
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/bankdb"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB Connection Error:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	fmt.Println("✅ Connected to MySQL successfully!")

	// Query account count
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM accounts").Scan(&count)
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	fmt.Printf("Total Accounts: %d\n", count)
}
