package main

import (
	"log"
	"tcp_server/pkg/tcp"
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
