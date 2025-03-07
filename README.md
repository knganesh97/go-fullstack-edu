# Go Fullstack Web Application

This is a simple todo list application built entirely with Go, using:

- Go's standard library for the backend API
- Vecty for the frontend (compiled to WebAssembly)

## Prerequisites

- Go 1.21 or later
- GOARCH=wasm GOOS=js support

## Project Structure

```
.
├── cmd
│   ├── frontend        # Frontend code (Vecty, compiles to WebAssembly)
│   │   └── components  # Vecty components
│   └── server          # Backend server code
├── pkg
│   └── models          # Shared data models
└── static              # Static assets
    └── css             # CSS styles
```

## Building the Application

You can use the included build script to build the entire application:

```bash
./build.sh
```

Or build manually:

1. First, copy the WebAssembly execution environment to the static directory:

```bash
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" static/
```

2. Build the frontend WebAssembly binary:

```bash
GOOS=js GOARCH=wasm go build -o static/main.wasm ./cmd/frontend
```

3. Build the backend server:

```bash
go build -o server ./cmd/server
```

## Running the Application

1. Start the server:

```bash
./server
```

2. Open your browser and navigate to:

```
http://localhost:8080
```

## Features

- Create, read, update, and delete todos
- Single-page application built with Go and WebAssembly
- RESTful API for todo management
- Responsive design
- Client-side routing with Vecty

## Notes on Vecty

Vecty is a library for building frontend web applications in Go that compiles to WebAssembly. It provides a React-like component model for building user interfaces.

Key benefits:
- Write both frontend and backend in Go
- Type safety across the entire application
- No need for JavaScript or HTML templates
- Efficient DOM updates
- Component-based architecture

## Architecture

This application uses a full Go stack:

1. **Backend**: Go HTTP server providing a RESTful API
2. **Frontend**: Go code compiled to WebAssembly using Vecty
3. **Communication**: HTTP requests between frontend and backend

The application no longer uses server-side HTML templates, as all UI rendering is handled by Vecty components in the browser. 