package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

var (
	clients = make(map[chan string]bool)
	mu      sync.Mutex
)

func sseHandler(writer http.ResponseWriter, request *http.Request) {
	// Set headers for SSE

	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "close")

	// Ensure the connection supports flushing
	flusher, ok := writer.(http.Flusher)
	if !ok {
		http.Error(writer, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Create a channel to send messages to the client
	clientChannel := make(chan string)

	// Lock the map to safely add this client
	mu.Lock()
	println("Client joined")
	clients[clientChannel] = true
	mu.Unlock()
	printClients()

	// To create a cancellable context that does not automatically timeout.
	ctx, cancel := context.WithCancel(request.Context())
	defer cancel()

	// Listen for context cancellation
	go func() {
		<-ctx.Done() // Wait for cancelation
		mu.Lock()
		delete(clients, clientChannel)
		println("Client left due to cancelation")
		mu.Unlock()
		printClients()
		close(clientChannel)
	}()

	for message := range clientChannel {
		fmt.Fprintf(writer, "data: %s\n\n", message)
		flusher.Flush()
	}
}

func printClients() {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("Current clients:")
	for client := range clients {
		fmt.Println(client)
	}
}

func main() {
	http.HandleFunc("/events", sseHandler)

	fmt.Println("Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
