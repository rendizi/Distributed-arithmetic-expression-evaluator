package accessible

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestGetAgent(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		now := time.Now()
		_ = GetAgent()
		log.Println(1, time.Now().Sub(now))
	}()
	go func() {
		defer wg.Done()
		now := time.Now()
		_ = GetAgent()
		log.Println(2, time.Now().Sub(now))
	}()
	go func() {
		defer wg.Done()
		now := time.Now()
		_ = GetAgent()
		log.Println(3, time.Now().Sub(now))
	}()
	wg.Wait()
}
