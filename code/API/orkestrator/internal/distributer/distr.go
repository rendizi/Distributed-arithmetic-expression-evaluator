package distributer

import (
	"context"
	"errors"
	"fmt"
	"github.com/rendizi/daee/code/API/orkestrator/internal/accessible"
	"github.com/rendizi/daee/code/API/orkestrator/internal/db"
	pb "github.com/rendizi/daee/proto"
	"log"
	"strconv"
	"strings"
)

// Данная функция делит выражение на операции и
// вычесляет их
func Do(expr db.Expression, id int64) {
	opId, err := db.InsertOperation(expr, id)
	if err != nil {
		return
	}
	symbols := strings.Fields(expr.Expression)
	n := len(symbols) - 1
	for i := 1; i < n; i += 2 {
		if symbols[i] == "*" || symbols[i] == "/" {
			subResult, err := subSolve(symbols[i-1], symbols[i], symbols[i+1], expr.Settings)
			if err != nil {
				db.UpdateResult(opId, id, err.Error())
				return
			}
			symbols[i-1] = subResult
			symbols = append(symbols[:i], symbols[i+2:]...)
			n -= 2
			i -= 2
			err = db.UpdateOperationState(strings.Join(symbols, " "), opId)
			if err != nil {
				db.UpdateResult(opId, id, err.Error())
				return
			}
		}
	}
	for i := 1; i < n; i += 2 {
		if symbols[i] == "+" || symbols[i] == "-" {
			subResult, err := subSolve(symbols[i-1], symbols[i], symbols[i+1], expr.Settings)
			if err != nil {
				db.UpdateResult(opId, id, err.Error())
				return
			}
			symbols[i-1] = subResult
			symbols = append(symbols[:i], symbols[i+2:]...)
			n -= 2
			i -= 2
			err = db.UpdateOperationState(strings.Join(symbols, " "), opId)
			if err != nil {
				db.UpdateResult(opId, id, err.Error())
				return
			}
		}
	}
	db.DeleteOperation(opId)
	db.UpdateResult(opId, id, symbols[0])
}

func subSolve(a, operator, b string, settings db.Settings) (string, error) {
	if b == "0" && operator == "/" {
		return "", errors.New("Division by zero")
	}

	inta, err := strconv.ParseFloat(a, 32)
	if err != nil {
		log.Println(err, inta)
	}
	intb, err := strconv.ParseFloat(b, 32)
	if err != nil {
		log.Println(err, intb)
	}
	conn := accessible.GetAgent() //Zdes
	defer conn.Close()
	grpcClient := pb.NewAgentServiceClient(conn)
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

	res, err := grpcClient.Op(context.Background(), &pb.OpRequest{
		A:        float32(inta),
		B:        float32(intb),
		Operator: operator,
		Time:     int64(time),
	})
	if err != nil {
		log.Println("failed invoking Area:", err)
	}
	return fmt.Sprintf("%.2f", res.Result), nil
}
