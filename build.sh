#!/bin/bash

# Exit on error
set -e

echo "=== Building Go Fullstack Application ==="

# Create directories if they don't exist
mkdir -p static/css

# Copy the WebAssembly execution environment
echo "Copying wasm_exec.js..."
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" static/

# Build the frontend WebAssembly binary
echo "Building frontend WebAssembly..."
GOOS=js GOARCH=wasm go build -o static/main.wasm ./cmd/frontend

# Build the backend server
echo "Building backend server..."
go build -o server ./cmd/server

echo "=== Build Complete ==="
echo "Run the server with: ./server"
echo "Then visit: http://localhost:8080" 