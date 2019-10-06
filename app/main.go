package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var tmpl = template.Must(template.ParseFiles("./template/base.html"))
var tmpl2 = template.Must(template.ParseFiles("./template/list.html"))

type Book struct {
	Id                     int
	Title                  string
	Author                 string
	Latest_Issue           float32
	Front_Cover_Image_Path string
}

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
	db, err := sql.Open("mysql", "racmy:racmy@tcp(db:3306)/book-management")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close() // 関数がリターンする直前に呼び出される

	rows, err := db.Query("SELECT * FROM book")
	if err != nil {
		panic(err.Error())
	}

	// Bookを格納するArray作成
	var books = []Book{}
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.Id, &book.Title, &book.Author, &book.Latest_Issue, &book.Front_Cover_Image_Path)
		log.Print(book.Title)
		if err != nil {
			panic(err)
		}
		log.Print(book.Title)
		books = append(books, book)
	}

	if err := tmpl2.ExecuteTemplate(w, "list.html", books); err != nil {
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
