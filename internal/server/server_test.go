package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandleHome(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{"root path", "/", http.StatusOK},
		{"non-root path", "/other", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			handleHome(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("handleHome() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandleAsciiArt(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		formData   url.Values
		wantStatus int
	}{
		{
			name:       "valid request",
			method:     http.MethodPost,
			formData:   url.Values{"input": {"Hello"}, "font": {"standard"}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "default banner",
			method:     http.MethodPost,
			formData:   url.Values{"input": {"Hello"}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "wrong method",
			method:     http.MethodGet,
			formData:   url.Values{},
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "empty text",
			method:     http.MethodPost,
			formData:   url.Values{"input": {""}, "font": {"standard"}},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "text too long",
			method:     http.MethodPost,
			formData:   url.Values{"input": {strings.Repeat("a", 1001)}, "font": {"standard"}},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid banner",
			method:     http.MethodPost,
			formData:   url.Values{"input": {"Hello"}, "font": {"invalid"}},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/ascii-art", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			handleAsciiArt(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("handleAsciiArt() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
