package main

import (
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/cmd/agent"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/cmd/orkestrator"
	"log"
	"os"
)

func main() {
	em := os.Getenv("EM")
	log.Println(em)
	if em == "O" {
		orkestrator.Main()
	} else if em == "A" {
		agent.Main()
	}
}
