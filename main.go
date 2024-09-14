package main

func main() {
	todos := Todos{}
	
	// Load todos from file
	storage := NewStorage[Todos]("todo.json")
	err := storage.Load(&todos)
	if err != nil {
		todos = Todos{}
	}

	// Start the TCP server
	StartTCPServer(&todos, "8080")

	// Save the todos periodically
	storage.Save(todos)
}
