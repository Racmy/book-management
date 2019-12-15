package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/docker_go_nginx/app/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type RegistValue struct {
	Title        string
	Author       string
	Latest_Issue float64
}

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
	var tmpl = template.Must(template.ParseFiles("./template/bookRegist.html"))

	tmpTitle := r.FormValue("Title")
	tmpAuthor := r.FormValue("Author")
	tmpLatest_Issue_String := r.FormValue("Latest_Issue")
	tmpLatest_Issue, strConvErr := strconv.ParseFloat(tmpLatest_Issue_String, 64)

	if strConvErr != nil {
		tmpLatest_Issue = 1
	}

	tmp := RegistValue{
		Title:        tmpTitle,
		Author:       tmpAuthor,
		Latest_Issue: tmpLatest_Issue,
	}

	if err := tmpl.ExecuteTemplate(w, "bookRegist.html", tmp); err != nil {
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
	tmpLatest_Issue, strConvErr := strconv.ParseFloat(tmpLatest_Issue_String, 64)

	if (tmpTitle == "") || (tmpAuthor == "") || (strConvErr != nil) {
		var url = "/regist"
		url += "?Title=" + r.Form["Title"][0] + "&Author=" + r.Form["Author"][0] + "&Latest_Issue=" + r.Form["Latest_Issue"][0]
		http.Redirect(w, r, url, http.StatusFound)
	}

	insertBook := db.Book{Title: tmpTitle, Author: tmpAuthor, Latest_Issue: tmpLatest_Issue}

	dbErr := db.InsertBook(insertBook)
	if dbErr != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	// テンプレートに埋め込むデータ作成
	dat := struct {
		Title        string
		Author       string
		Latest_Issue float64
	}{
		Title:        tmpTitle,
		Author:       tmpAuthor,
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
	r.HandleFunc("/regist", bookRegistHandler)
	r.HandleFunc("/regist/success", bookInsertHandler)
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/search", bookSearchHandler)
	r.HandleFunc("/detail", bookDetailHandler)
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules/"))))
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
