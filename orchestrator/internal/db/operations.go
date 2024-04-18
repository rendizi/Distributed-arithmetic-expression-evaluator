package db

import "database/sql"

type OperationJSON struct {
	Id        int    `json:"id"`
	Operation string `json:"operation"`
	Result    string `json:"result"`
}

//CREATE TABLE IF NOT EXISTS operations (
//			id SERIAL PRIMARY KEY,
//			operation TEXT NOT NULL,
//			result TEXT,
//    		expression_id INTEGER NOT NULL,
//    		FOREIGN KEY (expression_id) REFERENCES expressions(id)
//);

func InsertOperation(expression Expression, id int64) (int64, error) {
	query := `INSERT INTO operations (operation, expression_id) VALUES ($1, $2) RETURNING id`
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var opid int64
	err = stmt.QueryRow(expression.Expression, id).Scan(&opid)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateOperationState(operation string, opid int64) error {
	query := "UPDATE operations SET operation = $1 WHERE id = $2"

	_, err := db.Exec(query, operation, opid)
	if err != nil {
		return err
	}

	return nil
}

func DeleteOperation(id int64) error {
	stmt, err := db.Prepare("DELETE FROM operations WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func GetOperations() ([]int, []string, []string, error) {
	query := "SELECT id,operation,result FROM operations"
	rows, err := db.Query(query)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	operations := []string{}
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
		operations = append(operations, expression)
		ids = append(ids, id)
		if result.Valid {
			results = append(results, result.String)
		} else {
			results = append(results, "In progress...")
		}
	}

	if err = rows.Err(); err != nil {
		return nil, nil, nil, err
	}

	return ids, operations, results, nil
}
