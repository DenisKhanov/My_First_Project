package methods

import (
	"WWWgo/db"
	"WWWgo/structs"
	"fmt"
	"html/template"
	"net/http"
)

//Авторизация пользователя
func Check(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/autorisation.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmp.ExecuteTemplate(w, "autorisation", nil)
}
func Verification(w http.ResponseWriter, r *http.Request) {
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
		dataBase := db.DbConnect()
		var user structs.Users
		dataBase.Find(&user, "login", login)
		if user.Login == "" {
			status := fmt.Sprintf("Пользователь с таким email %s не зарегистрирован", login)
			page := "/autorisation"
			tmp.ExecuteTemplate(w, "error", struct{ Status, Page string }{Status: status, Page: page})
		} else {
			if user.Password == password {
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
