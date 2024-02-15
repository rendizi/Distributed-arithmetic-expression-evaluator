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
		//создаем новую машину
		First, err := machine.New("Superpupermachine")
		if err != nil {
			log.Println("Error creating machine number 1:", err.Error())
			return
		}
		fmt.Println("Machine is waked up")
		go func() {
			for {
				//если есть какое-нибудь задание
				if settingsString, resp, err := reqs.Task(); err == nil && len(resp) != 0 {
					//отделяем ответ на выражение и настройки, которые разделены символом &
					task := divideResp(resp)
					//проверяем настройки на длину и переводим в целочисленный список
					settings, err := settings(settingsString)
					if err != nil {
						//иначе отправляем ошибку
						_ = post.Task(string(resp[0]), "error", resp[1:], "")
						time.Sleep(5 * time.Second)
						continue
					}
					//решаем
					evaluate.Solve(task[1:], settings)
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
