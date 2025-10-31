package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// Alert represents an alert message produced by analytics.
type Alert struct {
	VehicleID string    `json:"vehicle_id"`
	Level     string    `json:"level"`   // INFO / WARN / CRITICAL
	Message   string    `json:"message"` // human message
	Ts        time.Time `json:"ts"`
	Source    string    `json:"source,omitempty"`
}

// Config via environment vars (12-factor)
var (
	redisAddr   = getenv("REDIS_ADDR", "redis:6379")
	redisPubSub = getenv("REDIS_ALERT_CHANNEL", "alerts")
	httpAddr    = getenv("HTTP_ADDR", ":8083")
	maxRecent   = getenvInt("MAX_RECENT_ALERTS", 200)
	pingPeriod  = 30 * time.Second
	writeWait   = 10 * time.Second
	readWait    = 60 * time.Second
)

// helper env functions
func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func getenvInt(k string, d int) int {
	if s := os.Getenv(k); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			return n
		}
	}
	return d
}

// RecentAlerts is a simple goroutine-safe ring/circular buffer
type RecentAlerts struct {
	mu    sync.RWMutex
	buf   []Alert
	cap   int
	start int
	len   int
}

func NewRecentAlerts(cap int) *RecentAlerts {
	return &RecentAlerts{buf: make([]Alert, cap), cap: cap}
}

func (r *RecentAlerts) Add(a Alert) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cap == 0 {
		return
	}
	if r.len < r.cap {
		r.buf[(r.start+r.len)%r.cap] = a
		r.len++
		return
	}
	// overwrite oldest
	r.buf[r.start] = a
	r.start = (r.start + 1) % r.cap
}

func (r *RecentAlerts) List() []Alert {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Alert, 0, r.len)
	for i := 0; i < r.len; i++ {
		idx := (r.start + i) % r.cap
		out = append(out, r.buf[idx])
	}
	return out
}

// WebSocket client management
type Client struct {
	conn *websocket.Conn
	send chan Alert
}

type Hub struct {
	clients    map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
	broadcast  chan Alert
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Alert, 256),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = struct{}{}
			h.mu.Unlock()
		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
			h.mu.Unlock()
		case a := <-h.broadcast:
			h.mu.RLock()
			for c := range h.clients {
				// non-blocking push to client
				select {
				case c.send <- a:
				default:
					// client send buffer full -> close it
					delete(h.clients, c)
					close(c.send)
				}
			}
			h.mu.RUnlock()
		case <-ctx.Done():
			// shutdown: close all clients
			h.mu.Lock()
			for c := range h.clients {
				close(c.send)
			}
			h.clients = nil
			h.mu.Unlock()
			return
		}
	}
}

// upgrader for websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// allow all origins for demo; in prod restrict origins
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWs(hub *Hub, recent *RecentAlerts, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ws] upgrade error: %v\n", err)
		return
	}
	client := &Client{conn: conn, send: make(chan Alert, 128)}
	hub.register <- client

	// send recent alerts on connect
	for _, a := range recent.List() {
		client.send <- a
	}

	// start writer and reader goroutines
	go wsWriter(client)
	go wsReader(client, hub)
}

func wsReader(c *Client, hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	_ = c.conn.SetReadDeadline(time.Now().Add(readWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(readWait)); return nil })
	for {
		// we only expect pings/pongs or client messages (ignored) — read to detect close
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
}

func wsWriter(c *Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case a, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// channel closed
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(a); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// NotificationService encapsulates Redis subscription and HTTP/WebSocket APIs
type NotificationService struct {
	rdb        *redis.Client
	recent     *RecentAlerts
	hub        *Hub
	logger     *log.Logger
	pubsub     *redis.PubSub
	cancelSub  context.CancelFunc
	subRunning chan struct{}
}

func NewNotificationService(rdb *redis.Client, logger *log.Logger) *NotificationService {
	return &NotificationService{
		rdb:        rdb,
		recent:     NewRecentAlerts(maxRecent),
		hub:        NewHub(),
		logger:     logger,
		subRunning: make(chan struct{}),
	}
}

// start subscription to Redis channel 'alerts'
func (n *NotificationService) StartSubscription(ctx context.Context, channel string) error {
	ctxSub, cancel := context.WithCancel(ctx)
	n.cancelSub = cancel
	ps := n.rdb.Subscribe(ctxSub, channel)
	n.pubsub = ps

	// wait for subscription to be ready
	_, err := ps.Receive(ctxSub)
	if err != nil {
		return err
	}

	go func() {
		defer close(n.subRunning)
		ch := ps.Channel()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					n.logger.Println("[sub] channel closed")
					return
				}
				// parse alert
				var a Alert
				if err := json.Unmarshal([]byte(msg.Payload), &a); err != nil {
					n.logger.Printf("[sub] invalid alert payload: %v\n", err)
					continue
				}
				// ensure timestamp
				if a.Ts.IsZero() {
					a.Ts = time.Now().UTC()
				}
				// 1) store in recent buffer
				n.recent.Add(a)
				// 2) publish to websocket clients
				select {
				case n.hub.broadcast <- a:
				default:
					// drop if hub busy
					n.logger.Println("[sub] hub busy, dropped alert broadcast")
				}
			case <-ctxSub.Done():
				n.logger.Println("[sub] stopping subscription")
				_ = ps.Close()
				return
			}
		}
	}()

	// run hub
	go n.hub.Run(ctxSub)
	return nil
}

func (n *NotificationService) Stop() {
	if n.cancelSub != nil {
		n.cancelSub()
		<-n.subRunning // wait for subscription loop to end
	}
}

// http handlers

// push test alert (also publishes to Redis channel so consumers see same)
func (n *NotificationService) handlePushAlert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var a Alert
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "invalid payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	if a.Ts.IsZero() {
		a.Ts = time.Now().UTC()
	}
	// publish to redis channel
	pb, _ := json.Marshal(a)
	ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
	defer cancel()
	if err := n.rdb.Publish(ctx, redisPubSub, pb).Err(); err != nil {
		n.logger.Printf("[push] redis publish err: %v", err)
		http.Error(w, "publish failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// list recent alerts
func (n *NotificationService) handleListAlerts(w http.ResponseWriter, r *http.Request) {
	list := n.recent.List()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

func (n *NotificationService) ServeHTTP(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/alerts", n.handleListAlerts)       // GET -> list
	mux.HandleFunc("/alert", n.handlePushAlert)         // POST -> push test alert
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { serveWs(n.hub, n.recent, w, r) }) // WS endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
		defer cancel()
		s := map[string]interface{}{"status": "ok", "redis": false}
		if err := n.rdb.Ping(ctx).Err(); err == nil {
			s["redis"] = true
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(s)
	})
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	// graceful shutdown handled by caller
	return server.ListenAndServe()
}

func main() {
	logger := log.New(os.Stdout, "[notification] ", log.LstdFlags|log.Lmsgprefix)

	// Redis client
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Fatalf("redis ping failed: %v", err)
	}

	// build service
	ns := NewNotificationService(rdb, logger)

	// start subscription
	rootCtx, rootCancel := context.WithCancel(context.Background())
	if err := ns.StartSubscription(rootCtx, redisPubSub); err != nil {
		logger.Fatalf("start subscription failed: %v", err)
	}
	logger.Println("subscription started on channel:", redisPubSub)

	// HTTP server in goroutine
	srvErr := make(chan error, 1)
	go func() {
		logger.Printf("http service listening on %s\n", httpAddr)
		srvErr <- ns.ServeHTTP(httpAddr)
	}()

	// graceful shutdown on signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	select {
	case sig := <-sigCh:
		logger.Printf("signal received: %v, shutting down", sig)
	case err := <-srvErr:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("http server error: %v", err)
		}
	}

	// stop subscription & hub
	rootCancel()
	ns.Stop()
	logger.Println("notification service stopped")
}

