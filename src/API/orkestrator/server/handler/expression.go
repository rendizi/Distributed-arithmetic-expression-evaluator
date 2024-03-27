package handler

import (
	"encoding/json"
	"github.com/rendizi/daee/src/API/orkestrator/internal/db"
	"github.com/rendizi/daee/src/API/orkestrator/internal/distributer"
	"github.com/rendizi/daee/src/API/orkestrator/server"
	"net/http"
	"strconv"
)

func Expression(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetExpression(w, r)
	} else if r.Method == http.MethodPost {
		PostExpression(w, r)
	} else {
		server.Error(map[string]interface{}{"message": "Method is not allowed", "status": 400}, w)
	}
}

func GetExpression(w http.ResponseWriter, r *http.Request) {
	expId := r.URL.Query().Get("id")
	if len(expId) == 0 {
		login := server.GetLogin(w, r)
		if login == "" {
			return
		}
		expressions, results, err := db.GetExpressions(login)
		if err != nil {
			server.Error(map[string]interface{}{"message": err.Error(), "status": 400}, w)
			return
		}
		if len(expressions) != len(results) {
			server.Error(map[string]interface{}{"message": "mismatched lengths of expressions and results slices", "status": 400}, w)
			return
		}

		var expressionsJSON []db.ExpressionJSON
		for i := 0; i < len(expressions); i++ {
			expressionJSON := db.ExpressionJSON{
				Expression: expressions[i],
				Result:     results[i],
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
		return
	}
	intExpId, err := strconv.Atoi(expId)
	if err != nil {
		server.Error(map[string]interface{}{"message": "Invalid id", "status": 400}, w)
		return
	}
	login := server.GetLogin(w, r)
	if login == "" {
		return
	}
	exp, res, err := db.GetExpression(intExpId, login)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error(), "status": 400}, w)
		return
	}
	server.Ok(map[string]interface{}{"expression": exp, "result": res}, w)
}

func PostExpression(w http.ResponseWriter, r *http.Request) {
	var expr db.Expression

	err := json.NewDecoder(r.Body).Decode(&expr)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error(), "status": 400}, w)
		return
	}

	login := server.GetLogin(w, r)

	if login == "" {
		return
	}

	id, err := db.InsertExpression(expr, login)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error(), "status": 400}, w)
		return
	}
	go distributer.Do(expr, id)
	server.Ok(map[string]interface{}{"message": "Expression added successfully", "id": id, "status": 200}, w)
}
