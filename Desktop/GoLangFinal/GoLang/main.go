package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Article{}

func home_page(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("html/home_page.html", "html/header.html", "html/footer.html")

	if err != nil {
		fmt.Print(w, err.Error())
	}
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT `name`, `age` FROM `users`")
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
		fmt.Println(fmt.Sprintf("Post: %s with id: %d", post.Title, post.Id))
	}

	t.ExecuteTemplate(w, "home_page", nil)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("html/create.html", "html/header.html", "html/footer.html")

	if err != nil {
		fmt.Print(err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "You Did Not Enter Any Data Yet!")
	} else {
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles`(`title`,`anons`, `full_text` ) VALUES('%s','%s','%s')", title, anons, full_text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprint(w, "ID: %v\n", vars["id"])
}

func HandleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", home_page).Methods("GET")
	rtr.HandleFunc("/create/", create).Methods("GET")
	rtr.HandleFunc("/save_article/", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")
	http.Handle("/", rtr)
	rtr.ListenAndServe(":8080", nil)
}

func main() {
	HandleFunc()

}
