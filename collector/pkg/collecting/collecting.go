package collecting

import (
	"github.com/TiNewman/LinuxMetricsCollector/pkg/cpu"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/process"
)

type Service interface {
	Collect()
}

type service struct {
	r Repository
	p process.Collector
	c cpu.Collector
}

type Metric struct {
	processes []process.Process
	cpu       []cpu.CPU
}

type Repository interface {
	PutMetric([]cpu.CPU, []process.Process)
}

func NewService(proc process.Collector, cpu cpu.Collector, repo Repository) service {
	return service{p: proc, c: cpu, r: repo}
}

func (s service) Collect() ([]cpu.CPU, []process.Process) {
	CPUInfo := s.c.Collect()
	// call cpu database code (remove database injection from cpu collector)?
	// or return new row id from cpu Collect

	// Old approach
	//s.r.PutNewCollector()

	processes := s.p.Collect()
	// call process database code (remove database injection from process collector)?

	// New approach
	s.r.PutMetric(CPUInfo, processes)

	return CPUInfo, processes
}
