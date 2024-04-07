package distributer

import (
	"context"
	"errors"
	"fmt"
	"github.com/rendizi/daee/internal/accessible"
	db2 "github.com/rendizi/daee/internal/db"
	"github.com/rendizi/daee/proto"
	"log"
	"strconv"
	"strings"
)

// Данная функция делит выражение на операции и
// вычесляет их
// Логика простая:
//   - Проходимся сначала по * и /
//   - Находим свободного агента
//   - Отправляем ему мини выражение
//   - Ждем ответа и заменяем 3 значение- a+b на result
//   - Проходимся по + и - с той же логикой
//   - Возвращаем единственное оставшиеся число- arr[0]
func Do(expr db2.Expression, id int64) {
	//добавляем операцию в бд
	opId, err := db2.InsertOperation(expr, id)
	if err != nil {
		return
	}
	//разделяем string в []string
	symbols := strings.Fields(expr.Expression)
	n := len(symbols) - 1
	for i := 1; i < n; i += 2 {
		if symbols[i] == "*" || symbols[i] == "/" {
			//Вычисляем a*b or a/b
			subResult, err := subSolve(symbols[i-1], symbols[i], symbols[i+1], expr.Settings)
			if err != nil {
				//если ошибка то так и записываем в бд и выходим
				db2.UpdateResult(opId, id, err.Error())
				return
			}
			//вставляем результат в []string
			symbols[i-1] = subResult
			symbols = append(symbols[:i], symbols[i+2:]...)
			n -= 2
			i -= 2
			//Обновляем операцию
			err = db2.UpdateOperationState(strings.Join(symbols, " "), opId)
			if err != nil {
				db2.UpdateResult(opId, id, err.Error())
				return
			}
		}
	}
	for i := 1; i < n; i += 2 {
		if symbols[i] == "+" || symbols[i] == "-" {
			subResult, err := subSolve(symbols[i-1], symbols[i], symbols[i+1], expr.Settings)
			if err != nil {
				db2.UpdateResult(opId, id, err.Error())
				return
			}
			symbols[i-1] = subResult
			symbols = append(symbols[:i], symbols[i+2:]...)
			n -= 2
			i -= 2
			err = db2.UpdateOperationState(strings.Join(symbols, " "), opId)
			if err != nil {
				db2.UpdateResult(opId, id, err.Error())
				return
			}
		}
	}
	//удаляем операцию и обновляем результат выражения
	db2.DeleteOperation(opId)
	db2.UpdateResult(opId, id, symbols[0])
}

func subSolve(a, operator, b string, settings db2.Settings) (string, error) {
	if b == "0" && operator == "/" {
		return "", errors.New("Division by zero")
	}

	//Парсим в флоут
	inta, err := strconv.ParseFloat(a, 32)
	if err != nil {
		log.Println(err, inta)
	}
	intb, err := strconv.ParseFloat(b, 32)
	if err != nil {
		log.Println(err, intb)
	}

	//Находим доступного агента
	conn := accessible.GetAgent()
	defer conn.Close()
	grpcClient := daee.NewAgentServiceClient(conn)
	time := 0
	switch operator {
	case "+":
		time = settings.Plus
	case "-":
		time = settings.Minus
	case "*":
		time = settings.Mult
	case "/":
		time = settings.Div
	}

	//Вычисляем выражение
	res, err := grpcClient.Op(context.Background(), &daee.daee{
		A:        float32(inta),
		B:        float32(intb),
		Operator: operator,
		Time:     int64(time),
	})
	if err != nil {
		log.Println("failed solving:", err)
	}
	return fmt.Sprintf("%.2f", res.Result), nil
}
