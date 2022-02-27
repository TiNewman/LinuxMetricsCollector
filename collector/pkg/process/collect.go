package process

import (
	"fmt"
	"time"

	"github.com/prometheus/procfs"
)

type Process struct {
	PID            int
	CPUUtilization float32
	RAMUtilization float32
	Status         string
	TimeStamp      time.Time
}

func Collect() {
	currentTime := time.Now()
	processList := []Process{}
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

	fmt.Printf("Current Time: %v\n", currentTime)
	fmt.Printf("Process: %v, CPU Time: %v, Mem Usage: %v, Status: %v\n", firstProcess.PID, cputime, float64(mem)/1000000, status)
	processList = append(processList, Process{PID: firstProcess.PID, CPUUtilization: float32(cputime), RAMUtilization: float32(mem) / 1000000, Status: status, TimeStamp: currentTime})

}

// use cpu times to calculate the percent utilization of a given process
func calcCPUUtilization() {
	// total time = utime + stime -> cputime
	/*
		sysinfo = &syscall.Sysinfo_t{}
		err := syscall.Sysinfo(sysinfo)
	*/
	// uptime -> sysinfo.Uptime
	// seconds = uptime - (starttime / hertz) -> uptime - procstat.StartTime()
	// usage = 100 * ((totaltime / hertz) / seconds)
}
