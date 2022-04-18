package disk

import (
	"fmt"
	"os/exec"
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

	out, err := exec.Command(c.command).Output()
	if err != nil {
		return result, err
	}

	fmt.Printf(string(out))

	return result, nil
}
