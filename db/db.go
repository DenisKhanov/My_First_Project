package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func DbConnect() *sql.DB {
	db, err := sql.Open("mysql", "denisk:02Denis1990@tcp(81.90.182.182:3306)/mydb")
	if err != nil {
		fmt.Println("Error connect to data base!!!")
		panic(err)
	}
	fmt.Println("Connected to data base!")
	return db
}
