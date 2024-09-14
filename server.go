package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"strconv"
	"time"
)

// Function to handle client commands over the socket connection
func handleConnection(conn net.Conn, todos *Todos) {
	defer conn.Close()

	// Reader for the client's command
	reader := bufio.NewReader(conn)

	for {
		// Send prompt to the client
		conn.Write([]byte("Enter command (e.g., add, list, edit, delete, toggle): "))

		// Read the client's input command
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}

		// Process the input
		command := strings.TrimSpace(input)

		// Parse and execute command
		switch {
		case strings.HasPrefix(command, "list"):
			conn.Write([]byte(todosList(todos)))

		case strings.HasPrefix(command, "add"):
			title := strings.TrimPrefix(command, "add ")
			todos.add(title)
			conn.Write([]byte("Todo added: " + title + "\n"))

		case strings.HasPrefix(command, "edit"):
			parts := strings.SplitN(command, " ", 3)
			if len(parts) != 3 {
				conn.Write([]byte("Invalid format for edit. Use: edit <index> <new_title>\n"))
				continue
			}
			index, err := strconv.Atoi(parts[1])
			if err != nil {
				conn.Write([]byte("Invalid index\n"))
				continue
			}
			newTitle := parts[2]
			if err := todos.edit(index, newTitle); err != nil {
				conn.Write([]byte(err.Error() + "\n"))
			} else {
				conn.Write([]byte("Todo edited successfully\n"))
			}

		case strings.HasPrefix(command, "delete"):
			parts := strings.Split(command, " ")
			if len(parts) != 2 {
				conn.Write([]byte("Invalid format for delete. Use: delete <index>\n"))
				continue
			}
			index, err := strconv.Atoi(parts[1])
			if err != nil {
				conn.Write([]byte("Invalid index\n"))
				continue
			}
			if err := todos.delete(index); err != nil {
				conn.Write([]byte(err.Error() + "\n"))
			} else {
				conn.Write([]byte("Todo deleted successfully\n"))
			}

		case strings.HasPrefix(command, "toggle"):
			parts := strings.Split(command, " ")
			if len(parts) != 2 {
				conn.Write([]byte("Invalid format for toggle. Use: toggle <index>\n"))
				continue
			}
			index, err := strconv.Atoi(parts[1])
			if err != nil {
				conn.Write([]byte("Invalid index\n"))
				continue
			}
			if err := todos.toggle(index); err != nil {
				conn.Write([]byte(err.Error() + "\n"))
			} else {
				conn.Write([]byte("Todo toggled successfully\n"))
			}

		case command == "exit":
			conn.Write([]byte("Goodbye!\n"))
			return

		default:
			conn.Write([]byte("Invalid command\n"))
		}
	}
}

// Function to convert the Todos list to a printable string
func todosList(todos *Todos) string {
	var sb strings.Builder
	for index, todo := range *todos {
		completed := "❌"
		if todo.Completed {
			completed = "✔️"
		}
		sb.WriteString(fmt.Sprintf("%d. %s [%s] (Created at: %s)\n", index, todo.Title, completed, todo.CreatedAt.Format(time.RFC1123)))
	}
	return sb.String()
}

// Start the TCP server on a specified port
func StartTCPServer(todos *Todos, port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on port", port)

	// Continuously listen for incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn, todos)
	}
}
