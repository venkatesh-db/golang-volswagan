package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Trade struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // default port
	}

	http.HandleFunc("/trade", func(w http.ResponseWriter, r *http.Request) {
		trade := Trade{Symbol: "AAPL", Price: 189.5}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trade)
	})

	log.Printf("Trading microservice running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}


