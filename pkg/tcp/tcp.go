package tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type ConnectionMap map[int]Connection

type Server struct {
	listenAddr    string
	listener      net.Listener
	Msgchan       chan []byte
	quitchan      chan struct{}
	connectionMap ConnectionMap
}

type Connection struct {
	stream net.Conn
}

func (s *Server) BroadcastMessage(message string) {
	for _, c := range s.connectionMap {
		go s.WriteMessage(c, []byte(message+"\n"))
	}
}

func (s *Server) WriteMessage(c Connection, b []byte) {
	_, err := c.stream.Write(b)
	if err != nil {
		log.Println(err)
	}
}

func CreateServer(host string, port string) *Server {
	return &Server{
		listenAddr:    host + ":" + port,
		Msgchan:       make(chan []byte, 10),
		quitchan:      make(chan struct{}),
		connectionMap: make(ConnectionMap),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	s.listener = listener

	go s.connLoop()

	<-s.quitchan
	close(s.Msgchan)

	return nil
}

func (s *Server) connLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(c net.Conn) {
	defer c.Close()

	s.connectionMap[len(s.connectionMap)] = Connection{stream: c}

	message := bufio.NewScanner(c)

	fmt.Printf("New connection from: %s\n", c.LocalAddr().String())

	for message.Scan() {
		s.Msgchan <- message.Bytes()
	}
}
