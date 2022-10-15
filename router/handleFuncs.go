package router

import (
	"WWWgo/methods"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func HandleFuncs() {
	port := os.Getenv("PORT")
	rout := mux.NewRouter().StrictSlash(true)

	rout.HandleFunc("/", methods.Index).Methods("GET")

	rout.HandleFunc("/autorisation", methods.Check).Methods("GET")
	rout.HandleFunc("/check_login", methods.Verification).Methods("POST")

	rout.HandleFunc("/login", methods.Register_new_user).Methods("GET")
	rout.HandleFunc("/add_user", methods.Add_user).Methods("POST")

	rout.HandleFunc("/create", methods.Create).Methods("GET")
	rout.HandleFunc("/save_article", methods.SaveArticle).Methods("POST")

	rout.HandleFunc("/post/{id:[0-9]+}", methods.ViewPosts).Methods("GET")

	rout.HandleFunc("/edit/post/{id:[0-9]+}", methods.Edit).Methods("GET", "POST")
	rout.HandleFunc("/edit_post/post/{id:[0-9]+}", methods.EditPost).Methods("GET", "POST")

	rout.HandleFunc("/delete/post/{id:[0-9]+}", methods.Delete).Methods("GET")

	http.Handle("/", rout)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":"+port, nil)
	//http.ListenAndServe(":8080", nil)
}
