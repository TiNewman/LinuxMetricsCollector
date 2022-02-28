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

	for _, proc := range p {
		procstat, err := proc.Stat()

		if err != nil {
			fmt.Printf("Could not get process status: %v\n", err.Error())
			return
		}

		// process schedule state: running, asleep, etc.
		status := procstat.State

		// resident memory in bytes
		mem := procstat.ResidentMemory()

		// calculate the CPU Utilization
		unixstarttime, err := procstat.StartTime()
		if err != nil {
			fmt.Printf("Could not get start time of process: %v\n", err.Error())
			return
		}

		starttime := time.Unix(int64(unixstarttime), 0)

		// total process cpu time in seconds
		cputime := procstat.CPUTime()

		cpuUtilization := calcCPUUtilization(cputime, currentTime, starttime)

		fmt.Printf("Current Time: %v\n", currentTime)
		fmt.Printf("Process: %v, CPU Utilization: %v, Mem Usage: %v, Status: %v, StartTime: %v\n", proc.PID, cpuUtilization, float64(mem)/1000000, status, starttime)
		processList = append(processList, Process{PID: proc.PID, CPUUtilization: float32(cputime), RAMUtilization: float32(mem) / 1000000, Status: status, TimeStamp: currentTime})
	}
	fmt.Printf("Processes in list: %v\n", len(processList))

}

// use cpu times to calculate the percent utilization of a given process
func calcCPUUtilization(cputime float64, currentTime time.Time, startTime time.Time) float64 {
	// total time = utime + stime -> cputime
	/*
		sysinfo = &syscall.Sysinfo_t{}
		err := syscall.Sysinfo(sysinfo)
	*/
	// uptime -> sysinfo.Uptime
	// seconds (how long the process has been running) = uptime - (starttime / hertz) -> uptime - procstat.StartTime() -> time.Now() - procstat.StartTime()
	// usage = 100 * ((totaltime / hertz) / seconds)
	// we have cpu time in seconds, not jiffies
	// -> usage = 100 * (totaltime / seconds)

	return 100 * (cputime / currentTime.Sub(startTime).Seconds())
}
