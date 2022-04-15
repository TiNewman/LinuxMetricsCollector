package process

import (
	"fmt"
	"os"
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
	r     Repository
	mount string
}

type Collector interface {
	Collect() ([]Process, error)
}

func NewProcessCollector(repo Repository) collector {
	return collector{repo, "/proc"}
}

func NewProcessCollectorWithoutRepo() collector {
	return collector{mount: "/proc"}
}

func NewTestCollector() collector {
	wd, _ := os.Getwd()
	mountpoint := wd + "/testdata"
	return collector{mount: mountpoint}
}

func (c collector) Collect() ([]Process, error) {
	currentTime := time.Now()
	processList := []Process{}
	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("Cannot determine current user: %v\n", err.Error())
		return processList, err
	}

	// collect data from the configured mount point
	// (collector/pkg/process/testdata when testing and /proc when in production)
	fs, err := procfs.NewFS(c.mount)
	if err != nil {
		fmt.Printf("Cannot locate proc mount %v", err.Error())
		return processList, err
	}
	p, err := fs.AllProcs()
	if err != nil {
		fmt.Printf("Could not get all processes: %v\n", err)
		return processList, err
	}

	// Calculate necessary values for each process and place them in a custom
	// Process struct
	for _, proc := range p {
		// filter by UID ( get processes for current user only )
		procStatus, err := proc.NewStatus()
		if err != nil {
			fmt.Printf("Could not get uids of process: %v\n", err.Error())
			return processList, err
		}
		uids := procStatus.UIDs
		if currentUser.Uid != uids[0] {
			continue
		}

		procstat, err := proc.Stat()

		if err != nil {
			fmt.Printf("Could not get process status: %v\n", err.Error())
			return processList, err
		}

		// get the process name
		pname, err := proc.Comm()
		if err != nil {
			fmt.Printf("Could not get process name: %v\n", err.Error())
			return processList, err
		}

		// process schedule state: running, asleep, etc.
		status := procstat.State

		// resident memory in bytes
		mem := procstat.ResidentMemory()

		// disk utilization (bytes read by process)
		io, err := proc.IO()

		var readTotal uint64
		if err != nil {
			fmt.Printf("Could not get IO metrics for process %v: %v\n", pname, err)
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
			return processList, err
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

		// fmt.Printf("%+v\n", nextprocess)
		processList = append(processList, nextprocess)

	}
	fmt.Printf("Processes in list: %v\n", len(processList))
	return processList, nil

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
	if executionTime <= 0 {
		return 0
	}

	return 100 * (cputime / executionTime)
}
