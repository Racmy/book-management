package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/docker_go_nginx/app/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var tmpl = template.Must(template.ParseFiles("./template/base.html"))

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// テンプレートに埋め込むデータ作成
	dat := struct {
		Title string
		Time  time.Time
	}{
		Title: "Test",
		Time:  time.Now(),
	}
	// テンプレートにデータを埋め込む
	if err := tmpl.ExecuteTemplate(w, "base.html", dat); err != nil {
		log.Fatal(err)
	}

}

func dbHandler(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("./template/list.html"))
	books := db.GetAllBooks()
	if err := tpl.ExecuteTemplate(w, "list.html", books); err != nil {
		log.Fatal(err)
	}
}
func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/db", dbHandler)
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
