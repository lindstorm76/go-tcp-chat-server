package main

import (
	"log"
	"net"

	"github.com/lindstorm76/go-tcp-chat-server/internal"
)

func main() {
	server := internal.NewServer()
	go server.Run()

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalf("failed to start tcp server: %s", err.Error())
	}

	log.Print("the server is up and running on :8080")

	for {
		clientConn, err := listener.Accept()

		if err != nil {
			log.Printf("failed to accept client connection: %s", err.Error())
		}

		log.Printf("a client has connected to the server: %s", clientConn.RemoteAddr())

		go internal.NewClient(clientConn, server.CommandChan)
	}
}
