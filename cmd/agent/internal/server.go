package internal

import (
	"context"
	"fmt"
	daee "github.com/rendizi/Distributed-arithmetic-expression-evaluator/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

// Создаем структуру сервера-агента, их у нас 3
type Server struct {
	daee.AgentServiceServer
	busy bool
	mu   sync.Mutex
}

// Функция для создания нового сервера
func NewServer() *Server {
	return &Server{}
}

//			service AgentService {
//	 		методы, которые можно будет реализовать и использовать
//	 		rpc Av (AvRequest) returns (AvResponse);
//	 		rpc Op (OpRequest) returns (OpResponse);
//			}
type AgentServiceServer interface {
	Av(context.Context, *daee.AvRequest) (*daee.AvResponse, error)
	Op(context.Context, *daee.OpRequest) (*daee.OpResponse, error)
	mustEmbedUnimplementedGeometryServiceServer()
}

// Функция для создания агента, передаем порт на котором будет работать
func CreateAgent(port int) (*grpc.Server, net.Listener) {
	addr := fmt.Sprintf("%s:%v", "localhost", port)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		return nil, nil
	}

	log.Println("tcp listener started at port: ", port)
	grpcServer := grpc.NewServer()
	// объект структуры, которая содержит реализацию
	// серверной части GeometryService
	geomServiceServer := NewServer()
	// зарегистрируем нашу реализацию сервера
	daee.RegisterAgentServiceServer(grpcServer, geomServiceServer)
	return grpcServer, lis
}
