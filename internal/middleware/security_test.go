package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	tests := []struct {
		name        string
		headerName  string
		headerValue string
	}{
		{
			name:        "X-Content-Type-Options",
			headerName:  "X-Content-Type-Options",
			headerValue: "nosniff",
		},
		{
			name:        "X-Frame-Options",
			headerName:  "X-Frame-Options",
			headerValue: "DENY",
		},
		{
			name:        "X-XSS-Protection",
			headerName:  "X-XSS-Protection",
			headerValue: "1; mode=block",
		},
		{
			name:        "Content-Security-Policy",
			headerName:  "Content-Security-Policy",
			headerValue: "default-src 'self'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			wrapped := SecurityHeaders(handler)

			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			wrapped.ServeHTTP(w, req)

			got := w.Header().Get(tt.headerName)
			if got != tt.headerValue {
				t.Errorf("header %s = %q, want %q", tt.headerName, got, tt.headerValue)
			}
		})
	}
}

func TestSecurityHeadersAllPresent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := SecurityHeaders(handler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	expectedHeaders := map[string]string{
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "DENY",
		"X-XSS-Protection":        "1; mode=block",
		"Content-Security-Policy": "default-src 'self'",
	}

	for name, want := range expectedHeaders {
		got := w.Header().Get(name)
		if got != want {
			t.Errorf("header %s = %q, want %q", name, got, want)
		}
	}
}

func TestSecurityHeadersWithDifferentMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			wrapped := SecurityHeaders(handler)

			req := httptest.NewRequest(method, "/", nil)
			w := httptest.NewRecorder()
			wrapped.ServeHTTP(w, req)

			if w.Header().Get("X-Content-Type-Options") != "nosniff" {
				t.Errorf("%s: security headers not set", method)
			}
		})
	}
}

func TestSecurityHeadersPreservesResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	})

	wrapped := SecurityHeaders(handler)

	req := httptest.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}

	if w.Body.String() != "Created" {
		t.Errorf("body = %q, want %q", w.Body.String(), "Created")
	}

	if w.Header().Get("X-Frame-Options") != "DENY" {
		t.Error("security headers not set")
	}
}
