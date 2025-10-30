package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// library packages
// lib methiods
// lib creates channel context
// in own goroutine we can cancel lib channel context
// main goroutine wait for lib channel context done
// main goroutine select for lib channel context done
// return value and error condition
// data tranformation

func Library() {

	// created channel context channel -> new Done channel
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	fmt.Println("library code with context:", ctx)

	go func() {
		fmt.Println("hi i am smiling")
		cancel() // stop
		fmt.Println("hi i am india no 1 coders")
	}()

	time.Sleep(5)

	select {

	case <-time.After(5 * time.Second):

		fmt.Println("Done")

	case <-ctx.Done():

		fmt.Println("main stopped", ctx.Err())

	}
}

// goroutine
// channel
// waitgroup mutex
// loop with select
// code execution

func goroutines(id int, jobs chan string, results chan string) {

	fmt.Println("having fun with channel")

	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		results <- fmt.Sprintf("worker %d finished job %s", id, j)
	}
}

//  resouce -mysql
//  connect to db -> Library
//  query execution -> Library
//  data mapping -> type strut
//  data tranformation ->type struct
//  logics --> Loop copy data of database in to simple variable or struct
//  error handling --> many error handling

func resources() {

	db, err := sql.Open("sqlite", "file:test.db?cache=shared&mode=memory")

	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()
	fmt.Println("Database connected successfully")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}
	fmt.Println("Table created successfully")

	_, err = db.Exec("INSERT INTO users (name) VALUES (?)", "Alice")
	if err != nil {
		fmt.Println("Error inserting data:", err)
		return
	}
	fmt.Println("Data inserted successfully")

	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		fmt.Println("Error querying data:", err)
		return
	}
	defer rows.Close()

	fmt.Println("Users:")
	for rows.Next() {

		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			fmt.Println("Error scanning data:", err)
			return
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println("Error with rows:", err)
		return
	}
}

// desiging
// struct
// interface methods
// composition

func main() {

	// Flow of code
	/*
		         1. two channel  jobs and results
				 2. three goroutines
				 3. three jobs sent to jobs channel

				 4. goroutine receive the jobs is done by goroutines
				 5. copy received job in to results channel

				 6. three results collected from results channel
	*/

	fmt.Println("calling monster!")

	jobs := make(chan string, 3)
	results := make(chan string, 3)

	for w := 1; w <= 3; w++ { // creating 3 goroutines
		go goroutines(w, jobs, results)
	}

	for j := 1; j <= 3; j++ { // sending 3 data in to  jobs
		jobs <- fmt.Sprintf("job%d", j)
	}

	close(jobs)

	for a := 1; a <= 3; a++ { // collecting 3 results from results channel
		result := <-results
		fmt.Println("Received:", result)
	}
	fmt.Println(" monster end!")
	Library()
	resources()

}
