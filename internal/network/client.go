package network

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn: conn,
		Send: make(chan []byte),
	}
}

func (c *Client) ReadPump() {
	defer func(Conn *websocket.Conn) {
		err := Conn.Close()
		if err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}(c.Conn)
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// TODO: Incoming messages processing
		c.Send <- message
	}
}
func (c *Client) WritePump() {
	defer func(Conn *websocket.Conn) {
		err := Conn.Close()
		if err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}(c.Conn)
	for {
		message := <-c.Send
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}
