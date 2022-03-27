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
	result := []CPU{}

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

	totalUsage := calculateUsage(stat.CPUTotal)

	result = append(result, CPU{Usage: totalUsage, Availability: 0})

	for _, cpu := range stat.CPU {
		coreUsage := calculateUsage(cpu)
		result = append(result, CPU{Usage: coreUsage, Availability: 0})
	}

	return result
}

func calculateUsage(stat procfs.CPUStat) float32 {

	active := stat.User + stat.System + stat.Iowait
	total := active + stat.Idle

	usage := (active / total) * 100

	return float32(usage)
}
