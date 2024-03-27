package db

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

	// Execute the insert query and get the inserted ID
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

func UpdateResult(opid, id int64, errr string) error {
	query := "UPDATE operations SET result = $1 WHERE id = $2"

	_, err := db.Exec(query, errr, opid)
	if err != nil {
		return err
	}

	query = "UPDATE expressions SET result = $1 WHERE id = $2"

	_, err = db.Exec(query, errr, id)
	if err != nil {
		return err
	}

	return nil
}
