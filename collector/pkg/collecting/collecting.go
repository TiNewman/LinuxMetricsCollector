package collecting

import (
	"fmt"
	"time"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/cpu"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/disk"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/logger"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/memory"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/process"
)

type Service interface {
	Collect() Metrics
	NewestHistory() History
}

type service struct {
	r Repository
	p process.Collector
	c cpu.Collector
	m memory.Collector
	d disk.Collector
}

type Metrics struct {
	Processes []process.Process
	CPU       cpu.CPU
	Memory    memory.Memory
	Disk      []disk.Disk
}

type History struct {
	Start           time.Time
	End             time.Time
	AverageCpuUsage float64
	AverageMemUsage float64
	AverageMemSize  float64
}

type Repository interface {
	BulkInsert(Metrics) bool
	GetNewestHistory() History
}

type ServiceOption func(*service)

func WithProcessCollector(proc process.Collector) ServiceOption {
	return func(s *service) {
		s.p = proc
	}
}

func WithCPUCollector(cpu cpu.Collector) ServiceOption {
	return func(s *service) {
		s.c = cpu
	}
}

func WithMemCollector(mem memory.Collector) ServiceOption {
	return func(s *service) {
		s.m = mem
	}
}

func WithDiskCollector(disk disk.Collector) ServiceOption {
	return func(s *service) {
		s.d = disk
	}
}

func WithRepository(repo Repository) ServiceOption {
	return func(s *service) {
		s.r = repo
	}
}

func NewService(opts ...ServiceOption) service {
	var s service
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

func NewServiceWithoutRepo(proc process.Collector, cpu cpu.Collector) service {
	return service{p: proc, c: cpu}
}

func (s service) Collect() Metrics {
	metrics := new(Metrics)

	if s.c != nil {
		CPUInfo, err := s.c.Collect()
		if err != nil {
			logger.Error(fmt.Sprintf("Error collecting CPU metrics: %v", err))
		}
		metrics.CPU = CPUInfo
	}

	if s.p != nil {
		processes, err := s.p.Collect()
		if err != nil {
			logger.Error(fmt.Sprintf("Error collecting Process metrics: %v", err))
		}
		metrics.Processes = processes
	}

	if s.m != nil {
		memInfo, err := s.m.Collect()
		if err != nil {
			logger.Error(fmt.Sprintf("Error collecting RAM metrics: %v", err))
		}
		metrics.Memory = memInfo
	}

	if s.d != nil {
		diskInfo, err := s.d.Collect()
		if err != nil {
			logger.Error(fmt.Sprintf("Error collecting disk metrics: %v", err))
		}
		metrics.Disk = diskInfo
	}

	m := *metrics

	if s.r != nil {
		s.r.BulkInsert(m)
	}

	return m
}

func (s service) NewestHistory() History {
	var history History
	history := s.r.GetNewestHistory()
	return history
}
