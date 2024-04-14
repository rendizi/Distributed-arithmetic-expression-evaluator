package accessible

import (
	"context"
	"errors"
	daee "github.com/rendizi/Distributed-arithmetic-expression-evaluator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

// структура агента
type Agent struct {
	Addr     string
	LastPing string
	IsBusy   bool
	m        sync.Mutex
}

// json структура агента, которую мы возвращаем при запросе
type AgentJson struct {
	Addr     string `json:"address"`
	IsBusy   bool   `json:"is_busy"`
	LastPing string `json:"last_ping"`
}

// новый агент, сразу даем последний пинг- сейчас
func newAgent(addr string) Agent {
	return Agent{
		Addr:     addr,
		LastPing: time.Now().Format("2006-01-02 15:04:05"),
		IsBusy:   false,
		m:        sync.Mutex{},
	}
}

var Agents = make([]Agent, 3)

// на старте программы создаем агентом
func init() {
	Agents[0] = newAgent("localhost:5000")
	Agents[1] = newAgent("localhost:5001")
	Agents[2] = newAgent("localhost:5002")
}

// данная функция нужна для получения свободного агента
func GetAgent() *grpc.ClientConn {
	for {
		select {
		//каждую секунду ищет и возвращает соединение с агентом
		case <-time.After(1 * time.Second):
			log.Println("In search...")
			conn := findAvailableAgent()
			if conn != nil {
				log.Println("Connection found")
				return conn
			}
		}
	}
}

func findAvailableAgent() *grpc.ClientConn {
	ch := make(chan *grpc.ClientConn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//запускаем цикл по всем агентам и проверяем доступны ли они
	log.Println(Agents)
	for i := range Agents {
		agent := &Agents[i]
		agent.m.Lock()
		go func(addr string) {
			log.Println("Pinging ", addr)
			conn, err := Ping(addr)
			//заодно обновляем его last ping
			agent.LastPing = time.Now().Format("2006-01-02 15:04:05")
			log.Println("there")
			if err == nil {
				ch <- conn
				cancel()
			} else {
				log.Println(err)
				agent.IsBusy = true
			}
		}(agent.Addr)
		agent.m.Unlock()
	}

	select {
	case conn := <-ch:
		return conn
	case <-ctx.Done():
		return nil
	}
}

func Ping(addr string) (*grpc.ClientConn, error) {
	//устанавливаем соединение и если ошибок нет- возвращаем соединение
	//Если av.Result- false, значит агент занят
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to dial grpc server")
	}
	grpcClient := daee.NewAgentServiceClient(conn)
	av, err := grpcClient.Av(context.Background(), &daee.AvRequest{})
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to check availability with Av RPC")
	}
	log.Println(av.Result)
	if av.Result {
		return conn, nil
	}
	return nil, errors.New("agent is busy")
}
