package process

import (
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

func newTestCollector(mp string) collector {
	wd, _ := os.Getwd()
	mountpoint := wd + "/testdata/" + mp
	return collector{mount: mountpoint}
}

func (c collector) Collect() ([]Process, error) {
	currentTime := time.Now()
	processList := []Process{}
	currentUser, err := user.Current()
	if err != nil {
		return processList, err
	}

	// collect data from the configured mount point
	// (collector/pkg/process/testdata when testing and /proc when in production)
	fs, err := procfs.NewFS(c.mount)
	if err != nil {
		return processList, err
	}
	p, err := fs.AllProcs()
	if err != nil {
		return processList, err
	}

	// Calculate necessary values for each process and place them in a custom
	// Process struct
	for _, proc := range p {
		// filter by UID ( get processes for current user only )
		procStatus, err := proc.NewStatus()
		if err != nil {
			return processList, err
		}
		uids := procStatus.UIDs
		if currentUser.Uid != uids[0] {
			continue
		}

		procstat, err := proc.Stat()

		if err != nil {
			return processList, err
		}

		pname := procstat.Comm

		// process schedule state: running, asleep, etc.
		status := procstat.State

		// resident memory in bytes
		mem := procstat.ResidentMemory()

		// disk utilization (bytes read by process)
		io, err := proc.IO()

		var readTotal uint64
		if err != nil {
			readTotal = 0
		} else {
			readTotal = io.ReadBytes
		}

		// calculate the CPU Utilization
		unixstarttime, err := procstat.StartTime()
		if err != nil {
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

		processList = append(processList, nextprocess)

	}

	return processList, nil

}

// use cpu times to calculate the percent utilization of a given process
func calcCPUUtilization(cputime float64, executionTime float64) float64 {
	if executionTime <= 0 {
		return 0
	}

	return 100 * (cputime / executionTime)
}
