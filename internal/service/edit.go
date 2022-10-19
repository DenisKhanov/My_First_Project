package service

import (
	"WWWgo/internal/db"
	"WWWgo/pkg/structs"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

//Редактирование статьи
func Edit(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("pkg/templates/edit.html",
		"pkg/templates/header.html", "pkg/templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	vars := mux.Vars(r)
	id := vars["id"]

	var post structs.Article
	dataBase := db.DbConnect()
	dataBase.Find(&post, id)
	tmp.ExecuteTemplate(w, "edit", post)

}
func EditPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	//
	dataBase := db.DbConnect()
	dataBase.Model(&structs.Article{}).Where("id", id).Updates(&structs.Article{Title: title,
		Anons: anons, FullText: full_text})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
