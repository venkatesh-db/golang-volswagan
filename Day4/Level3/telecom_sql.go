
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Plan struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Cost  float64 `json:"cost"`
}

type Customer struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Plans []Plan  `json:"plans"`
}

func main() {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/telecomdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT c.id, c.name, p.id, p.name, p.cost
		FROM customers c
		JOIN plans p ON c.id = p.customer_id
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	customerMap := make(map[int]*Customer)
	for rows.Next() {
		var cID, pID int
		var cName, pName string
		var cost float64

		if err := rows.Scan(&cID, &cName, &pID, &pName, &cost); err != nil {
			log.Fatal(err)
		}

		if _, exists := customerMap[cID]; !exists {
			customerMap[cID] = &Customer{ID: cID, Name: cName}
		}
		customerMap[cID].Plans = append(customerMap[cID].Plans, Plan{ID: pID, Name: pName, Cost: cost})
	}

	customers := []Customer{}
	for _, c := range customerMap {
		customers = append(customers, *c)
	}

	output, _ := json.MarshalIndent(customers, "", "  ")
	fmt.Println(string(output))
}
