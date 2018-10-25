package chat

import (
	"fmt"
	"log"
	"net"
)

// Client is a client side data structure
type Client struct {
	name   string
	socket net.Conn
	data   chan []byte
}

// Receive thread which reads data send by chat server
func (client *Client) Receive() {
	for {
		message := make([]byte, 4096)
		len, err := client.socket.Read(message)
		if err != nil {
			client.socket.Close()
			log.Fatalf("error reading data from chat server; shuttingdown: %v", err)
		}
		if len > 0 {
			fmt.Printf(">>: %v\n", string(message))
		}
	}
}
