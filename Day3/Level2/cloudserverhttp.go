package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type ServiceStatus struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	status := ServiceStatus{
		Name:      "CloudMonitor",
		Status:    "Healthy",
		Timestamp: time.Now(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func main() {
	http.HandleFunc("/status", statusHandler)
	http.ListenAndServe(":9090", nil)
}
