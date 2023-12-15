package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"gofr.dev/pkg/gofr"
)

var db *sql.DB

type Student struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Enrollment string `json:"enrollment"`
}

func createDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "students.db")
	if err != nil {
		fmt.Println("error opening the database")
	}
	createTableQuery := `CREATE TABLE students(id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR(255), age INTEGER, enrollment VARCHAR(255));`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		fmt.Println("Error creating database")
	}
}
func main() {
	app := gofr.New()
	createDatabase()
	app.GET("/", func(ctx *gofr.Context) (interface{}, error) {

		return "Hello This is the Student Management API", nil
	})
	app.POST("/add", func(ctx *gofr.Context) (interface{}, error) {
		var students Student
		err := AddStudent(student)
		if err != nil {
			return nil, err
		}
		return students, nil
	})
	app.Start()
}
