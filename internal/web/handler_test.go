package web

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHome_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		template       *template.Template
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success execution",
			template:       template.Must(template.New("index.html").Parse("<h1>HOME</h1>")),
			expectedStatus: http.StatusOK,
			expectedBody:   "<h1>HOME</h1>",
		},
		{
			name:           "execution failure",
			template:       template.Must(template.New("index.html").Parse(`{{call .}}`)),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to connect to the internal service",
		},
		{
			name:           "template missing",
			template:       nil,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "The template does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cache := make(map[string]*template.Template)

			// Only add template if it exists
			if tt.template != nil {
				cache["index.html"] = tt.template
			}

			app := &Application{
				TemplateCache: cache,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
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
