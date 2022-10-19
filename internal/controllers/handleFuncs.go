package controllers

import (
	"WWWgo/internal/service"
	"WWWgo/pkg/secure"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func HandleFuncs() {
	port := os.Getenv("PORT")
	rout := mux.NewRouter().StrictSlash(true)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	rout.HandleFunc("/", service.Index).Methods("GET")

	rout.HandleFunc("/autorisation", service.Check).Methods("GET")
	rout.HandleFunc("/check_login", service.Verification).Methods("POST")

	rout.HandleFunc("/login", service.Register_new_user).Methods("GET")
	rout.HandleFunc("/add_user", service.Add_user).Methods("POST")

	rout.HandleFunc("/create", secure.IsAuthorised(service.Create)).Methods("GET")
	rout.HandleFunc("/save_article", service.SaveArticle).Methods("POST")

	rout.HandleFunc("/post/{id:[0-9]+}", service.ViewPosts).Methods("GET")

	rout.HandleFunc("/edit/post/{id:[0-9]+}", service.Edit).Methods("GET", "POST")
	rout.HandleFunc("/edit_post/post/{id:[0-9]+}", service.EditPost).Methods("GET", "POST")

	rout.HandleFunc("/delete/post/{id:[0-9]+}", service.Delete).Methods("GET")

	http.Handle("/", rout)

	http.ListenAndServe(":"+port, nil)
	//http.ListenAndServe(":8080", nil)
}
