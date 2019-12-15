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

/*
	レスポンスデータ
*/
type ResponseData struct {
	Keyword string
	Books   []db.Book
}

/*
	本を登録画面へのハンドラ
*/
func bookRegistHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("./template/base.html"))

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

/*
	ホーム画面へのハンドラ
*/
func homeHandler(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("./template/list.html"))
	var responseData ResponseData
	responseData.Books = db.GetAllBooks()
	responseData.Keyword = ""
	if err := tpl.ExecuteTemplate(w, "list.html", responseData); err != nil {
		log.Fatal(err)
	}
}

/*
	本詳細画面へのハンドラ
*/
func bookDetailHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

/*
	本の検索のためのハンドラ
*/
func bookSearchHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	if keyword := query.Get("keyword"); query.Get("keyword") != "" {
		// keywordがnullの場合は、HOMEへリダイレクト
		if keyword == "" {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		var tpl = template.Must(template.ParseFiles("./template/list.html"))

		// ResponseDataの作成
		var responseData ResponseData
		responseData.Keyword = keyword
		responseData.Books = db.GetSearchedBooks(keyword)

		if err := tpl.ExecuteTemplate(w, "list.html", responseData); err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

/*
	ルーティング
*/
func main() {

	r := mux.NewRouter()
	r.HandleFunc("/book-regist", bookRegistHandler)
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/search", bookSearchHandler)
	r.HandleFunc("/detail", bookDetailHandler)
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules/"))))
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
