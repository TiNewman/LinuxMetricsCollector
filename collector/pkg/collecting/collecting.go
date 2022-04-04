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

type Metrics struct {
	Processes []process.Process
	CPU       cpu.CPU
}

type Repository interface {
	PutMetric(Metrics)
}

func NewService(proc process.Collector, cpu cpu.Collector, repo Repository) service {
	return service{p: proc, c: cpu, r: repo}
}

func (s service) Collect() Metrics {
	CPUInfo := s.c.Collect()
	// call cpu database code (remove database injection from cpu collector)?
	// or return new row id from cpu Collect

	// Old approach
	//s.r.PutNewCollector()

	processes := s.p.Collect()
	// call process database code (remove database injection from process collector)?
	m := Metrics{Processes: processes, CPU: CPUInfo}
	// New approach
	s.r.PutMetric(m)

	return m
}
