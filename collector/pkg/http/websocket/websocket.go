package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/process"
	"github.com/gorilla/websocket"
)

type Collector interface {
	Collect() []process.Process
}

type Repository interface {
	PutNewCollector() (int64, error)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Handler(process Collector, repository Repository) http.Handler {
	fmt.Printf("Websocket Handler\n")
	// repository.PutNewCollector()
	// process.Collect()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsEndpoint(process, repository))

	return mux
}

func wsEndpoint(process Collector, repository Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error upgrading connection: ", err.Error())
		}
		fmt.Println("Client Connected!")
		writeChan := make(chan string)
		go reader(ws, writeChan)
		writer(ws, writeChan, process, repository)
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
			fmt.Println("Error reading message: ", err.Error())
			return
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

func writer(conn *websocket.Conn, c chan string, process Collector, repository Repository) {
	var lastWrite time.Time
	count := 0
	publish := true
	for {
		now := time.Now()
		select {
		case m := <-c:
			if m == "process_list" {
				publish = true
			}
			if m == "quit" {
				fmt.Println("Stopping message stream...")
				publish = false
			}
			if m == "start" {
				fmt.Println("Start message stream")
				publish = true
			}
			fmt.Printf("writer received message: %v\n", m)
			lastWrite = now
		default:
			if publish && !lastWrite.IsZero() && now.Sub(lastWrite).Seconds() > 30 {
				repository.PutNewCollector()
				processes := process.Collect()

				response := make(map[string]interface{})

				response["process_list"] = processes

				jsonResponse, err := json.Marshal(response)
				if err != nil {
					fmt.Printf("Cannot Marshal Processes to JSON")
					return
				}
				err = conn.WriteMessage(websocket.TextMessage, jsonResponse)
				if err != nil {
					fmt.Println("Error writing message: ", err.Error())
					return
				}
				count = count + 1
				lastWrite = now
			}
		}
	}
}
