package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerOpenAPI(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rr := httptest.NewRecorder()
	Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "yaml") {
		t.Fatalf("Content-Type = %q, want yaml", ct)
	}
	body := rr.Body.String()
	for _, want := range []string{"openapi:", "/api/add", "/api/divide", "/api/history"} {
		if !strings.Contains(body, want) {
			t.Fatalf("spec missing %q", want)
		}
	}
}

func TestHandlerSwaggerUI(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rr := httptest.NewRecorder()
	Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if !strings.Contains(rr.Body.String(), "SwaggerUIBundle") {
		t.Fatalf("swagger html missing SwaggerUIBundle")
	}
}

func TestHandlerSwaggerRedirect(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/swagger", nil)
	rr := httptest.NewRecorder()
	Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusMovedPermanently {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusMovedPermanently)
	}
}
