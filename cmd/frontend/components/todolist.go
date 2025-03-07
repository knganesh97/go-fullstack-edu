//go:build js && wasm
// +build js,wasm

package components

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"
	"github.com/knganesh97/go-fullstack-edu/pkg/models"
)

// TodoList is a component that displays and manages todos
type TodoList struct {
	vecty.Core
	Todos        []models.Todo
	NewTodoTitle string
}

// Render implements the vecty.Component interface
func (t *TodoList) Render() vecty.ComponentOrHTML {
	return elem.Section(
		vecty.Markup(
			vecty.Class("todo-section"),
			prop.ID("todos"),
		),
		t.renderTodoForm(),
		t.renderTodoList(),
	)
}

// renderTodoForm renders the form for adding new todos
func (t *TodoList) renderTodoForm() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("todo-form"),
		),
		elem.Heading2(
			vecty.Text("Add New Todo"),
		),
		elem.Form(
			vecty.Markup(
				prop.ID("new-todo-form"),
				event.Submit(t.onAddTodo).PreventDefault(),
			),
			elem.Input(
				vecty.Markup(
					prop.ID("new-todo-title"),
					prop.Type(prop.TypeText),
					prop.Placeholder("What needs to be done?"),
					vecty.Property("required", true),
					prop.Value(t.NewTodoTitle),
					event.Input(func(e *vecty.Event) {
						t.NewTodoTitle = e.Target.Get("value").String()
					}),
				),
			),
			elem.Button(
				vecty.Markup(
					prop.Type(prop.TypeSubmit),
					vecty.Class("btn"),
				),
				vecty.Text("Add Todo"),
			),
		),
	)
}

// renderTodoList renders the list of todos
func (t *TodoList) renderTodoList() vecty.ComponentOrHTML {
	if t.Todos == nil {
		// Load todos if not already loaded
		go t.loadTodos()
		return elem.Div(
			vecty.Markup(
				vecty.Class("todo-list"),
			),
			elem.Heading2(
				vecty.Text("Loading Todos..."),
			),
		)
	}

	var todoItems []vecty.MarkupOrChild
	for _, todo := range t.Todos {
		todoItems = append(todoItems, t.renderTodoItem(todo))
	}

	return elem.Div(
		vecty.Markup(
			vecty.Class("todo-list"),
		),
		elem.Heading2(
			vecty.Text("Your Todos"),
		),
		elem.UnorderedList(
			append(
				[]vecty.MarkupOrChild{
					vecty.Markup(
						prop.ID("todos"),
					),
				},
				todoItems...,
			)...,
		),
	)
}

// renderTodoItem renders a single todo item
func (t *TodoList) renderTodoItem(todo models.Todo) vecty.ComponentOrHTML {
	return elem.ListItem(
		vecty.Markup(
			vecty.Class("todo-item"),
			vecty.MarkupIf(todo.Completed, vecty.Class("completed")),
			vecty.Property("data-id", todo.ID),
		),
		elem.Input(
			vecty.Markup(
				vecty.Class("todo-checkbox"),
				prop.Type(prop.TypeCheckbox),
				prop.Checked(todo.Completed),
				event.Change(func(e *vecty.Event) {
					todo.Completed = e.Target.Get("checked").Bool()
					go t.updateTodo(todo)
				}),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("todo-title"),
			),
			vecty.Text(todo.Title),
		),
		elem.Button(
			vecty.Markup(
				vecty.Class("delete-btn"),
				event.Click(func(e *vecty.Event) {
					go t.deleteTodo(todo)
				}),
			),
			vecty.Text("Delete"),
		),
	)
}

// loadTodos loads todos from the API
func (t *TodoList) loadTodos() {
	resp, err := http.Get("/api/todos")
	if err != nil {
		// Handle error
		return
	}
	defer resp.Body.Close()

	var todos []models.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		// Handle error
		return
	}

	t.Todos = todos
	vecty.Rerender(t)
}

// onAddTodo handles the form submission to add a new todo
func (t *TodoList) onAddTodo(e *vecty.Event) {
	if t.NewTodoTitle == "" {
		return
	}

	todo := models.Todo{
		Title:     t.NewTodoTitle,
		Completed: false,
	}

	// Reset the input field
	t.NewTodoTitle = ""

	// Send the new todo to the server
	go func() {
		jsonData, err := json.Marshal(todo)
		if err != nil {
			// Handle error
			return
		}

		req, err := http.NewRequest("POST", "/api/todos", strings.NewReader(string(jsonData)))
		if err != nil {
			// Handle error
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			// Handle error
			return
		}
		defer resp.Body.Close()

		var newTodo models.Todo
		if err := json.NewDecoder(resp.Body).Decode(&newTodo); err != nil {
			// Handle error
			return
		}

		t.Todos = append(t.Todos, newTodo)
		vecty.Rerender(t)
	}()
}

// updateTodo updates a todo on the server
func (t *TodoList) updateTodo(todo models.Todo) {
	jsonData, err := json.Marshal(todo)
	if err != nil {
		// Handle error
		return
	}

	req, err := http.NewRequest("PUT", "/api/todos", strings.NewReader(string(jsonData)))
	if err != nil {
		// Handle error
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Handle error
		return
	}
	defer resp.Body.Close()

	// Update the local todo list
	for i, existingTodo := range t.Todos {
		if existingTodo.ID == todo.ID {
			t.Todos[i] = todo
			break
		}
	}

	vecty.Rerender(t)
}

// deleteTodo deletes a todo from the server
func (t *TodoList) deleteTodo(todo models.Todo) {
	jsonData, err := json.Marshal(todo)
	if err != nil {
		// Handle error
		return
	}

	req, err := http.NewRequest("DELETE", "/api/todos", strings.NewReader(string(jsonData)))
	if err != nil {
		// Handle error
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Handle error
		return
	}
	defer resp.Body.Close()

	// Remove the todo from the local list
	for i, existingTodo := range t.Todos {
		if existingTodo.ID == todo.ID {
			t.Todos = append(t.Todos[:i], t.Todos[i+1:]...)
			break
		}
	}

	vecty.Rerender(t)
}
