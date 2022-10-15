package methods

import (
	"WWWgo/db"
	"WWWgo/structs"
	"fmt"
	"html/template"
	"net/http"
)

//Главная страница
func Index(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/index.html",
		"templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	//Выборка данных
	dataBase := db.DbConnect()
	var posts []structs.Article
	dataBase.Find(&posts)
	tmp.ExecuteTemplate(w, "index", posts)
	//json.NewEncoder(w).Encode(posts)
}
