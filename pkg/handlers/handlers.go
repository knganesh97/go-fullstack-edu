package handlers

import (
	"encoding/json"
	"html/template"
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
		{ID: 2, Title: "Build a web app", Completed: false},
		{ID: 3, Title: "Deploy to production", Completed: false},
	}
	todoID = 3
}

// HomeHandler renders the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TodosHandler renders the todos page
func TodosHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/todos.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todoMutex.Lock()
	data := struct {
		Todos []models.Todo
	}{
		Todos: todos,
	}
	todoMutex.Unlock()

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// APITodosHandler handles CRUD operations for todos via JSON API
func APITodosHandler(w http.ResponseWriter, r *http.Request) {
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
