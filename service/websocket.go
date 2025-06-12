package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketHandler handles incoming websocket connections
type WebsocketServer interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	PushData(data WebsocketMetricMessage) error
	IsRunning() bool
}

type websocketServer struct {
	upgrader websocket.Upgrader
	logger   *log.Logger
	dataChan chan WebsocketMetricMessage
	clients  map[*websocket.Conn]bool // Map of active clients
	mu       sync.RWMutex             // Mutex for thread-safe client access
	running  bool                     // Indicates if the server is running
}

// PushData sends data to the websocket server's channel
func (wss *websocketServer) PushData(data WebsocketMetricMessage) error {
	select {
	case wss.dataChan <- data:
		return nil
	default:
		// channel full, drop oldest message to make room
		select {
		case <-wss.dataChan:
		default:
		}
		wss.dataChan <- data
		return nil
	}
}

// IsRunning checks if the websocket server is currently running
func (wss *websocketServer) IsRunning() bool {
	return wss.running
}

// broadcastMessages reads from dataChan and sends messages to all clients
func (wss *websocketServer) broadcastMessages() {
	for data := range wss.dataChan {
		wss.mu.RLock()
		for client := range wss.clients {

			jsonData, jsonErr := json.Marshal(data)
			if jsonErr != nil {
				log.Fatalf("Error marshaling to JSON: %v", jsonErr)
			}

			err := client.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				// Mark client for removal (cannot modify map during iteration)
				client.Close()
			}
		}
		wss.mu.RUnlock()

		// Clean up closed clients
		wss.mu.Lock()
		for client := range wss.clients {
			if client.WriteMessage(websocket.PingMessage, nil) != nil {
				delete(wss.clients, client)
			}
		}
		wss.mu.Unlock()
	}
}

// ServeHTTP upgrades the HTTP connection to a WebSocket connection and handles incoming messages
func (wss *websocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !wss.running {
		wss.running = true
	}

	c, err := wss.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}

	// Register client
	wss.mu.Lock()
	wss.clients[c] = true
	wss.mu.Unlock()

	defer func() {
		log.Println("closing connection")
		// Unregister client
		wss.mu.Lock()
		delete(wss.clients, c)
		wss.mu.Unlock()
		c.Close()
	}()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error %s when reading message from client", err)
			return
		}
		if mt == websocket.BinaryMessage {
			err = c.WriteMessage(websocket.TextMessage, []byte("server doesn't support binary messages"))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
			}
			return
		}
		log.Printf("Receive message %s", string(message))
		if strings.Trim(string(message), "\n") != "start" {
			err = c.WriteMessage(websocket.TextMessage, []byte("You did not say the magic word!"))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				return
			}
			continue
		}
		// Client sent "start"; it will now receive broadcasted messages
		log.Println("client subscribed to messages")
	}
}

func NewWebsocketServer(logger *log.Logger) WebsocketServer {
	wss := &websocketServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for simplicity, adjust as needed
			},
		},
		logger:   logger,
		dataChan: make(chan WebsocketMetricMessage, 1000), // Buffered channel to hold messages
		clients:  make(map[*websocket.Conn]bool),
		running:  false,
	}
	// Start broadcasting goroutine
	go wss.broadcastMessages()
	return wss
}
