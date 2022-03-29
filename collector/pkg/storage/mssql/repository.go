/*
For now I have main left as a comment, as it allows for easy testing.
These functions will be exported.
As of 3/5/2022, the current functions mainly revolve around the Process Table.
More functions will be added for CPU/MEMORY/DISK tables.
There is insert for both COLLECTOR and PROCESS tables.
Will look into fully custom queries:
(tableName, column, field)...
*/

package mssql

//package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/cpu"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/process"
	_ "github.com/denisenkom/go-mssqldb"
)

// Database connection variables.
var server = "127.0.0.1"
var port = 1433
var user = "sa"
var password = "Password1_HOLDER"
var database = "MetricsCollectorDB"

type Storage struct {
	DB_CONNECTION *sql.DB
}

type DBProcess struct {
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

// Type for inserting into a Collector, as we don't need time or collectorID.
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

// ----------------------------- Connecting to Database Section -----------------------------

//  Opens a single database connection.
//  Doesn't need anything.
func NewStorage() (*Storage, error) {

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	var DB_CONNECTION *sql.DB
	var err error

	// Create connection pool
	DB_CONNECTION, err = sql.Open("sqlserver", connString)

	if err != nil {

		log.Fatal("Error creating connection pool: ", err.Error())
		return nil, err
	}

	ctx := context.Background()
	err = DB_CONNECTION.PingContext(ctx)

	if err != nil {

		log.Fatal(err.Error())
		return nil, err
	}

	s := new(Storage)
	s.DB_CONNECTION = DB_CONNECTION

	// Log connection here!
	fmt.Printf("Connected to DB!\n")
	return s, err
}

//  Close a single database connection.
//  Doesn't need anything, but a connection should be open before
//  this is called.
func (s *Storage) CloseDBConnection() {

	s.DB_CONNECTION.Close()
}

// ----------------------------- GPU Section Section -----------------------------

//  Get all GPUs from GPU Table.
//  Doesn't need anything, it just cycles through each gpu in the table.
//
//  Return:
//  	([]Cpu) all current CPUs.
func (s *Storage) GetCPUs() []cpu.CPU {

	//OpenDBConnection()

	ctx := context.Background()

	// For not we are just getting from the CPU table!
	singleQuery := fmt.Sprintf("SELECT * FROM CPU;")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	//CloseDBConnection()

	var toReturn []cpu.CPU

	// Iterate through the result set.
	for rows.Next() {

		var usage, availability float32
		var cpuID int

		// Get values from row.
		err := rows.Scan(&cpuID, &usage, &availability)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := cpu.CPU{Usage: usage}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

// ----------------------------- COLLECTOR Section -----------------------------

//  Get all Collectors from COLLECTOR Table.
//  !DONT USE THIS UNTIL WE ACTUALLY GET CPU/MEMORY/DISK tables running.!
//  Doesn't need anything, it just cycles through each collector in the table.
//
//  Return:
//  	([]Collector) all collectors.
/*func GetCollectors() []Collector {

	OpenDBConnection()

	ctx := context.Background()

	// Get all Collectors.
	singleQuery := fmt.Sprintf("SELECT * FROM COLLECTOR;")

	// Execute query
	rows, err := DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	CloseDBConnection()

	var toReturn []Collector

	// Iterate through the result set.
	for rows.Next() {

		var collectorID, cpuID, memoryID, diskID int
		var timeCollected time.Time

		// Get values from row.
		err := rows.Scan(&collectorID, &timeCollected, &cpuID, &memoryID, &diskID)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := Collector{collectorID, timeCollected, cpuID, memoryID, diskID}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}
*/

//  Get newest collector's ID from COLLECTOR table.
//  Doesn't need anything, just call it to get the newest ID.
//
//  Return:
//  	(int) collectorID.
func (s *Storage) GetCollectorIDNewest() int {

	//OpenDBConnection()

	ctx := context.Background()

	// Get newsest Processes, based off collectorID.
	singleQuery :=
		fmt.Sprintf("SELECT TOP 1 collectorID FROM COLLECTOR " +
			" ORDER BY timeCollected DESC;")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	//CloseDBConnection()

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

//  Insert for COLLECTOR Table.
//  Takes in a Collector, and uses its data to insert into the table.
//
//  Return:
//  	(int) rows inserted.
//  	(error) any error, this should be 'nil'.
func (s *Storage) PutNewCollector() (int64, error) {

	//OpenDBConnection()

	// These will be used once we get to CPU/MEMORY/DISK tables.
	// var cpuID = getCPUIDNewest()
	// var memoryID = getMemoryIDNewest()
	// var diskID = getDiskIDNewest()

	// Insert into Collector.
	// For now we only care about creating a timestamp and having a collectorID
	// for the PROCESS table.
	// CPU/MEMORY/DISK will be up later.
	singleInsert :=
		fmt.Sprintf("INSERT INTO COLLECTOR VALUES (GETDATE(), NULL, NULL, NULL);")

	/*
		// This will be used once we actually have to input CPU, etc..
			singleInsert :=
				fmt.Sprint("INSERT INTO COLLECT VALUES (GETDATE(), %d, %d, %d);",
				singleCollector.cpuID, singleCollector.memoryID, singleCollector.diskID)
	*/

	// Execute Insertion
	result, err := s.DB_CONNECTION.Exec(singleInsert)

	if err != nil {

		log.Fatal(err.Error())
	}

	//CloseDBConnection()

	return result.RowsAffected()
}

// ----------------------------- PROCESS Section -----------------------------

//  Get all Processes from PROCESS Table.
//  Doesn't need anything, it just cycles through each process in the table.
//
//  Return:
//  	([]Process) all processes.
func (s *Storage) GetProcesses() []process.Process {

	//OpenDBConnection()

	ctx := context.Background()

	// Get all Processes.
	singleQuery := fmt.Sprintf("SELECT processID, collectorID, PID, name, status," +
		" cpuUsage, memoryUsage, diskUsage, executionTime FROM PROCESS;")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	//CloseDBConnection()

	var toReturn []process.Process

	// Iterate through the result set.
	for rows.Next() {

		var processID, PID, collectorID int
		var name, status string
		var cpuUsage, memoryUsage, diskUsage, executionTime float32

		// Get values from row.
		err := rows.Scan(&processID, &collectorID, &PID, &name, &status, &cpuUsage,
			&memoryUsage, &diskUsage, &executionTime)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := process.Process{PID: PID, Name: name,
			CPUUtilization: cpuUsage, RAMUtilization: memoryUsage,
			DiskUtilization: diskUsage, Status: status, ExecutionTime: executionTime}

		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

//  Get all new Processes from PROCESS Table.
//  Doesn't need anything, it goes based off of the newest collectorID
//  which is taken from the COLLECTOR table.
//
//  Return:
//  	([]Process) newsest processes.
func (s *Storage) GetProcessesByNewest() []process.Process {

	//OpenDBConnection()

	ctx := context.Background()

	// Get newsest Processes, based off collectorID.
	singleQuery := fmt.Sprintf("SELECT processID, collectorID, PID, name, status," +
		" cpuUsage, memoryUsage, diskUsage, executionTime FROM PROCESS" +
		" WHERE collectorID IN " +
		"(SELECT TOP 1 collectorID FROM COLLECTOR ORDER BY timeCollected DESC);")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	//CloseDBConnection()

	var toReturn []process.Process

	// Iterate through the result set.
	for rows.Next() {

		var processID, PID, collectorID int
		var name, status string
		var cpuUsage, memoryUsage, diskUsage, executionTime float32

		// Get values from row.
		err := rows.Scan(&processID, &collectorID, &PID, &name, &status, &cpuUsage,
			&memoryUsage, &diskUsage, &executionTime)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := process.Process{PID: PID, Name: name,
			CPUUtilization: cpuUsage, RAMUtilization: memoryUsage,
			DiskUtilization: diskUsage, Status: status, ExecutionTime: executionTime}

		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

//  Get custom string searched Processes from PROCESS Table.
//  Given a column name, test it against a string field in the PROCESS table.
//	This will only work when searching columns that use 'string'/VARCHAR.
//
//  Return:
//  	([]Process) custom processes.
func (s *Storage) GetProcessesByCustomStringField(column string, field string) []process.Process {

	//OpenDBConnection()

	ctx := context.Background()

	// Get processes based custom column and string field.
	singleQuery := fmt.Sprintf("SELECT processID, collectorID, PID, name, status,"+
		" cpuUsage, memoryUsage, diskUsage, executionTime FROM PROCESS"+
		" WHERE %s = '%s';", column, field)

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	//CloseDBConnection()

	var toReturn []process.Process

	// Iterate through the result set.
	for rows.Next() {

		var processID, PID, collectorID int
		var name, status string
		var cpuUsage, memoryUsage, diskUsage, executionTime float32

		// Get values from row.
		err := rows.Scan(&processID, &collectorID, &PID, &name, &status, &cpuUsage,
			&memoryUsage, &diskUsage, &executionTime)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := process.Process{PID: PID, Name: name,
			CPUUtilization: cpuUsage, RAMUtilization: memoryUsage,
			DiskUtilization: diskUsage, Status: status, ExecutionTime: executionTime}

		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

//  Get custom float searched Processes from PROCESS Table.
//  Given a column name, test it against a float field in the PROCESS table.
// 	This will only work when searching columns that use float.
//
//  Return:
//  	([]Process) custom processes.
func (s *Storage) GetProcessesByCustomFloatField(column string, field float32) []process.Process {

	//OpenDBConnection()

	ctx := context.Background()

	// Get processes based custom column and float field.
	singleQuery := fmt.Sprintf("SELECT processID, collectorID, PID, name, status,"+
		" cpuUsage, memoryUsage, diskUsage, executionTime FROM PROCESS"+
		" WHERE %s = %.2f;", column, field)

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	//CloseDBConnection()

	var toReturn []process.Process

	// Iterate through the result set.
	for rows.Next() {

		var processID, PID, collectorID int
		var name, status string
		var cpuUsage, memoryUsage, diskUsage, executionTime float32

		// Get values from row.
		err := rows.Scan(&processID, &collectorID, &PID, &name, &status, &cpuUsage,
			&memoryUsage, &diskUsage, &executionTime)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := process.Process{PID: PID, Name: name,
			CPUUtilization: cpuUsage, RAMUtilization: memoryUsage,
			DiskUtilization: diskUsage, Status: status, ExecutionTime: executionTime}

		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

//  Get custom Integer searched Processes from PROCESS Table.
//  Given a column name, test it against an integer field in the PROCESS table.
//	This will only work when searching columns that use int/BIG INT.
//
//  Return:
//  	([]Process) custom processes.
func (s *Storage) GetProcessesByCustomIntField(column string, field int) []process.Process {

	//OpenDBConnection()

	ctx := context.Background()

	// Get processes based custom column and int field.
	singleQuery := fmt.Sprintf("SELECT processID, collectorID, PID, name, status,"+
		" cpuUsage, memoryUsage, diskUsage, executionTime FROM PROCESS"+
		" WHERE %s = %d;", column, field)

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		log.Fatal(err.Error())
	}

	defer rows.Close()

	//CloseDBConnection()

	var toReturn []process.Process

	// Iterate through the result set.
	for rows.Next() {

		var processID, PID, collectorID int
		var name, status string
		var cpuUsage, memoryUsage, diskUsage, executionTime float32

		// Get values from row.
		err := rows.Scan(&processID, &collectorID, &PID, &name, &status, &cpuUsage,
			&memoryUsage, &diskUsage, &executionTime)

		if err != nil {

			log.Fatal(err.Error())
		}

		singleInput := process.Process{PID: PID, Name: name,
			CPUUtilization: cpuUsage, RAMUtilization: memoryUsage,
			DiskUtilization: diskUsage, Status: status, ExecutionTime: executionTime}

		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

//  Insert for PROCESS Table
//  Takes in a Process, then checks for the newest collector,
//  and uses that collectorID (as you have to insert into collector first)
//  with the data in the Process to insert into the PROCESS table.
//
//  Return:
//  	(int) rows inserted.
//  	(error) any error, this should be 'nil'.
func (s *Storage) PutNewProcess(singleProcess process.Process) (int64, error) {

	//OpenDBConnection()

	var collectorID = s.GetCollectorIDNewest()

	// Insert into PROCESS based of singleProcess Data.
	singleInsert :=
		fmt.Sprintf("INSERT INTO PROCESS VALUES (%d, %d, '%s', '%s', %.2f, %.2f, "+
			"%.2f, %.2f);", collectorID, singleProcess.PID, singleProcess.Name,
			singleProcess.Status, singleProcess.CPUUtilization, singleProcess.RAMUtilization,
			singleProcess.DiskUtilization, singleProcess.ExecutionTime)

	// Execute Insertion
	result, err := s.DB_CONNECTION.Exec(singleInsert)

	if err != nil {

		log.Fatal(err.Error())
	}

	//CloseDBConnection()

	return result.RowsAffected()
}

// ------------------- Testing Section -------------------

//
func main() {

	fmt.Printf("Repository Implementation for mssql (Microsoft SQL Server)\n")

	/*var database, err = NewStorage()

	if err != nil {
	}
	*/

	// To start the connection, call 'databaseConnection'.
	//OpenDBConnection()

	// Test CPUs Get
	/*var answer []Cpu = getCPUs()

	for _, cpu := range answer {

		fmt.Printf("cpuID: %d, usage: %f, availability: %f\n", cpu.cpuID, cpu.usage, cpu.availability)
	}*/

	// Test Processes Get
	/*var answer []process.Process = database.GetProcesses()

	for _, process := range answer {

		fmt.Printf("PID: %d,  name: %s, status: %s, cpuUsage: %.2f, memoryUsage: %.2f, diskUsage: %.2f, executionTime: %.2f\n",
			 process.PID, process.Name, process.Status, process.CPUUtilization, process.RAMUtilization, process.DiskUtilization, process.ExecutionTime)
	}*/

	// Test Processes Get by newest collector == newest processes
	/*var answer []process.Process = database.GetProcessesByNewest()

	for _, process := range answer {

		fmt.Printf("PID: %d,  name: %s, status: %s, cpuUsage: %.2f, memoryUsage: %.2f, diskUsage: %.2f, executionTime: %.2f\n",
			process.PID, process.Name, process.Status, process.CPUUtilization, process.RAMUtilization, process.DiskUtilization, process.ExecutionTime)
	}*/

	// Test Processes Get by custom column and a string filed
	/*var answer []process.Process = database.GetProcessesByCustomStringField("name", "process2")

	for _, process := range answer {

		fmt.Printf("PID: %d,  name: %s, status: %s, cpuUsage: %.2f, memoryUsage: %.2f, diskUsage: %.2f, executionTime: %.2f\n",
			process.PID, process.Name, process.Status, process.CPUUtilization, process.RAMUtilization, process.DiskUtilization, process.ExecutionTime)
	}*/

	// Test Processes Get by custom column and a float filed
	/*var answer []process.Process = database.GetProcessesByCustomFloatField("diskUsage", 99.99)

	for _, process := range answer {

		fmt.Printf("PID: %d,  name: %s, status: %s, cpuUsage: %.2f, memoryUsage: %.2f, diskUsage: %.2f, executionTime: %.2f\n",
			process.PID, process.Name, process.Status, process.CPUUtilization, process.RAMUtilization, process.DiskUtilization, process.ExecutionTime)
	}*/

	// Test Processes Get by custom column and a int filed
	/*var answer []process.Process = database.GetProcessesByCustomIntField("collectorID", 1)
	for _, process := range answer {

		fmt.Printf("PID: %d,  name: %s, status: %s, cpuUsage: %.2f, memoryUsage: %.2f, diskUsage: %.2f, executionTime: %.2f\n",
			process.PID, process.Name, process.Status, process.CPUUtilization, process.RAMUtilization, process.DiskUtilization, process.ExecutionTime)
	}*/

	// Test Processes Put single
	/*var holderProcess = process.Process{PID: 5540, Name: "process0", CPUUtilization: 00.00, RAMUtilization: 00.00, DiskUtilization: 00.00, Status: "done", ExecutionTime: 00.00}

	var rowsInsertedCount, error1 = database.PutNewProcess(holderProcess)

	fmt.Printf("rowsInsertedCount: %d ", rowsInsertedCount)
	fmt.Println(error1)*/

	// Test Collector Put single
	/*var holderCollector = CollectorInsert{0, 0, 0}

	var rowsInsertedCount, error1 = putNewCollector(holderCollector)

	fmt.Printf("rowsInsertedCount: %d ", rowsInsertedCount)
	fmt.Println(error1)
	*/

	// For now I am closing it manually.
	// Not sure if we want it to stay open......
	//CloseDBConnection()

	// Test Collectors all
	// Dont run this, as we arent using ints for CPUID etc..
	/*var answer []Collector = GetCollectors()
	for _, collector := range answer {

		fmt.Printf("collectorID: %d, time: %t, CPUID: %d,  memoryID: %d, diskID: %d\n",
			collector.collectorID, collector.timeCollected.Day(), collector.cpuID, collector.memoryID, collector.diskID)
	}
	*/

	fmt.Print("DONE TEST")
}
