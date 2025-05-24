package ui

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

type (
	errMsg             error
	messageReceivedMsg string
)

// Model defines the structure for our application.
type Model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
	Conn        net.Conn
	ClientName  string
}

// messageChannel wraps a channel for incoming messages.
type messageChannel struct {
	ch chan string
}

func newMessageChannel() *messageChannel {
	return &messageChannel{ch: make(chan string)}
}

// startMessageListener launches a goroutine to read messages from the connection.
func startMessageListener(conn net.Conn, ch *messageChannel) {
	go func() {
		reader := bufio.NewReader(conn)
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				ch.ch <- "[error] " + err.Error()
				close(ch.ch)
				return
			}
			ch.ch <- strings.TrimSpace(msg)
		}
	}()
}

// waitForMessage returns a Bubble Tea command that waits for a message from the channel.
func waitForMessage(ch *messageChannel) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch.ch
		if !ok {
			return errMsg(fmt.Errorf("connection closed"))
		}
		return messageReceivedMsg(msg)
	}
}

var msgChan *messageChannel

// InitialModel creates the initial model for the application.
func InitialModel(conn net.Conn, clientName string) Model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 280
	ta.SetWidth(30)
	ta.SetHeight(3)
	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	// Disable insert newline
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	msgChan = newMessageChannel()
	// Send the client name as the first message to the server
	_, _ = conn.Write([]byte(clientName + "\n"))
	startMessageListener(conn, msgChan)

	return Model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		Conn:        conn,
		ClientName:  clientName,
	}
}

func (m Model) Init() tea.Cmd {
	return waitForMessage(msgChan)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.textarea, _ = m.textarea.Update(msg)
	m.viewport, _ = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)
		if len(m.messages) > 0 {
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			val := m.textarea.Value()
			if strings.TrimSpace(val) != "" {
				m.messages = append(m.messages, m.senderStyle.Render("You: ")+val)
				_, _ = m.Conn.Write([]byte(val + "\n"))
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
				m.textarea.Reset()
				m.viewport.GotoBottom()
			}
		}
	case messageReceivedMsg:
		// Try to extract sender name from the message, e.g. "x<name>: <msg>"
		msgStr := string(msg)
		if strings.HasPrefix(msgStr, "x") && strings.Contains(msgStr, ": ") {
			parts := strings.SplitN(msgStr[1:], ": ", 2)
			if len(parts) == 2 {
				name := parts[0]
				content := parts[1]
				m.messages = append(m.messages, m.senderStyle.Render(name+": ")+content)
			} else {
				m.messages = append(m.messages, m.senderStyle.Render(msgStr))
			}
		} else {
			m.messages = append(m.messages, m.senderStyle.Render(msgStr))
		}
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.textarea.Reset()
		m.viewport.GotoBottom()
		return m, waitForMessage(msgChan)
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf("%s%s%s", m.viewport.View(), gap, m.textarea.View())
}
