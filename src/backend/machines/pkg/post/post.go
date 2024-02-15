package post

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// функция для отправки ответа на задание
func Task(id, answer, task, time string) error {
	baseURL := "http://127.0.0.1:8080/task"
	query := url.Values{}
	query.Set("task", task)
	query.Set("answ", answer)
	query.Set("id", id)
	query.Set("time", time)
	url := baseURL + "?" + query.Encode()

	reqBody, err := json.Marshal(map[string]string{"key": "value"})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status: " + resp.Status)
	}

	return nil
}
