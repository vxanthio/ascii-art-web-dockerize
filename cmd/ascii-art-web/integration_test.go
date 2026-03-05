package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"ascii-art-web/internal/handlers"
)

// newTestServer builds a real server using the actual template cache and
// returns an httptest.Server. Tests must run from the project root so that
// NewTemplateCache can locate templates/ and static/ directories.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	cache, err := handlers.NewTemplateCache()
	if err != nil {
		t.Fatalf("NewTemplateCache() failed: %v\n(integration tests must run from project root)", err)
	}

	app := &handlers.Application{TemplateCache: cache}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/ascii-art", app.HandleASCIIArt)

	return httptest.NewServer(mux)
}

func TestIntegration_WebServer(t *testing.T) {
	if err := os.Chdir("../.."); err != nil {
		t.Fatalf("could not change to project root: %v", err)
	}

	srv := newTestServer(t)
	defer srv.Close()

	tests := []struct {
		name          string
		method        string
		path          string
		formData      url.Values
		wantStatus    int
		wantInBody    string
		wantNotInBody string
	}{
		{
			name:       "GET / returns 200 with HTML page",
			method:     http.MethodGet,
			path:       "/",
			wantStatus: http.StatusOK,
			wantInBody: "<form",
		},
		{
			name:       "GET /unknown returns 404",
			method:     http.MethodGet,
			path:       "/unknown",
			wantStatus: http.StatusNotFound,
			wantInBody: "404",
		},
		{
			name:       "POST / returns 405",
			method:     http.MethodPost,
			path:       "/",
			wantStatus: http.StatusMethodNotAllowed,
			wantInBody: "Method not allowed",
		},
		{
			name:       "GET /ascii-art returns 405",
			method:     http.MethodGet,
			path:       "/ascii-art",
			wantStatus: http.StatusMethodNotAllowed,
			wantInBody: "Method not allowed",
		},
		{
			name:       "POST /ascii-art empty text returns 400 with error in page",
			method:     http.MethodPost,
			path:       "/ascii-art",
			formData:   url.Values{"text": {""}, "banner": {"standard"}},
			wantStatus: http.StatusBadRequest,
			wantInBody: "text cannot be empty",
		},
		{
			name:       "POST /ascii-art invalid banner returns 404 with error in page",
			method:     http.MethodPost,
			path:       "/ascii-art",
			formData:   url.Values{"text": {"Hello"}, "banner": {"invalid"}},
			wantStatus: http.StatusNotFound,
			wantInBody: "invalid banner name",
		},
		{
			name:       "POST /ascii-art valid standard banner returns 200 with ASCII art",
			method:     http.MethodPost,
			path:       "/ascii-art",
			formData:   url.Values{"text": {"{123}"}, "banner": {"standard"}},
			wantStatus: http.StatusOK,
			wantInBody: "<pre>",
		},
		{
			name:       "POST /ascii-art valid shadow banner returns 200",
			method:     http.MethodPost,
			path:       "/ascii-art",
			formData:   url.Values{"text": {"Hi"}, "banner": {"shadow"}},
			wantStatus: http.StatusOK,
			wantInBody: "<pre>",
		},
		{
			name:       "POST /ascii-art valid thinkertoy banner returns 200",
			method:     http.MethodPost,
			path:       "/ascii-art",
			formData:   url.Values{"text": {"Go"}, "banner": {"thinkertoy"}},
			wantStatus: http.StatusOK,
			wantInBody: "<pre>",
		},
		{
			name:       "GET /static/style.css returns 200",
			method:     http.MethodGet,
			path:       "/static/style.css",
			wantStatus: http.StatusOK,
			wantInBody: "body",
		},
		{
			name:          "POST /ascii-art result page has no error message on success",
			method:        http.MethodPost,
			path:          "/ascii-art",
			formData:      url.Values{"text": {"Hi"}, "banner": {"standard"}},
			wantStatus:    http.StatusOK,
			wantNotInBody: "error-message",
		},
		{
			name:          "POST /ascii-art error page has no pre block",
			method:        http.MethodPost,
			path:          "/ascii-art",
			formData:      url.Values{"text": {""}, "banner": {"standard"}},
			wantStatus:    http.StatusBadRequest,
			wantNotInBody: "<pre>",
		},
	}

	client := srv.Client()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *http.Response
			var err error

			if tt.method == http.MethodPost {
				resp, err = client.PostForm(srv.URL+tt.path, tt.formData)
			} else {
				req, reqErr := http.NewRequest(tt.method, srv.URL+tt.path, nil)
				if reqErr != nil {
					t.Fatalf("failed to create request: %v", reqErr)
				}
				resp, err = client.Do(req)
			}

			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			bodyStr := string(body)

			if tt.wantInBody != "" && !strings.Contains(bodyStr, tt.wantInBody) {
				t.Errorf("body does not contain %q\nbody: %s", tt.wantInBody, bodyStr)
			}

			if tt.wantNotInBody != "" && strings.Contains(bodyStr, tt.wantNotInBody) {
				t.Errorf("body should not contain %q\nbody: %s", tt.wantNotInBody, bodyStr)
			}
		})
	}
}
