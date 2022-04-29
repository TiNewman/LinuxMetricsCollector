package disk

import (
	"regexp"

	"github.com/prometheus/procfs"
	"golang.org/x/sys/unix"
)

// Disk provides information and metrics
// related to a physical disk.
type Disk struct {
	Name       string
	MountPoint string
	Usage      float64
	Size       float64
}

// collector implements the Collector interface.
type collector struct {
	// regular expresion matching the directory to search for disk devices
	location string
}

// Collector is the interface wrapping the Collect method.
// Collect returns a new Disk slice, representing current physical disk
// information and metrics. Collect will return any errors
// encoutered during the collection process.
type Collector interface {
	Collect() ([]Disk, error)
}

// NewDiskCollector returns a new disk collector.
// The default collector will search the "/dev"
// mount point for physical disks.
func NewDiskCollector() collector {
	return collector{location: `/dev.*`}
}

// Collect collects physical disk information and metrics
// from the system, returning a Disk slice
// and any errors that occured during the
// collection process.
func (c collector) Collect() ([]Disk, error) {
	result := []Disk{}

	mountInfo, err := procfs.GetMounts()
	if err != nil {
		return result, err
	}

	for _, m := range mountInfo {
		if found, _ := regexp.MatchString(c.location, m.Source); found {
			var stat unix.Statfs_t
			unix.Statfs(m.MountPoint, &stat)
			usage := calculateUsage(float64(stat.Bfree), float64(stat.Blocks))
			size := float64(stat.Blocks*uint64(stat.Bsize)) / 1000000

			result = append(result, Disk{Name: m.Source, MountPoint: m.MountPoint, Usage: usage, Size: size})
		}
	}

	return result, nil
}

// calculateUsage calculates and returns the usage of a physical disk as
// a percentage.
func calculateUsage(free float64, total float64) float64 {
	if total <= 0 {
		return 0
	}
	return (1 - (free / total)) * 100
}
