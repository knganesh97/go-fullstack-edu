package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/knganesh97/go-fullstack-edu/pkg/models"
)

var (
	todos     []models.Todo
	todoMutex sync.Mutex
	todoID    int
)

// Initialize with some sample data
func init() {
	todos = []models.Todo{
		{ID: 1, Title: "Learn Go", Completed: true},
		{ID: 2, Title: "Learn Vecty", Completed: false},
		{ID: 3, Title: "Build a web app", Completed: false},
	}
	todoID = 3
}

func main() {
	// Define the port to listen on
	port := ":8080"

	// Set up file server for static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// API endpoints
	http.HandleFunc("/api/todos", handleTodos)

	// Serve the frontend - all routes serve the same index.html
	// which loads our WebAssembly application
	http.HandleFunc("/", serveIndex)

	// Start the server
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// serveIndex serves the index.html file which loads our WebAssembly
func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

// handleTodos handles CRUD operations for todos
func handleTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// Return all todos
		todoMutex.Lock()
		json.NewEncoder(w).Encode(todos)
		todoMutex.Unlock()

	case http.MethodPost:
		// Create a new todo
		var todo models.Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todoMutex.Lock()
		todoID++
		todo.ID = todoID
		todos = append(todos, todo)
		todoMutex.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(todo)

	case http.MethodPut:
		// Update an existing todo
		var updatedTodo models.Todo
		if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todoMutex.Lock()
		found := false
		for i, t := range todos {
			if t.ID == updatedTodo.ID {
				todos[i] = updatedTodo
				found = true
				break
			}
		}
		todoMutex.Unlock()

		if !found {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(updatedTodo)

	case http.MethodDelete:
		// Delete a todo
		var todoToDelete models.Todo
		if err := json.NewDecoder(r.Body).Decode(&todoToDelete); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todoMutex.Lock()
		found := false
		for i, t := range todos {
			if t.ID == todoToDelete.ID {
				todos = append(todos[:i], todos[i+1:]...)
				found = true
				break
			}
		}
		todoMutex.Unlock()

		if !found {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
