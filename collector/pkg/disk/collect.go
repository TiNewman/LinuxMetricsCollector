package disk

import (
	"fmt"
	"regexp"

	"github.com/prometheus/procfs"
)

type Disk struct {
	Usage float64
	Size  float64
}

type collector struct {
	command string
}

type Collector interface {
	Collect() ([]Disk, error)
}

func NewDiskCollector() collector {
	return collector{}
}

func (c collector) Collect() ([]Disk, error) {
	result := []Disk{}

	/*
		out, err := exec.Command(c.command).Output()
		if err != nil {
			return result, err
		}

		fmt.Printf(string(out))
	*/

	mountInfo, err := procfs.GetMounts()
	if err != nil {
		return result, err
	}

	for _, m := range mountInfo {
		if found, _ := regexp.MatchString(`/dev.*`, m.Source); found {
			fmt.Printf("mount info: %+v\n\n", m)
		}
	}

	return result, nil
}
