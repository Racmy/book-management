package main

import (
	"html/template"
	"log"
	"net/http"
	"github.com/docker_go_nginx/app/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"strconv"
)

/*
	本を登録画面へのハンドラ
*/
func bookRegistHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("./template/bookRegist.html"))

	if err := tmpl.Execute(w, "bookRegist.html"); err != nil {
		log.Fatal(err)
	}

}
/*
	本を登録画面へのハンドラ
*/
func bookInsertHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("./template/bookRegistResult.html"))
	r.ParseForm()
	tmpTitle := r.Form["Title"][0]
	tmpAuthor := r.Form["Author"][0]
	tmpLatest_Issue_String := r.Form["Latest_Issue"][0]
	tmpLatest_Issue , strConvErr := strconv.ParseFloat(tmpLatest_Issue_String,64)

	if strConvErr !=nil{
		tmpLatest_Issue = 1;
	}

	insertBook := db.Book{Title: tmpTitle,Author: tmpAuthor,Latest_Issue: tmpLatest_Issue}
	
	dbErr := db.InsertBook(insertBook)
	if dbErr != nil{
		http.Redirect(w,r,"/",http.StatusFound)
	}
	// テンプレートに埋め込むデータ作成
	dat := struct {
		Title string
		Author string
		Latest_Issue float64
	}{
		Title: tmpTitle,
		Author: tmpAuthor,
		Latest_Issue: tmpLatest_Issue,
	}

	// テンプレートにデータを埋め込む
	if err := tmpl.ExecuteTemplate(w, "bookRegistResult.html", dat); err != nil {
		log.Fatal(err)
	}

}
/*
	ホーム画面へのハンドラ
*/
func homeHandler(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("./template/list.html"))
	books := db.GetAllBooks()
	if err := tpl.ExecuteTemplate(w, "list.html", books); err != nil {
		log.Fatal(err)
	}
}

/*
	本詳細画面へのハンドラ
*/
func bookDetailHandler(w http.ResponseWriter, r *http.Request) {}

/*
	ルーティング
*/
func main() {

	r := mux.NewRouter()
	r.HandleFunc("/regist", bookRegistHandler)
	r.HandleFunc("/regist/success", bookInsertHandler)
	r.HandleFunc("/", homeHandler)
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules/"))))
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
