//	For now I have *main* left as a comment, as it allows for easy testing.
//	As of 3/24/2022, there are custom searches (based on table name) and inserts for
//	MEMORY/DISk tables.
//	There are fully custom (tableName, column, field) for the PROCESS table.
//	CPU has it's own functions as it only holds usage now.
//	Use the BULK insert Function to insert everything together.

//package mssql

package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/collecting"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/cpu"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/disk"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/logger"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/memory"
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

type Collector struct {
	collectorID   int
	timeCollected time.Time
	cpuID         int
	memoryID      int
	diskID        int
}

type IndividualComponent struct {
	usage float64
	size  float64
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

		logger.Error("Error creating connection pool: " + err.Error())
		return nil, err
	}

	ctx := context.Background()
	err = DB_CONNECTION.PingContext(ctx)

	if err != nil {

		logger.Error(err.Error())
		return nil, err
	}

	s := new(Storage)
	s.DB_CONNECTION = DB_CONNECTION

	// Log connection here!
	logger.Info("Connected to DB!\n")
	return s, err
}

//  Close a single database connection.
//  Doesn't need anything, but a connection should be open before
//  this is called.
func (s *Storage) CloseDBConnection() {

	s.DB_CONNECTION.Close()

	logger.Info("Closed DB connection.")
}

// ----------------------------- CPU Section Section -----------------------------

//  Get all CPUs from the CPU Table.
//  Doesn't need anything, it just cycles through each cpu in the table.
//
//  Return:
//  	([]cpu.CPU) all CPUs.
func (s *Storage) GetCPUs() []cpu.CPU {

	ctx := context.Background()

	// For not we are just getting from the CPU table!
	singleQuery := fmt.Sprintf("SELECT usage FROM CPU;")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

	var toReturn []cpu.CPU

	// Iterate through the result set.
	for rows.Next() {

		var usage float32

		// Get values from row.
		err := rows.Scan(&usage)

		if err != nil {
			logger.Error(err.Error())
		}

		singleInput := cpu.CPU{Usage: usage}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

//  Get newest CPU ID from the CPU Table
//  Nothing needs to be passed, just call te function.
//
//  Return:
//  	(int) single ID from the CPU table.
func (s *Storage) GetNewestCPUID() int {

	ctx := context.Background()

	// For not we are just getting from the single table!
	singleQuery := fmt.Sprintf("SELECT TOP 1 cpuID FROM CPU ORDER BY cpuID DESC;")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

	var toReturn int

	// Iterate through the result set.
	for rows.Next() {

		var id int

		// Get values from row.
		err := rows.Scan(&id)

		if err != nil {

			logger.Error(err.Error())
		}

		toReturn = id
	}

	return toReturn
}

//  Get newest from CPU
//  Nothing needs to be passed, just call te function.
//
//  Return:
//  	(cpu.CPU) single from the CPU table.
func (s *Storage) GetNewestCPU() cpu.CPU {

	ctx := context.Background()

	// For not we are just getting from the single table!
	singleQuery := fmt.Sprintf("SELECT usage FROM CPU WHERE cpuID IN " +
		"(SELECT TOP 1 cpuID FROM COLLECTOR ORDER BY timeCollected DESC);")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

	var toReturn cpu.CPU

	// Iterate through the result set.
	for rows.Next() {

		var usage float32

		// Get values from row.
		err := rows.Scan(&usage)

		if err != nil {

			logger.Error(err.Error())
		}

		toReturn = cpu.CPU{Usage: usage}
	}

	return toReturn
}

//  Get a single CPU from the CPU Table based off it's ID.
//  Only needs the ID that is being searched for.
//
//  Return:
//  	([]cpu.CPU) CPUs.
func (s *Storage) GetCPUByID(cpuID int) []cpu.CPU {

	ctx := context.Background()

	// For not we are just getting from the CPU table!
	singleQuery := fmt.Sprintf("SELECT usage FROM CPU WHERE cpuID = %d;", cpuID)

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

	var toReturn []cpu.CPU

	// Iterate through the result set.
	for rows.Next() {

		var usage float32

		// Get values from row.
		err := rows.Scan(&usage)

		if err != nil {

			logger.Error(err.Error())
		}

		singleInput := cpu.CPU{Usage: usage}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

//  Insert into CPU.
//  Takes in the data to be inserted.
//
//  Return:
//  	(int) rows inserted.
//  	(error) any error, this should be 'nil'.
func (s *Storage) PutNewCPU(singleInput cpu.CPU) (int64, error) {

	// Insert into a single component.
	singleInsert :=
		fmt.Sprintf("INSERT INTO CPU VALUES (%.2f);", singleInput.Usage)

	// Execute Insertion
	result, err := s.DB_CONNECTION.Exec(singleInsert)

	if err != nil {

		logger.Error(err.Error())
	}

	return result.RowsAffected()
}

// ----------------------------- MEMORY Section Section -----------------------------

//  Get all memories from MEMORY Table.
//  Doesn't need anything, it just cycles through each memory in the table.
//
//  Return:
//  	([]IndividualComponent) all Memories.
func (s *Storage) GetMemories() []IndividualComponent {

	ctx := context.Background()

	// For not we are just getting from the MEMORY table!
	singleQuery := fmt.Sprintf("SELECT usage, size FROM MEMORY;")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

	var toReturn []IndividualComponent

	// Iterate through the result set.
	for rows.Next() {

		var usage, size float64

		// Get values from row.
		err := rows.Scan(&usage, &size)

		if err != nil {

			logger.Error(err.Error())
		}

		singleInput := IndividualComponent{usage, size}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

// ----------------------------- DISK Section Section -----------------------------

//  Get all disks from DISK Table.
//  Doesn't need anything, it just cycles through each disk in the table.
//
//  Return:
//  	([]IndividualComponent) all disks.
func (s *Storage) GetDisks() []IndividualComponent {

	ctx := context.Background()

	// For not we are just getting from the DISK table!
	singleQuery := fmt.Sprintf("SELECT usage, size FROM DISK;")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

	var toReturn []IndividualComponent

	// Iterate through the result set.
	for rows.Next() {

		var usage, size float64

		// Get values from row.
		err := rows.Scan(&usage, &size)

		if err != nil {

			logger.Error(err.Error())
		}

		singleInput := IndividualComponent{usage, size}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

// ------------------------ INDIVIDUAL COMPONENT Section -----------------------

//  Get all from either MEMORY/DISK or _AVERAGE's.
//  You need to give it the name of what table you want to get all from.
//	You can only use this method for MEMORY/DISK and their _AVERAGE tables!
//
//  Return:
//  	([]IndividualComponent) all from one of the 4 tables.
func (s *Storage) GetIndivComponents(tableName string) []IndividualComponent {

	tableName = strings.ToUpper(tableName)

	ctx := context.Background()

	var singleQuery string

	if strings.Contains(tableName, "_AVERAGE") {

		// _AVEGAGE tables
		singleQuery = fmt.Sprintf("SELECT averageUsage, averageSize FROM %s;", tableName)
	} else {

		// Normal tables
		singleQuery = fmt.Sprintf("SELECT usage, size FROM %s;", tableName)
	}

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

	var toReturn []IndividualComponent

	// Iterate through the result set.
	for rows.Next() {

		var usage, size float64

		// Get values from row.
		err := rows.Scan(&usage, &size)

		if err != nil {

			logger.Error(err.Error())
		}

		singleInput := IndividualComponent{usage, size}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

//  Get newest from either MEMORY/DISK and _AVERAGE.
//  You need to give it the name of what table you want to get all from.
//	You can only use this method for MEMORY/DISK and their _AVERAGE tables!
//
//  Return:
//  	(IndividualComponent) single from one of the 4 tables.
func (s *Storage) GetNewestIndivComponent(tableName string) IndividualComponent {

	tableName = strings.ToUpper(tableName)
	var IdName string
	var singleQuery string

	if strings.Contains(tableName, "_AVERAGE") {

		if tableName == "MEMORY_AVERAGE" {

			IdName = "memoryAverageID"
		} else {

			IdName = "diskAverageID"
		}

		// For not we are just getting from the single table!
		singleQuery = fmt.Sprintf("SELECT averageUsage, averageSize FROM %s WHERE %s IN "+
			"(SELECT TOP 1 %s FROM COLLECTOR_HISTORY ORDER BY collectorHistoryID DESC);",
			tableName, IdName, IdName)

	} else {

		if tableName == "MEMORY" {

			IdName = "memoryID"
		} else {

			IdName = "diskID"
		}

		// For not we are just getting from the single table!
		singleQuery = fmt.Sprintf("SELECT usage, size FROM %s WHERE %s IN "+
			"(SELECT TOP 1 %s FROM COLLECTOR ORDER BY timeCollected DESC);",
			tableName, IdName, IdName)
	}

	ctx := context.Background()

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

	var toReturn IndividualComponent

	// Iterate through the result set.
	for rows.Next() {

		var usage, size float64

		// Get values from row.
		err := rows.Scan(&usage, &size)

		if err != nil {

			logger.Error(err.Error())
		}

		toReturn = IndividualComponent{usage, size}
	}

	return toReturn
}

//  Get newest ID from either MEMORY/DISK and _AVERAGE.
//  You need to give it the name of what table you want to get all from.
//	You can only use this method for MEMORY/DISK and their _AVERAGE tables!
//
//  Return:
//  	(int) single ID from one of the 4 tables.
func (s *Storage) GetNewestIndivComponentID(tableName string) int {

	tableName = strings.ToUpper(tableName)
	var IdName string

	if tableName == "MEMORY" {

		IdName = "memoryID"
	} else if tableName == "MEMORY_AVERAGE" {

		IdName = "memoryAverageID"
	} else if tableName == "DISK_AVERAGE" {

		IdName = "diskAverageID"
	} else if tableName == "DISK" {

		IdName = "diskID"
	}

	ctx := context.Background()

	// For not we are just getting from the single table!
	singleQuery := fmt.Sprintf("SELECT TOP 1 %s FROM %s ORDER BY %s DESC;",
		IdName, tableName, IdName)

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

	var toReturn int

	// Iterate through the result set.
	for rows.Next() {

		var id int

		// Get values from row.
		err := rows.Scan(&id)

		if err != nil {

			logger.Error(err.Error())
		}

		toReturn = id
	}

	return toReturn
}

//  Insert for either MEMORY/DISK.
//  Takes in a table name, and the data to be inserted.
//
//  Return:
//  	(int) rows inserted.
//  	(error) any error, this should be 'nil'.
func (s *Storage) PutNewSingleComponent(
	tableName string, singleInput IndividualComponent) (int64, error) {

	tableName = strings.ToUpper(tableName)

	// Insert into a single component.
	singleInsert :=
		fmt.Sprintf("INSERT INTO %s VALUES (%f, %f);",
			tableName, singleInput.usage, singleInput.size)

	// Execute Insertion
	result, err := s.DB_CONNECTION.Exec(singleInsert)

	if err != nil {

		logger.Error(err.Error())
	}

	return result.RowsAffected()
}

// ----------------------------- COLLECTOR Section -----------------------------

//  Get all Collectors from COLLECTOR Table.
//  Doesn't need anything, it just cycles through each collector in the table.
//
//  Return:
//  	([]Collector) all collectors.
func (s *Storage) GetCollectors() []Collector {

	ctx := context.Background()

	// Get all Collectors.
	singleQuery := fmt.Sprintf("SELECT collectorID, timeCollected, cpuID, " +
		"memoryID, diskID FROM COLLECTOR;")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

	var toReturn []Collector

	// Iterate through the result set.
	for rows.Next() {

		var collectorID, cpuID, memoryID, diskID int
		var timeCollected time.Time

		// Get values from row.
		err := rows.Scan(&collectorID, &timeCollected, &cpuID, &memoryID, &diskID)

		if err != nil {

			logger.Error(err.Error())
		}

		singleInput := Collector{collectorID, timeCollected, cpuID, memoryID, diskID}
		toReturn = append(toReturn, singleInput)
	}

	return toReturn
}

//  Get newest collector from COLLECTOR table.
//  Doesn't need anything, just call it to get the newest collector.
//
//  Return:
//  	(Collector) COLLECTOR.
func (s *Storage) GetCollectorNewest() Collector {

	ctx := context.Background()

	// Get newsest Processes, based off collectorID.
	singleQuery :=
		fmt.Sprintf("SELECT TOP 1 collectorID, timeCollected, cpuID, memoryID, diskID" +
			" FROM COLLECTOR ORDER BY timeCollected DESC;")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

	var toReturn Collector

	// Iterate through the result set.
	for rows.Next() {

		var collectorID, cpuID, memoryID, diskID int
		var timeCollected time.Time

		// Get values from row.
		err := rows.Scan(&collectorID, &timeCollected, &cpuID, &memoryID, &diskID)

		if err != nil {

			logger.Error(err.Error())
		}

		toReturn = Collector{collectorID: collectorID, timeCollected: timeCollected,
			cpuID: cpuID, memoryID: memoryID, diskID: diskID}
	}

	return toReturn
}

//  Get newest collector's ID from COLLECTOR table.
//  Doesn't need anything, just call it to get the newest ID.
//
//  Return:
//  	(int) collectorID.
func (s *Storage) GetCollectorIDNewest() int {

	ctx := context.Background()

	// Get newsest Processes, based off collectorID.
	singleQuery :=
		fmt.Sprintf("SELECT TOP 1 collectorID FROM COLLECTOR " +
			"ORDER BY timeCollected DESC;")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

	var toReturnInt int

	// Iterate through the result set.
	for rows.Next() {

		var collectorID int

		// Get values from row.
		err := rows.Scan(&collectorID)

		if err != nil {

			logger.Error(err.Error())
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

	var cpuID = s.GetNewestCPUID()
	var memoryID = s.GetNewestIndivComponentID("MEMORY")
	var diskID = s.GetNewestIndivComponentID("DISK")

	// Insert into Collector.
	// CPU/MEMORY/DISK will be up later.
	singleInsert :=
		fmt.Sprintf("INSERT INTO COLLECTOR VALUES (GETDATE(), %v, %v, %v);",
			cpuID, memoryID, diskID)

	// Execute Insertion
	result, err := s.DB_CONNECTION.Exec(singleInsert)

	if err != nil {

		logger.Error(err.Error())
	}

	return result.RowsAffected()
}

// ----------------------------- PROCESS Section -----------------------------

//  Get all Processes from PROCESS Table.
//  Doesn't need anything, it just cycles through each process in the table.
//
//  Return:
//  	([]Process) all processes.
func (s *Storage) GetProcesses() []process.Process {

	ctx := context.Background()

	// Get all Processes.
	singleQuery := fmt.Sprintf("SELECT processID, collectorID, PID, name, status," +
		" cpuUsage, memoryUsage, diskUsage, executionTime FROM PROCESS;")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

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

			logger.Error(err.Error())
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

	ctx := context.Background()

	// Get newsest Processes, based off collectorID.
	singleQuery := fmt.Sprintf("SELECT processID, collectorID, PID, name, status," +
		" cpuUsage, memoryUsage, diskUsage, executionTime FROM PROCESS" +
		" WHERE collectorID IN " +
		"(SELECT TOP 1 collectorID FROM COLLECTOR ORDER BY timeCollected DESC) " +
		"ORDER BY processID DESC;")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

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

			logger.Error(err.Error())
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

	ctx := context.Background()

	// Get processes based custom column and string field.
	singleQuery := fmt.Sprintf("SELECT processID, collectorID, PID, name, status,"+
		" cpuUsage, memoryUsage, diskUsage, executionTime FROM PROCESS"+
		" WHERE %s = '%s';", column, field)

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

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

			logger.Error(err.Error())
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

	ctx := context.Background()

	// Get processes based custom column and float field.
	singleQuery := fmt.Sprintf("SELECT processID, collectorID, PID, name, status,"+
		" cpuUsage, memoryUsage, diskUsage, executionTime FROM PROCESS"+
		" WHERE %s = %.2f;", column, field)

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

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

			logger.Error(err.Error())
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

	ctx := context.Background()

	// Get processes based custom column and int field.
	singleQuery := fmt.Sprintf("SELECT processID, collectorID, PID, name, status,"+
		" cpuUsage, memoryUsage, diskUsage, executionTime FROM PROCESS"+
		" WHERE %s = %d;", column, field)

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

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

			logger.Error(err.Error())
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

	var collectorID = s.GetCollectorIDNewest()
	var repeatedProcess = false

	// Test if this is already inserted in the DB for this collector time
	processes := s.GetProcessesByNewest()

	for _, singleCheckProcess := range processes {

		if singleCheckProcess == singleProcess {

			logger.Info("Tried (skipped) inserting the same process, PID: " + fmt.Sprint(singleProcess.PID))
			repeatedProcess = true
		}
	}

	var err error

	if !(repeatedProcess) {
		// Insert into PROCESS based of singleProcess Data.
		singleInsert :=
			fmt.Sprintf("INSERT INTO PROCESS VALUES (%v, %v, '%v', '%v', %.2f, %.2f, "+
				"%.2f, %.2f);", collectorID, singleProcess.PID, singleProcess.Name,
				singleProcess.Status, singleProcess.CPUUtilization, singleProcess.RAMUtilization,
				singleProcess.DiskUtilization, singleProcess.ExecutionTime)

		// Execute Insertion
		result, err := s.DB_CONNECTION.Exec(singleInsert)

		if err != nil {

			logger.Error(err.Error())
		}

		return result.RowsAffected()
	}

	return 0, err
}

// ------------------- Latest HISTORY Section -------------------

//  Get History data from COLLECTOR_HISTORY/CPU_AVERAGE/MEMORY_AVERAGE.
//  Nothing needs to be passed, as the query gets the data needed.
//
//  Return:
//  	(History) single History data structure.
func (s *Storage) GetNewestHistory() collecting.History {

	ctx := context.Background()

	// Get processes based custom column and float field.
	singleQuery := fmt.Sprintf("SELECT TOP 1 CL.timeCollectedStart, " +
		"CL.timeCollectedEnd, CA.averageUsage, MA.averageUsage, MA.averageSize " +
		"FROM COLLECTOR_HISTORY AS CL, CPU_AVERAGE AS CA, MEMORY_AVERAGE AS MA " +
		"WHERE CA.cpuAverageID = CL.cpuAverageID AND " +
		"MA.memoryAverageID = CL.memoryAverageID ORDER BY collectorHistoryID DESC;")

	// Execute query
	rows, err := s.DB_CONNECTION.QueryContext(ctx, singleQuery)

	if err != nil {

		logger.Error(err.Error())
	}

	defer rows.Close()

	var toReturn collecting.History

	// Iterate through the result set.
	for rows.Next() {

		var timeCollectedStart, timeCollectedEnd time.Time
		var averageCpuUsage, averageMemUsage, averageMemSize float64

		// Get values from row.
		err := rows.Scan(&timeCollectedStart, &timeCollectedEnd, &averageCpuUsage,
			&averageMemUsage, &averageMemSize)

		if err != nil {

			logger.Error(err.Error())
		}

		singleInput := collecting.History{Start: timeCollectedStart, End: timeCollectedEnd,
			AverageCpuUsage: averageCpuUsage, AverageMemUsage: averageMemUsage,
			AverageMemSize: averageMemSize}

		toReturn = singleInput
	}

	return toReturn
}

// ------------------- BULK INSERT Section -------------------

//  Insert for All the tables
//  Takes in the Metrics struct, which should hold all the data needed to
//	be inserted into the database (CPU, MEMORY, DISK, and PROCESS).
//  Purge stored procedure will be called from here as well.
// 	The procedure checks if it needs to move/purge data, and if it does it will.
//  The rest of the insertion will wait for it to complete.
//
//  Return:
//  	(bool) true if an error occurred.
func (s *Storage) BulkInsert(totalMetrics collecting.Metrics) bool {

	errorHappened := false

	// Purge Stored Procedure.

	singleInsert := fmt.Sprintf("EXEC PRC_PurgeData;")

	// Execute Insertion
	result, err := s.DB_CONNECTION.Exec(singleInsert)

	affectedRows, err := result.RowsAffected()

	if err != nil {

		logger.Error("Error PurgeData: " + err.Error())
	}

	if !(affectedRows >= 1) {

		logger.Error("Error PurgeData.")
	}

	// Insert into CPU/MEMORY/DISK
	rowsAffected, err := s.PutNewCPU(totalMetrics.CPU)

	if err != nil {

		logger.Error("Error in adding in CPU Table -- Bulk Insert Function: " + err.Error())
		errorHappened = true
	}
	if !(rowsAffected >= 1) {

		logger.Error("Error in adding in CPU Table -- Bulk Insert Function.\n")
		errorHappened = true
	}

	memoryHolder := IndividualComponent{usage: totalMetrics.Memory.Usage, size: totalMetrics.Memory.Size}

	rowsAffected, err = s.PutNewSingleComponent("MEMORY", memoryHolder)
	if err != nil {

		logger.Error("Error in adding in MEMORY Table" +
			" -- Bulk Insert Function:" + err.Error())
		errorHappened = true
	}
	if !(rowsAffected >= 1) {

		logger.Error("Error in adding in MEMORY Table" +
			" -- Bulk Insert Function.\n")
		errorHappened = true
	}

	var singleDiskHolder disk.Disk

	for _, singleDisk := range totalMetrics.Disk {

		if singleDisk.MountPoint == "/" {

			singleDiskHolder = singleDisk
		}
	}

	diskHolder := IndividualComponent{usage: singleDiskHolder.Usage, size: singleDiskHolder.Size}

	rowsAffected, err = s.PutNewSingleComponent("DISK", diskHolder)
	if err != nil {

		logger.Error("Error in adding in DISK Table" +
			" -- Bulk Insert Function: " + err.Error())
		errorHappened = true
	}
	if !(rowsAffected >= 1) {

		logger.Error("Error in adding in DISK Table" +
			" -- Bulk Insert Function.\n")
		errorHappened = true
	}

	// Insert into Collector
	rowsAffected, err = s.PutNewCollector()
	if err != nil {

		logger.Error("Error in adding in COLLECTOR Table" +
			"-- Bulk Insert Function: " + err.Error())
		errorHappened = true
	}
	if !(rowsAffected >= 1) {

		logger.Error("Error in adding in COLLECTOR Table" +
			" -- Bulk Insert Function.\n")
		errorHappened = true
	}

	// Insert into PROCESS
	for iteration, singleProcess := range totalMetrics.Processes {

		rowsAffected, err = s.PutNewProcess(singleProcess)
		if err != nil {

			logger.Error("Error in adding in PROCESS Table" +
				" -- Bulk Insert Function, iteration: " + fmt.Sprint(iteration) +
				", error: " + err.Error())
			errorHappened = true
		}
		if !(rowsAffected >= 1) {

			logger.Error("Error in adding in PROCESS Table" +
				" -- Bulk Insert Function, iteration: " + fmt.Sprint(iteration))
			errorHappened = true
		}
	}

	if !errorHappened {

		logger.Info("BULK insertion completed successfully.")
	}

	return errorHappened
}

// ------------------- Testing Section -------------------

//
func main() {

	fmt.Printf("Repository Implementation for mssql (Microsoft SQL Server)\n")

	logger.Init()

	var database, err = NewStorage()

	if err != nil {
	}

	//newestCollector := database.GetCollectorNewest()

	//fmt.Printf("collector: %v, %v\n", newestCollector.collectorID, newestCollector.timeCollected)

	cpuHolder1 := cpu.CPU{Usage: 10.10}
	memoryHolder := memory.Memory{Usage: 10.10, Size: 1000.00}

	diskHolder := []disk.Disk{}

	diskHolder = append(diskHolder, disk.Disk{Name: "NA", MountPoint: "hehehe", Usage: 40.00, Size: 40.00})
	diskHolder = append(diskHolder, disk.Disk{Name: "/", MountPoint: "/", Usage: 60.00, Size: 1245.67})

	listProcess := []process.Process{}

	listProcess = append(listProcess, process.Process{PID: 5540, Name: "process0", CPUUtilization: 00.00, RAMUtilization: 00.00, DiskUtilization: 00.00, Status: "done", ExecutionTime: 00.00})
	listProcess = append(listProcess, process.Process{PID: 999, Name: "process1", CPUUtilization: 01.00, RAMUtilization: 01.00, DiskUtilization: 00.10, Status: "running", ExecutionTime: 01.00})
	listProcess = append(listProcess, process.Process{PID: 666, Name: "process2", CPUUtilization: 02.20, RAMUtilization: 22.00, DiskUtilization: 00.22, Status: "failed", ExecutionTime: 22.00})
	listProcess = append(listProcess, process.Process{PID: 999, Name: "process1", CPUUtilization: 00.00, RAMUtilization: 00.00, DiskUtilization: 00.10, Status: "done", ExecutionTime: 00.00})

	metricsHolder1 := collecting.Metrics{Processes: listProcess, CPU: cpuHolder1, Memory: memoryHolder, Disk: diskHolder}

	database.BulkInsert(metricsHolder1)

	// Test Get Individual Components
	/*var answer []IndividualComponent = database.GetIndivComponents("MEMORY")

	for _, singleComponent := range answer {

		fmt.Printf("usage: %.2f, size: %.2f\n",
			singleComponent.usage, singleComponent.size)
	}*/

	// Test Get Newest Individual Component
	/*var answer IndividualComponent = database.GetNewestIndivComponent("DISK")

	fmt.Printf("usage: %.2f, size: %.2f\n",
		answer.usage, answer.size)
	*/
	// Test Get Newest Individual Component
	/*var id = database.GetNewestIndivComponentID("DISK")

	fmt.Printf("ID: %v", id)
	*/

	database.CloseDBConnection()

	//database.PutNewProcess(process.Process{PID: 5540, Name: "process0", CPUUtilization: 00.00, RAMUtilization: 00.00, DiskUtilization: 00.00, Status: "done", ExecutionTime: 00.00})

	// To start the connection, call 'databaseConnection'.

	// Test CPUs Get
	/*var answer []CPU = database.GetCPUs()

	for _, cpu := range answer {

		fmt.Printf("usage: %.2f\n", cpu.usage)
	}*/

	// Test Newest CPUs Get
	/*var answer []CPU = database.GetNewestCPU()

	for _, cpu := range answer {

		fmt.Printf("usage: %.2f\n", cpu.usage)
	}*/

	// Test CPU Get by ID
	/*var answer []CPU = database.GetCPUByID(1)

	for _, cpu := range answer {

		fmt.Printf("usage: %.2f\n", cpu.usage)
	}*/

	// Test CPU Put single
	/*var holderProcess = CPU{usage: 11.11}

	var rowsInsertedCount, error1 = database.PutNewCPU(holderProcess)

	fmt.Printf("rowsInsertedCount: %d ", rowsInsertedCount)
	fmt.Println(error1)
	*/

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

	// Test Collectors all
	// Dont run this, as we arent using ints for CPUID etc..
	/*var answer []Collector = GetCollectors()
	for _, collector := range answer {

		fmt.Printf("collectorID: %d, time: %t, CPUID: %d,  memoryID: %d, diskID: %d\n",
			collector.collectorID, collector.timeCollected.Day(), collector.cpuID, collector.memoryID, collector.diskID)
	}
	*/

	// Test Get Individual Components
	/*var answer []IndividualComponent = database.GetIndivComponents("MEMORY")

	for _, singleComponent := range answer {

		fmt.Printf("usage: %.2f, size: %.2f\n",
			singleComponent.usage, singleComponent.size)
	}*/

	// Test Get Newest Individual Component
	/*var answer []IndividualComponent = database.GetNewestIndivComponent("DISK")

	for _, singleComponent := range answer {

		fmt.Printf("usage: %.2f, size: %.2f\n",
			singleComponent.usage, singleComponent.size)
	}
	*/

	fmt.Print("DONE TEST")
}
