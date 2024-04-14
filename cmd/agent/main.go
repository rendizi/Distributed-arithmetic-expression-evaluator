package agent

import (
	"context"
	"fmt"
	daee "github.com/rendizi/Distributed-arithmetic-expression-evaluator/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
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

//Данная функция отвечает на запросы
//занят ли сервер
//		message AvRequest {}

//			Сообщение для описания результата доступности агента(да/нет)
//			message AvResponse {
//	 		bool result = 1;
//			}
func (s *Server) Av(
	ctx context.Context,
	in *daee.AvRequest,
) (*daee.AvResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &daee.AvResponse{
		Result: !s.busy,
	}, nil
}

// Данная функция вычисляет передаваемое
// выражение
//
//			Сообщение для описания запроса на вычисление выражения
//			message OpRequest {
//	 		float a = 1;
//	 		float b = 2;
//	 		string operator = 3;
//	 		int64 time = 4;
//			}
//			Сообщение для описания результата вычисления выражения
//			message OpResponse {
//	 		float result = 1;
//			}
func (s *Server) Op(
	ctx context.Context,
	in *daee.OpRequest,
) (*daee.OpResponse, error) {
	s.mu.Lock()
	s.busy = !s.busy
	//теперь сервер занят
	defer s.mu.Unlock()
	var res float32
	//делаем действия в зависимости от оператора
	log.Println(in.A, in.Operator, in.B)
	switch in.Operator {
	case "+":
		res = in.A + in.B
	case "-":
		res = in.A - in.B
	case "*":
		res = in.A * in.B
	case "/":
		if in.B == 0 {
			//На ноль делить нельзя, результат 0, но
			//в орекстраторе это также проверяется и
			//по идеи данное выражение не должно добраться
			//до агента и выйдет ошибка вместо 0
			res = 0.0
		} else {
			res = in.A / in.B
		}
	}
	//спим необходимое время
	time.Sleep(time.Duration(in.Time) * time.Second)
	log.Println(res)
	//теперь агент не занят
	s.busy = !s.busy
	return &daee.OpResponse{
		Result: res,
	}, nil
}

// Функция для создания агента, передаем порт на котором будет работать
func createAgent(port int) (*grpc.Server, net.Listener) {
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

func Main() {
	log.Println("Agent is running")
	var wg sync.WaitGroup
	port := 5000
	i := 0
	//создаем три агента
	for i < 3 {
		wg.Add(1)
		grpcServer, lis := createAgent(port)
		port++
		i++
		//в горутинах слушаем каждый адрес
		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				log.Println("error serving grpc: ", err)
			}
		}()
	}
	wg.Wait()
}
