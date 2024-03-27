package accessible

import (
	"context"
	"errors"
	pb "github.com/rendizi/daee/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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

	for _, agent := range Agents {
		log.Println(agent)
		go func(addr string) {
			conn, err := ping(addr)
			if err == nil {
				ch <- conn
				cancel()
			}
		}(agent.Addr)
	}

	// Wait for any of the goroutines to finish
	select {
	case conn := <-ch:
		return conn
	case <-ctx.Done():
		log.Println("No accessible agents")
		return nil
	}
}

func ping(addr string) (*grpc.ClientConn, error) {
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
	log.Println(addr, "Addr is busy")
	return nil, errors.New("agent is busy")
}
