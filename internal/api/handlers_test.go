package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"simple-go-api/internal/cache"
	"sync"
	"testing"
)

// testCase is a common struct for defining test cases.
// testCase is a common struct for defining test cases.
type testCase struct {
	name             string
	path             string
	expectedStatus   int
	expectedResponse string
	expectedMethod   string
	expectedBody     string // This field can be empty for cases that don't use it
}

// Factorial test cases
var factorialTestCases = []testCase{
	{
		name:             "Factorial route with valid input returns OK",
		path:             "/factorial/5",
		expectedStatus:   http.StatusOK,
		expectedResponse: `{"message":"120"}`,
		expectedMethod:   http.MethodGet,
	},
	{
		name:             "Factorial route with negative input returns BadRequest",
		path:             "/factorial/-5",
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: `{"error":"Input must be a non-negative integer"}`,
		expectedMethod:   http.MethodGet,
	},
	{
		name:             "Factorial route with non-numeric input returns BadRequest",
		path:             "/factorial/abc",
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: `{"error":"Input must be a non-negative integer"}`,
		expectedMethod:   http.MethodGet,
	},
	{
		name:             "Factorial route with POST method returns MethodNotAllowed",
		path:             "/factorial/5",
		expectedStatus:   http.StatusMethodNotAllowed,
		expectedResponse: `{"error":"Method Not Allowed"}`,
		expectedMethod:   http.MethodPost,
	},
}

// Hello test cases
var helloTestCases = []testCase{
	{
		name:             "Hello route with GET returns OK",
		path:             "/hello",
		expectedStatus:   http.StatusOK,
		expectedResponse: `{"message":"Hello, World!"}`,
		expectedMethod:   http.MethodGet,
	},
	{
		name:             "Hello route with POST returns MethodNotAllowed",
		path:             "/hello",
		expectedStatus:   http.StatusMethodNotAllowed,
		expectedResponse: `{"error":"Method Not Allowed"}`,
		expectedMethod:   http.MethodPost,
	},
}

// Reverse test cases
var reverseTestCases = []testCase{
	{
		name:             "Reverse route with GET returns MethodNotAllowed",
		path:             "/reverse",
		expectedStatus:   http.StatusMethodNotAllowed,
		expectedResponse: `{"error":"Method Not Allowed"}`,
		expectedMethod:   http.MethodGet,
	},
	{
		name:             "Reverse route with POST empty body returns BadRequest",
		path:             "/reverse",
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: `{"error":"Text field cannot be empty"}`,
		expectedMethod:   http.MethodPost,
		expectedBody:     `{"text":""}`,
	},
	{
		name:             "Reverse route with POST invalid payload returns BadRequest",
		path:             "/reverse",
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: `{"error":"Invalid request payload"}`,
		expectedMethod:   http.MethodPost,
		expectedBody:     `{"text": 12345}`,
	},
	{
		name:             "Reverse single character with POST returns OK",
		path:             "/reverse",
		expectedStatus:   http.StatusOK,
		expectedResponse: `{"reversed":"a"}`,
		expectedMethod:   http.MethodPost,
		expectedBody:     `{"text":"a"}`,
	},
	{
		name:             "Reverse long sentence with POST returns OK",
		path:             "/reverse",
		expectedStatus:   http.StatusOK,
		expectedResponse: `{"reversed":"stleb elprup ydwor esoht lla ta kool ereht yeh"}`,
		expectedMethod:   http.MethodPost,
		expectedBody:     `{"text":"hey there look at all those rowdy purple belts"}`,
	},
	{
		name:             "Reverse short word with POST returns OK",
		path:             "/reverse",
		expectedStatus:   http.StatusOK,
		expectedResponse: `{"reversed":"sso"}`,
		expectedMethod:   http.MethodPost,
		expectedBody:     `{"text":"oss"}`,
	},
}

// Home test cases
var homeTestCases = []testCase{
	{
		name:           "Home route with GET returns OK",
		path:           "/",
		expectedStatus: http.StatusOK,
		expectedResponse: `{
			"message": "Welcome to the API! Here are some useful endpoints:",
			"endpoints": [
				{"path": "/hello", "method": "GET", "description": "Returns a greeting message."},
				{"path": "/reverse", "method": "POST", "description": "Reverses the provided text."},
				{"path": "/factorial/{n}", "method": "GET", "description": "Calculates the factorial of n."}
			],
			"project_info": "This project demonstrates building a simple API in Go. Key strengths include the use of guard clauses for error handling, organized test cases for maintainability, and effective concurrency handling."
		}`,
		expectedMethod: http.MethodGet,
	},
	{
		name:             "Home route with POST returns MethodNotAllowed",
		path:             "/",
		expectedStatus:   http.StatusMethodNotAllowed,
		expectedResponse: `{"error":"Method Not Allowed"}`,
		expectedMethod:   http.MethodPost,
	},
}

// Invalid route test case
var invalidRouteTestCase = testCase{
	name:             "Invalid route returns BadRequest",
	path:             "/invalid",
	expectedStatus:   http.StatusBadRequest,
	expectedResponse: `{"error":"Invalid Route"}`,
	expectedMethod:   http.MethodGet,
}

// TestNewHandler tests the creation of a new Handler.
func TestNewHandler(t *testing.T) {
	cache := cache.NewInMemCacheMap()
	handler := NewHandler(cache)
	if handler.cache != cache {
		t.Errorf("NewHandler() = %v, want %v", handler.cache, cache)
	}
}

// TestRouteHandlers tests the API route handlers for various endpoints.
func TestRouteHandlers(t *testing.T) {
	// Init cache for passing to RouteHandler test cases
	testCache := cache.NewInMemCacheMap()

	// Combine all test cases
	allTestCases := append(homeTestCases, factorialTestCases...)
	allTestCases = append(allTestCases, helloTestCases...)
	allTestCases = append(allTestCases, reverseTestCases...)
	allTestCases = append(allTestCases, invalidRouteTestCase)

	for _, tc := range allTestCases {
		t.Run(tc.name, func(t *testing.T) {
			var reqBody io.Reader
			if tc.expectedBody != "" {
				reqBody = bytes.NewBufferString(tc.expectedBody)
			}

			req := httptest.NewRequest(tc.expectedMethod, tc.path, reqBody)
			rec := httptest.NewRecorder()

			RouteHandler(rec, req, testCache)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.expectedStatus {
				t.Errorf("expected status %v, got %v", tc.expectedStatus, res.StatusCode)
			}

			bodyBytes, _ := io.ReadAll(res.Body)
			bodyString := string(bodyBytes)

			var expected, actual map[string]interface{}
			if err := json.Unmarshal([]byte(tc.expectedResponse), &expected); err != nil {
				t.Fatalf("invalid expected response JSON: %v", err)
			}
			if err := json.Unmarshal(bodyBytes, &actual); err != nil {
				t.Fatalf("invalid actual response JSON: %v", err)
			}

			if !equalJSON(expected, actual) {
				t.Errorf("expected response %v, got %v", tc.expectedResponse, bodyString)
			}
		})
	}
}

func TestFactorialRequestsConcurrency(t *testing.T) {
	// Initialize the cache
	factorialCache := cache.NewInMemCacheMap()

	// Create a new request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		RouteHandler(w, r, factorialCache)
	})

	// Start a test server
	server := httptest.NewServer(mux)
	defer server.Close()

	// Define inputs for factorial calculation
	inputs := []int{5, 6, 7, 8, 9, 10}
	var wg sync.WaitGroup

	// Function to call the factorial endpoint
	callFactorial := func(n int) {
		defer wg.Done()
		resp, err := http.Get(fmt.Sprintf("%s/factorial/%d", server.URL, n))
		if err != nil {
			t.Errorf("Error requesting factorial of %d: %v", n, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK for %d, got %v", n, resp.StatusCode)
		}
	}

	// Launch concurrent requests
	for _, n := range inputs {
		wg.Add(1)
		go callFactorial(n)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func equalJSON(a, b map[string]interface{}) bool {
	aBytes, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bBytes, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return string(aBytes) == string(bBytes)
}
