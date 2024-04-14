package main

import (
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/cmd/agent/internal"
	"log"
	"sync"
)

func main() {
	log.Println("Agent is running")
	var wg sync.WaitGroup
	var mutex sync.Mutex // Add a mutex to synchronize access to port
	port := 5000
	i := 0
	// Create three agents
	for i < 3 {
		wg.Add(1)
		grpcServer, lis := internal.CreateAgent(port)
		mutex.Lock()
		port++
		mutex.Unlock()
		i++
		// Goroutines to serve each agent
		go func() {
			defer wg.Done() // Decrement the WaitGroup when done
			if err := grpcServer.Serve(lis); err != nil {
				log.Println("error serving grpc: ", err)
			}
		}()
	}
	wg.Wait()
	log.Println("bye bye")
}
