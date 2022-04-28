package websocket

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/TiNewman/LinuxMetricsCollector/pkg/collecting"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/cpu"
	"github.com/TiNewman/LinuxMetricsCollector/pkg/process"
	"github.com/gorilla/websocket"
)

var done chan interface{}
var interrupt chan os.Signal

func TestMain(m *testing.M) {
	// start websocket server
	pcollector := process.NewProcessCollectorWithoutRepo()
	cpuCollector := cpu.NewDefaultCPUCollector()
	collectingService := collecting.NewServiceWithoutRepo(pcollector, cpuCollector)

	// serve endpoints
	fmt.Println("Starting Service")
	router := Handler(collectingService)
	go http.ListenAndServe(":8080", router)

	// run test cases
	exitVal := m.Run()

	os.Exit(exitVal)
}

func receiveHandler(connection *websocket.Conn, t *testing.T) {
	defer close(done)
	for {
		_, _, err := connection.ReadMessage()
		if err != nil {
			t.Errorf("Error in receive: %v", err.Error())
			done <- "done"
			return
		}
		// log.Printf("Received: %s\n", msg)
		done <- "done"
		return
	}
}

func receiveErrorHandler(connection *websocket.Conn, t *testing.T) {
	defer close(done)
	for {
		_, _, err := connection.ReadMessage()
		if err == nil {
			t.Errorf("Expected error message from server")
			done <- "done"
			return
		}
		// log.Printf("Received: %s\n", msg)
		done <- "done"
		return
	}
}

func TestMalformedReq(t *testing.T) {
	done = make(chan interface{})    // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	// connect client to server
	socketUrl := "ws://localhost:8080" + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	// handle server responses
	go receiveErrorHandler(conn, t)

	// send request to websocket server
	response := make(map[string]interface{})

	response["this isn't right"] = "this should not work"

	err = writeSocketResponse(conn, response)
	if err != nil {
		t.Errorf("Error: %v", err.Error())
		return
	}

	// wait for read handler to finish
	for {
		select {
		case <-interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			log.Println("Received SIGINT interrupt signal. Closing all pending connections")

			// Close our websocket connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
		case <-done:
			// log.Println("Receiver Channel Closed! Exiting....")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
			return
		case <-time.After(time.Duration(10) * time.Second):
			// not an error -> server should ignore the request
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
			return
		}
	}
}

func TestStopReq(t *testing.T) {
	done = make(chan interface{})    // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	// connect client to server
	socketUrl := "ws://localhost:8080" + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	// handle server responses
	go receiveErrorHandler(conn, t)

	// send request to websocket server
	response := make(map[string]interface{})

	response["request"] = "stop"

	err = writeSocketResponse(conn, response)
	if err != nil {
		t.Errorf("Error: %v", err.Error())
		return
	}

	// wait for read handler to finish
	for {
		select {
		case <-interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			log.Println("Received SIGINT interrupt signal. Closing all pending connections")

			// Close our websocket connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
		case <-done:
			// log.Println("Receiver Channel Closed! Exiting....")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
			return
		case <-time.After(time.Duration(10) * time.Second):
			// not an error -> server should ignore the request
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
			return
		}
	}
}

func TestCPUReq(t *testing.T) {
	done = make(chan interface{})    // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	// connect client to server
	socketUrl := "ws://localhost:8080" + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	// handle server responses
	go receiveHandler(conn, t)

	// send request to websocket server
	response := make(map[string]interface{})

	response["request"] = "cpu"

	err = writeSocketResponse(conn, response)
	if err != nil {
		t.Errorf("Error: %v", err.Error())
		return
	}

	// wait for read handler to finish
	for {
		select {
		case <-interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			log.Println("Received SIGINT interrupt signal. Closing all pending connections")

			// Close our websocket connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
		case <-done:
			// log.Println("Receiver Channel Closed! Exiting....")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
			return
		case <-time.After(time.Duration(10) * time.Second):
			t.Errorf("Timeout in closing receiving channel. Exiting....")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
			return
		}
	}
}

func TestProcessReq(t *testing.T) {
	done = make(chan interface{})    // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	// connect client to server
	socketUrl := "ws://localhost:8080" + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	// handle server responses
	go receiveHandler(conn, t)

	// send request to websocket server
	response := make(map[string]interface{})

	response["request"] = "process_list"

	err = writeSocketResponse(conn, response)
	if err != nil {
		t.Errorf("Error: %v", err.Error())
		return
	}

	// wait for read handler to finish
	for {
		select {
		case <-interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			log.Println("Received SIGINT interrupt signal. Closing all pending connections")

			// Close our websocket connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
		case <-done:
			// log.Println("Receiver Channel Closed! Exiting....")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
			return
		case <-time.After(time.Duration(10) * time.Second):
			t.Errorf("Timeout in closing receiving channel. Exiting....")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
			return
		}
	}
}
