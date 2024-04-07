package db

import (
	"database/sql"
	"errors"
)

type UserJson struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

//CREATE TABLE IF NOT EXISTS users (
//	id SERIAL PRIMARY KEY,
//	login TEXT NOT NULL UNIQUE,
//	password TEXT NOT NULL
//);

// Добавляем юзера в бд
func InsertUser(user UserJson) error {
	if user.Password == "" || user.Login == "" {
		return errors.New("user data can't be an empty string")
	}
	insertQuery := `INSERT INTO users (login,password) VALUES ($1, $2)`
	_, err := db.Exec(insertQuery, user.Login, user.Password)
	if err != nil {
		return err
	}
	return nil
}

// Данную функцию используем во время логина, проверяя переданные данные с данными в бд
func ValidateUser(user UserJson) error {
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE login = $1", user.Login).Scan(&storedPassword)
	if err != nil {
		//Если не имеются записи- нет юзера с переданным логином
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		//Иначе просто вернем ошибку бд
		return err
	}

	//Если пароль, сохраненный в бд и переданный юзером не совадают, то возвращаем ошибку
	if storedPassword != user.Password {
		return errors.New("invalid password")
	}

	//Все сходится, можно пропускать
	return nil
}
