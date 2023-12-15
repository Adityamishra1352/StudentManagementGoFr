package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

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
func addStudent(student Student) error {
	_, err := db.Exec("INSERT INTO students(name,age,enrollment)VALUES(?,?,?)", student.Name, student.Age, student.Enrollment)
	return err
}
func viewStudents() ([]Student, error) {
	rows, err := db.Query("SELECT * FROM students")
	if err != nil {
		fmt.Println("Error while executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var a Student
		err := rows.Scan(&a.ID, &a.Name, &a.Age, &a.Enrollment)
		if err != nil {
			fmt.Println("Error while scanning row:", err)
			return nil, err
		}
		students = append(students, a)
	}

	return students, nil
}
func deleteStudent(id int) error {
	_, err := db.Exec("DELETE FROM students WHERE id=?", id)
	return err
}
func updateStudent(id int, student Student) error {
	_, err := db.Exec("UPDATE students SET name=?, age=?, enrollment=? WHERE id=?", student.Name, student.Age, student.Enrollment, id)
	return err
}
func main() {
	app := gofr.New()
	createDatabase()
	app.GET("/", func(ctx *gofr.Context) (interface{}, error) {

		return "Hello This is the Student Management API", nil
	})
	app.POST("/add", func(ctx *gofr.Context) (interface{}, error) {
		var students Student
		if err := json.NewDecoder(ctx.Request().Body).Decode(&students); err != nil {
			return nil, err
		}
		err := addStudent(students)
		if err != nil {
			return nil, err
		}
		return students, nil
	})
	app.GET("/view", func(ctx *gofr.Context) (interface{}, error) {
		students, err := viewStudents()
		if err != nil {
			fmt.Println("COuldnt view")
		}
		return students, nil
	})
	app.GET("/d/:id", func(ctx *gofr.Context) (interface{}, error) {
		idParam := ctx.Param("id")
		if idParam == "" {
			return nil, fmt.Errorf("ID not provided")
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			return nil, fmt.Errorf("invalidformat")
		}

		deletedStudent, err := viewStudents()
		if err != nil {
			return nil, err
		}

		err = deleteStudent(id)
		if err != nil {
			fmt.Println("Couldn't delete student:", err)
			return nil, err
		}

		return deletedStudent, nil
	})
	app.PUT("/update/:id", func(ctx *gofr.Context) (interface{}, error) {
		idParam := ctx.Param("id")
		if idParam == "" {
			return nil, fmt.Errorf("idnotprovided")
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			return nil, fmt.Errorf("invalidid")
		}

		var updatedStudent Student
		if err := json.NewDecoder(ctx.Request().Body).Decode(&updatedStudent); err != nil {
			return nil, err
		}

		err = updateStudent(id, updatedStudent)
		if err != nil {
			fmt.Println("Couldn't update student:", err)
			return nil, err
		}

		return updatedStudent, nil
	})
	app.Start()
}
