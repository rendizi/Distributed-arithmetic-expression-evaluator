package main

import (
	"fmt"
	"net/http"
	"os"

	handle "github.com/rendizi/Distributed-arithmetic-expression-evaluator/src/backend/daee/server/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/expression", handle.Expression)
	mux.HandleFunc("/operations", handle.Operations)
	mux.HandleFunc("/reg", handle.Register)
	mux.HandleFunc("/task", handle.Task)
	mux.HandleFunc("/machines", handle.Machines)

	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*") // Allow requests from any origin
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}

	err := http.ListenAndServe(":8080", corsHandler(mux))
	if err != nil {
		if err == http.ErrServerClosed {
			fmt.Println("server closed")
		} else {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
	}
}
