package process

import (
	"fmt"
	"os/user"
	"time"

	"github.com/prometheus/procfs"
)

type Process struct {
	PID             int
	Name            string
	CPUUtilization  float32
	RAMUtilization  float32
	DiskUtilization float32
	Status          string
	ExecutionTime   float32
}

type collector struct {
	r Repository
}

func NewProcessCollector(repo Repository) collector {
	return collector{repo}
}

func NewProcessCollectorWithoutRepo() collector {
	return collector{}
}

func (c collector) Collect() {
	currentTime := time.Now()
	processList := []Process{}
	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("Cannot determine current user: %v\n", err.Error())
		return
	}

	// read process list from the proc file system
	p, err := procfs.AllProcs()
	if err != nil {
		fmt.Printf("Could not get all processes: %v\n", err)
	}
	fmt.Printf("Number of processes: %v\n", p.Len())

	// Calculate necessary values for each process and place them in a custom
	// Process struct
	for _, proc := range p {
		// filter by UID ( get processes for current user only )
		procStatus, err := proc.NewStatus()
		if err != nil {
			fmt.Printf("Could not get uids of process: %v\n", err.Error())
			return
		}
		uids := procStatus.UIDs
		if currentUser.Uid != uids[0] {
			continue
		}

		procstat, err := proc.Stat()

		if err != nil {
			fmt.Printf("Could not get process status: %v\n", err.Error())
			return
		}

		// get the process name
		pname, err := proc.Comm()
		if err != nil {
			fmt.Printf("Could not get process name: %v\n", err.Error())
			return
		}

		// process schedule state: running, asleep, etc.
		status := procstat.State

		// resident memory in bytes
		mem := procstat.ResidentMemory()

		// disk utilization (bytes read by process)
		io, err := proc.IO()

		var readTotal uint64
		if err != nil {
			fmt.Printf("Could not get IO metrics: %v\n", err)
			//return
			// set a negative value to signify N/A?
			readTotal = 0
		} else {
			readTotal = io.ReadBytes
		}

		// calculate the CPU Utilization
		unixstarttime, err := procstat.StartTime()
		if err != nil {
			fmt.Printf("Could not get start time of process: %v\n", err.Error())
			return
		}

		// time the process started
		starttime := time.Unix(int64(unixstarttime), 0)

		// total process cpu time in seconds
		cputime := procstat.CPUTime()

		// total time the process has been running in seconds
		executionTime := currentTime.Sub(starttime).Seconds()

		cpuUtilization := calcCPUUtilization(cputime, executionTime)

		nextprocess := Process{PID: proc.PID, Name: pname, CPUUtilization: float32(cpuUtilization), RAMUtilization: float32(mem) / 1000000, DiskUtilization: float32(readTotal) / 1000000, Status: status, ExecutionTime: float32(executionTime)}

		if c.r != nil {
			c.r.PutNewProcess(nextprocess)
		}

		fmt.Printf("%+v\n", nextprocess)
		processList = append(processList, nextprocess)

	}
	fmt.Printf("Processes in list: %v\n", len(processList))

}

// use cpu times to calculate the percent utilization of a given process
func calcCPUUtilization(cputime float64, executionTime float64) float64 {
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

	return 100 * (cputime / executionTime)
}
