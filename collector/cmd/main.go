package main

import (
	"fmt"
	"net/http"

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

	// serve endpoints
	fmt.Println("Starting Service")
	http.ListenAndServe(":8080", nil)
}
