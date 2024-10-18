package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"simple-go-api/internal/cache"
	"simple-go-api/internal/services"
	"strconv"
)

// Handler is a struct that holds a cache for factorial calculations.
type Handler struct {
	cache *cache.InMemCacheMap
}

// NewHandler creates a new Handler with the provided cache.
func NewHandler(cache *cache.InMemCacheMap) *Handler {
	return &Handler{cache: cache}
}

// RouteHandler routes requests to the appropriate handler based on the path.
func RouteHandler(w http.ResponseWriter, r *http.Request, factorialCache *cache.InMemCacheMap) {
	switch {
	case r.URL.Path == "/":
		HomeHandler(w, r)
	case r.URL.Path == "/hello":
		HelloHandler(w, r)
	case r.URL.Path == "/reverse":
		ReverseHandler(w, r)
	case r.URL.Path == "/factorial/" || len(r.URL.Path) > len("/factorial/"):
		FactorialHandler(w, r, factorialCache)
	default:
		respondWithError(w, http.StatusBadRequest, "Invalid Route")
	}
}

// HomeHandler handles GET requests to the home page.
// It provides an overview of the API endpoints and project information.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Create a response with API overview and project information
	response := map[string]interface{}{
		"message": "Welcome to the API! Here are some useful endpoints:",
		"endpoints": []map[string]string{
			{"path": "/hello", "method": "GET", "description": "Returns a greeting message."},
			{"path": "/reverse", "method": "POST", "description": "Reverses the provided text."},
			{"path": "/factorial/{n}", "method": "GET", "description": "Calculates the factorial of n."},
		},
		"project_info": "This project demonstrates building a simple API in Go. " +
			"Key strengths include the use of guard clauses for error handling, organized test cases for maintainability, " +
			"and effective concurrency handling.",
	}

	respondWithJSON(w, http.StatusOK, response)
}

// HelloHandler handles GET requests to the /hello endpoint.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	message := services.GetHelloMessage()
	response := MessageResponse{Message: message}
	respondWithJSON(w, http.StatusOK, response)
}

// ReverseHandler processes POST requests to /reverse,
// reversing the "text" field from the JSON payload and returning it as "reversed".
func ReverseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var req ReverseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Text == "" {
		respondWithError(w, http.StatusBadRequest, "Text field cannot be empty")
		return
	}

	reversedText := services.GetReverseMessage(req.Text)
	response := ReverseResponse{Reversed: reversedText}
	respondWithJSON(w, http.StatusOK, response)
}

// FactorialHandler handles GET requests to calculate the factorial of a given number.
// It extracts the number from the URL path, validates it, and calculates the factorial using the cache.
func FactorialHandler(w http.ResponseWriter, r *http.Request, fcache *cache.InMemCacheMap) {
	fmt.Println("Received request for factorial calculation")

	// Ensure the request method is GET
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Extract the 'n' param from the URL path /factorial/{n}
	path := r.URL.Path
	param := path[len("/factorial/"):]

	// Parse to int and check for non-numeric or negative input
	num, err := strconv.Atoi(param)
	if err != nil || num < 0 || len(path) <= len("/factorial/") {
		respondWithError(w, http.StatusBadRequest, "Input must be a non-negative integer")
		return
	}

	// Calculate Factorials starting with Cache and base case {0:1} via Memoization:
	factorial := services.CalculateFactorial(num, fcache)

	// Respond with the result
	fmt.Printf("Responding with factorial: %d\n", factorial)

	response := MessageResponse{Message: factorial.String()}
	respondWithJSON(w, http.StatusOK, response)
}
