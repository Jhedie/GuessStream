package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strings"
	"sync"
)

var (
	clients     = make(map[chan string]bool)
	mu          sync.Mutex
	guesses     []string
	secretWord  = "golang"
	randomWords = []string{"golang", "race", "yellow"}
)

func sseHandler(writer http.ResponseWriter, request *http.Request) {
	// Set headers for SSE
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "close")

	writer.Header().Set("Access-Control-Allow-Origin", "*")

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

func broadcast(message string) {
	println("Broadcasting")
	mu.Lock()
	defer mu.Unlock()
	for clientChannel := range clients {
		clientChannel <- message
	}
}

func guessHandler(writer http.ResponseWriter, request *http.Request) {
	// Set headers for SSE
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "close")

	writer.Header().Set("Access-Control-Allow-Origin", "*")

	err := request.ParseForm()
	if err != nil {
		http.Error(writer, "Unable to parse form", http.StatusBadRequest)
		return
	}
	playerGuess := strings.ToLower(request.FormValue("guess"))

	fmt.Println("Player guess:", playerGuess)

	mu.Lock()
	guesses = append(guesses, playerGuess)
	mu.Unlock()

	if playerGuess == secretWord {
		broadcast(fmt.Sprintf("🎉 Correct! The word was '%s'.", secretWord))
		broadcast("Game Over! Reset to play again.")
	} else {
		broadcast(fmt.Sprintf("❌ Guess: '%s'", playerGuess))
	}
	// Optionally, you can send a response back to the client
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Guess received"))
}

func getNewGuessWord() {
	println("new Random word")
	secretWord = randomWords[rand.IntN(len(randomWords))]
	println("Secret word", secretWord)
}

func resetHandler(writer http.ResponseWriter, request *http.Request) {

	// Set headers for SSE
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "close")

	writer.Header().Set("Access-Control-Allow-Origin", "*")

	println("Cleaning up")
	mu.Lock()
	guesses = nil
	getNewGuessWord()
	mu.Unlock()

	broadcast("RESET")
}
func main() {
	http.HandleFunc("/events", sseHandler)
	http.HandleFunc("/guess", guessHandler)
	http.HandleFunc("/reset", resetHandler)

	fmt.Println("Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
