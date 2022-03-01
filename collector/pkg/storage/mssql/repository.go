/*
For now I have main left as a comment, as it allows for easy testing.
Will need to change function names etc, but for now we can see that a connection
to the SQL Server works.
To start the connection, call 'databaseConnection'.
*/
package mssql

//package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

var DB_CONNECTION *sql.DB

// Database connection variables.
var server = "127.0.0.1"
var port = 1433
var user = "sa"
var password = "Password1_HOLDER"
var database = "MetricsCollectorDB"

type Storage struct {
	DB_CONNECTION *sql.DB
}

type Process struct {
	PID           int
	collectorID   int
	name          string
	status        string
	cpuUsage      float32
	memoryUsage   float32
	diskUsage     float32
	executionTime time.Time
}

type Collector struct {
	collectorID   int
	timeCollected time.Time
	cpuID         int
	memoryID      int
	diskID        int
}

type Cpu struct {
	cpuID        int
	usage        float32
	availability float32
}

func openDBConnection() {

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	var err error

	// Create connection pool
	DB_CONNECTION, err = sql.Open("sqlserver", connString)

	if err != nil {

		log.Fatal("Error creating connection pool: ", err.Error())
	}

	ctx := context.Background()
	err = DB_CONNECTION.PingContext(ctx)

	if err != nil {

		log.Fatal(err.Error())
	}

	// Log connection here!
	fmt.Printf("Connected to DB!\n")
}

func closeDBConnection() {

	DB_CONNECTION.Close()
}

func getAllCPU() []Cpu {

	ctx := context.Background()

	// For not we are just getting from the CPU table!
	singleQuery := fmt.Sprintf("SELECT * FROM CPU;")

	// Execute query
	rows, err := DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	var toReturn []Cpu

	// Iterate through the result set.
	for rows.Next() {

		var usage, availability float32
		var cpuID int

		// Get values from row.
		err := rows.Scan(&cpuID, &usage, &availability)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := Cpu{cpuID, usage, availability}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

func main() {

	fmt.Printf("Repository Implementation for mssql (Microsoft SQL Server)\n")

	// To start the connection, call 'databaseConnection'.
	openDBConnection()

	/*
		var answer []Cpu = getAllCPU()

		for _, cpu := range answer {

			fmt.Printf("cpuID: %d, usage: %f, availability: %f\n", cpu.cpuID, cpu.usage, cpu.availability)
		}
	*/

	// For now I am closing it manually.
	// Not sure if we want it to stay open......
	closeDBConnection()

	fmt.Print("Look, im done!")
}
