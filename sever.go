package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

// Global map of connected agents
var agents = make(map[string]*Agent)
var agentsMutex = sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Agent represents a connected tunnel agent
type Agent struct {
	ID   string
	Conn *websocket.Conn
}

func main() {
	// Use Railway PORT or default 8080
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	// WebSocket endpoint for agents
	http.HandleFunc("/agent", agentHandler)

	fmt.Println("Tunnel server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// agentHandler handles agent WebSocket connections
func agentHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	agentID := r.URL.Query().Get("id")
	if agentID == "" {
		agentID = "agent1"
	}

	agent := &Agent{ID: agentID, Conn: conn}

	agentsMutex.Lock()
	agents[agentID] = agent
	agentsMutex.Unlock()

	fmt.Println("Agent connected:", agentID)

	defer func() {
		conn.Close()
		agentsMutex.Lock()
		delete(agents, agentID)
		agentsMutex.Unlock()
		fmt.Println("Agent disconnected:", agentID)
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// For now, just print received messages from agent
		fmt.Printf("Received from %s: %d bytes\n", agentID, len(msg))
	}
}

// ====== TCP client handling ======
// This function shows how you would handle a TCP client
// For Railway, raw TCP may not work on extra ports, so we can wrap TCP over WebSocket
func handleTCPClient(agentID string, client net.Conn) {
	defer client.Close()

	agentsMutex.Lock()
	agent, ok := agents[agentID]
	agentsMutex.Unlock()

	if !ok {
		fmt.Println("No agent found for ID:", agentID)
		return
	}

	fmt.Println("Forwarding client", client.RemoteAddr(), "to agent", agentID)

	buf := make([]byte, 1024)
	for {
		n, err := client.Read(buf)
		if err != nil {
			break
		}
		// Send to agent via WebSocket
		err = agent.Conn.WriteMessage(websocket.BinaryMessage, buf[:n])
		if err != nil {
			break
		}
	}
}
