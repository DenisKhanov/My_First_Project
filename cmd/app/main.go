package main

import (
	"WWWgo/internal/controllers"
	"WWWgo/internal/db"
	"WWWgo/internal/service"
	"fmt"
)

func main() {
	db.InitialMigration()
	controllers.HandleFuncs()
	fmt.Println(service.Token)
}

//id := r.URL.Query().Get("id") Парсинг URL при передачи данных в конце URL после знака ?
