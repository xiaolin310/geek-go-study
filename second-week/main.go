package main

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"os"
)

type Student struct {
	Id   int
	Name string
	Class string
}

// controller
func getStudentInfoById(id int) (*Student, error) {
	student, err := queryStudentById(id)
	return student, errors.WithMessage(err, "could not get student info by id")


}

// dao: wrap errors
func queryStudentById(id int) (*Student, error) {
	db, err := sql.Open("mysql", "xl:123456@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	student := &Student{}
	err = db.QueryRow("select id, name, class from students where id = ? ", id).Scan(&student.Id, &student.Name, &student.Class)
	if err != nil {
		msg := fmt.Sprintf("queryStudentById:%d 查询学生信息异常", id)
		return student, errors.Wrap(err, msg)
	}
	return student, nil
}



func main() {
	student, err := getStudentInfoById(666)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("query student info err: %+v\n", err)
			os.Exit(1)
		}
		fmt.Printf("打印其他错误信息: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("学生信息：%v\n", student)
}
