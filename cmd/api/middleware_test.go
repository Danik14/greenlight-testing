package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"greenlight.bcc/internal/assert"
)

func TestRecoverPanicMiddleware(t *testing.T) {
	// Create a new application instance
	app := newTestApplication(t)

	// Create a new request that will trigger a panic
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	recorder := httptest.NewRecorder()

	// Create a test handler that will panic
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})

	// Wrap the handler with the middleware
	middleware := app.recoverPanic(handler)

	// Call the middleware
	middleware.ServeHTTP(recorder, req)

	// Check that the response status code is 500
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// // Check that the response body contains the error message
	// expectedBody := `{"error":"the server encountered a problem and could not process your request"}`
	// assert.Equal(t, expectedBody, recorder.Body.String())

	expectedJSON := `{"error":"the server encountered a problem and could not process your request"}`
	actualJSON := strings.TrimSpace(recorder.Body.String())

	if !json.Valid([]byte(actualJSON)) {
		t.Fatalf("invalid JSON response: %s", actualJSON)
	}

	var expected interface{}
	var actual interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(actualJSON), &actual); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected response body:\nexpected: %v\nactual: %v", expected, actual)
	}
}

// func TestRateLimitMiddleware(t *testing.T) {
// 	app := newTestApplication(t)
// 	rr := httptest.NewRecorder()
// 	req, err := http.NewRequest("GET", "/v1/healthcheck", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Make several requests within the rate limit
// 	for i := 0; i < 2; i++ {
// 		app.rateLimit(http.HandlerFunc(app.healthcheckHandler)).ServeHTTP(rr, req)
// 		assert.Equal(t, rr.Code, http.StatusOK)
// 		rr.Body.Reset()
// 	}

// 	// Make a request that exceeds the rate limit
// 	app.rateLimit(http.HandlerFunc(app.healthcheckHandler)).ServeHTTP(rr, req)
// 	assert.Equal(t, rr.Code, http.StatusTooManyRequests)
// 	assert.Equal(t, rr.Body.String(), `{"error": "rate limit exceeded"}`)
// }

// func TestAuthenticateMiddleware(t *testing.T) {
// 	// Create a new application instance and set up the necessary dependencies
// 	app := newTestApplication(t)

// 	// Create a new user and generate a bearer token for them

// 	users := []*data.User{
// 		{
// 			Email: "test@test.com",
// 		},
// 		{
// 			Email: "test2@test2.com",
// 		},
// 	}
// 	token1, err := app.models.Tokens.New(users[0])
// 	// token, err := user.GenerateToken(data.ScopeAuthentication, time.Hour)

// 	// Create a mock request with the bearer token in the Authorization header
// 	req, err := http.NewRequest("GET", "/test", nil)
// 	if err != nil {
// 	}
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

// 	// Create a mock response recorder and a handler that just writes the user ID to the response
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		user := app.contextGetUser(r)
// 		fmt.Fprintf(w, "User ID: %d", user.ID)
// 	})

// 	// Create a new request context with the mock request and a fresh database connection
// 	ctx := context.WithValue(context.Background(), key("db"), app.db)

// 	// Call the Authenticate middleware with the test handler
// 	middleware := app.authenticate(handler)
// 	middleware.ServeHTTP(rr, req.WithContext(ctx))

// 	// Verify that the response code is 200 OK and that the response body contains the expected user ID
// 	assert.Equal(t, http.StatusOK, rr.Code)
// 	assert.Equal(t, "User ID: 1", rr.Body.String())
// }
