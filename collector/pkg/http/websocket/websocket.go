package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/collecting"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/cpu"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/disk"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/memory"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/process"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Handler(collector collecting.Service) http.Handler {
	fmt.Printf("Websocket Handler\n")

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsEndpoint(collector))

	return mux
}

func wsEndpoint(collector collecting.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error upgrading connection: ", err.Error())
		}
		// fmt.Println("Client Connected!")
		writeChan := make(chan string)
		go reader(ws, writeChan)
		writer(ws, writeChan, collector)
		// fmt.Printf("go routines: %v\n", runtime.NumGoroutine())
	}
}

type clientreq struct {
	Request string `json:"request"`
}

func reader(conn *websocket.Conn, writeChan chan string) {
	for {
		_, p, err := conn.ReadMessage()

		if err != nil {
			if ce, ok := err.(*websocket.CloseError); ok {
				switch ce.Code {
				case websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived:
					// fmt.Println("Connection closed by client")
					writeChan <- fmt.Sprint("stop")
					conn.Close()
					return
				}
			} else {
				fmt.Printf("Error reading message: %v\n", err.Error())
				continue
			}
		}

		// fmt.Println("Message Received: ", string(p))
		var req clientreq
		err = json.Unmarshal(p, &req)
		if err != nil {
			fmt.Printf("Error Decoding JSON Request: %v\n", err.Error())
		}
		// fmt.Printf("%+v\n", req)

		writeChan <- string(req.Request)

	}
}

func writer(conn *websocket.Conn, c chan string, collector collecting.Service) {
	var lastWrite time.Time
	publish := false
	metric := ""
	for {
		now := time.Now()
		select {
		case m := <-c:
			if m == "process_list" {
				publish = true
				metric = "process_list"
				data := collector.Collect()
				sendProcessList(conn, data.Processes)
			}
			if m == "cpu" {
				publish = true
				metric = "cpu"
				data := collector.Collect()
				sendCPUInfo(conn, data.CPU)
			}
			if m == "memory" {
				publish = true
				metric = "memory"
				data := collector.Collect()
				sendMemInfo(conn, data.Memory)
			}
			if m == "disk" {
				publish = true
				metric = "disk"
				data := collector.Collect()
				sendDiskInfo(conn, data.Disk)
			}
			if m == "stop" {
				// fmt.Println("Stopping message stream...")
				publish = false
			}
			// fmt.Printf("writer received message: %v\n", m)
			lastWrite = now
		default:
			if publish && !lastWrite.IsZero() && now.Sub(lastWrite).Seconds() > 5 {
				switch metric {
				case "process_list":
					data := collector.Collect()
					sendProcessList(conn, data.Processes)
				case "cpu":
					data := collector.Collect()
					sendCPUInfo(conn, data.CPU)
				case "memory":
					data := collector.Collect()
					sendMemInfo(conn, data.Memory)
				case "disk":
					data := collector.Collect()
					sendDiskInfo(conn, data.Disk)
				}
				lastWrite = now
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func sendProcessList(conn *websocket.Conn, processes []process.Process) {
	response := make(map[string]interface{})

	response["process_list"] = processes

	err := writeSocketResponse(conn, response)
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
	}
}

func sendCPUInfo(conn *websocket.Conn, cpuList cpu.CPU) {
	response := make(map[string]interface{})

	response["cpu"] = cpuList

	err := writeSocketResponse(conn, response)
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
	}
}

func sendMemInfo(conn *websocket.Conn, mem memory.Memory) {
	response := make(map[string]interface{})

	response["memory"] = mem

	err := writeSocketResponse(conn, response)
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
	}
}

func sendDiskInfo(conn *websocket.Conn, disks []disk.Disk) {
	response := make(map[string]interface{})

	response["disk"] = disks

	err := writeSocketResponse(conn, response)
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
	}
}

func writeSocketResponse(conn *websocket.Conn, res map[string]interface{}) error {
	jsonResponse, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("Cannot Marshal Processes to JSON")
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, jsonResponse)
	if err != nil {
		fmt.Println("Error writing message: ", err.Error())
		return err
	}
	return nil
}
