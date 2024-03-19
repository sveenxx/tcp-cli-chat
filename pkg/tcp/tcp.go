package tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type ConnectionMap map[int]Connection

type Connection struct {
	stream   net.Conn
	nickname string
}

func broadcastMessage(message string, ConnectionMap ConnectionMap) {
	for _, connection := range ConnectionMap {
		go func(connection Connection) {
			_, err := connection.stream.Write([]byte(message + "\n"))
			if err != nil {
				return
			}
		}(connection)
	}
}

func CreateTCPServer(port string) {
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	connectionMap := make(ConnectionMap)

	connId := 0

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		connId++

		connectionMap[connId] = Connection{stream: conn, nickname: ""}

		go handleConnection(conn, connectionMap)
	}
}

func handleConnection(conn net.Conn, connectionMap ConnectionMap) {
	defer conn.Close()

	message := bufio.NewScanner(conn)

	fmt.Printf("New connection from: %s\n", conn.LocalAddr().String())

	for message.Scan() {
		broadcastMessage(message.Text(), connectionMap)
	}
}
