package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New incoming Connection From Client: ", ws.RemoteAddr())
	// maps are not concurrent(thread) safe so mutex can be used for production
	s.conns[ws] = true
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				// Connection from other side is closed!
				break
			}
			fmt.Println("Read Error: ", err)
			continue
		}
		msg := buf[:n]
		fmt.Println("Message:", string(msg))
		ws.Write([]byte("Thank you for the Message!!"))
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.ListenAndServe(":3550", nil)
}
