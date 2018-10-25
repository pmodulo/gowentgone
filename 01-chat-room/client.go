package chat

import (
	"fmt"
	"log"
	"net"
)

// Client is a client side data structure
type Client struct {
	Name   string
	Socket net.Conn
	Data   chan []byte
}

// Receive thread which reads data send by chat server
func (client *Client) Receive() {
	for {
		message := make([]byte, 4096)
		len, err := client.Socket.Read(message)
		if err != nil {
			client.Socket.Close()
			log.Fatalf("error reading data from chat server; shuttingdown: %v", err)
		}
		if len > 0 {
			fmt.Printf(">>: %v\n", string(message))
		}
	}
}
