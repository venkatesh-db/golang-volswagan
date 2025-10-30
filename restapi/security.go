package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// curl -H "Authorization: Bearer SECRET123" http://localhost:9090/secure

// http://localhost:9090/sessionid

var (
	sessions = map[string]time.Time{}
	mu       sync.Mutex // to prevent concurrent map writes
)

func secureapi(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(map[string]string{"message": "access is granted"})
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")

		if token != "Bearer SECRET123" {
			http.Error(res, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func main() {
	http.Handle("/secure", Middleware(http.HandlerFunc(secureapi)))

	http.HandleFunc("/sessionid", func(res http.ResponseWriter, req *http.Request) {
		sessionId := fmt.Sprintf("guest-%d", time.Now().UnixNano())

		mu.Lock()
		sessions[sessionId] = time.Now()
		mu.Unlock()

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(map[string]string{"session_id": sessionId})
	})

	fmt.Println("✅ Server is running on port 9090...")
	http.ListenAndServe(":9090", nil)
}
