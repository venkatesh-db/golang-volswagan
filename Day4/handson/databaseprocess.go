
package main 

import (
		"database/sql"	
 	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

//  resouce -mysql
//  connect to db -> Library
//  query execution -> Library
//  data mapping -> type strut
//  data tranformation ->type struct
//  logics --> Loop copy data of database in to simple variable or struct
//  error handling --> many error handling




func main() {
	// ✅ Replace 'your_new_password' with your actual MySQL root password
	dsn := "root:Th36under!@tcp(127.0.0.1:3306)/testdb?tls=false"
	//dsn := "root:Th36under!@tcp(127.0.0.1:3306)/testdb"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// ✅ Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("✅ Database connection established and alive!")
}
