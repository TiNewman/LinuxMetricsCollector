package memory

import (
	"os"

	"github.com/prometheus/procfs"
)

// Memory provides information and metrics
// related to the system memory.
type Memory struct {
	Usage float64 // amount of ram used as a percentage
	Size  float64 // total ram in Megabytes
}

// collector implements the Collector interface.
type collector struct {
	mount string
}

// Collector is the interface wrapping the Collect method.
// Collect returns a Memory struct, representing current memory
// information and metrics. Collect will return any errors
// encoutered during the collection process.
type Collector interface {
	Collect() (Memory, error)
}

// NewMemoryCollector returns a new memory collector.
// The default collector will search the "/proc"
// mount point for memory information and metrics.
func NewMemoryCollector() collector {
	return collector{mount: "/proc"}
}

func newTestCollector(mp string) collector {
	wd, _ := os.Getwd()
	mountpoint := wd + "/testdata/" + mp
	return collector{mount: mountpoint}
}

// Collect collects memory information and metrics
// from the system, returning a Memory struct
// and any errors that occured during the
// collection process.
func (c collector) Collect() (Memory, error) {
	result := Memory{}

	fs, err := procfs.NewFS(c.mount)
	if err != nil {
		return result, err
	}

	info, err := fs.Meminfo()
	if err != nil {
		return result, err
	}

	total := info.MemTotal

	// convert total from kibibytes to Megabytes
	kibToMBRatio := float64(0.001024)
	totalMB := float64(*total) * kibToMBRatio
	available := info.MemAvailable

	usage := calculateUsage(total, available)

	result = Memory{Usage: usage, Size: totalMB}

	return result, nil
}

// calculateUsate calculates and returns the memory usage as a percentage.
func calculateUsage(total *uint64, available *uint64) float64 {
	ftotal := float64(*total)
	favailable := float64(*available)
	used := ftotal - favailable
	usage := (used / ftotal) * 100
	return usage
}
