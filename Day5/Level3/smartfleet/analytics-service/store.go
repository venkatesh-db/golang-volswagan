package analyticsservice


import (
	"context"
	"encoding/json"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Store is the DB wrapper.
type Store struct {
	db     *gorm.DB
	logger *log.Logger
}

// NewStore constructs a Store.
func NewStore(db *gorm.DB, logger *log.Logger) *Store {
	return &Store{db: db, logger: logger}
}

// InsertTelemetry persists raw telemetry.
func (s *Store) InsertTelemetry(ctx context.Context, ev TelemetryEvent) error {
	raw := TelemetryRaw{
		VehicleID: ev.VehicleID,
		Timestamp: time.UnixMilli(ev.Ts).UTC(),
		Speed:     ev.Speed,
		Fuel:      ev.FuelLevel,
		Latitude:  ev.Lat,
		Longitude: ev.Lon,
	}
	if err := s.db.WithContext(ctx).Create(&raw).Error; err != nil {
		s.logger.Printf("InsertTelemetry err: %v", err)
		return err
	}
	return nil
}

// UpsertAggregate updates per-minute aggregates using DB upsert.
func (s *Store) UpsertAggregate(ctx context.Context, ev TelemetryEvent) error {
	ts := time.UnixMilli(ev.Ts).UTC()
	bucket := bucketMinute(ts)

	// For incoming single event: we upsert with ON CONFLICT combining counts and computing simple aggregates
	agg := Aggregate{
		VehicleID:  ev.VehicleID,
		Bucket:     bucket,
		AvgSpeed:   ev.Speed,
		MinFuel:    ev.FuelLevel,
		MaxSpeed:   ev.Speed,
		EventCount: 1,
	}

	// On conflict: update aggregate fields:
	// new_avg = (avg*count + ev.Speed) / (count+1)
	// min_fuel = least(min_fuel, ev.FuelLevel), max_speed = greatest(max_speed, ev.Speed), event_count = event_count+1
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "vehicle_id"}, {Name: "bucket"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"avg_speed":   clause.Expr("((coalesce(avg_speed,0) * coalesce(event_count,0) + ?) / (coalesce(event_count,0) + 1))", ev.Speed),
			"min_fuel":    clause.Expr("LEAST(coalesce(min_fuel,9999), ?)", ev.FuelLevel),
			"max_speed":   clause.Expr("GREATEST(coalesce(max_speed,0), ?)", ev.Speed),
			"event_count": clause.Expr("coalesce(event_count,0) + 1"),
			"updated_at":  time.Now(),
		}),
	}).Create(&agg).Error
}

// SaveOrUpdateTrip persists trip start/end and updates distance/avg speed.
func (s *Store) SaveOrUpdateTrip(ctx context.Context, trip *Trip) error {
	// If trip has ID (existing), update; else create.
	if trip.ID == 0 {
		return s.db.WithContext(ctx).Create(trip).Error
	}
	return s.db.WithContext(ctx).Save(trip).Error
}

// GetActiveTrip fetches an active (no EndedAt) trip for vehicle
func (s *Store) GetActiveTrip(ctx context.Context, vehicleID string) (*Trip, error) {
	var t Trip
	err := s.db.WithContext(ctx).Where("vehicle_id = ? AND ended_at IS NULL", vehicleID).Order("started_at DESC").First(&t).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// For debug: dump aggregates for vehicle
func (s *Store) DumpAggregates(ctx context.Context, vehicleID string) ([]Aggregate, error) {
	var out []Aggregate
	if err := s.db.WithContext(ctx).Where("vehicle_id = ?", vehicleID).Order("bucket DESC").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

// helper to log JSON structure (for debugging)
func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
