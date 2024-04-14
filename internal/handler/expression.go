package handler

import (
	"encoding/json"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/internal/db"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/internal/distributer"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/internal/server"
	"log"
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
	//берем айди. Если есть айди- ищем в бд данное выражение,
	//иначе все выражения
	expId := r.URL.Query().Get("id")
	if len(expId) == 0 {
		//из JWT ключа берем логин
		login := server.GetLogin(w, r)
		//если нет логина то просто возвращаем, ошибка обрабатывается в функции
		if login == "" {
			return
		}
		//берем выражения
		ids, expressions, results, err := db.GetExpressions(login)
		if err != nil {
			server.Error(map[string]interface{}{"message": err.Error(), "status": 400}, w)
			return
		}
		if len(expressions) != len(results) {
			server.Error(map[string]interface{}{"message": "mismatched lengths of expressions and results slices", "status": 400}, w)
			return
		}

		//изменияем время
		var expressionsJSON []db.ExpressionJSON
		for i := 0; i < len(expressions); i++ {
			time, err := db.Time(int64(ids[i]))
			if err != nil {
				time = err.Error()
			}
			expressionJSON := db.ExpressionJSON{
				Id:         ids[i],
				Expression: expressions[i],
				Result:     results[i],
				Time:       time,
			}
			expressionsJSON = append(expressionsJSON, expressionJSON)
		}

		//маршаллим и возвращаем
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
	//Берем выражение по конкретному id
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
	//берем время
	time, _ := db.Time(int64(intExpId))
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error(), "status": 400}, w)
		return
	}
	//возвращаем json
	server.Ok(map[string]interface{}{"expression": exp, "result": res, "time": time}, w)
}

func PostExpression(w http.ResponseWriter, r *http.Request) {
	var expr db.Expression

	err := json.NewDecoder(r.Body).Decode(&expr)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error(), "status": 400}, w)
		return
	}
	if len(expr.Expression) == 0 {
		server.Error(map[string]interface{}{"message": "Expression can not be an empty string", "status": 400}, w)
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
	log.Println(expr)
	go distributer.Do(expr, id)
	server.Ok(map[string]interface{}{"message": "Expression added successfully", "id": id, "status": 200}, w)
}
