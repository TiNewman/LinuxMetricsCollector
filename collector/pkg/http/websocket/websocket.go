package websocket

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

type Collector interface {
	Collect()
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

func reader(conn *websocket.Conn, writeChan chan string) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message: ", err.Error())
			return
		}

		fmt.Println("Message Received: ", string(p))

		writeChan <- string(p)

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
			if m == "quit" {
				fmt.Println("Stopping message stream...")
				publish = false
			}
			if m == "start" {
				fmt.Println("Start message stream")
				publish = true
			}
			fmt.Printf("writer received message: %v\n", m)
			err := conn.WriteMessage(websocket.TextMessage, []byte("Hello from server"))
			if err != nil {
				fmt.Println("Error writing message: ", err.Error())
				return
			}
			lastWrite = now
		default:
			if publish && !lastWrite.IsZero() && now.Sub(lastWrite).Seconds() > 30 {
				repository.PutNewCollector()
				process.Collect()
				err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", count)))
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
