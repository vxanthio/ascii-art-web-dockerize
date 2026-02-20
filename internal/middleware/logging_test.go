package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogging(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		handlerStatus  int
		wantLogContain []string
	}{
		{
			name:          "GET request",
			method:        "GET",
			path:          "/",
			handlerStatus: http.StatusOK,
			wantLogContain: []string{"GET", "/", "200"},
		},
		{
			name:          "POST request",
			method:        "POST",
			path:          "/ascii-art",
			handlerStatus: http.StatusOK,
			wantLogContain: []string{"POST", "/ascii-art", "200"},
		},
		{
			name:          "404 error",
			method:        "GET",
			path:          "/notfound",
			handlerStatus: http.StatusNotFound,
			wantLogContain: []string{"GET", "/notfound", "404"},
		},
		{
			name:          "500 error",
			method:        "POST",
			path:          "/error",
			handlerStatus: http.StatusInternalServerError,
			wantLogContain: []string{"POST", "/error", "500"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			logger := log.New(&buf, "", 0)

			// Mock handler that returns the test status
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.handlerStatus)
			})

			// Wrap with logging middleware
			wrapped := Logging(logger)(handler)

			// Make request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			wrapped.ServeHTTP(w, req)

			// Check log output
			logOutput := buf.String()
			for _, want := range tt.wantLogContain {
				if !strings.Contains(logOutput, want) {
					t.Errorf("log missing %q, got: %s", want, logOutput)
				}
			}
		})
	}
}

func TestLoggingDuration(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := Logging(logger)(handler)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	logOutput := buf.String()
	if !strings.Contains(logOutput, "ms") && !strings.Contains(logOutput, "µs") && !strings.Contains(logOutput, "ns") {
		t.Errorf("log should contain duration, got: %s", logOutput)
	}
}
