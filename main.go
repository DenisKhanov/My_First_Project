package main

import (
	"WWWgo/db"
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

//type Users struct {
//	Id              int
//	Login, Password string
//}

var posts = []Article{}
var showPost = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	//Выборка данных
	defer db.DbConnect().Close()
	res, err := db.DbConnect().Query("SELECT * FROM `articles`")
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
		//Добавление данных
		defer db.DbConnect().Close()
		insert, err := db.DbConnect().Query(fmt.Sprintf("INSERT INTO `articles` (`title`,`anons`,`full_text`) VALUES('%s','%s','%s')", title, anons, full_text))
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
	defer db.DbConnect().Close()
	res, err := db.DbConnect().Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s'", vars["id"]))
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
func reduct_story(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Yes, it's working!")
}

func register_new_user(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/login.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmp.ExecuteTemplate(w, "login", nil)
}
func add_user(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	if login == "" || password == "" {
		fmt.Fprintf(w, "Логин и пароль не могут быть пустыми")
	} else {
		//Добавление данных
		defer db.DbConnect().Close()
		insert, err := db.DbConnect().Query(fmt.Sprintf("INSERT INTO `autentification` (`login`,`password`) VALUES('%s','%s')", login, password))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/ok", http.StatusSeeOther)
	}
}

func status_registration(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/ok.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmp.ExecuteTemplate(w, "ok", nil)
}

func handleFuncs() {
	port := os.Getenv("PORT")
	rout := mux.NewRouter()
	rout.HandleFunc("/", index).Methods("GET")

	rout.HandleFunc("/login", register_new_user).Methods("GET")
	rout.HandleFunc("/add_user", add_user).Methods("POST")
	rout.HandleFunc("/ok", status_registration).Methods("GET")
	rout.HandleFunc("/create", create).Methods("GET")
	rout.HandleFunc("/save_article", save_article).Methods("POST")
	rout.HandleFunc("/post/{id:[0-9]+}", wiewPost).Methods("GET")

	http.Handle("/", rout)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":"+port, nil)
	//http.ListenAndServe(":8080", nil)
}
func main() {
	handleFuncs()
}

//
//package main
//
//import (
//	"WWWgo/db"
//	"bufio"
//	"fmt"
//	_ "github.com/go-sql-driver/mysql"
//	"os"
//)
//
//type articles struct {
//	id                      int
//	title, anons, full_text string
//}
//
//func reduct_title() {
//	//Редактирование данных в таблице
//	defer db.DbConnect().Close()
//	var numId int
//	fmt.Println("Enter id and his new title...")
//	fmt.Scan(&numId)
//	bio := bufio.NewReader(os.Stdin)
//	newTitle, _, _ := bio.ReadLine()
//	res, err := db.DbConnect().Exec(fmt.Sprintf("UPDATE articles SET `title`='%s' WHERE `id`= '%d'", newTitle, numId))
//	if err != nil {
//		panic(err)
//	}
//	res.LastInsertId()
//}
//
//func wiew_db_info() {
//	//Выводим данные всей базы данных
//	defer db.DbConnect().Close()
//	res, err := db.DbConnect().Query("SELECT * FROM `articles`")
//	if err != nil {
//		panic(err)
//	}
//	defer res.Close()
//	informations := []articles{}
//
//	for res.Next() {
//		inf := articles{}
//		err := res.Scan(&inf.id, &inf.title, &inf.anons, &inf.full_text)
//		if err != nil {
//			fmt.Println(err, "Attention,it's error!")
//			continue
//		}
//		informations = append(informations, inf)
//	}
//	for _, inf := range informations {
//		fmt.Printf(" Id        %d\n Title     %s\n Anons     %s\n Full_text %s\n\n",
//			inf.id, inf.title, inf.anons, inf.full_text)
//	}
//}
//
//func main() {
//	reduct_title()
//	wiew_db_info()
//}

//Выводим только одну строку из базы данных
//	res := db.QueryRow("SELECT * FROM articles WHERE id=2")
//	inf := articles{}
//	err = res.Scan(&inf.id, &inf.title, &inf.anons, &inf.full_text)
//	if err != nil {
//		panic(err)
//	} else {
//		fmt.Printf(" Id        %d\n Title     %s\n Anons     %s\n Full_text %s\n\n",
//			inf.id, inf.title, inf.anons, inf.full_text)
//	}
//}
