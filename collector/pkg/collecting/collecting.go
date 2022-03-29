package collecting

import (
	"github.com/TiNewman/LinuxMetricsCollector/pkg/cpu"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/process"
)

type Service interface {
	Collect() Metrics
}

type service struct {
	r Repository
	p process.Collector
	c cpu.Collector
}

type Metrics struct {
	Processes []process.Process
	CPU       []cpu.CPU
}

type Repository interface {
	PutMetric([]cpu.CPU, []process.Process)
}

func NewService(proc process.Collector, cpu cpu.Collector, repo Repository) service {
	return service{p: proc, c: cpu, r: repo}
}

func NewServiceWithoutRepo(proc process.Collector, cpu cpu.Collector) service {
	return service{p: proc, c: cpu}
}

func (s service) Collect() Metrics {
	CPUInfo := s.c.Collect()
	// call cpu database code (remove database injection from cpu collector)?
	// or return new row id from cpu Collect

	// Old approach
	//s.r.PutNewCollector()

	processes := s.p.Collect()
	// call process database code (remove database injection from process collector)?

	// New approach
	if s.r != nil {
		s.r.PutMetric(CPUInfo, processes)
	}

	return Metrics{CPU: CPUInfo, Processes: processes}
}
