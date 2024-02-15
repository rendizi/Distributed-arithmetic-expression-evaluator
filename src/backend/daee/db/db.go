package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1"
	dbname   = "golang_microservice"
)

var db *sql.DB

// На запуске создаем таблицу если не существует и подключаемся к бд
func init() {
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

	fmt.Println("Successfully connected to PostgreSQL!")

	createTableQuery := `
        CREATE TABLE IF NOT EXISTS tasks (
            id SERIAL PRIMARY KEY,
            userId TEXT NOT NULL,
            task TEXT NOT NULL,
            answer TEXT DEFAULT 'no',
            time TEXT DEFAULT 'not finished yet',
            settings TEXT NOT NULL
        );
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table 'tasks' created successfully!")
}

// функция, добавляющая в таблицу новое задание и на столбец ответ ставит значение нет
func Insert(user, task, settings string) error {
	insertQuery := `INSERT INTO tasks (userId, task,answer,settings) VALUES ($1, $2,$3,$4);`
	_, err := db.Exec(insertQuery, user, task, "no", settings)
	if err != nil {
		return err
	}
	return nil
}

// функция обновляющая таблицу, добавляет ответ и время выполнения
func Update(user, task, answer, time string) error {
	if len(task) > 0 && task[len(task)-1] == '\n' {
		task = task[:len(task)-1]
	}
	updateQuery := `UPDATE tasks SET answer = $1, time = $2 WHERE userId = $3 AND task = $4;`
	_, err := db.Exec(updateQuery, answer, time, user, task)
	if err != nil {
		return err
	}
	return nil
}

//Функция ниже по дефолту просто берет все записи с бд, но если подается userId, который не является
// ns , берет только записи , принадлежащие к некоторому юзеру, иначе если userId равен ns он выведет
// все записи у которых нет ответа.

// Получается что при подаче в функцию пустую строку вернутся все записи, при подаче ns операции(ведь
// у них нет ответа) , при подаче юзерайди всю историю запросов юзера
func Get(userId string) (map[int][]string, error) {
	query := "SELECT id, userId, task,time, answer FROM tasks"
	if len(userId) != 0 && userId != "ns" {
		query = fmt.Sprintf("SELECT userId , task, answer,time FROM tasks WHERE userId = '%s'", userId)
	} else if userId == "ns" {
		query = "SELECT id, userId , task,settings FROM tasks WHERE answer='no'"
	}
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[int][]string)
	id := 0
	for rows.Next() {
		var user, task, answer, time string
		err := rows.Scan(&user, &task, &answer, &time)
		if err != nil {
			return nil, err
		}
		result[id] = []string{user, task, answer, time}
		id++
	}
	return result, nil
}

// Данная функция проверяет на импатентность, тоесть если юзер уже отправлял такое задание
// и оно решено - просто вернуть ответ
func Solved(userId, task string) (string, string, error) {
	query := "SELECT answer FROM tasks WHERE userId = $1 AND task = $2 AND answer != $3"
	fmt.Println(query)
	rows, err := db.Query(query, userId, task, "no")
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	for rows.Next() {
		var answer string
		err = rows.Scan(&answer)
		if err != nil {
			return "", "", err
		}
		return answer, time.Microsecond.String(), nil
	}
	return "", "", fmt.Errorf("no record found")
}
