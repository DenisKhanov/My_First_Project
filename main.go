package main

import (
	"WWWgo/db"
	//"encoding/json"
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

//Эксперименты с Json
//type UserRequest struct {
//	Name  string `json:"name"`
//	Email string `json:"email"`
//	Age   int    `json:"age"`
//	Sex   string `json:"sex"`
//}
//
//var user1 = Users{2, "Denis", "123123124"}
//var userRec = UserRequest{"Ivan", "fafaasad@mail.ru", 32, "male"}
//
//type Test_struct struct {
//	Test   string `json:"test"`
//	Number int    `json:"number"`
//}
//
//var testt = Test_struct{}
//
//func outputJson(w http.ResponseWriter, r *http.Request) {
//	w.Header().Add("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(userRec)
//	json.NewEncoder(w).Encode(user1)
//}
//func inputJson(w http.ResponseWriter, r *http.Request) {
//fmt.Fprintf(w, r.Body)
//json.NewDecoder(r.Body).Decode(&testt)
//w.Header().Add("Content-Type", "application/json")
//json.NewEncoder(w).Encode(testt)
//jsonString := `{"test":"asdasds","number":123}`
//err := json.Unmarshal([]byte(jsonString), &testt)
//if err != nil {
//	fmt.Println(err)
//}

//}

//------------------------------------------------------------------------------------------------

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
		tmp.ExecuteTemplate(w, "error", struct{ Status string }{Status: status})
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
		tmp.ExecuteTemplate(w, "error", struct{ Status string }{Status: status})
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
		fmt.Fprintf(w, "Это все не то!")
		status := "Поля логин или пароль не могут быть пустыми"
		tmp.ExecuteTemplate(w, "error", struct{ Status string }{Status: status})
	} else {
		//Вытягивание строки из БД по логину и проверка с введенным паролем
		res := db.DbConnect().QueryRow(fmt.Sprintf("SELECT * FROM `autentification` WHERE `login`='%s'", login))
		inf := Users{}
		err := res.Scan(&inf.Id, &inf.Login, &inf.Password)
		if err != nil {
			status := fmt.Sprintf("Пользователь %s не зарегистрирован", login)
			tmp.ExecuteTemplate(w, "error", struct{ Status string }{Status: status})
		} else {
			if inf.Password == password {
				fmt.Println("Complete")
				status := fmt.Sprintf("%s, мы вас узнали!", login)
				tmp.ExecuteTemplate(w, "ok", struct{ Status string }{Status: status})
			} else {
				status := "Сочетание логина и пароля не верны!"
				tmp.ExecuteTemplate(w, "error", struct{ Status string }{Status: status})
			}
		}
	}
}

func handleFuncs() {
	port := os.Getenv("PORT")
	rout := mux.NewRouter()

	//rout.HandleFunc("/output", outputJson)
	//rout.HandleFunc("/input", inputJson)

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
	//http.ListenAndServe(":8080", nil)
}
func main() {
	handleFuncs()
}

//Статусы выполнения
//func status_registration(w http.ResponseWriter, r *http.Request) {
//	tmp, err := template.ParseFiles("templates/ok.html", "templates/header.html", "templates/footer.html")
//	if err != nil {
//		fmt.Fprintf(w, err.Error())
//	}
//	tmp.ExecuteTemplate(w, "ok", nil)
//}
//func error_autorisatin(w http.ResponseWriter, r *http.Request) {
//	tmp, err := template.ParseFiles("templates/error.html", "templates/header.html", "templates/footer.html")
//	if err != nil {
//		fmt.Fprintf(w, err.Error())
//	}
//	tmp.ExecuteTemplate(w, "error", nil)
//}
//rout.HandleFunc("/ok", status_registration).Methods("GET")
//rout.HandleFunc("/error", error_autorisatin).Methods("GET")

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
