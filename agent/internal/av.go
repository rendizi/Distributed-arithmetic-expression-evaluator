package internal

import daee "github.com/rendizi/Distributed-arithmetic-expression-evaluator/proto"
import "context"

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
