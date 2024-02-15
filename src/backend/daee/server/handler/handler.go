package handle

import (
	"fmt"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/src/backend/daee/db"
	"log"
	"net/http"
	"sync"
	"time"

	server "github.com/rendizi/Distributed-arithmetic-expression-evaluator/src/backend/daee/server"
)

// Virazheniya
var (
	compId        = 0
	machinesList  *MachinesList
	machinesCount = 0
)

type MachinesList struct {
	List  [][]string
	mutex sync.Mutex
}

// на запуске создаем список вычислительных машин
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
		//Данная часть добавляет ответ на некое выражение в бд
		id := r.URL.Query().Get("id")
		answ := r.URL.Query().Get("answ")
		task := r.URL.Query().Get("task")
		time := r.URL.Query().Get("time")

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
	//ниже мы добавляем новую вычислительную машину в список
	compId++
	info := []string{name, port, time.Now().String(), "alive"}
	machinesList.mutex.Lock()
	machinesList.List = append(machinesList.List, info)
	machinesList.mutex.Unlock()
	machinesCount++
	//если это первая машине, то запускаем горутину, переодически пингующая вычислительные машины
	if machinesCount == 1 {
		go pingMachine()
	}

}

func Machines(w http.ResponseWriter, r *http.Request) {
	//возвращает инфу о машинах
	machinesList.mutex.Lock()
	w.WriteHeader(http.StatusOK)
	for _, machine := range machinesList.List {
		fmt.Fprintln(w, machine[0]+": "+machine[1]+".Last ping: "+machine[2]+".Status: "+machine[3])
	}
	machinesList.mutex.Unlock()
}

func pingMachine() {
	//отправляет запрос в соответствии которым меняет статус
	for {
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
}
