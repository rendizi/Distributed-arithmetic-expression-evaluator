package handler

import (
	"encoding/json"
	"github.com/rendizi/daee/internal/db"
	"github.com/rendizi/daee/server"
	"net/http"
)

//Здесь тоже самое что и в ./expression.go

func Operations(w http.ResponseWriter, r *http.Request) {
	ids, operations, results, err := db.GetOperations()
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error(), "status": 400}, w)
		return
	}
	if len(operations) != len(results) {
		server.Error(map[string]interface{}{"message": "mismatched lengths of operations and results slices", "status": 400}, w)
		return
	}

	var expressionsJSON []db.OperationJSON
	for i := 0; i < len(operations); i++ {
		expressionJSON := db.OperationJSON{
			Id:        ids[i],
			Operation: operations[i],
			Result:    results[i],
		}
		expressionsJSON = append(expressionsJSON, expressionJSON)
	}

	jsonData, err := json.Marshal(expressionsJSON)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error(), "status": 400}, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
