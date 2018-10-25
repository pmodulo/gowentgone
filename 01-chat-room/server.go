package chat

import (
	"fmt"
)

// ClientManager is a server side data structure
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// Start the server worker thread
func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			fmt.Println("Accepted a new conection")
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.data)
				delete(manager.clients, conn)
				fmt.Println("Closed a connection")
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.data <- message:
				default:
					close(conn.data)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

// Receive reads data from client. This is one per client
func (manager *ClientManager) Receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Printf("Received: %v", string(message))
			manager.broadcast <- message
		}
	}
}

// Send writes data to client. This is one per client
func (manager *ClientManager) Send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			client.socket.Write(message)
		}
	}
}
