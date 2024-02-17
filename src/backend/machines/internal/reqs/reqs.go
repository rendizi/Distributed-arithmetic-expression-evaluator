package reqs

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// код нижу отправляет запрос на главный сервер, возвращая новое задание для выполнения
func Task() (string, string, error) {
	response, err := http.Get("http://127.0.0.1:8080/task?id=1")
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", "", errors.New("no task at the moment")
	}
	task, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}
	if len(string(task)) == 0 {
		return "", "", err
	}
	index := len(string(task)) - 1
	fmt.Println(string(task))
	for i, v := range string(task) {
		if v == '&' {
			index = i
			break
		}
	}

	return string(task[index:]), string(task[:index+1]), nil
}
