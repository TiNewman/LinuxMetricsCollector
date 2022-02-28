/*
For now I have main left as a comment, as it allows for easy testing.
Will need to change function names etc, but for now we can see that a connection
to the SQL Server works.
*/
package mssql

//package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var db *sql.DB

// Database connection variables.
var server = "127.0.0.1"
var port = 1433
var user = "sa"
var password = "Password1_HOLDER"
var database = "MetricsCollectorDB"

func main() {

	fmt.Printf("Repository Implementation for mssql (Microsoft SQL Server)\n")

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	var err error

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)

	if err != nil {

		log.Fatal("Error creating connection pool: ", err.Error())
	}

	ctx := context.Background()
	err = db.PingContext(ctx)

	if err != nil {

		log.Fatal(err.Error())
	}

	fmt.Printf("Connected to DB!\n")

	// For not we are just getting from the CPU table!
	tsql := fmt.Sprintf("SELECT * FROM CPU;")

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	var count int

	// Iterate through the result set.
	for rows.Next() {

		var usage, availability string
		var cpuID int

		// Get values from row.
		err := rows.Scan(&cpuID, &usage, &availability)

		if err != nil {

			log.Fatal(err.Error())
		}

		fmt.Printf("cpuID: %d, usage: %s, availability: %s\n", cpuID, usage, availability)
		count++
	}

	// For now I am closing it manually.
	// Not sure if we want it to stay open......
	db.Close()

	fmt.Print("Look, im done!")
}
