package cpu

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
	return result
}
