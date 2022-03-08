package main

import (
	"fmt"
	"net/http"

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
	pcollector := process.NewProcessCollectorWithoutRepo()
	//pcollector.Collect()

	// serve endpoints
	fmt.Println("Starting Service")
	router := websocket.Handler(pcollector)
	http.ListenAndServe(":8080", router)
}
