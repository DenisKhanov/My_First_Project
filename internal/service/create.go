package service

import (
	"WWWgo/internal/db"
	"WWWgo/pkg/structs"
	"fmt"
	"html/template"
	"net/http"
)

//Создание статьи
func Create(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("pkg/templates/create.html",
		"pkg/templates/header.html", "pkg/templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmp.ExecuteTemplate(w, "create", Token)
}
func SaveArticle(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("pkg/templates/blank.html",
		"pkg/templates/header.html", "pkg/templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		tmp.ExecuteTemplate(w, "blank", nil)
	} else {
		//Добавление данных
		dataBase := db.DbConnect()
		dataBase.Create(&structs.Article{Title: title, Anons: anons, FullText: full_text})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
