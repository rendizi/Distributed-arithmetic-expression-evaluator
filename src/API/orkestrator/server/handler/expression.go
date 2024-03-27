package handler

import (
	"encoding/json"
	"github.com/rendizi/daee/src/API/orkestrator/internal/db"
	"github.com/rendizi/daee/src/API/orkestrator/internal/distributer"
	"github.com/rendizi/daee/src/API/orkestrator/server"
	"net/http"
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

}

func PostExpression(w http.ResponseWriter, r *http.Request) {
	var expr db.Expression

	err := json.NewDecoder(r.Body).Decode(&expr)
	if err != nil {
		server.Error(map[string]interface{}{"message": "Data is not provided", "status": 400}, w)
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
	distributer.Do(expr, id)
	server.Ok(map[string]interface{}{"message": "Expression added successfully", "id": id, "status": 200}, w)
}
