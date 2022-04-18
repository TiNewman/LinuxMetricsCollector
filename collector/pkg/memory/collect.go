package memory

import "os"

type Memory struct {
	Usage float32
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
	return result, nil
}
