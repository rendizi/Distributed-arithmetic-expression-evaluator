package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/src/backend/machines/internal/evaluate"
	machine "github.com/rendizi/Distributed-arithmetic-expression-evaluator/src/backend/machines/internal/machine"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/src/backend/machines/internal/reqs"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/src/backend/machines/pkg/post"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		First, err := machine.New("Superpupermachine")
		if err != nil {
			log.Println("Error creating machine number 1:", err.Error())
			return
		}
		fmt.Println("Machine is waked up")
		go func() {
			for {
				if settingsString, resp, err := reqs.Task(); err == nil && len(resp) != 0 {
					task := divideResp(resp)

					settings, err := settings(settingsString)
					if err != nil {
						_ = post.Task(string(resp[0]), "error", resp[1:], "")
						time.Sleep(5 * time.Second)
						continue
					}
					result, time := evaluate.Solve(task[1:], settings)

					if !isInt(result) {
						result = "error"
					}

					err = post.Task(string(task[0]), result, task[1:], time)
					if err != nil {
						log.Println("Error posting task result:", err.Error())
					}

				}
				time.Sleep(5 * time.Second)
			}
		}()
		err = First.Server.ListenAndServe()
		if err != nil {
			log.Println("Error starting server on machine number 1:", err.Error())
			return
		}
	}()
	wg.Wait()
}

func isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func settings(settings string) ([]int, error) {
	settings = settings[1:]
	result := strings.Split(settings, ",")
	if len(result) != 4 {
		return nil, errors.New("length of settings should be 4")
	}
	output := make([]int, 4)
	var err error
	for i, val := range result {
		val = strings.TrimSpace(val)
		output[i], err = strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}

func divideResp(resp string) string {
	output := ""
	for i := 0; i < len(resp); i++ {
		if resp[i] == '&' {
			break
		}
		output += string(resp[i])
	}
	return output
}
