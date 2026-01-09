package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

// Map of connected agents
var agents = make(map[string]*websocket.Conn)
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func main() {
	// WebSocket endpoint for agents
	http.HandleFunc("/agent", agentHandler)

	go startTCPServer(4000) // Public TCP port for clients

	fmt.Println("Tunnel server running on :8080 for agents")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func agentHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// For demo, simple agent ID from query
	agentID := r.URL.Query().Get("id")
	if agentID == "" {
		agentID = "agent1"
	}

	agents[agentID] = conn
	fmt.Println("Agent connected:", agentID)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Agent disconnected:", agentID)
			delete(agents, agentID)
			break
		}
		// Just print messages for now
		fmt.Println("Message from agent:", string(msg))
	}
}

// TCP server for clients
func startTCPServer(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Public TCP server listening on port", port)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println("Client accept error:", err)
			continue
		}
		go handleClient(clientConn)
	}
}

// Handle TCP client
func handleClient(client net.Conn) {
	defer client.Close()
	fmt.Println("Client connected:", client.RemoteAddr())

	// Demo: just echo back for now
	buf := make([]byte, 1024)
	for {
		n, err := client.Read(buf)
		if err != nil {
			fmt.Println("Client disconnected:", client.RemoteAddr())
			break
		}
		client.Write(buf[:n])
	}
}
