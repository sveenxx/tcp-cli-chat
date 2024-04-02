package main

import (
	"github.com/sveenxx/tcp-cli-chat/pkg/tcp"
	"log"
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
