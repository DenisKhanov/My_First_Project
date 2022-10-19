package db

import (
	structs2 "WWWgo/pkg/structs"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DbConnect() *gorm.DB {
	db, err := gorm.Open(mysql.Open("denisk:Qwerty123@tcp(81.90.182.182:3320)/orm?parseTime=true"))
	if err != nil {
		fmt.Println("Could not connect to Data Base")
		panic(err)
	}
	return db
}

func InitialMigration() {
	dataBase := DbConnect()

	dataBase.AutoMigrate(&structs2.Users{}, &structs2.Article{})

}
