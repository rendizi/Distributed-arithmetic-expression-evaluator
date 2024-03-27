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
	LastPing time.Duration
	IsBusy   bool
}

func newAgent(addr string) Agent {
	return Agent{
		Addr:   addr,
		IsBusy: false,
	}
}

var Agents = make([]Agent, 3)

func Init() {
	Agents[0] = newAgent("localhost:5000")
	Agents[1] = newAgent("localhost:5001")
	Agents[2] = newAgent("localhost:5002")
}

func ping(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	grpcClient := pb.NewAgentServiceClient(conn)
	av, err := grpcClient.Av(context.Background(), &pb.AvRequest{})
	if err != nil {
		log.Println("failed invoking Area:", err)
	}
	if av.Result {
		return conn, nil
	}
	return nil, errors.New("busy")
}

func GetAgent() *grpc.ClientConn {
	var wg sync.WaitGroup
	ch := make(chan *grpc.ClientConn)

	for {
		for _, agent := range Agents {
			wg.Add(1)
			go func(addr string) {
				defer wg.Done()
				conn, err := ping(addr)
				if err == nil {
					ch <- conn
				}
			}(agent.Addr)
			wg.Wait()

			select {
			case conn := <-ch:
				return conn
			default:
				continue
			}
		}
	}
}
