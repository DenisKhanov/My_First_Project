package main

import (
	"WWWgo/db"
	"WWWgo/router"
)

func main() {
	db.InitialMigration()
	router.HandleFuncs()
}

//id := r.URL.Query().Get("id") Парсинг URL при передачи данных в конце URL после знака ?
