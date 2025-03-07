//go:build js && wasm
// +build js,wasm

package main

import (
	"github.com/hexops/vecty"
	"github.com/knganesh97/go-fullstack-edu/cmd/frontend/components"
)

func main() {
	// Create a new application
	app := &components.Application{}

	// Render the application to the page
	vecty.SetTitle("Go Fullstack App with Vecty")
	vecty.RenderBody(app)
}
