package cpu

import (
	"fmt"
	"time"

	"github.com/prometheus/procfs"
)

type CPU struct {
	Model string
	Cores int
	Usage float32
}

type collector struct {
	r Repository
}

type Collector interface {
	Collect() []CPU
}

func NewCPUCollector(repo Repository) collector {
	return collector{r: repo}
}

func NewCPUCollectorWithoutRepo() collector {
	return collector{}
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
	cores := info[0].CPUCores
	model := info[0].ModelName

	startStat, err := fs.Stat()
	if err != nil {
		fmt.Printf("Could not get CPU stat: %v\n", err)
	}

	time.Sleep(time.Second)

	endStat, err := fs.Stat()
	if err != nil {
		fmt.Printf("Could not get CPU stat: %v\n", err)
	}
	// fmt.Printf("%+v\n", stat)

	totalUsage := calculateUsage(startStat.CPUTotal, endStat.CPUTotal)

	result = append(result, CPU{Usage: totalUsage, Model: model, Cores: int(cores)})

	for i := range startStat.CPU {
		coreUsage := calculateUsage(startStat.CPU[i], endStat.CPU[i])
		result = append(result, CPU{Usage: coreUsage})
	}

	fmt.Printf("%+v\n", result)

	return result
}

func calculateUsage(start procfs.CPUStat, end procfs.CPUStat) float32 {

	userDiff := end.User - start.User
	sysDiff := end.System - start.System
	ioDiff := end.Iowait - start.Iowait
	idleDiff := end.Idle - start.Idle

	active := userDiff + sysDiff + ioDiff
	total := active + idleDiff

	usage := (active / total) * 100

	return float32(usage)
}
