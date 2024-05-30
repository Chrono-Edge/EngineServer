package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Chrono-Edge/EngineServer/internal/game"
	"github.com/Chrono-Edge/EngineServer/internal/network"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var gameState = game.NewGameState()

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}(ws)

	client := network.NewClient(ws)
	go client.WritePump()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("read:", err)
			break
		}
		log.Printf("recv: %s", msg)

		// TODO: Commands processing
		switch string(msg) {
		case "join":
			client.Send <- []byte("Welcome to ChronoEdge!")
		case "move":
			gameState.Update()
			client.Send <- []byte("Player moved")
		default:
			client.Send <- []byte("Unknown command")

		}
	}
}

func main() {
	// Open log file
	logFile, err := os.OpenFile("logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Could not open log file: %v", err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			log.Fatalf("Failed to close log file: %v", err)
		}
	}(logFile)

	// Create multi-writer to write to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Set log output to multi-writer
	log.SetOutput(multiWriter)
	log.Println("Server is starting...")
	http.HandleFunc("/ws", handleConnections)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
