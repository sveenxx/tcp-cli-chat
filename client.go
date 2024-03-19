package main

import (
	"bufio"
	"fmt"
	screen "github.com/aditya43/clear-shell-screen-golang"
	"net"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func receiveServerOutput(conn net.Conn) {
	reader := bufio.NewScanner(conn)
	messages := make([]string, 0)

	for reader.Scan() {
		screen.Clear()

		messages = append(messages, reader.Text())

		if len(messages) > 10 {
			messages = messages[len(messages)-10:]
		}

		for _, message := range messages {
			fmt.Println(message)
		}

		fmt.Print("> Enter text: ")
	}
}

func setUpNickName(n *string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your nickname: ")
	nickname, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	*n = strings.TrimRight(nickname, "\r\n")

	screen.Clear()
}

func sendClientInput(conn net.Conn, nickName string) {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n> Enter text: ")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}

		_, err = writer.WriteString(nickName + ": " + message)

		if err != nil {
			fmt.Println("Error writing:", err.Error())
		}

		err = writer.Flush()
		if err != nil {
			fmt.Println("Error flushing:", err.Error())
		}
	}
}

const (
	ConnHost = "localhost"
	ConnPort = "2000"
)

func main() {
	var nickName string
	setUpNickName(&nickName)

	conn, err := net.Dial("tcp", ConnHost+":"+ConnPort)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	go receiveServerOutput(conn)
	sendClientInput(conn, nickName)
}
