package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
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
			name:           "GET request",
			method:         "GET",
			path:           "/",
			handlerStatus:  http.StatusOK,
			wantLogContain: []string{"GET", "/", "200"},
		},
		{
			name:           "POST request",
			method:         "POST",
			path:           "/ascii-art",
			handlerStatus:  http.StatusOK,
			wantLogContain: []string{"POST", "/ascii-art", "200"},
		},
		{
			name:           "404 error",
			method:         "GET",
			path:           "/notfound",
			handlerStatus:  http.StatusNotFound,
			wantLogContain: []string{"GET", "/notfound", "404"},
		},
		{
			name:           "500 error",
			method:         "POST",
			path:           "/error",
			handlerStatus:  http.StatusInternalServerError,
			wantLogContain: []string{"POST", "/error", "500"},
		},
		{
			name:           "PUT request",
			method:         "PUT",
			path:           "/update",
			handlerStatus:  http.StatusOK,
			wantLogContain: []string{"PUT", "/update", "200"},
		},
		{
			name:           "DELETE request",
			method:         "DELETE",
			path:           "/delete",
			handlerStatus:  http.StatusNoContent,
			wantLogContain: []string{"DELETE", "/delete", "204"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := log.New(&buf, "", 0)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.handlerStatus)
			})

			wrapped := Logging(logger)(handler)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			wrapped.ServeHTTP(w, req)

			logOutput := buf.String()
			for _, want := range tt.wantLogContain {
				if !strings.Contains(logOutput, want) {
					t.Errorf("log missing %q, got: %s", want, logOutput)
				}
			}
		})
	}
}

func TestLoggingFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := Logging(logger)(handler)
	req := httptest.NewRequest("POST", "/ascii-art", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	logOutput := buf.String()
	// Format: [2024-02-20 19:54:20] POST /ascii-art 200 15ms
	pattern := `^\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\] POST /ascii-art 200 \d+ms\n$`
	matched, err := regexp.MatchString(pattern, logOutput)
	if err != nil {
		t.Fatalf("regex error: %v", err)
	}
	if !matched {
		t.Errorf("log format incorrect, got: %s", logOutput)
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
	if !strings.Contains(logOutput, "ms") {
		t.Errorf("log should contain duration in ms, got: %s", logOutput)
	}
}

func TestLoggingWrapsAllHandlers(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		method  string
		path    string
		status  int
	}{
		{
			name: "home handler",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Home"))
			},
			method: "GET",
			path:   "/",
			status: http.StatusOK,
		},
		{
			name: "ascii-art handler",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ASCII Art"))
			},
			method: "POST",
			path:   "/ascii-art",
			status: http.StatusOK,
		},
		{
			name: "error handler",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			},
			method: "GET",
			path:   "/error",
			status: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := log.New(&buf, "", 0)

			wrapped := Logging(logger)(tt.handler)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			wrapped.ServeHTTP(w, req)

			if w.Code != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, w.Code)
			}

			logOutput := buf.String()
			if logOutput == "" {
				t.Error("expected log output, got empty string")
			}
		})
	}
}

func TestResponseWriterPreservesStatus(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	})

	wrapped := Logging(logger)(handler)
	req := httptest.NewRequest("POST", "/create", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "201") {
		t.Errorf("log should contain status 201, got: %s", logOutput)
	}
}
