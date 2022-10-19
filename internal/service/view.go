package service

import (
	"WWWgo/internal/db"
	"WWWgo/pkg/structs"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

//Развернутый просмотр статьи
func ViewPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.WriteHeader(http.StatusOK)

	tmp, err := template.ParseFiles("pkg/templates/show.html",
		"pkg/templates/header.html", "pkg/templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	dataBase := db.DbConnect()
	var post structs.Article
	dataBase.Find(&post, id)
	tmp.ExecuteTemplate(w, "show", post)
}
