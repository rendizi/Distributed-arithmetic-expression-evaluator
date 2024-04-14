package main

import (
	"fmt"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/internal/db"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/internal/handler"
	"log"
	"net/http"
	"os"

	"github.com/MadAppGang/httplog"
)

var (
	registerHandler   http.Handler = http.HandlerFunc(handler.Register)
	loginHandler      http.Handler = http.HandlerFunc(handler.Login)
	expressionHandler http.Handler = http.HandlerFunc(handler.Expression)
	agentsHandler     http.Handler = http.HandlerFunc(handler.Agents)
	operationsHandler http.Handler = http.HandlerFunc(handler.Operations)
)

func main() {
	db.Init()
	log.Println("Orkestrator is running")
	//mux нам нужен для cors
	mux := http.NewServeMux()

	//loggerWithFormatter - красивый логгер запросов
	loggerWithFormatter := httplog.LoggerWithFormatter(httplog.DefaultLogFormatterWithRequestHeader)
	mux.Handle("/register", loggerWithFormatter(registerHandler))
	mux.Handle("/login", loggerWithFormatter(loginHandler))
	mux.Handle("/expression", loggerWithFormatter(expressionHandler))
	mux.Handle("/agents", loggerWithFormatter(agentsHandler))
	mux.Handle("/operations", loggerWithFormatter(operationsHandler))

	//устанавливаем cors
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

	//слушаем на порту 8080
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
