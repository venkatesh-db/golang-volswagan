package analyticsservice

package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

// TelemetryEvent matches producer payload
type TelemetryEvent struct {
	VehicleID string  `json:"vehicle_id"`
	Speed     float64 `json:"speed"`
	FuelLevel float64 `json:"fuel_level"`
	Lat       float64 `json:"latitude"`
	Lon       float64 `json:"longitude"`
	Ts        int64   `json:"ts"`
}

// TripState holds transient state for trip detection per vehicle
type TripState struct {
	LastLat     float64
	LastLon     float64
	LastTs      int64
	AccumDistKm float64
	AccumSpeed  float64
	EventCount  int64
	Moving      bool
	StartedAt   time.Time
	mutex       sync.Mutex
}

// TripStateMap stores per-vehicle states
type TripStateMap struct {
	m map[string]*TripState
	l sync.RWMutex
}

func NewTripStateMap() *TripStateMap {
	return &TripStateMap{m: make(map[string]*TripState)}
}

func (tsm *TripStateMap) GetOrCreate(v string) *TripState {
	tsm.l.Lock()
	defer tsm.l.Unlock()
	if s, ok := tsm.m[v]; ok {
		return s
	}
	s := &TripState{}
	tsm.m[v] = s
	return s
}

func (tsm *TripStateMap) Delete(v string) {
	tsm.l.Lock()
	defer tsm.l.Unlock()
	delete(tsm.m, v)
}

// consumer loop
func runConsumerLoop(ctx context.Context, reader *kafka.Reader, store *Store, logger *log.Logger) error {
	tripStates := NewTripStateMap()
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				// normal shutdown
				return nil
			}
			return err
		}
		var ev TelemetryEvent
		if err := json.Unmarshal(m.Value, &ev); err != nil {
			logger.Printf("invalid message: %v", err)
			continue
		}
		// process synchronously (for simplicity); for performance use worker pool
		if err := processTelemetryEvent(ctx, store, tripStates, ev, logger); err != nil {
			logger.Printf("processTelemetryEvent err: %v", err)
		}
	}
}

// processTelemetryEvent persists telemetry, updates aggregate, manages trip state
func processTelemetryEvent(ctx context.Context, store *Store, tsm *TripStateMap, ev TelemetryEvent, logger *log.Logger) error {
	// 1. persist raw telemetry
	if err := store.InsertTelemetry(ctx, ev); err != nil {
		logger.Printf("Insert telemetry err: %v", err)
		// continue to process aggregates/trips — don't return fatal
	}

	// 2. update per-minute aggregate
	if err := store.UpsertAggregate(ctx, ev); err != nil {
		logger.Printf("Upsert aggregate err: %v", err)
	}

	// 3. handle trip FSM
	ts := tsm.GetOrCreate(ev.VehicleID)
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	// compute dt and distance from last point
	var distKm float64
	if ts.LastTs > 0 {
		distKm = haversineKm(ts.LastLat, ts.LastLon, ev.Lat, ev.Lon)
	}
	// update accumulators
	ts.AccumDistKm += distKm
	ts.AccumSpeed += ev.Speed
	ts.EventCount++

	// trip detection thresholds
	const movingSpeedThreshold = 5.0       // km/h
	const tripEndIdleSeconds = 120         // if idle for 120s -> end trip
	now := time.UnixMilli(ev.Ts)
	isMoving := ev.Speed >= movingSpeedThreshold

	if !ts.Moving && isMoving {
		// trip started
		ts.Moving = true
		ts.StartedAt = time.UnixMilli(ev.Ts).UTC()
		// create Trip record
		trip := &Trip{
			VehicleID:  ev.VehicleID,
			StartedAt:  ts.StartedAt,
			DistanceKm: 0,
			EventCount: 0,
		}
		if err := store.SaveOrUpdateTrip(ctx, trip); err != nil {
			logger.Printf("create trip err: %v", err)
		}
	} else if ts.Moving && !isMoving {
		// possible trip end - check idle duration since last moving
		// if last event timestamp exists and gap exceeds threshold -> end
		if ts.LastTs > 0 {
			idleSeconds := (now - ts.LastTs) / 1000
			if idleSeconds >= tripEndIdleSeconds {
				// finish trip: fetch active trip, update stats and close
				active, err := store.GetActiveTrip(ctx, ev.VehicleID)
				if err != nil {
					logger.Printf("GetActiveTrip err: %v", err)
				} else if active != nil {
					active.DistanceKm = ts.AccumDistKm
					if ts.EventCount > 0 {
						active.AvgSpeedKmph = ts.AccumSpeed / float64(ts.EventCount)
					}
					active.EventCount = ts.EventCount
					nowT := time.UnixMilli(ev.Ts).UTC()
					active.EndedAt = &nowT
					if err := store.SaveOrUpdateTrip(ctx, active); err != nil {
						logger.Printf("close trip err: %v", err)
					}
				}
				// reset trip state
				ts.Moving = false
				ts.AccumDistKm = 0
				ts.AccumSpeed = 0
				ts.EventCount = 0
			}
		}
	}

	// update last point
	ts.LastLat = ev.Lat
	ts.LastLon = ev.Lon
	ts.LastTs = ev.Ts

	// 4. Quick alerts example (speed)
	if ev.Speed > 140 {
		store.Alert(ev.VehicleID, "CRITICAL", "Overspeed > 140 km/h")
	}

	return nil
}
