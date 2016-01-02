package main

import (
	_ "github.com/orchestrate-io/dvr"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// setUp() equivalent code would go here

	// this call actually run the tests
	r := m.Run()

	// tearDown() equivalent code here
	os.Exit(r)
}

// unconfigured ApiRouter
var ar = ApiRouter{}

// a list of handlers and their respective expected http status code
var handlerTests = []struct {
	method          string
	path            string
	handler         http.Handler
	expected_status int
}{
	{
		method:          "GET",
		path:            "/health",
		handler:         http.HandlerFunc(ar.HealthHandler),
		expected_status: http.StatusOK,
	},
	{
		method:          "GET",
		path:            "/httpclient",
		handler:         http.HandlerFunc(ar.HttpClientHandler),
		expected_status: http.StatusOK,
	},
}

// one test for all handlers
func TestHandlers(t *testing.T) {
	for _, tt := range handlerTests {
		req, _ := http.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()

		tt.handler.ServeHTTP(w, req)
		if w.Code != tt.expected_status {
			t.Errorf("Request to path %s, answered by handler %s didn't return the expected %d status code.",
				tt.path, tt.handler, tt.expected_status)
		}
	}
}

// Benchmark a handler -- the way this is written is not particularly useful,
// as here we create a client and a server instance for each iteration of the loop.
// at the same time it seems this can be considered the worst case for this handler.
func BenchmarkHelloHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/hello", nil)
		http.HandlerFunc(ar.HelloHandler).ServeHTTP(httptest.NewRecorder(), req)
	}
}
