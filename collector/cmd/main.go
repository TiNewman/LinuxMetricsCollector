//go:build linux && amd64

package main

import (
	"fmt"
	"net/http"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/collecting"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/cpu"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/disk"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/http/websocket"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/logger"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/memory"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/process"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/storage/mssql"
)

func main() {
	// initialize logging
	logger.Init()

	logger.Info("Linux Metrics Collector Started")

	// initialize storage
	s, err := mssql.NewStorage()
	if err != nil {
		logger.Error(fmt.Sprintf("Could not initialize persistent storage: %v", err.Error()))
	}

	// initialize collectors
	pcollector := process.NewDefaultProcessCollector()
	cpuCollector := cpu.NewDefaultCPUCollector()
	memCollector := memory.NewMemoryCollector()
	diskCollector := disk.NewDiskCollector()

	collectingService := collecting.NewService(
		collecting.WithProcessCollector(pcollector),
		collecting.WithCPUCollector(cpuCollector),
		collecting.WithMemCollector(memCollector),
		collecting.WithDiskCollector(diskCollector),
		collecting.WithRepository(s))

	// serve endpoints
	router := websocket.Handler(collectingService)
	http.ListenAndServe(":8080", router)
}
