package web

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHome_TableDriven tests the Home handler with multiple scenarios.
func TestHome_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		method         string
		template       *template.Template
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success execution",
			path:           "/",
			method:         http.MethodGet,
			template:       template.Must(template.New("index.html").Parse("<h1>HOME</h1>")),
			expectedStatus: http.StatusOK,
			expectedBody:   "<h1>HOME</h1>",
		},
		{
			name:           "execution failure",
			path:           "/",
			method:         http.MethodGet,
			template:       template.Must(template.New("index.html").Parse(`{{call .}}`)),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to connect to the internal service",
		},
		{
			name:           "template missing",
			path:           "/",
			method:         http.MethodGet,
			template:       nil,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "The template does not exist",
		},
		{
			name:           "unknown path returns 404",
			path:           "/random-page",
			method:         http.MethodGet,
			template:       template.Must(template.New("index.html").Parse("<h1>HOME</h1>")),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found",
		},
		{
			name:           "wrong method returns 405",
			path:           "/",
			method:         http.MethodPost,
			template:       template.Must(template.New("index.html").Parse("<h1>HOME</h1>")),
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := make(map[string]*template.Template)

			if tt.template != nil {
				cache["index.html"] = tt.template
			}

			app := &Application{
				TemplateCache: cache,
			}

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			app.Home(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			body := rr.Body.String()
			if !strings.Contains(body, tt.expectedBody) {
				t.Fatalf("expected body to contain %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

// TestNewTemplateCache tests the template cache initialization.
func TestNewTemplateCache(t *testing.T) {
	t.Run("returns error on missing files", func(t *testing.T) {
		_, err := NewTemplateCache()
		if err == nil {
			t.Fatal("expected error when template files are missing, got nil")
		}
	})
}
