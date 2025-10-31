package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Order represents a trading order
type Order struct {
	ID     string  `json:"id"`
	Symbol string  `json:"symbol"`
	Side   string  `json:"side"`
	Price  float64 `json:"price"`
}

var orders = []Order{
	{"ORD001", "AAPL", "BUY", 190.0},
	{"ORD002", "GOOG", "SELL", 125.5},
}

func main() {
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	})

	log.Println("Order service running on :8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}
