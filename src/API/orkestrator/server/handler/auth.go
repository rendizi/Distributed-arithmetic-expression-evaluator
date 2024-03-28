package handler

import (
	"encoding/json"
	"github.com/rendizi/daee/src/API/orkestrator/internal/db"
	"github.com/rendizi/daee/src/API/orkestrator/server"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Секретный ключ для генерации auth токенов
var jwtSecret = []byte("secret_key")

func Register(w http.ResponseWriter, r *http.Request) {
	//Берем из body json
	var creds db.UserJson

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		//Если ошибка- не предоставлена дата про юзера
		server.Error(map[string]interface{}{"message": "Data is not provided", "status": 400}, w)
		return
	}
	//добавляем в бд
	err = db.InsertUser(creds)
	if err != nil {
		//выводим ошибку в виде json
		server.Error(map[string]interface{}{"message": err.Error(), "status": 400}, w)
		return
	}
	//выводим сообщение о успехе в виде json
	server.Ok(map[string]interface{}{"message": "Signed-up successful", "status": 200}, w)
}

func Login(w http.ResponseWriter, r *http.Request) {
	//Получаем дату из body
	var creds db.UserJson

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		server.Error(map[string]interface{}{"message": "Data is not provided", "status": 400}, w)
		return
	}

	//проверяем совпадают ли данные с данными из бд
	err = db.ValidateUser(creds)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error(), "status": 400}, w)
		return
	}

	//генерируем новый токен, через который можно будет получить логин
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": creds.Login,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		server.Error(map[string]interface{}{"message": "Error generating token", "status": 400}, w)
		return
	}

	//возвращаем токен
	server.Ok(map[string]interface{}{"message": "Signed-in successful", "token": tokenString, "status": 200}, w)
}
