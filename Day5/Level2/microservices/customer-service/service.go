
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Customer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	http.HandleFunc("/customer", func(w http.ResponseWriter, r *http.Request) {
		customer := Customer{ID: 1, Name: "Ravi Bala"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customer)
	})

	log.Println("Customer service running on :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
