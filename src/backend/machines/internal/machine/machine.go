package machine

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Machine struct {
	Server http.Server
	Addr   string
	Id     string
	Status bool
}

type MachinesStruct struct {
	Map   map[string]Machine
	mutex sync.Mutex
}

var Machines = &MachinesStruct{
	Map:   make(map[string]Machine),
	mutex: sync.Mutex{},
}

type port struct {
	port  int
	mutex sync.Mutex
}

var portId = port{
	port:  8081,
	mutex: sync.Mutex{},
}

// Создавая новую машину задаем ему порт, который далее увеличиваем и отправляем запрос на
// регистрацию машины с некоторым именем и портом . Задаем ему хэндлер по которому он будет слушать запросы
func New(name string) (*Machine, error) {
	portId.mutex.Lock()
	machinePort := portId.port
	portId.port++
	portId.mutex.Unlock()
	response, err := http.Get("http://127.0.0.1:8080/reg?name=" + name + "&port=" + fmt.Sprintf("%d", machinePort))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	id, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	machine := &Machine{
		Server: http.Server{
			Addr: fmt.Sprintf(":%d", machinePort),

			Handler: http.HandlerFunc(Handle),
		},
		Addr:   fmt.Sprintf(":%d", machinePort),
		Id:     string(id),
		Status: false,
	}

	Machines.mutex.Lock()
	Machines.Map[string(id)] = *machine
	Machines.mutex.Unlock()

	return machine, nil
}

func Handle(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "server is running")
}
