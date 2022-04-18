package collecting

import (
	"fmt"

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
	CPU       cpu.CPU
}

type Repository interface {
	BulkInsert(Metrics) bool
}

func NewService(proc process.Collector, cpu cpu.Collector, repo Repository) service {
	return service{p: proc, c: cpu, r: repo}
}

func NewServiceWithoutRepo(proc process.Collector, cpu cpu.Collector) service {
	return service{p: proc, c: cpu}
}

func (s service) Collect() Metrics {
	CPUInfo, err := s.c.Collect()
	if err != nil {
		fmt.Printf("Error collecting CPU metrics: %v\n", err)
	}
	// call cpu database code (remove database injection from cpu collector)?
	// or return new row id from cpu Collect

	// Old approach
	//s.r.PutNewCollector()

	processes, err := s.p.Collect()
	if err != nil {
		fmt.Printf("Error collecting CPU metrics: %v\n", err)
	}
	// call process database code (remove database injection from process collector)?

	m := Metrics{CPU: CPUInfo, Processes: processes}

	// New approach
	if s.r != nil {
		s.r.BulkInsert(m)
	}

	return m
}
