# chapp

A terminal-based chat application written in Go, featuring a simple UI using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Bubbles](https://github.com/charmbracelet/bubbles) libraries.

## Features
- Real-time chat between multiple clients
- Simple, modern terminal UI
- Message history and sender highlighting
- Easy to run locally

## Project Structure
```
go.mod         # Go module definition
go.sum         # Go dependencies
README.md      # Project documentation
client/        # Client application code
server/        # Server application code
ui/            # Terminal UI components
```

## Getting Started

### Prerequisites
- Go 1.18 or newer

### Install dependencies
```
go mod tidy
```

### Running the Server
```
go run ./server
```

### Running the Client
In a new terminal window:

Enter a username as an argument when running the client, e.g.:
```
go run ./client your_username
```

## How It Works
- The server listens for incoming TCP connections and broadcasts messages to all connected clients.
- Each client connects to the server, sends messages, and receives updates in real time.
- The UI is built with Bubble Tea and Bubbles, providing a responsive and user-friendly terminal interface.

## Customization
- UI code is located in `ui/ui.go` and can be modified to change the look and feel.
- Message handling and networking logic can be extended for more features (e.g., private messages, authentication).
