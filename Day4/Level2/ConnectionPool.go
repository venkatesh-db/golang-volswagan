package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Order struct {
	ID     int
	Symbol string
	Side   string
	Price  float64
}

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/tradingdb"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB Connection failed:", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Insert order
	insertStmt, err := db.Prepare("INSERT INTO orders(symbol, side, price) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer insertStmt.Close()

	_, err = insertStmt.Exec("AAPL", "BUY", 189.45)
	if err != nil {
		log.Fatal("Insert failed:", err)
	}

	// Select multiple orders
	rows, err := db.Query("SELECT id, symbol, side, price FROM orders")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.Symbol, &o.Side, &o.Price); err != nil {
			log.Fatal(err)
		}
		orders = append(orders, o)
	}
	fmt.Println("Orders:", orders)
}
