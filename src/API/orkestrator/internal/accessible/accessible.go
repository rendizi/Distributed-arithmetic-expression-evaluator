package accessible

import (
	"context"
	"errors"
	pb "github.com/rendizi/daee/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

type Agent struct {
	Addr     string
	LastPing string
	IsBusy   bool
	m        sync.Mutex
}

type AgentJson struct {
	Addr     string `json:"address"`
	IsBusy   bool   `json:"is_busy"`
	LastPing string `json:"last_ping"`
}

func newAgent(addr string) Agent {
	return Agent{
		Addr:     addr,
		LastPing: time.Now().Format("2006-01-02 15:04:05"),
		IsBusy:   false,
		m:        sync.Mutex{},
	}
}

var Agents = make([]Agent, 3)

func init() {
	Agents[0] = newAgent("localhost:5000")
	Agents[1] = newAgent("localhost:5001")
	Agents[2] = newAgent("localhost:5002")
}

func GetAgent() *grpc.ClientConn {
	for {
		select {
		case <-time.After(1 * time.Second):
			log.Println("Searching for available agents...")
			if conn := findAvailableAgent(); conn != nil {
				return conn
			}
		}
	}
}

func findAvailableAgent() *grpc.ClientConn {
	ch := make(chan *grpc.ClientConn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := range Agents {
		agent := &Agents[i]
		agent.m.Lock()
		go func(addr string) {
			conn, err := Ping(addr)
			agent.LastPing = time.Now().Format("2006-01-02 15:04:05")
			log.Println(agent.LastPing)
			if err == nil {
				ch <- conn
				cancel()
			} else {
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
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.New("failed to dial grpc server")
	}
	grpcClient := pb.NewAgentServiceClient(conn)
	av, err := grpcClient.Av(context.Background(), &pb.AvRequest{})
	if err != nil {
		return nil, errors.New("failed to check availability with Av RPC")
	}
	if av.Result {
		return conn, nil
	}
	return nil, errors.New("agent is busy")
}
