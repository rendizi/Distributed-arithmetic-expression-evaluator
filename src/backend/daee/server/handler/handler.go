package handle

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	server "github.com/rendizi/Distributed-arithmetic-expression-evaluator/src/backend/daee/server"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/src/backend/database/db"
)

// Virazheniya
var (
	compId       = 0
	machinesList *MachinesList
)

type MachinesList struct {
	List  [][]string
	mutex sync.Mutex
}

func init() {
	machinesList = &MachinesList{List: make([][]string, 0), mutex: sync.Mutex{}}
}

func Expression(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		server.PostExp(w, r)
	} else if r.Method == http.MethodGet {
		if len(r.URL.Query().Get("id")) == 0 {
			http.Error(w, "id not found", http.StatusBadRequest)
			return
		}
		server.GetExpList(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func Operations(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		server.GetOps(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func Task(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := r.URL.Query().Get("id")
		if len(id) == 0 {
			http.Error(w, "id not found", http.StatusBadRequest)
		} else {
			server.GetTask(w, r)
		}
	} else if r.Method == http.MethodPost {
		id := r.URL.Query().Get("id")
		answ := r.URL.Query().Get("answ")
		task := r.URL.Query().Get("task")
		time := r.URL.Query().Get("time")
		fmt.Println(id, answ, task, time)

		if len(id) == 0 || len(answ) == 0 || len(task) == 0 {
			http.Error(w, "Empty fields", http.StatusBadRequest)
			return
		}

		err := db.Update(id, task, answ, time)
		if err != nil {
			log.Println("Failed to update:", err)
			http.Error(w, "Failed to update", http.StatusInternalServerError)
			return
		}

		fmt.Println("Updated")
		w.WriteHeader(http.StatusOK)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		http.Error(w, "no name provided", http.StatusBadRequest)
		return
	}
	port := r.URL.Query().Get("port")
	if len(port) == 0 {
		http.Error(w, "no port provided", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, compId)
	compId++
	info := []string{name, port, time.November.String(), "alive"}
	machinesList.mutex.Lock()
	machinesList.List = append(machinesList.List, info)
	machinesList.mutex.Unlock()
	go pingMachine()
}

func Machines(w http.ResponseWriter, r *http.Request) {
	machinesList.mutex.Lock()
	w.WriteHeader(http.StatusOK)
	for _, machine := range machinesList.List {
		fmt.Fprintln(w, machine[0]+": "+machine[1]+".Last ping: "+machine[2]+".Status: "+machine[3])
	}
	machinesList.mutex.Unlock()
}

func pingMachine() {
	machinesList.mutex.Lock()
	for _, machine := range machinesList.List {
		resp, err := http.Get("http://127.0.0.1:" + machine[1] + "/")
		if err != nil {
			log.Println(err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			machine[3] = "dead"
		}
		machine[2] = time.Now().String()
	}
	machinesList.mutex.Unlock()
	time.Sleep(1 * time.Minute)
}
