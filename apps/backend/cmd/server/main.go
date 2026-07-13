package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/zigglib/webapp-monorepo/backend/internal/calc"
)

type operands struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

type resultResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /api/add", handleOp(func(a, b float64) (float64, error) {
		return calc.Add(a, b), nil
	}))
	mux.HandleFunc("POST /api/subtract", handleOp(func(a, b float64) (float64, error) {
		return calc.Subtract(a, b), nil
	}))
	mux.HandleFunc("POST /api/multiply", handleOp(func(a, b float64) (float64, error) {
		return calc.Multiply(a, b), nil
	}))
	mux.HandleFunc("POST /api/divide", handleOp(calc.Divide))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	addr := ":" + port
	log.Printf("backend listening on %s", addr)
	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleOp(op func(a, b float64) (float64, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in operands
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		result, err := op(in.A, in.B)
		if err != nil {
			if errors.Is(err, calc.ErrDivideByZero) {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			writeError(w, http.StatusInternalServerError, "operation failed")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resultResponse{Result: result})
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: msg})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
