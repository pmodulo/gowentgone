package chat

import (
	"fmt"
)

// ClientManager is a server side data structure
type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

// Start the server worker thread
func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.Register:
			manager.Clients[conn] = true
			fmt.Println("Accepted a new conection")
		case conn := <-manager.Unregister:
			if _, ok := manager.Clients[conn]; ok {
				close(conn.Data)
				delete(manager.Clients, conn)
				fmt.Println("Closed a connection")
			}
		case message := <-manager.Broadcast:
			for conn := range manager.Clients {
				select {
				case conn.Data <- message:
				default:
					close(conn.Data)
					delete(manager.Clients, conn)
				}
			}
		}
	}
}

// Receive reads data from client. This is one per client
func (manager *ClientManager) Receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.Socket.Read(message)
		if err != nil {
			manager.Unregister <- client
			client.Socket.Close()
			break
		}
		if length > 0 {
			fmt.Printf("Received: %v", string(message))
			manager.Broadcast <- message
		}
	}
}

// Send writes data to client. This is one per client
func (manager *ClientManager) Send(client *Client) {
	defer client.Socket.Close()
	for {
		select {
		case message, ok := <-client.Data:
			if !ok {
				return
			}
			client.Socket.Write(message)
		}
	}
}
