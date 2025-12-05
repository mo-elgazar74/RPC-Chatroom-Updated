package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Client struct {
	id   string
	conn net.Conn
	out  chan string
}

type Server struct {
	mu      sync.Mutex
	clients map[string]*Client
}

func NewServer() *Server {
	return &Server{
		clients: make(map[string]*Client),
	}
}

func (s *Server) broadcast(senderID, msg string, includeSender bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, c := range s.clients {
		if !includeSender && id == senderID {
			continue
		}
		select {
		case c.out <- msg:
		default:
			log.Printf("dropping message for client %s (out channel full)", id)
		}
	}
}

func (s *Server) addClient(c *Client) {
	s.mu.Lock()
	s.clients[c.id] = c
	s.mu.Unlock()
}

func (s *Server) removeClient(id string) {
	s.mu.Lock()
	if c, ok := s.clients[id]; ok {
		delete(s.clients, id)
		close(c.out)
		c.conn.Close()
	}
	s.mu.Unlock()
}

func handleClient(srv *Server, conn net.Conn) {
	reader := bufio.NewReader(conn)

	nameLine, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("failed to read client name: %v", err)
		conn.Close()
		return
	}
	id := strings.TrimSpace(nameLine)
	if id == "" {
		id = "Anonymous"
	}

	client := &Client{
		id:   id,
		conn: conn,
		out:  make(chan string, 16),
	}
	srv.addClient(client)

	joinMsg := fmt.Sprintf("User [%s] joined\n", client.id)
	log.Printf(joinMsg)
	srv.broadcast(client.id, joinMsg, false)

	go func(c *Client) {
		writer := bufio.NewWriter(c.conn)
		for msg := range c.out {
			_, err := writer.WriteString(msg)
			if err != nil {
				log.Printf("write error to %s: %v", c.id, err)
				return
			}
			if err := writer.Flush(); err != nil {
				log.Printf("flush error to %s: %v", c.id, err)
				return
			}
		}
	}(client)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("client %s disconnected: %v", client.id, err)
			break
		}
		text := strings.TrimSpace(line)
		if text == "" {
			continue
		}

		if strings.EqualFold(text, "exit") {
			log.Printf("client %s requested exit", client.id)
			break
		}

		msg := fmt.Sprintf("[%s]: %s\n", client.id, text)
		srv.broadcast(client.id, msg, false) 
	}

	srv.removeClient(client.id)
	leaveMsg := fmt.Sprintf("User [%s] left\n", client.id)
	log.Printf(leaveMsg)
	srv.broadcast(client.id, leaveMsg, false)
}

func main() {
	port := "1234"
	addr := ":" + port
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen error on %s: %v", addr, err)
	}
	log.Printf("Chat server running on %s ...", addr)

	server := NewServer()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go handleClient(server, conn) 
	}
}
