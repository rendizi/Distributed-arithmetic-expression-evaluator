package distributer

import (
	"context"
	"fmt"
	pb "github.com/rendizi/daee/proto"
	"github.com/rendizi/daee/src/API/orkestrator/internal/accessible"
	"github.com/rendizi/daee/src/API/orkestrator/internal/db"
	"log"
	"os"
	"strconv"
	"strings"
)

// Данная функция делит выражение на операции и
// вычесляет их
func Do(expr db.Expression, id int64) {
	log.Println("Started")
	opId, err := db.InsertOperation(expr, id)
	if err != nil {
		log.Println(err)
	}
	symbols := strings.Fields(expr.Expression)
	n := len(symbols) - 1
	for i := 1; i < n; i += 2 {
		if symbols[i] == "*" || symbols[i] == "/" {
			subResult, err := subSolve(symbols[i-1], symbols[i], symbols[i+1], expr.Settings)
			if err != nil {
				db.UpdateResult(opId, id, err.Error())
			}
			symbols[i-1] = subResult
			symbols = append(symbols[:i], symbols[i+2:]...)
			n -= 2
			i -= 2
			err = db.UpdateOperationState(strings.Join(symbols, " "), opId)
			if err != nil {
				db.UpdateResult(opId, id, err.Error())
			}
			log.Println(symbols)
		}
	}
	for i := 1; i < n; i += 2 {
		if symbols[i] == "+" || symbols[i] == "-" {
			subResult, err := subSolve(symbols[i-1], symbols[i], symbols[i+1], expr.Settings)
			if err != nil {
				db.UpdateResult(opId, id, err.Error())
			}
			symbols[i-1] = subResult
			symbols = append(symbols[:i], symbols[i+2:]...)
			n -= 2
			i -= 2
			err = db.UpdateOperationState(strings.Join(symbols, " "), opId)
			if err != nil {
				db.UpdateResult(opId, id, err.Error())
			}
			log.Println(symbols)
		}
	}
	db.UpdateResult(opId, id, symbols[0])
}

func subSolve(a, operator, b string, settings db.Settings) (string, error) {
	inta, _ := strconv.Atoi(a)
	intb, _ := strconv.Atoi(b)
	log.Println("before")
	conn := accessible.GetAgent() //Zdes
	log.Println("after")
	defer conn.Close()
	grpcClient := pb.NewAgentServiceClient(conn)
	time := 0
	switch b {
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
		Operator: b,
		Time:     int64(time),
	})
	log.Println(res.Result)
	if err != nil {
		log.Println("failed invoking Area:", err)
		os.Exit(1)
	}
	return fmt.Sprintf("%.2f", res.Result), nil
}
