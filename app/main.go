package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/handler/bookHandler"
	"text/template"
	"net/http"
	"github.com/gorilla/mux"
)

// ルーティング
func main() {
	bookhandler.Tpl, _ = template.ParseGlob("./template/parts/*")
	r := mux.NewRouter()
	r.HandleFunc(appconst.BookURL, bookhandler.BookListHandler)
	r.HandleFunc(appconst.BookRegistURL, bookhandler.BookRegistHandler)
	r.HandleFunc(appconst.BookRegistProcessURL, bookhandler.BookInsertHandler)
	r.HandleFunc(appconst.BookRegistResultURL, bookhandler.BookInsertResultHandler)
	r.HandleFunc(appconst.BookSearchURL, bookhandler.BookSearchHandler)
	r.HandleFunc(appconst.BookDetailLURL, bookhandler.BookDetailHandler)
	r.HandleFunc(appconst.BookUpdatehURL, bookhandler.BookUpdateHandler)
	r.HandleFunc(appconst.BookDeleteURL, bookhandler.BookDeleteHandler)
	// cssフレームワーク読み込み
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules/"))))
	// 画像フォルダ
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.Handle(appconst.RootURL, r)
	http.ListenAndServe(":3000", nil)
}