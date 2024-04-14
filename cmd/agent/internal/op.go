package internal

import (
	"context"
	daee "github.com/rendizi/Distributed-arithmetic-expression-evaluator/proto"
	"log"
	"time"
)

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
