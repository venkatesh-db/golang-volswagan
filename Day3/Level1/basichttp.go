package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "IRCTC Server Running: All trains operational ✅")
	})

	http.HandleFunc("/train", func(w http.ResponseWriter, r *http.Request) {
		train := r.URL.Query().Get("name")
		if train == "" {
			http.Error(w, "Train name missing", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Train %s — On Time 🚆", train)
	})

	fmt.Println("Starting IRCTC HTTP Server on :8080 ...")
	http.ListenAndServe(":8080", nil)
}
