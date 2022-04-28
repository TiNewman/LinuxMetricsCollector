package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/collecting"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/cpu"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/disk"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/logger"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/memory"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/process"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Handler(collector collecting.Service) http.Handler {
	logger.Debug("Websocket Handler Started")

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsEndpoint(collector))

	return mux
}

func wsEndpoint(collector collecting.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error(fmt.Sprintf("Error upgradiing connection: %s", err.Error()))
		}

		logger.Debug(fmt.Sprintf("Client Connected!"))

		writeChan := make(chan string)
		go reader(ws, writeChan)
		writer(ws, writeChan, collector)
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
					// logger.Info("Connection closed by client")
					writeChan <- fmt.Sprint("stop")
					conn.Close()
					return
				}
			} else {
				// logger.Error(fmt.Sprintf("Error reading message: %v", err.Error()))
				continue
			}
		}

		// logger.Debug(fmt.Sprintf("Message Received: %s", string(p)))
		var req clientreq
		err = json.Unmarshal(p, &req)
		if err != nil {
			// logger.Error(fmt.Sprintf("Error Decoding JSON Request: %s\n", err.Error()))
		}

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
			if m == "all" {
				publish = true
				metric = "all"
				data := collector.Collect()
				sendAllMetrics(conn, data)
			}
			if m == "stop" {
				logger.Debug(fmt.Sprintf("Stopping message stream..."))
				publish = false
			}
			logger.Debug(fmt.Sprintf("writer received message: %s", m))
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
				case "all":
					data := collector.Collect()
					sendAllMetrics(conn, data)
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
		logger.Error(fmt.Sprintf("Error: %s", err.Error()))
	}
}

func sendCPUInfo(conn *websocket.Conn, cpuList cpu.CPU) {
	response := make(map[string]interface{})

	response["cpu"] = cpuList

	err := writeSocketResponse(conn, response)
	if err != nil {
		logger.Error(fmt.Sprintf("Error: %s", err.Error()))
	}
}

func sendMemInfo(conn *websocket.Conn, mem memory.Memory) {
	response := make(map[string]interface{})

	response["memory"] = mem

	err := writeSocketResponse(conn, response)
	if err != nil {
		logger.Error(fmt.Sprintf("Error: %v", err.Error()))
	}
}

func sendDiskInfo(conn *websocket.Conn, disks []disk.Disk) {
	response := make(map[string]interface{})

	response["disk"] = disks

	err := writeSocketResponse(conn, response)
	if err != nil {
		logger.Error(fmt.Sprintf("Error: %v", err.Error()))
	}
}

func sendAllMetrics(conn *websocket.Conn, metrics collecting.Metrics) {
	response := make(map[string]interface{})

	response["process_list"] = metrics.Processes
	response["cpu"] = metrics.CPU
	response["memory"] = metrics.Memory
	response["disk"] = metrics.Disk

	err := writeSocketResponse(conn, response)
	if err != nil {
		logger.Error(fmt.Sprintf("Error: %v", err.Error()))
	}
}

func writeSocketResponse(conn *websocket.Conn, res map[string]interface{}) error {
	jsonResponse, err := json.Marshal(res)
	if err != nil {
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, jsonResponse)
	if err != nil {
		return err
	}
	return nil
}
