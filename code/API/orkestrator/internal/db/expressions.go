package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

type Settings struct {
	Plus  int `json:"plus"`
	Minus int `json:"minus"`
	Mult  int `json:"mult"`
	Div   int `json:"div"`
}

type Expression struct {
	Expression string   `json:"expression"`
	Settings   Settings `json:"settings"`
}

type ExpressionJSON struct {
	Id         int    `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
	Time       string `json:"time"`
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
	query := `INSERT INTO expressions (expression, settings,createdAt, user_id) VALUES ($1, $2, $3,$4) RETURNING id`
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	currentTime := time.Now()

	// Define the layout for date and time only
	layout := "2006-01-02 15:04:05"

	// Format the current time using the layout
	formattedTime := currentTime.Format(layout)

	// Execute the insert query and get the inserted ID
	var id int64
	err = stmt.QueryRow(expression.Expression, string(settingsJSON), formattedTime, userID).Scan(&id)
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

func GetExpressions(authorLogin string) ([]int, []string, []string, error) {
	var userID int64
	err := db.QueryRow("SELECT id FROM users WHERE login = $1", authorLogin).Scan(&userID)
	if err != nil {
		return nil, nil, nil, err
	}
	query := "SELECT id,expression,result FROM expressions WHERE user_id = $1"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	// Store the expressions
	expressions := []string{}
	results := []string{}
	ids := []int{}
	for rows.Next() {
		var expression string
		var id int
		var result sql.NullString
		err = rows.Scan(&id, &expression, &result)
		if err != nil {
			return nil, nil, nil, err
		}
		expressions = append(expressions, expression)
		ids = append(ids, id)
		if result.Valid {
			results = append(results, result.String)
		} else {
			results = append(results, "In progress...")
		}
	}

	// Check for errors during row iteration
	if err = rows.Err(); err != nil {
		return nil, nil, nil, err
	}

	return ids, expressions, results, nil
}

func UpdateResult(opid, id int64, errr string) error {
	currentTime := time.Now()
	query := "UPDATE expressions SET result = $1 , endTime = $2 WHERE id = $3"

	// Define the layout for date and time only
	layout := "2006-01-02 15:04:05"

	// Format the current time using the layout
	formattedTime := currentTime.Format(layout)

	_, err := db.Exec(query, errr, formattedTime, id)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func Time(id int64) (string, error) {
	var start string
	var end sql.NullString
	err := db.QueryRow(`SELECT createdAt, endTime FROM expressions WHERE id = $1`, id).Scan(&start, &end)
	if err != nil {
		return "", err
	}

	createdAt, err := time.Parse("2006-01-02 15:04:05", start)
	if err != nil {
		return "", err
	}

	if !end.Valid {
		return "In progress...", nil
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", end.String)
	if err != nil {
		return "", err
	}
	log.Println(start, end.String)
	duration := endTime.Sub(createdAt)
	return duration.String(), nil
}
