package main

import (
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/cmd/agent"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/cmd/orkestrator"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	em := os.Getenv("EM")
	log.Println(em)
	if em == "O" {
		orkestrator.Main()
	} else if em == "A" {
		agent.Main()
	}
}
