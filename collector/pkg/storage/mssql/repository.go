/*
For now I have main left as a comment, as it allows for easy testing.
Will need to change function names etc, but for now we can see that a connection
to the SQL Server works.
To start the connection, call 'databaseConnection'.
*/
//package mssql

package main

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
	processID     int
	collectorID   int
	PID           int
	name          string
	status        string
	cpuUsage      float32
	memoryUsage   float32
	diskUsage     float32
	executionTime float32
}

type Collector struct {
	collectorID   int
	timeCollected time.Time
	cpuID         int
	memoryID      int
	diskID        int
}

type CollectorInsert struct {
	cpuID    int
	memoryID int
	diskID   int
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

func getCPUs() []Cpu {

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

func getCollectorIDNewest() int {

	ctx := context.Background()

	// Get newsest Processes, based off collectorID.
	singleQuery :=
		fmt.Sprintf("SELECT TOP 1 collectorID FROM COLLECTOR " +
			" ORDER BY timeCollected DESC;")

	// Execute query
	rows, err := DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	var toReturnInt int

	// Iterate through the result set.
	for rows.Next() {

		var collectorID int

		// Get values from row.
		err := rows.Scan(&collectorID)

		if err != nil {

			log.Fatal(err.Error())
		}

		toReturnInt = collectorID
	}

	return toReturnInt
}

// Insert for COLLECTOR Table
// Takes in a Collector, and uses its data to insert into the table.
//
// Return:
//	(int) rows inserted.
//	(error) any error, this should be 'nil'.
func putNewCollector(singleCollector CollectorInsert) (int64, error) {

	// These will be used once we get to CPU/MEMORY/DISK tables.
	// var cpuID = getCPUIDNewest()
	// var memoryID = getMemoryIDNewest()
	// var diskID = getDiskIDNewest()

	// Insert into Collector.
	// For now we only care about creating a timestamp and having a collectorID
	// for the PROCESS table.
	// CPU/MEMORY/DISK will be up later.
	singleInsert :=
		fmt.Sprint("INSERT INTO COLLECTOR VALUES (GETDATE(), NULL, NULL, NULL);")

	/*
		// This will be used once we actually have to input CPU, etc..
			singleInsert :=
				fmt.Sprint("INSERT INTO COLLECT VALUES (GETDATE(), %d, %d, %d);",
				singleCollector.cpuID, singleCollector.memoryID, singleCollector.diskID)
	*/

	// Execute Insertion
	result, err := DB_CONNECTION.Exec(singleInsert)

	if err != nil {

		log.Fatal(err.Error())
	}

	return result.RowsAffected()
}

func getProcesses() []Process {

	ctx := context.Background()

	// Get all Processes.
	singleQuery := fmt.Sprintf("SELECT * FROM PROCESS;")

	// Execute query
	rows, err := DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	var toReturn []Process

	// Iterate through the result set.
	for rows.Next() {

		var processID, PID, collectorID int
		var name, status string
		var cpuUsage, memoryUsage, diskUsage, executionTime float32

		// Get values from row.
		err := rows.Scan(&processID, &PID, &collectorID, &name, &status, &cpuUsage,
			&memoryUsage, &diskUsage, &executionTime)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := Process{processID, PID, collectorID, name, status, cpuUsage,
			memoryUsage, diskUsage, executionTime}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

func getProcessesByCollector(collectorID int) []Process {

	ctx := context.Background()

	// Based off of a collectorID (which is where the timeStamp is held),
	// get PROCESSES.
	singleQuery := fmt.Sprintf("SELECT * FROM PROCESS WHERE collectorID = %d;", collectorID)

	// Execute query
	rows, err := DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	var toReturn []Process

	// Iterate through the result set.
	for rows.Next() {

		var processID, PID, collectorID int
		var name, status string
		var cpuUsage, memoryUsage, diskUsage, executionTime float32

		// Get values from row.
		err := rows.Scan(&processID, &PID, &collectorID, &name, &status, &cpuUsage,
			&memoryUsage, &diskUsage, &executionTime)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := Process{processID, PID, collectorID, name, status, cpuUsage,
			memoryUsage, diskUsage, executionTime}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

func getProcessesByNewest() []Process {

	ctx := context.Background()

	// Get newsest Processes, based off collectorID.
	singleQuery := fmt.Sprintf("SELECT * FROM PROCESS WHERE collectorID IN " +
		"(SELECT TOP 1 collectorID FROM COLLECTOR ORDER BY timeCollected DESC);")

	// Execute query
	rows, err := DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	var toReturn []Process

	// Iterate through the result set.
	for rows.Next() {

		var processID, PID, collectorID int
		var name, status string
		var cpuUsage, memoryUsage, diskUsage, executionTime float32

		// Get values from row.
		err := rows.Scan(&processID, &PID, &collectorID, &name, &status, &cpuUsage,
			&memoryUsage, &diskUsage, &executionTime)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := Process{processID, PID, collectorID, name, status, cpuUsage,
			memoryUsage, diskUsage, executionTime}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

func getProcessesByPID(PID int) []Process {

	ctx := context.Background()

	// Get processes based off PID.
	singleQuery := fmt.Sprintf("SELECT * FROM PROCESS WHERE PID = %d;", PID)

	// Execute query
	rows, err := DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	var toReturn []Process

	// Iterate through the result set.
	for rows.Next() {

		var processID, PID, collectorID int
		var name, status string
		var cpuUsage, memoryUsage, diskUsage, executionTime float32

		// Get values from row.
		err := rows.Scan(&processID, &PID, &collectorID, &name, &status, &cpuUsage,
			&memoryUsage, &diskUsage, &executionTime)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := Process{processID, PID, collectorID, name, status, cpuUsage,
			memoryUsage, diskUsage, executionTime}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

// Given a column name, test it against a string field.
func getProcessesByCustomStringField(column string, field string) []Process {

	ctx := context.Background()

	// Get processes based custom column and string field.
	singleQuery := fmt.Sprintf("SELECT * FROM PROCESS WHERE %s = '%s';", column, field)

	// Execute query
	rows, err := DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	var toReturn []Process

	// Iterate through the result set.
	for rows.Next() {

		var processID, PID, collectorID int
		var name, status string
		var cpuUsage, memoryUsage, diskUsage, executionTime float32

		// Get values from row.
		err := rows.Scan(&processID, &PID, &collectorID, &name, &status, &cpuUsage,
			&memoryUsage, &diskUsage, &executionTime)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := Process{processID, PID, collectorID, name, status, cpuUsage,
			memoryUsage, diskUsage, executionTime}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

// Given a column name, test it against a float field.
func getProcessesByCustomFloatField(column string, field float32) []Process {

	ctx := context.Background()

	// Get processes based custom column and float field.
	singleQuery := fmt.Sprintf("SELECT * FROM PROCESS WHERE %s = %.2f;", column, field)

	// Execute query
	rows, err := DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	var toReturn []Process

	// Iterate through the result set.
	for rows.Next() {

		var processID, PID, collectorID int
		var name, status string
		var cpuUsage, memoryUsage, diskUsage, executionTime float32

		// Get values from row.
		err := rows.Scan(&processID, &PID, &collectorID, &name, &status, &cpuUsage,
			&memoryUsage, &diskUsage, &executionTime)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := Process{processID, PID, collectorID, name, status, cpuUsage,
			memoryUsage, diskUsage, executionTime}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

func getProcessesByStatus(field string) []Process {

	ctx := context.Background()

	// Get processes based off status string.
	singleQuery := fmt.Sprintf("SELECT * FROM PROCESS WHERE status = '%s';", field)

	// Execute query
	rows, err := DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	var toReturn []Process

	// Iterate through the result set.
	for rows.Next() {

		var processID, PID, collectorID int
		var name, status string
		var cpuUsage, memoryUsage, diskUsage, executionTime float32

		// Get values from row.
		err := rows.Scan(&processID, &PID, &collectorID, &name, &status, &cpuUsage,
			&memoryUsage, &diskUsage, &executionTime)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := Process{processID, PID, collectorID, name, status, cpuUsage,
			memoryUsage, diskUsage, executionTime}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

// Insert for PROCESS Table
// Takes in a Process, then checks for the newest collector,
// and uses that collectorID (as you have to insert into collector first)
// with the data in the Process to insert into the PROCESS table.
//
// Return:
//	(int) rows inserted.
//	(error) any error, this should be 'nil'.
func putNewProcess(singleProcess Process) (int64, error) {

	var collectorID = getCollectorIDNewest()

	// Insert into PROCESS based of singleProcess Data.
	singleInsert :=
		fmt.Sprintf("INSERT INTO PROCESS VALUES (%d, %d, '%s', '%s', %.2f, %.2f, "+
			"%.2f, %.2f);", collectorID, singleProcess.PID, singleProcess.name,
			singleProcess.status, singleProcess.cpuUsage, singleProcess.memoryUsage,
			singleProcess.diskUsage, singleProcess.executionTime)

	// Execute Insertion
	result, err := DB_CONNECTION.Exec(singleInsert)

	if err != nil {

		log.Fatal(err.Error())
	}

	return result.RowsAffected()
}

func main() {

	fmt.Printf("Repository Implementation for mssql (Microsoft SQL Server)\n")

	// To start the connection, call 'databaseConnection'.
	openDBConnection()

	// Test CPUs Get
	/*var answer []Cpu = getCPUs()

	for _, cpu := range answer {

		fmt.Printf("cpuID: %d, usage: %f, availability: %f\n", cpu.cpuID, cpu.usage, cpu.availability)
	}*/

	// Test Processes Get
	/*var answer []Process = getProcesses()

	for _, process := range answer {

		fmt.Printf("processID: %d, collectorID: %d, PID: %d,  name: %s, status: %s, cpuUsage: %f, memoryUsage: %f, diskUsage: %f, executionTime: %f\n",
			process.processID, process.PID, process.collectorID, process.name, process.status, process.cpuUsage, process.memoryUsage, process.diskUsage, process.executionTime)
	}*/

	// Test Processes Get by Collector
	/*var answer []Process = getProcessesByCollector(1)

	for _, process := range answer {

		fmt.Printf("processID: %d, collectorID: %d, PID: %d,  name: %s, status: %s, cpuUsage: %f, memoryUsage: %f, diskUsage: %f, executionTime: %f\n",
			process.processID, process.PID, process.collectorID, process.name, process.status, process.cpuUsage, process.memoryUsage, process.diskUsage, process.executionTime)
	}*/

	// Test Processes Get by newest collector == newest processes
	/*var answer []Process = getProcessesByNewest()

	for _, process := range answer {

		fmt.Printf("processID: %d, collectorID: %d, PID: %d,  name: %s, status: %s, cpuUsage: %f, memoryUsage: %f, diskUsage: %f, executionTime: %f\n",
			process.processID, process.PID, process.collectorID, process.name, process.status, process.cpuUsage, process.memoryUsage, process.diskUsage, process.executionTime)
	}
	*/

	// Test Processes Get by PID
	/*var answer []Process = getProcessesByPID(6640)

	for _, process := range answer {

		fmt.Printf("processID: %d, collectorID: %d, PID: %d,  name: %s, status: %s, cpuUsage: %f, memoryUsage: %f, diskUsage: %f, executionTime: %f\n",
			process.processID, process.PID, process.collectorID, process.name, process.status, process.cpuUsage, process.memoryUsage, process.diskUsage, process.executionTime)
	}
	*/

	// Test Processes Get by custom column and a string filed
	/*var answer []Process = getProcessesByCustomStringField("name", "process2")

	for _, process := range answer {

		fmt.Printf("processID: %d, collectorID: %d, PID: %d,  name: %s, status: %s, cpuUsage: %f, memoryUsage: %f, diskUsage: %f, executionTime: %f\n",
			process.processID, process.PID, process.collectorID, process.name, process.status, process.cpuUsage, process.memoryUsage, process.diskUsage, process.executionTime)
	}
	*/

	// Test Processes Get by custom column and a string filed
	/*var answer []Process = getProcessesByStatus("done")

	for _, process := range answer {

		fmt.Printf("processID: %d, collectorID: %d, PID: %d,  name: %s, status: %s, cpuUsage: %f, memoryUsage: %f, diskUsage: %f, executionTime: %f\n",
			process.processID, process.PID, process.collectorID, process.name, process.status, process.cpuUsage, process.memoryUsage, process.diskUsage, process.executionTime)
	}
	*/

	// Test Processes Get by custom column and a float filed
	/*var answer []Process = getProcessesByCustomFloatField("diskUsage", 99.99)

	for _, process := range answer {

		fmt.Printf("processID: %d, collectorID: %d, PID: %d,  name: %s, status: %s, cpuUsage: %.2f, memoryUsage: %.2f, diskUsage: %.2f, executionTime: %.2f\n",
			process.processID, process.PID, process.collectorID, process.name, process.status, process.cpuUsage, process.memoryUsage, process.diskUsage, process.executionTime)
	}
	*/

	// Test Processes Put single
	/*var holderProcess = Process{0, 0, 5540, "process0", "done", 00.00, 00.00, 00.00, 00.00}

	var rowsInsertedCount, error1 = putNewProcess(holderProcess)

	fmt.Printf("rowsInsertedCount: %d ", rowsInsertedCount)
	fmt.Println(error1)
	*/

	// Test Collector Put single
	/*var holderCollector = CollectorInsert{0, 0, 0}

	var rowsInsertedCount, error1 = putNewCollector(holderCollector)

	fmt.Printf("rowsInsertedCount: %d ", rowsInsertedCount)
	fmt.Println(error1)
	*/

	// For now I am closing it manually.
	// Not sure if we want it to stay open......
	closeDBConnection()

	fmt.Print("DONE TEST")
}
