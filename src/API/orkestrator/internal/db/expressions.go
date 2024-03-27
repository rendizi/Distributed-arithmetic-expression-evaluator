package db

import (
	"encoding/json"
	"time"
)

type Settings struct {
	Plus  time.Duration `json:"plus"`
	Minus time.Duration `json:"minus"`
	Mult  time.Duration `json:"mul"`
	Div   time.Duration `json:"div"`
}

type Expression struct {
	Expression string   `json:"expression"`
	Settings   Settings `json:"settings"`
}

//CREATE TABLE IF NOT EXISTS expressions (
////    id INTEGER PRIMARY KEY,
////    expression TEXT NOT NULL,
////    result TEXT NOT NULL,
////    user_id INTEGER NOT NULL,
////    FOREIGN KEY (user_id) REFERENCES users(id)
////);

func InsertExpression(expression Expression, author string) (int64, error) {
	settingsJSON, err := json.Marshal(expression.Settings)
	if err != nil {
		return 0, err
	}
	query := `INSERT INTO expressions (expression, settings, user_id) VALUES ($1, $2, $3) RETURNING id`
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(expression.Expression, string(settingsJSON), author).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
