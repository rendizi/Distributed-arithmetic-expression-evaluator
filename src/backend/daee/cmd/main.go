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

	//Для того чтобы отправлять запросы с вебсайта на сервер нужно настроить cors, что тут и делаю
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
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
