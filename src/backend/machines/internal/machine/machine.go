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

func New(name string) (*Machine, error) {
	//увеличиваем порт
	portId.mutex.Lock()
	machinePort := portId.port
	portId.port++
	portId.mutex.Unlock()
	//регистрируем вычислительную машину
	response, err := http.Get("http://127.0.0.1:8080/reg?name=" + name + "&port=" + fmt.Sprintf("%d", machinePort))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	id, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	//хэндлим и сздаем новую машину
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handle)

	// Adding CORS middleware to allow requests from a web page
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	// Wrap the mux with CORS middleware
	handler := corsHandler(mux)

	machine := &Machine{
		Server: http.Server{
			Addr:    fmt.Sprintf(":%d", machinePort),
			Handler: handler,
		},
		Addr:   fmt.Sprintf(":%d", machinePort),
		Id:     string(id),
		Status: false,
	}

	// Start the server
	//go func() {
	//	if err := machine.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		log.Fatalf("HTTP server ListenAndServe: %v", err)
	//	}
	//}()
	// You may also want to handle graceful shutdown of the server
	// Here's an example of how you can do it:
	// shutdown := make(chan os.Signal, 1)
	// signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	// <-shutdown
	// log.Println("Shutting down server...")
	// if err := machine.Server.Shutdown(context.Background()); err != nil {
	//     log.Fatalf("Server shutdown failed: %v", err)
	// }
	// log.Println("Server gracefully stopped")

	Machines.mutex.Lock()
	Machines.Map[string(id)] = *machine
	Machines.mutex.Unlock()

	return machine, nil
}

func Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "server is running")
}
