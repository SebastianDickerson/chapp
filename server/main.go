package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type broadcastMsg struct {
	msg    string
	sender net.Conn
}

var (
	clients   = make(map[net.Conn]string)
	mu        sync.Mutex
	broadcast = make(chan broadcastMsg)
)

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Server listening on :9000")

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to accept connection: %v\n", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	name, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	name = strings.TrimSpace(name)
	if name == "" {
		name = conn.RemoteAddr().String()
	}

	mu.Lock()
	clients[conn] = name
	mu.Unlock()

	broadcast <- broadcastMsg{msg: fmt.Sprintf("🟢 %s has joined the chat", name)}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		broadcast <- broadcastMsg{msg: fmt.Sprintf("%s: %s", name, line), sender: conn}
	}

	mu.Lock()
	delete(clients, conn)
	mu.Unlock()
	broadcast <- broadcastMsg{msg: fmt.Sprintf("🔴 %s has left the chat", name)}
}

func broadcaster() {
	for bmsg := range broadcast {
		fmt.Println(bmsg.msg)
		mu.Lock()
		for conn := range clients {
			if bmsg.sender != nil && conn == bmsg.sender {
				continue
			}
			fmt.Fprintln(conn, bmsg.msg)
		}
		mu.Unlock()
	}
}
