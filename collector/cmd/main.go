//go:build linux && amd64

package main

import (
	"fmt"
	"net/http"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/collecting"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/cpu"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/http/websocket"
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

	// without persistent storage
	pcollector := process.NewProcessCollectorWithoutRepo()
	cpuCollector := cpu.NewCPUCollectorWithoutRepo()
	collectingService := collecting.NewServiceWithoutRepo(pcollector, cpuCollector)

	// with persistent storage
	/*
		pcollector := process.NewProcessCollector(s)
		cpuCollector := cpu.NewCPUCollector(s)
		collectingService := collecting.NewService(pcollector, cpuCollector, s)
	*/

	// serve endpoints
	fmt.Println("Starting Service")
	router := websocket.Handler(collectingService)
	http.ListenAndServe(":8080", router)
}
