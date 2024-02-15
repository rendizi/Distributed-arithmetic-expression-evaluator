package evaluate

import (
	"errors"
	"fmt"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/src/backend/machines/pkg/post"
	"log"
	"strconv"
	"strings"
	"time"
)

//Идея такова: проходимся первый раз, решаем все * и / и вместо 3 значений (первое число, знак
// и второй число) ставлю ответ. Далее прохожусь во второй раз и решаю + и -, с той же логикой заменой
// 2 + 2 * 2 -> 2 * 2 = 4 ; 2 + 4 - > 2 + 4 = 6 - > 6
//В конце просто результат = единственному значению

// Проверяю сколько прошло времени и отправляю запрос на пост ответа
func Solve(task string, settings []int) {
	start := time.Now()
	symbols := strings.Fields(task)
	n := len(symbols) - 1
	result := ""
	isErr := false

	for i := 1; i < n; i += 2 {
		if symbols[i] == "*" || symbols[i] == "/" {
			subResult, err := subSolve(symbols[i-1], symbols[i], symbols[i+1], settings)
			if err != nil {
				result = err.Error()
				isErr = true
				break
			}
			symbols[i-1] = subResult
			symbols = append(symbols[:i], symbols[i+2:]...)
			n -= 2
			i -= 2
		}
	}

	for i := 1; i < n; i += 2 {
		if !isErr {
			break
		}
		if symbols[i] == "+" || symbols[i] == "-" {
			subResult, err := subSolve(symbols[i-1], symbols[i], symbols[i+1], settings)
			if err != nil {
				result = err.Error()
				isErr = true
				break
			}
			symbols[i-1] = subResult
			symbols = append(symbols[:i], symbols[i+2:]...)
			n -= 2
			i -= 2
		}
	}
	if !isErr {
		result = symbols[0]
	} else {
		result = "error"
	}
	time := time.Since(start).String()
	err := post.Task(string(task[0]), result, task[1:], time)
	if err != nil {
		log.Println("Error posting task result:", err.Error())
	}
}

func subSolve(first, znak, second string, settings []int) (string, error) {
	firstNum, err := strconv.Atoi(first)
	if err != nil {
		return "", err
	}
	secondNum, err := strconv.Atoi(second)
	if err != nil {
		return "", err
	}

	if znak == "*" {
		time.Sleep(time.Duration(settings[0]) * time.Second)
		return fmt.Sprint(firstNum * secondNum), nil
	}
	if znak == "/" && second != "0" {
		time.Sleep(time.Duration(settings[1]) * time.Second)
		return fmt.Sprint(firstNum / secondNum), nil
	}
	if znak == "+" {
		time.Sleep(time.Duration(settings[2]) * time.Second)
		return fmt.Sprint(firstNum + secondNum), nil
	}
	if znak == "-" {
		time.Sleep(time.Duration(settings[3]) * time.Second)
		return fmt.Sprint(firstNum - secondNum), nil
	}
	return "", errors.New("not valid operator for sequnce")
}
