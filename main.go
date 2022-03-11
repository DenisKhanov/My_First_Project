package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
)

type Article struct {
	Id                     int
	Title, Anons, FullText string
}

var posts = []Article{}
var showPost = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/myproject")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//Выборка данных
	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}
	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	tmp.ExecuteTemplate(w, "index", posts)

}
func create(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmp.ExecuteTemplate(w, "create", nil)
}
func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не должно быть пустых строк!")
	} else {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/myproject")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		//Добавление данных
		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`,`anons`,`full_text`) VALUES('%s','%s','%s')", title, anons, full_text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
func wiewPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	tmp, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/myproject")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}

	showPost = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		showPost = post
	}
	tmp.ExecuteTemplate(w, "show", showPost)
}
func handleFunc() {
	port := os.Getenv("PORT")
	rout := mux.NewRouter()
	rout.HandleFunc("/", index)
	rout.HandleFunc("/create", create)
	rout.HandleFunc("/save_article", save_article)
	rout.HandleFunc("/post/{id:[0-9]+}", wiewPost)

	http.Handle("/", rout)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":"+port, nil)
}
func main() {
	handleFunc()
}
