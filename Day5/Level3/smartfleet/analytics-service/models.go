package analyticsservice

import (
	"time"

	"gorm.io/gorm"
)

// TelemetryRaw stores incoming telemetry events (persisted)
type TelemetryRaw struct {
	ID        uint      `gorm:"primaryKey"`
	VehicleID string    `gorm:"index:idx_vehicle_ts,priority:1;index"`
	Timestamp time.Time `gorm:"index:idx_vehicle_ts,priority:2"`
	Speed     float64
	Fuel      float64
	Latitude  float64
	Longitude float64
	CreatedAt time.Time
}

// Aggregate represents per-minute aggregate metrics per vehicle
type Aggregate struct {
	ID         uint      `gorm:"primaryKey"`
	VehicleID  string    `gorm:"index:idx_agg_vehicle_bucket,priority:1"`
	Bucket     time.Time `gorm:"index:idx_agg_vehicle_bucket,priority:2"` // bucket start (UTC minute)
	AvgSpeed   float64
	MinFuel    float64
	MaxSpeed   float64
	EventCount int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Trip summary for detected trips
type Trip struct {
	ID           uint      `gorm:"primaryKey"`
	VehicleID    string    `gorm:"index"`
	StartedAt    time.Time `gorm:"index"`
	EndedAt      *time.Time
	DistanceKm   float64
	AvgSpeedKmph float64
	EventCount   int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// MigrateSchemas runs AutoMigrate
func MigrateSchemas(db *gorm.DB) error {
	return db.AutoMigrate(&TelemetryRaw{}, &Aggregate{}, &Trip{})
}
