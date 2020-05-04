package main

import (
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/handler/bookHandler"
	"github.com/docker_go_nginx/app/handler/loginHandler"
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

func baseHandlerFunc(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return baseHandler(http.HandlerFunc(handler))
}

func baseHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL, r.Method)
		handler.ServeHTTP(w, r)

	})
}

// ルーティング
func main() {
	bookhandler.Tpl, _ = template.ParseGlob("./template/parts/*")

	// ホーム画面のハンドラ
	http.Handle(appconst.RootURL, baseHandlerFunc(homeHandler))
	// ユーザ登録のハンドラ
	http.Handle(appconst.UserRegistURL, baseHandlerFunc(userHandler.UserRegistHandler))
	// ユーザ登録情報の更新ハンドラ
	http.Handle(appconst.UserEditURL, baseHandlerFunc(userHandler.UserEditHandler))
	// ユーザパスワード再発行申込ハンドラ
	http.Handle(appconst.UserPassWordOrderURL, baseHandlerFunc(userHandler.UserPassWordOrderHandler))
	// ユーザパスワード再登録ハンドラ
	http.Handle(appconst.UserPassWordRegistURL, baseHandlerFunc(userHandler.UserPasswordRegist))
	// 本一覧画面表示ハンドラ
	http.Handle(appconst.BookURL, baseHandlerFunc(bookhandler.BookListHandler))
	// 本登録画面表示ハンドラ
	http.Handle(appconst.BookRegistURL, baseHandlerFunc(bookhandler.BookRegistHandler))
	// 本登録処理ハンドラ
	http.Handle(appconst.BookRegistProcessURL, baseHandlerFunc(bookhandler.BookInsertHandler))
	// 本登録結果画面表示ハンドラ
	http.Handle(appconst.BookRegistResultURL, baseHandlerFunc(bookhandler.BookInsertResultHandler))
	// 本検索画面表示ハンドラ
	http.Handle(appconst.BookSearchURL, baseHandlerFunc(bookhandler.BookSearchHandler))
	// 本詳細画面表示ハンドラ
	http.Handle(appconst.BookDetailLURL, baseHandlerFunc(bookhandler.BookDetailHandler))
	// 本更新処理ハンドラ
	http.Handle(appconst.BookUpdatehURL, baseHandlerFunc(bookhandler.BookUpdateHandler))
	// 本削除処理ハンドラ
	http.Handle(appconst.BookDeleteURL, baseHandlerFunc(bookhandler.BookDeleteHandler))
	// ログイン処理ハンドラ
	http.Handle(appconst.LoginURL, baseHandlerFunc(loginHandler.LoginHandler))
	// cssフレームワーク読み込み
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules/"))))
	// 画像フォルダ
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.ListenAndServe(":3000", nil)
}
