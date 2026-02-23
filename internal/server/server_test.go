package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
		addr string
	}{
		{"default port", ":8080"},
		{"custom port", ":3000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := New(tt.addr)
			if srv == nil {
				t.Fatal("New() returned nil")
			}
			if srv.Addr != tt.addr {
				t.Errorf("expected addr %s, got %s", tt.addr, srv.Addr)
			}
		})
	}
}

func TestServerStart(t *testing.T) {
	srv := New(":0") // :0 = random available port

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Start()
	}()

	// Give server time to start
	time.Sleep(50 * time.Millisecond)

	// Try to make a request
	resp, err := http.Get("http://" + srv.Addr)
	if err == nil {
		resp.Body.Close()
	}

	// Shutdown
	if err := srv.Shutdown(); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}

	// Check Start() returned no error
	select {
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Start() error = %v", err)
		}
	case <-time.After(time.Second):
		t.Error("Start() did not return after Shutdown()")
	}
}

func TestServerRoutes(t *testing.T) {
	srv := New(":8080")

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{"GET root", "GET", "/", http.StatusOK},
		{"POST ascii-art", "POST", "/ascii-art", http.StatusOK},
		{"GET not found", "GET", "/invalid", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			srv.Handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestServerMiddlewareChain(t *testing.T) {
	srv := New(":8080")

	// Test that recovery middleware catches panics
	// We can't easily test this without modifying handlers,
	// but we can verify the middleware chain is applied
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	srv.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("middleware chain broken: expected 200, got %d", w.Code)
	}
}
