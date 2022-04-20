package memory

import (
	"os"

	"github.com/prometheus/procfs"
)

type Memory struct {
	Usage float64 // amount of ram used as a percentage
	Size  float64 // total ram in Megagytes
}

type collector struct {
	mount string
}

type Collector interface {
	Collect() (Memory, error)
}

func NewMemoryCollector() collector {
	return collector{mount: "/proc"}
}

func newTestCollector(mp string) collector {
	wd, _ := os.Getwd()
	mountpoint := wd + "/testdata/" + mp
	return collector{mount: mountpoint}
}

func (c collector) Collect() (Memory, error) {
	result := Memory{}

	fs, err := procfs.NewFS(c.mount)
	if err != nil {
		// fmt.Printf("Cannot locate proc mount %v", err.Error())
		return result, err
	}

	info, err := fs.Meminfo()
	if err != nil {
		// fmt.Printf("Could not get CPU info: %v\n", err)
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

func calculateUsage(total *uint64, available *uint64) float64 {
	ftotal := float64(*total)
	favailable := float64(*available)
	used := ftotal - favailable
	usage := (used / ftotal) * 100
	return usage
}
