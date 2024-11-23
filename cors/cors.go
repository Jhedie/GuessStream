package main

import (
	"fmt"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")                   // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Allowed methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Allowed headers
		w.Header().Set("Access-Control-Allow-Credentials", "true")           // Allow credentials

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
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Add your handlers
	mux.HandleFunc("/events", sseHandler)
	mux.HandleFunc("/guess", guessHandler)

	// Wrap the ServeMux with the CORS middleware
	handlerWithCORS := corsMiddleware(mux)

	// Start the server
	server := &http.Server{
		Addr:    ":8080",
		Handler: handlerWithCORS,
	}

	fmt.Println("Server running at http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
