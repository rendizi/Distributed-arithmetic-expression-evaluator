package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1"
	dbname   = "daee"
)

var db *sql.DB

func init() {
	//creating db
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to PostgreSQL!")

	//Table Users
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login TEXT NOT NULL UNIQUE, 
			password TEXT NOT NULL
);
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Table 'users' created successfully!")

	//Table Expressions
	createTableQuery = `
        CREATE TABLE IF NOT EXISTS expressions (
    		id SERIAL PRIMARY KEY,
    		expression TEXT NOT NULL,
    		settings TEXT NOT NULL,
    		result TEXT,
    		createdAt TEXT NOT NULL,
    		endTime TEXT,
    		user_id INTEGER NOT NULL,
    		FOREIGN KEY (user_id) REFERENCES users(id)
);

    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Table 'urls' created successfully!")

	//Table Operations
	createTableQuery = `
        CREATE TABLE IF NOT EXISTS operations (
			id SERIAL PRIMARY KEY,
			operation TEXT NOT NULL,
			result TEXT,
    		expression_id INTEGER NOT NULL,
    		FOREIGN KEY (expression_id) REFERENCES expressions(id)
);
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Table 'operations' created successfully!")
}
