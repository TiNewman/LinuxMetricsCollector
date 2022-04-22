package disk

import (
	"regexp"

	"github.com/prometheus/procfs"
	"golang.org/x/sys/unix"
)

type Disk struct {
	Name       string
	MountPoint string
	Usage      float64
	Size       float64
}

type collector struct {
	location string // regular expresion matching the directory to search for disk devices
}

type Collector interface {
	Collect() ([]Disk, error)
}

func NewDiskCollector() collector {
	return collector{location: `/dev.*`}
}

func (c collector) Collect() ([]Disk, error) {
	result := []Disk{}

	mountInfo, err := procfs.GetMounts()
	if err != nil {
		return result, err
	}

	for _, m := range mountInfo {
		if found, _ := regexp.MatchString(c.location, m.Source); found {
			// calculate disk usage and size
			// then append to the list of disks
			var stat unix.Statfs_t
			// fmt.Printf("%v:\n", m.MountPoint)
			unix.Statfs(m.MountPoint, &stat)
			usage := calculateUsage(float64(stat.Bfree), float64(stat.Blocks))
			// fmt.Printf("Used: %v%%\n", usage)
			size := float64(stat.Blocks*uint64(stat.Bsize)) / 1000000
			// fmt.Printf("Size: %v\n", size)

			result = append(result, Disk{Name: m.Source, MountPoint: m.MountPoint, Usage: usage, Size: size})
		}
	}

	return result, nil
}

func calculateUsage(free float64, total float64) float64 {
	return (1 - (free / total)) * 100
}
