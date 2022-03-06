package main

import (
	"fmt"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/process"
)

func main() {
	fmt.Printf("Linux Metrics Collector\n")

	// initialize storage

	// initialize collectors
	pcollector := process.NewProcessCollectorWithoutRepo()
	pcollector.Collect()

	// serve endpoints
	//fmt.Println("Starting Service")
	//router := websocket.Handler()
	//http.ListenAndServe(":8080", router)
}
