package main

import (
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/agent/internal"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	port := 5000
	i := 0

	for i < 3 {
		wg.Add(1)
		grpcServer, lis := internal.CreateAgent(port)
		mutex.Lock()
		port++
		mutex.Unlock()
		i++

		go func() {
			defer wg.Done()
			if err := grpcServer.Serve(lis); err != nil {
				log.Println("error serving grpc: ", err)
			}
		}()
	}
	wg.Wait()
}
