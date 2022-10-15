package methods

import (
	"WWWgo/db"
	"WWWgo/structs"
	"fmt"
	"html/template"
	"net/http"
)

//Регистрация пользователя
func Register_new_user(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/login.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmp.ExecuteTemplate(w, "login", nil)
}
func Add_user(w http.ResponseWriter, r *http.Request) {
	tmp, errr := template.ParseFiles("templates/error.html", "templates/ok.html", "templates/header.html", "templates/footer.html")
	if errr != nil {
		fmt.Fprintf(w, errr.Error())
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	dataBase := db.DbConnect()
	var user structs.Users
	dataBase.Find(&user, "login", login)

	if len(login) < 4 || len(password) < 4 {
		status := "Логин и пароль не могут быть менее 4 символов"
		page := "/login"
		tmp.ExecuteTemplate(w, "error", struct{ Status, Page string }{Status: status, Page: page})
	} else if user.Login == "" {
		//Добавление данных
		dataBase.Create(&structs.Users{Login: login, Password: password})
		status := fmt.Sprintf("Пользователь %s успешно зарегистрирован", login)
		tmp.ExecuteTemplate(w, "ok", struct{ Status string }{Status: status})
	} else {
		status := fmt.Sprintf("Пользователь %s уже зарегистрирован", login)
		page := "/login"
		tmp.ExecuteTemplate(w, "error", struct{ Status, Page string }{Status: status, Page: page})
	}
}
