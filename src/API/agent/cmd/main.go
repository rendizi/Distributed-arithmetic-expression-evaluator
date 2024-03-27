package main

import (
	"context"
	"fmt"
	pb "github.com/rendizi/daee/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"sync"
)

type Server struct {
	pb.AgentServiceServer
	busy bool
	mu   sync.Mutex
}

func NewServer() *Server {
	return &Server{}
}

type AgentServiceServer interface {
	Av(context.Context, *pb.AvRequest) (*pb.AvResponse, error)
	Op(context.Context, *pb.OpRequest) (*pb.OpResponse, error)
	mustEmbedUnimplementedGeometryServiceServer()
}

func (s *Server) Av(
	ctx context.Context,
	in *pb.AvRequest,
) (*pb.AvResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &pb.AvResponse{
		Result: !s.busy,
	}, nil
}

func (s *Server) Op(
	ctx context.Context,
	in *pb.OpRequest,
) (*pb.OpResponse, error) {
	log.Println("invoked Perimeter: ", in)
	// вычислим периметр и вернём ответ
	s.mu.Lock()
	defer s.mu.Unlock()
	var res float32
	switch in.Operator {
	case "+":
		res = in.A + in.B
	case "-":
		res = in.A - in.B
	case "*":
		res = in.A * in.B
	case "/":
		if in.B == 0 {
			res = 0.0
		} else {
			res = in.A / in.B
		}
	}
	return &pb.OpResponse{
		Result: res,
	}, nil
}

func createAgent(port int) (*grpc.Server, net.Listener) {
	addr := fmt.Sprintf("%s:%v", "localhost", port)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Println("tcp listener started at port: ", port)
	// создадим сервер grpc
	grpcServer := grpc.NewServer()
	// объект структуры, которая содержит реализацию
	// серверной части GeometryService
	geomServiceServer := NewServer()
	// зарегистрируем нашу реализацию сервера
	pb.RegisterAgentServiceServer(grpcServer, geomServiceServer)
	return grpcServer, lis
}

func main() {
	var wg sync.WaitGroup
	port := 5000
	i := 0
	for i < 3 {
		wg.Add(1)
		grpcServer, lis := createAgent(port)
		port++
		i++
		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				log.Println("error serving grpc: ", err)
				os.Exit(1)
			}
		}()
	}
	wg.Wait()
}
