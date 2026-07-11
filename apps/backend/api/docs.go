package api

import (
	"embed"
	"net/http"
)

//go:embed openapi.yaml swagger.html
var files embed.FS

// Handler serves OpenAPI spec and Swagger UI.
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		data, err := files.ReadFile("openapi.yaml")
		if err != nil {
			http.Error(w, "spec not found", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(data)
	})
	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
	mux.HandleFunc("GET /swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data, err := files.ReadFile("swagger.html")
		if err != nil {
			http.Error(w, "swagger ui not found", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(data)
	})
	return mux
}
