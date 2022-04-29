package cpu

import (
	"math"
	"os"
	"time"

	"github.com/prometheus/procfs"
)

// CPU provides information and metrics
// related to the system CPU.
type CPU struct {
	Model string
	Cores int
	Usage float32
}

// collector implements the Collector interface.
type collector struct {
	mount string
}

// Collector is the interface wrapping the Collect method.
// Collect returns a new CPU, representing current CPU
// information and metrics. Collect will return any errors
// encoutered during the collection process.
type Collector interface {
	Collect() (CPU, error)
}

// NewDefaultCPUCollector returns a new cpu collector.
// The default collector will search the "/proc"
// mount point for CPU information and metrics.
func NewDefaultCPUCollector() collector {
	return collector{mount: "/proc"}
}

func newTestCollector(mp string) collector {
	wd, _ := os.Getwd()
	mountpoint := wd + "/testdata/" + mp
	return collector{mount: mountpoint}
}

// Collect collects CPU information and metrics
// from the system, returning a CPU struct
// and any errors that occured during the
// collection process.
func (c collector) Collect() (CPU, error) {
	result := CPU{}

	fs, err := procfs.NewFS(c.mount)
	if err != nil {
		return result, err
	}

	info, err := fs.CPUInfo()
	if err != nil {
		return result, err
	}
	cores := info[0].CPUCores
	model := info[0].ModelName

	startStat, err := fs.Stat()
	if err != nil {
		return result, err
	}

	time.Sleep(time.Second)

	endStat, err := fs.Stat()
	if err != nil {
		return result, err
	}

	totalUsage := calculateUsage(startStat.CPUTotal, endStat.CPUTotal)

	result = CPU{Usage: totalUsage, Model: model, Cores: int(cores)}

	return result, nil
}

// calculateUsage calculates and returns the total CPU usage as a percentage.
func calculateUsage(start procfs.CPUStat, end procfs.CPUStat) float32 {

	userDiff := end.User - start.User
	sysDiff := end.System - start.System
	ioDiff := end.Iowait - start.Iowait
	idleDiff := end.Idle - start.Idle

	active := userDiff + sysDiff + ioDiff
	total := active + idleDiff

	usage := (active / total) * 100

	if math.IsNaN(usage) {
		return 0
	}

	return float32(usage)
}
