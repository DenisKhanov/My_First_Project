package methods

import (
	"WWWgo/db"
	"WWWgo/structs"
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

	tmp, err := template.ParseFiles("templates/show.html",
		"templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	dataBase := db.DbConnect()
	var post structs.Article
	dataBase.Find(&post, id)
	tmp.ExecuteTemplate(w, "show", post)
}
