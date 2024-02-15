package server

import (
	"fmt"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/src/backend/daee/db"
	"net/http"
	"strings"
)

func PostExp(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("uid")
	if len(userId) == 0 {
		http.Error(w, "no user id provided", http.StatusBadRequest)
		return
	}
	task := r.URL.Query().Get("task")
	if len(task) == 0 {
		http.Error(w, "no task provided", http.StatusBadRequest)
		return
	}
	settingsString := r.URL.Query().Get("settings")
	if len(settingsString) == 0 {
		http.Error(w, "no settings provided", http.StatusBadRequest)
		return
	}
	settings := strings.Split(settingsString, ",")
	if len(settings) != 4 {
		http.Error(w, "length of settings is not 4:"+settingsString, http.StatusBadRequest)
		return
	}

	//проверяем на импатентность
	answer, time, _ := db.Solved(userId, task)

	if answer != "" {
		//если ответ не пустой- задание было уже решено и мы обновляем бд
		err := db.Update(userId, task, answer, time)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	//иначу добавляем в бд и дальше он будет добавлен в очередь
	err := db.Insert(userId, task, settingsString)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "expression uploaded")

}

func GetExpList(w http.ResponseWriter, r *http.Request) {
	//отправляет запрос в бд на получение списка всех выражений отправленные юзером
	expMap, err := db.Get(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	result := ""

	for _, value := range expMap {
		fmt.Println(value, "value")
		result += "userId: " + value[0] + ",expression: " + value[1] + ",answer: " + value[2] + ",tken time: " + value[3] + "\n"
	}
	fmt.Fprintln(w, result)
}

func GetOps(w http.ResponseWriter, r *http.Request) {
	//отправляет запрос в бд где нет ответов
	opMap, err := db.Get("ns")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	result := ""
	for _, value := range opMap {
		result += value[0] + "," + value[1] + "," + value[2] + "\n"
	}
	fmt.Fprintln(w, result)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	//из списка выражений у которых нет ответов берет 1 и возвращает
	tasks, err := db.Get("ns")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, value := range tasks {
		fmt.Println(value)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, value[1]+value[2]+"&"+value[3])
		break
	}
}
