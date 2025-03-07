//go:build js && wasm
// +build js,wasm

package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/knganesh97/go-fullstack-edu/pkg/models"
)

// Application is the main component of our frontend
type Application struct {
	vecty.Core
	Todos       []models.Todo
	CurrentView string // "home" or "todos"
}

// Render implements the vecty.Component interface
func (a *Application) Render() vecty.ComponentOrHTML {
	// Initialize current view if not set
	if a.CurrentView == "" {
		a.CurrentView = "home"
	}

	return elem.Body(
		vecty.Markup(
			vecty.Class("container"),
		),
		a.renderHeader(),
		a.renderMain(),
		a.renderFooter(),
	)
}

// renderHeader renders the header section
func (a *Application) renderHeader() vecty.ComponentOrHTML {
	return elem.Header(
		vecty.Markup(
			vecty.Class("header"),
		),
		elem.Heading1(
			vecty.Text("Go Fullstack App with Vecty"),
		),
		elem.Navigation(
			elem.UnorderedList(
				elem.ListItem(
					elem.Anchor(
						vecty.Markup(
							vecty.MarkupIf(a.CurrentView == "home", vecty.Class("active")),
							event.Click(func(e *vecty.Event) {
								a.CurrentView = "home"
								vecty.Rerender(a)
								e.Call("preventDefault")
							}),
							vecty.Property("href", "#"),
						),
						vecty.Text("Home"),
					),
				),
				elem.ListItem(
					elem.Anchor(
						vecty.Markup(
							vecty.MarkupIf(a.CurrentView == "todos", vecty.Class("active")),
							event.Click(func(e *vecty.Event) {
								a.CurrentView = "todos"
								vecty.Rerender(a)
								e.Call("preventDefault")
							}),
							vecty.Property("href", "#"),
						),
						vecty.Text("Todos"),
					),
				),
			),
		),
	)
}

// renderMain renders the main content section
func (a *Application) renderMain() vecty.ComponentOrHTML {
	if a.CurrentView == "todos" {
		return elem.Main(
			&TodoList{},
		)
	}

	// Home view
	return elem.Main(
		elem.Section(
			vecty.Markup(
				vecty.Class("hero"),
			),
			elem.Heading2(
				vecty.Text("A Simple Todo Application"),
			),
			elem.Paragraph(
				vecty.Text("Built with Go and Vecty - both frontend and backend!"),
			),
			elem.Button(
				vecty.Markup(
					vecty.Class("btn"),
					event.Click(func(e *vecty.Event) {
						a.CurrentView = "todos"
						vecty.Rerender(a)
					}),
				),
				vecty.Text("View Todos"),
			),
		),
		elem.Section(
			vecty.Markup(
				vecty.Class("features"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("feature"),
				),
				elem.Heading3(
					vecty.Text("Go Backend"),
				),
				elem.Paragraph(
					vecty.Text("Powered by Go's standard library for HTTP handling"),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("feature"),
				),
				elem.Heading3(
					vecty.Text("Vecty Frontend"),
				),
				elem.Paragraph(
					vecty.Text("Frontend built with Vecty and WebAssembly"),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("feature"),
				),
				elem.Heading3(
					vecty.Text("Single Language"),
				),
				elem.Paragraph(
					vecty.Text("Both frontend and backend written in Go"),
				),
			),
		),
	)
}

// renderFooter renders the footer section
func (a *Application) renderFooter() vecty.ComponentOrHTML {
	return elem.Footer(
		elem.Paragraph(
			vecty.Text("© 2023 Go Fullstack App with Vecty"),
		),
	)
}
