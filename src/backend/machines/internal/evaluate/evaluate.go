package evaluate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Solve(task string, settings []int) (string, string) {
	start := time.Now()
	symbols := strings.Fields(task)
	n := len(symbols) - 1
	for i := 1; i < n; i += 2 {
		if symbols[i] == "*" || symbols[i] == "/" {
			subResult, err := subSolve(symbols[i-1], symbols[i], symbols[i+1], settings)
			if err != nil {
				return err.Error(), ""
			}
			symbols[i-1] = subResult
			symbols = append(symbols[:i], symbols[i+2:]...)
			n -= 2
			i -= 2
		}
	}
	for i := 1; i < n; i += 2 {
		if symbols[i] == "+" || symbols[i] == "-" {
			subResult, err := subSolve(symbols[i-1], symbols[i], symbols[i+1], settings)
			if err != nil {
				return err.Error(), ""
			}
			symbols[i-1] = subResult
			symbols = append(symbols[:i], symbols[i+2:]...)
			n -= 2
			i -= 2
		}
	}

	return symbols[0], time.Since(start).String()
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
