
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Vehicle struct {
	ID        int       `json:"id"`
	Model     string    `json:"model"`
	Location  string    `json:"location"`
	Timestamp time.Time `json:"timestamp"`
	Sensors   []Sensor  `json:"sensors"`
}

type Sensor struct {
	ID          int     `json:"id"`
	VehicleID   int     `json:"vehicle_id"`
	Temperature float64 `json:"temperature"`
	Speed       float64 `json:"speed"`
}

func main() {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/iotfleet")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(15)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Update sensor data (simulation)
	updateStmt, _ := db.Prepare("UPDATE sensors SET temperature=?, speed=? WHERE id=?")
	defer updateStmt.Close()
	updateStmt.Exec(32.8, 90.5, 1)

	rows, err := db.Query(`
		SELECT v.id, v.model, v.location, v.timestamp, 
		       s.id, s.temperature, s.speed
		FROM vehicles v
		JOIN sensors s ON v.id = s.vehicle_id
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fleetMap := make(map[int]*Vehicle)
	for rows.Next() {
		var vID, sID int
		var model, location string
		var ts time.Time
		var temp, speed float64

		if err := rows.Scan(&vID, &model, &location, &ts, &sID, &temp, &speed); err != nil {
			log.Fatal(err)
		}

		if _, exists := fleetMap[vID]; !exists {
			fleetMap[vID] = &Vehicle{ID: vID, Model: model, Location: location, Timestamp: ts}
		}
		fleetMap[vID].Sensors = append(fleetMap[vID].Sensors, Sensor{
			ID: sID, VehicleID: vID, Temperature: temp, Speed: speed,
		})
	}

	var vehicles []Vehicle
	for _, v := range fleetMap {
		vehicles = append(vehicles, *v)
	}

	data, _ := json.MarshalIndent(vehicles, "", "  ")
	fmt.Println(string(data))
}

