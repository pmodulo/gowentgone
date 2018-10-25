package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	chat "github.com/pmodulo/gowentgone/01-chat-room"
)

func main() {
	mode := flag.String("mode", "server", "start in client or server mode")
	port := flag.Int("port", 6741, "port number for chat server")
	name := flag.String("name", "chat-server", "name of client or server")

	if strings.ToLower(*mode) == "server" {
		startServerMode(*port)
	} else {
		startClientMode(*port, *name)
	}
}

func startServerMode(port int) {
	fmt.Println("Starting chat-room server...")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("could not start chat server: %v", err)
	}
	manager := chat.ClientManager{
		clients:    make(map[*chat.Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *chat.Client),
		unregister: make(chan *chat.Client),
	}
	go manager.Start()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panicf("could not accept connection: %v", err)
		}
		client := &chat.Client{
			socket: conn,
			data:   make(chan []byte),
		}
		manager.register <- client
		go manager.Receive(client)
		go manager.Send(client)
	}
}

func startClientMode(port int, name string) {
	fmt.Printf("Starting client %v...\n", name)
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("cannot connect to chat server; shuttingdown...: %v", err)
	}
	client := &chat.Client{
		name:   name,
		socket: conn,
	}
	go client.Receive()
	reader := bufio.NewReader(os.Stdin)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("cannot read message from std input; shuttingdown...: %v", err)
		}
		if len(msg) > 0 {
			conn.Write([]byte(fmt.Sprintf("%v: %v", name, strings.TrimRight(msg, "\n"))))
		}
	}
}
