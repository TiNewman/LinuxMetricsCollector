//	Unable to test the following as DISK and MEMORY are not implemented yet:
//	Inserting only into CPU/MEMORY and COLLECTOR tables: Must fail successfully
//	Inserting only into CPU/DISK and COLLECTOR tables: Must fail successfully
//	Inserting only into MEMORY/DISK and COLLECTOR tables: Must fail successfully.

package mssql

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/collecting"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/cpu"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/process"
)

var globalcollectorID int = 0

var globalcpuID int = 0
var globalmemoryID int = 0
var globaldiskID int = 0

var globalprocess0ID int = 5540
var globalprocess1ID int = 422
var globalprocess2ID int = 643
var globalprocess3ID int = 8432

//	As We have not implemented RAM or MEMORY yet, these tests only work with:
//	CPU, COLLECTOR, and PROCESS.
//	This also cycles through the process 10 times, where on the 5th cycle,
//	the database is 'cleared' to simulate purging
//	(data is not moved to long term tables yet).
func TestAllCorrectInputs(t *testing.T) {

	var database, err = NewStorage()

	if err != nil {
	}

	var numberOfRounds int = 10

	for numberOfRounds != 0 {

		var cpuHolder1 cpu.CPU
		// var memoryHolder1 memory.MEMORY
		// var ramHolder1 ram.RAM

		if globalcollectorID != 0 {

			cpuHolder1 = cpu.CPU{Usage: (float32(globalcollectorID*10) + 0.01)}
			// memoryHolder1 := memory.MEMORY{}
			// ramHolder1 := ram.RAM{}
		}
		cpuHolder1 = cpu.CPU{Usage: 10.10}
		// memoryHolder1 := memory.MEMORY{}
		// ramHolder1 := ram.RAM{}

		collectorHolder1 := Collector{collectorID: globalcollectorID, cpuID: globalcpuID}
		/*, memoryID: memoryID, diskID: diskID*/

		globalcollectorID++
		globalcpuID++
		globalmemoryID++
		globaldiskID++

		listProcess := []process.Process{}
		listProcess = append(listProcess, process.Process{PID: globalprocess0ID, Name: "process0", CPUUtilization: 54.00, RAMUtilization: 50.00, DiskUtilization: 00.00, Status: "done", ExecutionTime: 00.00})
		listProcess = append(listProcess, process.Process{PID: globalprocess1ID, Name: "process1", CPUUtilization: 01.00, RAMUtilization: 01.00, DiskUtilization: 00.10, Status: "running", ExecutionTime: 01.00})
		listProcess = append(listProcess, process.Process{PID: globalprocess2ID, Name: "process2", CPUUtilization: 02.20, RAMUtilization: 22.00, DiskUtilization: 77.22, Status: "failed", ExecutionTime: 22.00})
		listProcess = append(listProcess, process.Process{PID: globalprocess3ID, Name: "process3", CPUUtilization: 70.00, RAMUtilization: 00.00, DiskUtilization: 00.10, Status: "done", ExecutionTime: 99.99})

		globalprocess0ID++
		globalprocess1ID++
		globalprocess2ID++
		globalprocess3ID++

		metricsHolder1 := collecting.Metrics{Processes: listProcess, CPU: cpuHolder1}

		completedPart1 := database.BulkInsert(metricsHolder1)

		// BulkInsert returns true if it didn't insert correctly.
		if completedPart1 {

			t.Errorf("Error inserting data..")
		}

		singleCPU := database.GetNewestCPU()

		if singleCPU != cpuHolder1 {

			t.Errorf("CPU insert--> Returned: %v, Wanted: %v\n", singleCPU, cpuHolder1)
		}

		// As MEMORY and DISK isn't implemented yet, these are left out.
		/*
			singleMemory := database.GetNewestIndivComponent("MEMORY")

			if singleMemory != memoryHolder1 {

				t.Errorf("MEMORY insert--> Returned: %v, Wanted: %v\n", singleMemory, memoryHolder1)
			}
		*/

		/*
			singleDisk := database.GetNewestIndivComponent("DISK")

			if singleDisk != diskHolder1 {

				t.Errorf("CPU insert--> Returned: %v, Wanted: %v\n", singleDisk, diskHolder1)
			}
		*/

		newestCollector := database.GetCollectorNewest()

		// As we only look for CPU right now, we don't have memory of disk.
		if newestCollector.collectorID != collectorHolder1.collectorID || newestCollector.cpuID != collectorHolder1.cpuID {

			t.Errorf("COLLECTOR insert--> (ReturnedCollectorID: %v, WantedCollectorID: %v), (ReturnedCpuID: %v, WantedCpuID: %v)\n",
				newestCollector.collectorID, collectorHolder1.collectorID, newestCollector.cpuID, collectorHolder1.cpuID)
		}

		newProcesses := database.GetProcessesByNewest()

		for iteration, process := range newProcesses {

			if process != listProcess[iteration] {

				t.Errorf("PROCESS insert--> GotPID:%v, WanetdPID: %v\n",
					process.PID, listProcess[iteration].PID)
			}
		}

		if numberOfRounds == 5 {

			ctx := context.Background()

			// For not we are just getting from the single table!
			singleQuery := fmt.Sprintf("DELETE FROM PROCESS; DELETE FROM COLLECTOR;" +
				" DELETE FROM CPU; DELETE FROM MEMORY; DELETE FROM DISK;")

			// Execute query
			rows, err := database.DB_CONNECTION.QueryContext(ctx, singleQuery)

			if err != nil {

				log.Fatal(err.Error())
			}

			defer rows.Close()
		}

		numberOfRounds--

	}

	database.CloseDBConnection()

}

// When RAM and DISK is fully implemented, this test will be used to test:
// CPU/RAM/DISK and COLLECTOR without PROCESS insertion.
// This tests for correct insertion without inserting into PROCESS, and
// Makes sure that everything is saved correctly.
func TestMainThreeCorrectInputs(t *testing.T) {

	var database, err = NewStorage()

	if err != nil {
	}

	cpuHolder1 := cpu.CPU{Usage: 22.22}
	//memoryHolder1 := memory.MEMORY{}
	//ramHolder1 := ram.RAM{}

	collectorHolder1 := Collector{collectorID: globalcollectorID, cpuID: globalcpuID /*, memoryID: globalmemoryID, diskID: globaldiskID*/}

	// 'BulkInsert' needs to take in Metrics type, so an empty process list is needed.
	listProcess := []process.Process{}

	metricsHolder1 := collecting.Metrics{Processes: listProcess, CPU: cpuHolder1}

	completedPart1 := database.BulkInsert(metricsHolder1)

	// BulkInsert returns true if it didn't insert correctly.
	if completedPart1 {

		t.Errorf("Error inserting data..")
	}

	singleCPU := database.GetNewestCPU()

	if singleCPU != cpuHolder1 {

		t.Errorf("CPU insert--> Returned: %v, Wanted: %v\n", singleCPU, cpuHolder1)
	}

	// As MEMORY and DISK isn't implemented yet, these are left out.
	/*
		singleMemory := database.GetNewestIndivComponent("MEMORY")

		if singleMemory != memoryHolder1 {

			t.Errorf("MEMORY insert--> Returned: %v, Wanted: %v\n", singleMemory, memoryHolder1)
		}
	*/

	/*
		singleDisk := database.GetNewestIndivComponent("DISK")

		if singleDisk != diskHolder1 {

			t.Errorf("CPU insert--> Returned: %v, Wanted: %v\n", singleDisk, diskHolder1)
		}
	*/

	newestCollector := database.GetCollectorNewest()

	// As we only look for CPU right now, we don't have memory of disk.
	if newestCollector.collectorID != collectorHolder1.collectorID || newestCollector.cpuID != collectorHolder1.cpuID {

		t.Errorf("COLLECTOR insert--> (ReturnedCollectorID: %v, WantedCollectorID: %v), (ReturnedCpuID: %v, WantedCpuID: %v)\n",
			newestCollector.collectorID, collectorHolder1.collectorID, newestCollector.cpuID, collectorHolder1.cpuID)
	}

	database.CloseDBConnection()
}

//	This tests sending a new process list to be inserted, however there is a duplicate.
// 	This means the process insert function should return an error.
//	If an error is returned, then it was successful
//	(and the duplicate was not inserted).
func TestDuplicateProcess(t *testing.T) {

	var database, err = NewStorage()

	if err != nil {
	}

	// 'BulkInsert' needs to take in Metrics type, so an empty process list is needed.
	listProcess := []process.Process{}
	listProcess = append(listProcess, process.Process{PID: globalprocess0ID, Name: "process0", CPUUtilization: 54.00, RAMUtilization: 50.00, DiskUtilization: 00.00, Status: "done", ExecutionTime: 00.00})
	listProcess = append(listProcess, process.Process{PID: globalprocess1ID, Name: "process1", CPUUtilization: 01.00, RAMUtilization: 01.00, DiskUtilization: 00.10, Status: "running", ExecutionTime: 01.00})
	listProcess = append(listProcess, process.Process{PID: globalprocess2ID, Name: "process2", CPUUtilization: 02.20, RAMUtilization: 22.00, DiskUtilization: 77.22, Status: "failed", ExecutionTime: 22.00})
	listProcess = append(listProcess, process.Process{PID: globalprocess3ID, Name: "process3", CPUUtilization: 70.00, RAMUtilization: 00.00, DiskUtilization: 00.10, Status: "done", ExecutionTime: 99.99})
	// Duplicate is here
	listProcess = append(listProcess, process.Process{PID: globalprocess0ID, Name: "process0", CPUUtilization: 54.00, RAMUtilization: 50.00, DiskUtilization: 00.00, Status: "done", ExecutionTime: 00.00})

	errorHappened := false

	for _, singleProcess := range listProcess {

		rowsAffected, err := database.PutNewProcess(singleProcess)
		if err != nil {
		}
		if rowsAffected == 0 {
			errorHappened = true
		}
	}

	if !errorHappened {

		t.Errorf("Duplicate PROCESS insert--> Process insert function should have " +
			"returned an error as a duplicate was found, it did not.\n")
	}

	database.CloseDBConnection()
}
