package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/collecting"
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
		fmt.Println("Client Connected!")
		writeChan := make(chan string)
		go reader(ws, writeChan)
		writer(ws, writeChan, collector)
		fmt.Printf("go routines: %v\n", runtime.NumGoroutine())
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
					fmt.Println("Connection closed by client")
					writeChan <- fmt.Sprint("stop")
					conn.Close()
					return
				}
			} else {
				fmt.Printf("Error reading message: %v\n", err.Error())
				continue
			}
		}

		fmt.Println("Message Received: ", string(p))
		var req clientreq
		err = json.Unmarshal(p, &req)
		if err != nil {
			fmt.Printf("Error Decoding JSON Request: %v\n", err.Error())
		}
		fmt.Printf("%+v\n", req)

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
			if m == "stop" {
				fmt.Println("Stopping message stream...")
				publish = false
			}
			fmt.Printf("writer received message: %v\n", m)
			lastWrite = now
		default:
			if publish && !lastWrite.IsZero() && now.Sub(lastWrite).Seconds() > 30 {
				switch metric {
				case "process_list":
					data := collector.Collect()
					sendProcessList(conn, data.Processes)
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

func collectAndSendProcessList(conn *websocket.Conn, process process.Collector) {
	response := make(map[string]interface{})

	processes := process.Collect()
	response["process_list"] = processes

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
