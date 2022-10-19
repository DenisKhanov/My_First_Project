package service

import (
	"WWWgo/internal/db"
	"WWWgo/pkg/secure"
	"WWWgo/pkg/structs"
	"fmt"
	"html/template"
	"net/http"
)

var Token string

//Авторизация пользователя
func Check(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("pkg/templates/autorisation.html", "pkg/templates/header.html", "pkg/templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmp.ExecuteTemplate(w, "autorisation", nil)
}
func Verification(w http.ResponseWriter, r *http.Request) {
	tmp, errr := template.ParseFiles("pkg/templates/error.html", "pkg/templates/ok.html", "pkg/templates/header.html", "pkg/templates/footer.html")
	if errr != nil {
		fmt.Fprintf(w, errr.Error())
	}
	login := r.FormValue("login")
	password := secure.HashPassword(r.FormValue("password"))

	if login == "" || password == "" {
		status := "Поля логин или пароль не могут быть пустыми"
		page := "/autorisation"
		tmp.ExecuteTemplate(w, "error", struct{ Status, Page string }{Status: status, Page: page})
	} else {
		//Вытягивание строки из БД по логину и проверка с введенным паролем
		dataBase := db.DbConnect()
		var user structs.Users
		dataBase.Find(&user, "login", login)
		if user.Login == "" {
			status := fmt.Sprintf("Пользователь с таким email %s не зарегистрирован", login)
			page := "/autorisation"
			tmp.ExecuteTemplate(w, "error", struct{ Status, Page string }{Status: status, Page: page})
		} else {
			if user.Password == password {
				status := fmt.Sprintf("%s, мы вас узнали!", login)
				Token, _ = secure.GenerateJWT()
				tmp.ExecuteTemplate(w, "ok", struct{ Status string }{Status: status})
			} else {
				status := "Сочетание логина и пароля не верны!"
				page := "/autorisation"
				tmp.ExecuteTemplate(w, "error", struct{ Status, Page string }{Status: status, Page: page})
			}
		}
	}
	secure.Token = Token
}
