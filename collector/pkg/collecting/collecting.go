package collecting

import (
	"fmt"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/cpu"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/disk"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/logger"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/memory"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/process"
)

// Service is the interface wrapping the Collect method.
// Collect returns a new Metrics measurement, holding
// different system metrics collected from the system.
type Service interface {
	Collect() Metrics
}

// service imlements the Service interface,
// collecting and persisting Metrics using
// the assigned repository and collectors.
type service struct {
	r Repository
	p process.Collector
	c cpu.Collector
	m memory.Collector
	d disk.Collector
}

// Metrics provides metrics data points
// collected from the system.
type Metrics struct {
	Processes []process.Process
	CPU       cpu.CPU
	Memory    memory.Memory
	Disk      []disk.Disk
}

// Repository is the interface wrapping methods to
// interact with persisten metrics storage.
type Repository interface {
	BulkInsert(Metrics) bool
}

// ServiceOption is an option that can be used to
// configure a new collecting service.
type ServiceOption func(*service)

// WithProcessCollector returns a ServiceOption that
// assigns the given process Collector to a collecting Service.
func WithProcessCollector(proc process.Collector) ServiceOption {
	return func(s *service) {
		s.p = proc
	}
}

// WithCPUCollector returns a ServiceOption that
// assigns the given cpu Collector to a collecting Service.
func WithCPUCollector(cpu cpu.Collector) ServiceOption {
	return func(s *service) {
		s.c = cpu
	}
}

// WithMemCollector returns a ServiceOption that
// assigns the given memory Collector to a collecting Service.
func WithMemCollector(mem memory.Collector) ServiceOption {
	return func(s *service) {
		s.m = mem
	}
}

// WithDiskCollector returns a ServiceOption that
// assigns the given disk Collector to a collecting Service.
func WithDiskCollector(disk disk.Collector) ServiceOption {
	return func(s *service) {
		s.d = disk
	}
}

// WithRepository returns a ServiceOption that
// assigns the given Repository to a collecting Service.
func WithRepository(repo Repository) ServiceOption {
	return func(s *service) {
		s.r = repo
	}
}

// NewService configures a new collecting service with
// the given ServiceOption parameters and returns
// the new service.
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

// Collect collects system metrics using the
// configured collectors and persists the
// Metrics if a repository was configured,
// then returns the collected Metrics.
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
