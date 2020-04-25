package main

import (
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/handler/bookHandler"
	"github.com/docker_go_nginx/app/handler/loginHandler"
	"text/template"
	"net/http"
	"github.com/gorilla/mux"
)

var rootTemplatePath = "./template/"
var homeTemplatePath = rootTemplatePath + "home/"
var homeHTMLName = "index.html"

/*
	ホーム画面を表示するハンドラ
*/
func homeHandler(w http.ResponseWriter, r *http.Request) {
	Tpl, _ := template.ParseGlob("./template/parts/*")
	Tpl.New(homeHTMLName).ParseFiles(homeTemplatePath + homeHTMLName)
	if err := Tpl.ExecuteTemplate(w, homeHTMLName, nil); err != nil {
		log.Fatal(err)
	}
}

// ルーティング
func main() {
	bookhandler.Tpl, _ = template.ParseGlob("./template/parts/*")
	r := mux.NewRouter()
	r.HandleFunc(appconst.RootURL, homeHandler)
	r.HandleFunc(appconst.BookURL, bookhandler.BookListHandler)
	r.HandleFunc(appconst.BookRegistURL, bookhandler.BookRegistHandler)
	r.HandleFunc(appconst.BookRegistProcessURL, bookhandler.BookInsertHandler)
	r.HandleFunc(appconst.BookRegistResultURL, bookhandler.BookInsertResultHandler)
	r.HandleFunc(appconst.BookSearchURL, bookhandler.BookSearchHandler)
	r.HandleFunc(appconst.BookDetailLURL, bookhandler.BookDetailHandler)
	r.HandleFunc(appconst.BookUpdatehURL, bookhandler.BookUpdateHandler)
	r.HandleFunc(appconst.BookDeleteURL, bookhandler.BookDeleteHandler)
	r.HandleFunc(appconst.LoginURL,loginHandler.LoginHandler).Methods("GET")
	r.HandleFunc(appconst.LoginURL,loginHandler.LoginCheckHandler).Methods("POST")
	// cssフレームワーク読み込み
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules/"))))
	// 画像フォルダ
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.Handle(appconst.RootURL, r)
	http.ListenAndServe(":3000", nil)
}