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
	players     = make(map[string]bool) // Track players by their identifiers (e.g., IPs)
	guesses     []string
	secretWord  = "golang"
	randomWords = []string{"golang", "race", "yellow"}
)

func sseHandler(writer http.ResponseWriter, request *http.Request) {

	// Ensure the connection supports flushing
	flusher, ok := writer.(http.Flusher)
	if !ok {
		http.Error(writer, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Create a channel to send messages to the client
	clientChannel := make(chan string, 10) // Buffered channel with capacity of 10
	clientID := request.RemoteAddr         // Use IP address as an identifier

	// Lock the map to safely add this client
	mu.Lock()
	println("Client joined")
	clients[clientChannel] = true
	players[clientID] = true
	mu.Unlock()
	printClients()

	broadcastPlayerList() // Broadcast updated player list to all clients

	// To create a cancellable context that does not automatically timeout.
	ctx, cancel := context.WithCancel(request.Context())
	defer cancel()

	// Listen for context cancellation
	go func() {
		<-ctx.Done() // Wait for cancelation

		mu.Lock()
		delete(clients, clientChannel)
		delete(players, clientID)
		println("Client left due to cancellation")
		mu.Unlock()

		printClients()

		close(clientChannel)
		broadcastPlayerList() // Update player list on leave

	}()

	for message := range clientChannel {
		fmt.Fprintf(writer, "data: %s\n\n", message)
		flusher.Flush()
		println("Flushed message to client:", message)
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
	println("Broadcasting message:", message)
	mu.Lock()
	defer mu.Unlock()
	for clientChannel := range clients {
		select {
		case clientChannel <- message:
		default:
			println("Client channel is unresponsive, skipping")
		}
	}
}

func broadcastPlayerList() {
	playerList := make([]string, 0, len(players))
	mu.Lock()
	for player := range players {
		playerList = append(playerList, player)
	}
	mu.Unlock()

	message := fmt.Sprintf("PLAYERS: %s", strings.Join(playerList, ", "))
	broadcast(message)
}

func guessHandler(writer http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		http.Error(writer, "Unable to parse form", http.StatusBadRequest)
		return
	}
	playerGuess := strings.ToLower(request.FormValue("guess"))

	fmt.Println("Player guess:", playerGuess)

	mu.Lock()
	fmt.Println("appending guesses message")
	guesses = append(guesses, playerGuess)
	mu.Unlock()

	if playerGuess == secretWord {
		broadcast(fmt.Sprintf("🎉 Correct! The word was '%s'.", secretWord))
		broadcast("Game Over! Reset to play again.")
	} else {
		broadcast(fmt.Sprintf("❌ Guess: '%s'", playerGuess))
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Guess received"))
}

func getNewGuessWord() {
	println("new Random word")
	secretWord = randomWords[rand.IntN(len(randomWords))]
	println("Secret word", secretWord)
}

func resetHandler(writer http.ResponseWriter, request *http.Request) {

	println("Cleaning up")
	mu.Lock()
	guesses = nil
	getNewGuessWord()
	mu.Unlock()

	broadcast("RESET")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")                   // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Allowed methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Allowed headers
		w.Header().Set("Access-Control-Allow-Credentials", "true")           // Allow credentials

		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "close")
		// Handle preflight requests (OPTIONS method)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	//Create a new ServeMux
	mux := http.NewServeMux()

	// Wrap the ServeMux with the CORS middleware
	handlerWithCORS := corsMiddleware(mux)

	mux.HandleFunc("/events", sseHandler)
	mux.HandleFunc("/guess", guessHandler)
	mux.HandleFunc("/reset", resetHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handlerWithCORS,
	}

	fmt.Println("Server running at http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
