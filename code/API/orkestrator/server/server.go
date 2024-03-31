package server

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

var jwtSecret = []byte("secret_key")

// Json со статусом 400
func Error(data map[string]interface{}, w http.ResponseWriter) {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(jsonResponse)
	return
}

// Json со статусом ok
func Ok(data map[string]interface{}, w http.ResponseWriter) {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	return
}

// Берем логин из Auth token
func GetLogin(w http.ResponseWriter, r *http.Request) string {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		Error(map[string]interface{}{"message": "auth token is not provided", "status": 400}, w)
		return ""
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		Error(map[string]interface{}{"message": "Token is not valid", "status": 400}, w)
		return ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		Error(map[string]interface{}{"message": "Error extracting token claims", "status": 400}, w)
		return ""
	}
	login, ok := claims["login"].(string)
	if !ok {
		Error(map[string]interface{}{"message": "Error extracting email from token", "status": 400}, w)
		return ""
	}
	return login
}
