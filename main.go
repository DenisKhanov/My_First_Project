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

type Users struct {
	Id              int
	Login, Password string
}

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

//Создание статьи
func create(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmp.ExecuteTemplate(w, "create", nil)
}
func save_article(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/blank.html", "templates/header.html", "templates/footer.html")
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

//Регистрация пользователя
func register_new_user(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/login.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmp.ExecuteTemplate(w, "login", nil)
}
func add_user(w http.ResponseWriter, r *http.Request) {
	tmp, errr := template.ParseFiles("templates/error.html", "templates/ok.html", "templates/header.html", "templates/footer.html")
	if errr != nil {
		fmt.Fprintf(w, errr.Error())
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	res := db.DbConnect().QueryRow(fmt.Sprintf("SELECT * FROM `autentification` WHERE `login`='%s'", login))
	inf := Users{}
	err := res.Scan(&inf.Id, &inf.Login, &inf.Password)
	defer db.DbConnect().Close()
	if len(login) < 4 || len(password) < 4 {
		status := "Логин и пароль не могут быть менее 4 символов"
		page := "/login"
		tmp.ExecuteTemplate(w, "error", struct{ Status, Page string }{Status: status, Page: page})
	} else if err != nil {
		//Добавление данных
		insert, erro := db.DbConnect().Query(fmt.Sprintf("INSERT INTO `autentification` (`login`,`password`) VALUES('%s','%s')", login, password))
		if erro != nil {
			panic(err)
		}
		defer insert.Close()
		status := fmt.Sprintf("Пользователь %s успешно зарегистрирован", login)
		tmp.ExecuteTemplate(w, "ok", struct{ Status string }{Status: status})
	} else {
		status := fmt.Sprintf("Пользователь %s уже зарегистрирован", login)
		page := "/login"
		tmp.ExecuteTemplate(w, "error", struct{ Status, Page string }{Status: status, Page: page})
	}
}

//Авторизация пользователя
func check(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/autorisation.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmp.ExecuteTemplate(w, "autorisation", nil)
}
func verification(w http.ResponseWriter, r *http.Request) {
	tmp, errr := template.ParseFiles("templates/error.html", "templates/ok.html", "templates/header.html", "templates/footer.html")
	if errr != nil {
		fmt.Fprintf(w, errr.Error())
	}
	login := r.FormValue("login")
	password := r.FormValue("password")

	if login == "" || password == "" {
		status := "Поля логин или пароль не могут быть пустыми"
		page := "/autorisation"
		tmp.ExecuteTemplate(w, "error", struct{ Status, Page string }{Status: status, Page: page})
	} else {
		//Вытягивание строки из БД по логину и проверка с введенным паролем
		res := db.DbConnect().QueryRow(fmt.Sprintf("SELECT * FROM `autentification` WHERE `login`='%s'", login))
		inf := Users{}
		err := res.Scan(&inf.Id, &inf.Login, &inf.Password)
		if err != nil {
			status := fmt.Sprintf("Пользователь %s не зарегистрирован", login)
			page := "/autorisation"
			tmp.ExecuteTemplate(w, "error", struct{ Status, Page string }{Status: status, Page: page})
		} else {
			if inf.Password == password {
				fmt.Println("Complete")
				status := fmt.Sprintf("%s, мы вас узнали!", login)
				tmp.ExecuteTemplate(w, "ok", struct{ Status string }{Status: status})
			} else {
				status := "Сочетание логина и пароля не верны!"
				page := "/autorisation"
				tmp.ExecuteTemplate(w, "error", struct{ Status, Page string }{Status: status, Page: page})
			}
		}
	}
}

func handleFuncs() {
	port := os.Getenv("PORT")
	rout := mux.NewRouter()

	rout.HandleFunc("/", index).Methods("GET")

	rout.HandleFunc("/autorisation", check).Methods("GET")
	rout.HandleFunc("/check_login", verification).Methods("POST")

	rout.HandleFunc("/login", register_new_user).Methods("GET")
	rout.HandleFunc("/add_user", add_user).Methods("POST")

	rout.HandleFunc("/create", create).Methods("GET")
	rout.HandleFunc("/save_article", save_article).Methods("POST")
	rout.HandleFunc("/post/{id:[0-9]+}", wiewPost).Methods("GET")

	http.Handle("/", rout)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":"+port, nil)
}
func main() {
	handleFuncs()
}
