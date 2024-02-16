package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/src/backend/database/db"
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
		http.Error(w, "length of settings is not 4", http.StatusBadRequest)
		return
	}

	err := db.Insert(userId, task, settingsString)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "expression uploaded")
}

func GetExpList(w http.ResponseWriter, r *http.Request) {
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
