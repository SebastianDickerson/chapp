package main

import (
	"log"
	"net"
	"os"

	"chapp/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <username>", os.Args[0])
	}

	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer conn.Close()

	p := tea.NewProgram(ui.InitialModel(conn, os.Args[1]), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
