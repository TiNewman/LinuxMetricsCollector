//go:build linux && amd64

package main

import (
	"fmt"
	"net/http"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/collecting"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/cpu"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/disk"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/http/websocket"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/memory"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/process"
)

func main() {
	fmt.Printf("Linux Metrics Collector\n")

	// initialize storage
	/*
		s, err := mssql.NewStorage()
		if err != nil {
			fmt.Printf("Could not initialize persistent storage: %v", err.Error())
		}
	*/

	// initialize collectors
	pcollector := process.NewProcessCollectorWithoutRepo()
	cpuCollector := cpu.NewCPUCollectorWithoutRepo()
	memCollector := memory.NewMemoryCollector()
	diskCollector := disk.NewDiskCollector()
	collectingService := collecting.NewService(
		collecting.WithProcessCollector(pcollector),
		collecting.WithCPUCollector(cpuCollector),
		collecting.WithMemCollector(memCollector),
		collecting.WithDiskCollector(diskCollector))
	// collectingService := collecting.NewServiceWithoutRepo(pcollector, cpuCollector)
	// collectingService := collecting.NewService(pcollector, cpuCollector, s)

	// serve endpoints
	fmt.Println("Starting Service")
	router := websocket.Handler(collectingService)
	http.ListenAndServe(":8080", router)
}
