package main

import (
	"fmt"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/process"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	fmt.Printf("Linux Metrics Collector\n")

	// initialize storage

	// initialize collectors
	process.Collect()

	// serve endpoints
	//fmt.Println("Starting Service")
	//http.ListenAndServe(":8080", nil)
}
