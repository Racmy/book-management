package main

import (
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/common/processTemplate"
	"github.com/docker_go_nginx/app/handler/bookHandler"
	"github.com/docker_go_nginx/app/handler/authHandler"
	"github.com/docker_go_nginx/app/handler/userHandler"
	"github.com/docker_go_nginx/app/utility/ulogin"
	_ "github.com/go-sql-driver/mysql"
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

	// ホーム画面のハンドラ
	http.Handle(appconst.RootURL, processTemplate.BaseHandlerFunc(homeHandler, 0))
	// ユーザ登録のハンドラ
	http.Handle(appconst.UserRegistURL, processTemplate.BaseHandlerFunc(userHandler.UserRegistHandler, 0))
	// ユーザ登録情報の更新ハンドラ
	http.Handle(appconst.UserEditURL, processTemplate.BaseHandlerFunc(userHandler.UserEditHandler, 1))
	// ユーザパスワード再発行申込ハンドラ
	http.Handle(appconst.UserPassWordOrderURL, processTemplate.BaseHandlerFunc(userHandler.UserPassWordOrderHandler, 0))
	// ユーザパスワード再登録ハンドラ
	http.Handle(appconst.UserPassWordRegistURL, processTemplate.BaseHandlerFunc(userHandler.UserPasswordRegist, 0))
	// 本一覧画面表示ハンドラ
	http.Handle(appconst.BookURL, processTemplate.BaseHandlerFunc(bookhandler.BookListHandler, 1))
	// 本登録画面表示ハンドラ
	http.Handle(appconst.BookRegistURL, processTemplate.BaseHandlerFunc(bookhandler.BookRegistHandler, 1))
	// 本登録処理ハンドラ
	http.Handle(appconst.BookRegistProcessURL, processTemplate.BaseHandlerFunc(bookhandler.BookInsertHandler, 1))
	// 本登録結果画面表示ハンドラ
	http.Handle(appconst.BookRegistResultURL, processTemplate.BaseHandlerFunc(bookhandler.BookInsertResultHandler, 1))
	// 本検索画面表示ハンドラ
	http.Handle(appconst.BookSearchURL, processTemplate.BaseHandlerFunc(bookhandler.BookSearchHandler, 1))
	// 本詳細画面表示ハンドラ
	http.Handle(appconst.BookDetailLURL, processTemplate.BaseHandlerFunc(bookhandler.BookDetailHandler, 1))
	// 本更新処理ハンドラ
	http.Handle(appconst.BookUpdatehURL, processTemplate.BaseHandlerFunc(bookhandler.BookUpdateHandler, 1))
	// 本削除処理ハンドラ
	http.Handle(appconst.BookDeleteURL, processTemplate.BaseHandlerFunc(bookhandler.BookDeleteHandler, 1))
	// ログイン処理ハンドラ
	http.Handle(appconst.LoginURL, processTemplate.BaseHandlerFunc(authHandler.LoginHandler, 0))
	// ログアウト処理ハンドラ
	http.Handle(appconst.LogoutURL, processTemplate.BaseHandlerFunc(authHandler.LogoutHandler, 0))
	// cssフレームワーク読み込み
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules/"))))
	// 画像フォルダ
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.ListenAndServe(":3000", nil)
}
