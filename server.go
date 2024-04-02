package main

import (
	"log"
	"tcp_chat/pkg/tcp"
)

func main() {
	server := tcp.CreateServer("localhost", "2000")

	go func() {
		for msg := range server.Msgchan {
			server.BroadcastMessage(string(msg))
		}
	}()

	log.Fatal(server.Start())
}
