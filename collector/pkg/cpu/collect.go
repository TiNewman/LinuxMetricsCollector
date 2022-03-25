package cpu

import (
	"fmt"

	"github.com/prometheus/procfs"
)

type CPU struct {
	Usage        float32
	Availability float32
}

type collector struct {
	r Repository
}

type Collector interface {
	Collect() []CPU
}

func NewCPUCollector(repo Repository) collector {
	return collector{repo}
}

func (c collector) Collect() []CPU {
	fs, err := procfs.NewDefaultFS()
	if err != nil {
		fmt.Printf("Cannot locate proc mount %v", err.Error())
	}

	info, err := fs.CPUInfo()
	if err != nil {
		fmt.Printf("Could not get CPU info: %v\n", err)
	}
	fmt.Printf("%v\n", len(info))
	fmt.Printf("%+v\n", info[0])

	stat, err := fs.Stat()
	if err != nil {
		fmt.Printf("Could not get CPU stat: %v\n", err)
	}
	fmt.Printf("%+v\n", stat)

	/*
		for _, cpu := range info {
			fmt.Printf("%+v\n", cpu)
		}
	*/

	return []CPU{}
}
