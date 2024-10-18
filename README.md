## Simple Toy API
### Description
Simple toy api that implements 3 endpoints, a simple in-memory cache, is thread safe, and handles multiple requests concurrently.

### Endpoints
1. GET /
    - home route that describes/showcases the project's endpoints
2. GET /hello
    - Returns a JSON response with a message: `{"message": "Hello, World!"}`.
2. POST /reverse:
    - Accepts a JSON payload with a single field text (string).
    - Returns a JSON response with the reversed string: `{"reversed": "<reversed_text>"}`.
3. GET /factorial/{n}:
    - Returns the factorial of a given integer n passed as a URL parameter.
    - If n is negative, return an error message: `{"error": "Input must be a non-negative integer"}`.
