package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecovery(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		wantStatus     int
		wantLogContain string
		wantBody       string
	}{
		{
			name: "panic with string",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic("something went wrong")
			},
			wantStatus:     http.StatusInternalServerError,
			wantLogContain: "panic recovered: something went wrong",
			wantBody:       "Internal Server Error",
		},
		{
			name: "panic with error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic(http.ErrAbortHandler)
			},
			wantStatus:     http.StatusInternalServerError,
			wantLogContain: "panic recovered:",
			wantBody:       "Internal Server Error",
		},
		{
			name: "panic with nil",
			handler: func(w http.ResponseWriter, r *http.Request) {
				var ptr *int
				_ = *ptr
			},
			wantStatus:     http.StatusInternalServerError,
			wantLogContain: "panic recovered:",
			wantBody:       "Internal Server Error",
		},
		{
			name: "no panic",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			},
			wantStatus:     http.StatusOK,
			wantLogContain: "",
			wantBody:       "OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := log.New(&buf, "", 0)

			wrapped := Recovery(logger)(tt.handler)

			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			wrapped.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}

			if tt.wantLogContain != "" {
				logOutput := buf.String()
				if !strings.Contains(logOutput, tt.wantLogContain) {
					t.Errorf("log missing %q, got: %s", tt.wantLogContain, logOutput)
				}
			}

			if !strings.Contains(w.Body.String(), tt.wantBody) {
				t.Errorf("body should contain %q, got: %s", tt.wantBody, w.Body.String())
			}
		})
	}
}

func TestRecoveryPreventsCrash(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("crash!")
	})

	wrapped := Recovery(logger)(handler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// This should not crash the test
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestRecoveryWithMultipleRequests(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	callCount := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 2 {
			panic("second request panics")
		}
		w.WriteHeader(http.StatusOK)
	})

	wrapped := Recovery(logger)(handler)

	// First request succeeds
	req1 := httptest.NewRequest("GET", "/", nil)
	w1 := httptest.NewRecorder()
	wrapped.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Errorf("first request: expected 200, got %d", w1.Code)
	}

	// Second request panics but is recovered
	req2 := httptest.NewRequest("GET", "/", nil)
	w2 := httptest.NewRecorder()
	wrapped.ServeHTTP(w2, req2)
	if w2.Code != http.StatusInternalServerError {
		t.Errorf("second request: expected 500, got %d", w2.Code)
	}

	// Third request succeeds (server still running)
	req3 := httptest.NewRequest("GET", "/", nil)
	w3 := httptest.NewRecorder()
	wrapped.ServeHTTP(w3, req3)
	if w3.Code != http.StatusOK {
		t.Errorf("third request: expected 200, got %d", w3.Code)
	}
}
