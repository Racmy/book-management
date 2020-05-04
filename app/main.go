package main

import (
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/handler/bookHandler"
	"github.com/docker_go_nginx/app/handler/loginHandler"
	"github.com/docker_go_nginx/app/handler/userHandler"
	"github.com/docker_go_nginx/app/utility/ulogin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
)

var rootTemplatePath = "./template/"
var homeTemplatePath = rootTemplatePath + "home/"
var homeHTMLName = "index.html"

// ホーム画面用の画面データ構造
type HomeResponseData struct {
	ViewData map[string]string
	Message  map[string][]string
}

/*
	ホーム画面を表示するハンドラ
*/
func homeHandler(w http.ResponseWriter, r *http.Request) {
	Tpl, _ := template.ParseGlob("./template/parts/*")
	Tpl.New(homeHTMLName).ParseFiles(homeTemplatePath + homeHTMLName)
	session, _ := ulogin.GetSession(r)

	// 画面表示データ構造作成
	responseData := HomeResponseData{
		ViewData: map[string]string{},
		Message:  nil,
	}

	if message := session.Flashes(appconst.SessionMsg); len(message) > 0 {
		castedMessage := message[0].(map[string][]string)
		viewData := session.Flashes(appconst.SessionViewData)[0].(map[string]string)
		session.Save(r, w)
		// 画面表示データ構造作成
		responseData = HomeResponseData{
			ViewData: viewData,
			Message:  castedMessage,
		}
	}

	if err := Tpl.ExecuteTemplate(w, homeHTMLName, responseData); err != nil {
		log.Fatal(err)
	}
}

// ルーティング
func main() {
	bookhandler.Tpl, _ = template.ParseGlob("./template/parts/*")
	r := mux.NewRouter()
	// ホーム画面のハンドラ
	r.HandleFunc(appconst.RootURL, homeHandler)
	// ユーザ登録のハンドラ
	r.HandleFunc(appconst.UserRegistURL, userHandler.UserRegistHandler)
	// ユーザ登録情報の更新ハンドラ
	r.HandleFunc(appconst.UserEditURL, userHandler.UserEditHandler)
	// ユーザパスワード再発行申込ハンドラ
	r.HandleFunc(appconst.UserPassWordOrderURL, userHandler.UserPassWordOrderHandler)
	// ユーザパスワード再登録ハンドラ
	r.HandleFunc(appconst.UserPassWordRegistURL, userHandler.UserPasswordRegist)
	// 本一覧画面表示ハンドラ
	r.HandleFunc(appconst.BookURL, bookhandler.BookListHandler)
	// 本登録画面表示ハンドラ
	r.HandleFunc(appconst.BookRegistURL, bookhandler.BookRegistHandler)
	// 本登録処理ハンドラ
	r.HandleFunc(appconst.BookRegistProcessURL, bookhandler.BookInsertHandler)
	// 本登録結果画面表示ハンドラ
	r.HandleFunc(appconst.BookRegistResultURL, bookhandler.BookInsertResultHandler)
	// 本検索画面表示ハンドラ
	r.HandleFunc(appconst.BookSearchURL, bookhandler.BookSearchHandler)
	// 本詳細画面表示ハンドラ
	r.HandleFunc(appconst.BookDetailLURL, bookhandler.BookDetailHandler)
	// 本更新処理ハンドラ
	r.HandleFunc(appconst.BookUpdatehURL, bookhandler.BookUpdateHandler)
	// 本削除処理ハンドラ
	r.HandleFunc(appconst.BookDeleteURL, bookhandler.BookDeleteHandler)
	//r.HandleFunc(appconst.LoginURL,loginHandler.LoginHandler).Methods("GET")
	r.HandleFunc(appconst.LoginURL, loginHandler.LoginHandler).Methods("POST")
	// cssフレームワーク読み込み
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules/"))))
	// 画像フォルダ
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.Handle(appconst.RootURL, r)
	http.ListenAndServe(":3000", nil)
}
