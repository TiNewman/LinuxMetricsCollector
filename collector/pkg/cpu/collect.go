package cpu

import (
	"math"
	"os"
	"time"

	"github.com/prometheus/procfs"
)

type CPU struct {
	Model string
	Cores int
	Usage float32
}

type collector struct {
	r     Repository
	mount string
}

type Collector interface {
	Collect() (CPU, error)
}

func NewCPUCollector(repo Repository) collector {
	return collector{r: repo, mount: "/proc"}
}

func NewCPUCollectorWithoutRepo() collector {
	return collector{mount: "/proc"}
}

func newTestCollector(mp string) collector {
	wd, _ := os.Getwd()
	mountpoint := wd + "/testdata/" + mp
	return collector{mount: mountpoint}
}

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
