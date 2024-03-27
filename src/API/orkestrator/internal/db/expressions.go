package db

import (
	"database/sql"
	"encoding/json"
)

type Settings struct {
	Plus  int `json:"plus"`
	Minus int `json:"minus"`
	Mult  int `json:"mul"`
	Div   int `json:"div"`
}

type Expression struct {
	Expression string   `json:"expression"`
	Settings   Settings `json:"settings"`
}

type ExpressionJSON struct {
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

//CREATE TABLE IF NOT EXISTS expressions (
////    id  PRIMARY KEY,
////    expression TEXT NOT NULL,
////    result TEXT NOT NULL,
////    user_id INTEGER NOT NULL,
////    FOREIGN KEY (user_id) REFERENCES users(id)
////);

func InsertExpression(expression Expression, authorLogin string) (int64, error) {
	// Retrieve user ID based on login
	var userID int64
	err := db.QueryRow("SELECT id FROM users WHERE login = $1", authorLogin).Scan(&userID)
	if err != nil {
		return 0, err
	}

	// Marshal settings to JSON
	settingsJSON, err := json.Marshal(expression.Settings)
	if err != nil {
		return 0, err
	}

	// Prepare the insert query
	query := `INSERT INTO expressions (expression, settings, user_id) VALUES ($1, $2, $3) RETURNING id`
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Execute the insert query and get the inserted ID
	var id int64
	err = stmt.QueryRow(expression.Expression, string(settingsJSON), userID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetExpression(id int, authorLogin string) (string, string, error) {
	var userID int64
	err := db.QueryRow("SELECT id FROM users WHERE login = $1", authorLogin).Scan(&userID)
	if err != nil {
		return "", "", err
	}

	var expression string
	var result sql.NullString
	err = db.QueryRow(`SELECT expression,result FROM expressions WHERE id = $1 and user_id = $2`, id, userID).Scan(&expression, &result)
	if err != nil {
		return "", "", err
	}
	if result.Valid {
		// result is not null
		return expression, result.String, nil
	} else {
		// result is null
		return expression, "In progress...", nil
	}
}

func GetExpressions(authorLogin string) ([]string, []string, error) {
	var userID int64
	err := db.QueryRow("SELECT id FROM users WHERE login = $1", authorLogin).Scan(&userID)
	if err != nil {
		return nil, nil, err
	}
	query := "SELECT expression,result FROM expressions WHERE user_id = $1"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	// Store the expressions
	expressions := []string{}
	results := []string{}
	for rows.Next() {
		var expression string
		var result sql.NullString
		err := rows.Scan(&expression, &result)
		if err != nil {
			return nil, nil, err
		}
		expressions = append(expressions, expression)
		if result.Valid {
			results = append(results, result.String)
		} else {
			results = append(results, "In progress...")
		}
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return expressions, results, nil
}
