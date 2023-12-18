package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func (s *Server) handleWSStock(ws *websocket.Conn) {
	fmt.Println("New Incoming Connection from client to Stock Feed:", ws.RemoteAddr())
	price := strings.Split(ws.Request().URL.RawQuery, "=")[1]

	if price == "" {
		payload := "N/A"
		ws.Write([]byte(payload))
	} else {
		for {
			parsedPrice, err := strconv.Atoi(price)
			if err != nil {
				fmt.Println("Error: Parsing Price error")
			}

			min := parsedPrice - 5
			max := parsedPrice + 5
			randomNumber := rand.Intn(max-min+1)%(max-min+1) + min

			payload := strconv.Itoa(randomNumber)
			ws.Write([]byte(payload))
			time.Sleep(time.Second * 1)
		}
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
		s.broadcast(msg)
		fmt.Println("Message:", string(msg))
		ws.Write([]byte("Thank you for the Message!!"))
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("Write Error:", err)
			}
		}(ws)
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/stock", websocket.Handler(server.handleWSStock))
	http.ListenAndServe(":3550", nil)
}
