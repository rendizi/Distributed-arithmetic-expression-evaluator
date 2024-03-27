package distributer

import (
	"encoding/json"
	"github.com/rendizi/daee/src/API/orkestrator/internal/db"
	"strings"
)

// Данная функция делит выражение на операции и
// вычесляет их
func Do(expr db.Expression, id int64) error {
	symbols := strings.Fields(expr.Expression)
	n := len(symbols) - 1
	for i := 1; i < n; i += 2 {
		if symbols[i] == "*" || symbols[i] == "/" {
			sets, _ := json.Marshal(expr.Settings)
			subResult, err := subSolve(symbols[i-1], symbols[i], symbols[i+1], string(sets))
			if err != nil {
				return err
			}
			symbols[i-1] = subResult
			symbols = append(symbols[:i], symbols[i+2:]...)
			n -= 2
			i -= 2
		}
	}
	return nil
}

func subSolve(a, operator, b, settings string) (string, error) {
	//inta, _ := strconv.Atoi(a)
	//intb, _ := strconv.Atoi(b)
	//conn := accessible.GetAgent()
	//defer conn.Close()
	//grpcClient := pb.NewAgentServiceClient(conn)
	//sets := 0
	//switch b {
	//case "+":
	//	sets = settings[0]
	//
	//}
	//
	//res, err := grpcClient.Op(context.Background(), &pb.OpRequest{
	//	A: int64(inta),
	//	B: int64(intb),
	//	Operator: b,
	//	Time: int64(),
	//})
	//if err != nil {
	//	log.Println("failed invoking Area:", err)
	//	os.Exit(1)
	//}
	return "", nil
}
