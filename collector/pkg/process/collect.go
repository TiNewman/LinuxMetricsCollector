package process

import (
	"fmt"

	"github.com/prometheus/procfs"
)

type Process struct {
	PID            int
	CPUUtilization float32
	RAMUtilization float32
	Status         string
}

func Collect() {
	p, err := procfs.AllProcs()
	if err != nil {
		fmt.Printf("Could not get all processes: %v\n", err)
	}
	fmt.Printf("Number of processes: %v\n", p.Len())

	firstProcess := p[0]

	procstat, err := firstProcess.Stat()

	if err != nil {
		fmt.Printf("Could not get process status: %v\n", err.Error())
		return
	}

	// total cpu time in seconds
	cputime := procstat.CPUTime()

	// process schedule state: running, asleep, etc.
	status := procstat.State

	// resident memory in bytes
	mem := procstat.ResidentMemory()

	fmt.Printf("Process: %v, CPU Time: %v, Mem Usage: %v, Status: %v\n", firstProcess.PID, cputime, mem, status)

}

func calcCPUUtilization() {
	// use cpu times to calculate the percent utilization of a given process

}
