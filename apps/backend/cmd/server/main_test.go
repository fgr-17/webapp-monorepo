package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	newHandler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status body = %q, want ok", body["status"])
	}
}

func TestCalculatorOps(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		a, b       float64
		want       float64
		wantStatus int
		wantErr    string
	}{
		{name: "add", path: "/api/add", a: 10, b: 2, want: 12, wantStatus: http.StatusOK},
		{name: "subtract", path: "/api/subtract", a: 10, b: 2, want: 8, wantStatus: http.StatusOK},
		{name: "multiply", path: "/api/multiply", a: 10, b: 2, want: 20, wantStatus: http.StatusOK},
		{name: "divide", path: "/api/divide", a: 10, b: 2, want: 5, wantStatus: http.StatusOK},
		{name: "divide_by_zero", path: "/api/divide", a: 10, b: 0, wantStatus: http.StatusBadRequest, wantErr: "division by zero"},
	}

	handler := newHandler()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, _ := json.Marshal(operands{A: tt.a, B: tt.b})
			req := httptest.NewRequest(http.MethodPost, tt.path, bytes.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rr.Code, tt.wantStatus, rr.Body.String())
			}
			if tt.wantErr != "" {
				var errBody errorResponse
				if err := json.NewDecoder(rr.Body).Decode(&errBody); err != nil {
					t.Fatalf("decode error: %v", err)
				}
				if errBody.Error != tt.wantErr {
					t.Fatalf("error = %q, want %q", errBody.Error, tt.wantErr)
				}
				return
			}
			var body resultResponse
			if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if body.Result != tt.want {
				t.Fatalf("result = %v, want %v", body.Result, tt.want)
			}
		})
	}
}

func TestInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/add", strings.NewReader("{not-json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	newHandler().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestCORSPreflight(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/api/add", nil)
	rr := httptest.NewRecorder()
	newHandler().ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("Allow-Origin = %q, want *", got)
	}
}

func TestSwaggerAndOpenAPI(t *testing.T) {
	handler := newHandler()

	t.Run("openapi", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
		body := rr.Body.String()
		if !strings.Contains(body, "openapi:") {
			t.Fatalf("body missing openapi header")
		}
		if !strings.Contains(body, "/api/history") {
			t.Fatalf("body missing /api/history")
		}
	})

	t.Run("swagger_ui", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
		if !strings.Contains(rr.Body.String(), "swagger-ui") {
			t.Fatalf("body missing swagger-ui")
		}
	})

	t.Run("swagger_redirect", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusMovedPermanently {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusMovedPermanently)
		}
		if loc := rr.Header().Get("Location"); loc != "/swagger/" {
			t.Fatalf("Location = %q, want /swagger/", loc)
		}
	})
}

func TestHistoryEmpty(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/history", nil)
	rr := httptest.NewRecorder()
	newHandler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	var entries []map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("len = %d, want 0", len(entries))
	}
}

func TestHistoryRecordsSuccessfulOps(t *testing.T) {
	handler := newHandler()

	post := func(path string, a, b float64) *httptest.ResponseRecorder {
		payload, _ := json.Marshal(operands{A: a, B: b})
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		return rr
	}

	if rr := post("/api/add", 10, 2); rr.Code != http.StatusOK {
		t.Fatalf("add status = %d", rr.Code)
	}
	if rr := post("/api/multiply", 3, 4); rr.Code != http.StatusOK {
		t.Fatalf("multiply status = %d", rr.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/history", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("history status = %d", rr.Code)
	}

	var entries []struct {
		Op     string  `json:"op"`
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(entries) < 2 {
		t.Fatalf("len = %d, want at least 2", len(entries))
	}
	ops := map[string]bool{}
	for _, e := range entries {
		ops[e.Op] = true
	}
	if !ops["add"] || !ops["multiply"] {
		t.Fatalf("history missing expected ops: %+v", entries)
	}
}

func TestHistoryHasFiniteCap(t *testing.T) {
	handler := newHandler()
	for i := 0; i < 20; i++ {
		payload, _ := json.Marshal(operands{A: float64(i), B: 1})
		req := httptest.NewRequest(http.MethodPost, "/api/add", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("add #%d status = %d", i, rr.Code)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/history", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var entries []struct {
		A float64 `json:"a"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(entries) == 0 || len(entries) >= 20 {
		t.Fatalf("len = %d, want a finite cap below 20", len(entries))
	}
}
