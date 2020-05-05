package main

import (
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/common/message"
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

///baseHandlerFunc...ハンドラに共通処理を付与する関数
/*
@param handler func(w http.ResponseWriter, r *http.Request) 固有処理関数
@param sessionFlg int セッションのログイン情報を用いる場合:1　それ以外:0
@return http.Handler 共通処理を付与したハンドラ
*/
func baseHandlerFunc(handler func(w http.ResponseWriter, r *http.Request), sessionFlg int) http.Handler {
	return baseHandler(http.HandlerFunc(handler), sessionFlg)
}

///baseHandler...ハンドラの共通処理
/*
@param handler http.Handler 固有処理関数
@param sessionFlg int セッションのログイン情報を用いる場合:1　それ以外:0
@return http.Handler 共通処理を付与したハンドラ
*/
func baseHandler(handler http.Handler, sessionFlg int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL, r.Method)
		if sessionFlg == 1 {
			if ulogin.IsLogined(r) {
				//ログイン済み
				handler.ServeHTTP(w, r)
			} else {
				//未ログインorセッション切れ
				log.Println("session err")
				errMsgMap := map[string][]string{}
				viewData := map[string]string{}
				noSessionMsg := []string{}
				noSessionMsg = append(noSessionMsg, message.ErrMsgNoSession)
				errMsgMap["nosession"] = noSessionMsg

				// セッションにエラーメッセージと画面データをつめる
				session, _ := ulogin.GetSession(r)
				session.AddFlash(errMsgMap, appconst.SessionMsg)

				viewData["mail"] = ""
				session.AddFlash(viewData, appconst.SessionViewData)

				session.Save(r, w)
				http.Redirect(w, r, appconst.RootURL, http.StatusFound)
			}
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

// ルーティング
func main() {
	bookhandler.Tpl, _ = template.ParseGlob("./template/parts/*")

	// ホーム画面のハンドラ
	http.Handle(appconst.RootURL, baseHandlerFunc(homeHandler, 0))
	// ユーザ登録のハンドラ
	http.Handle(appconst.UserRegistURL, baseHandlerFunc(userHandler.UserRegistHandler, 0))
	// ユーザ登録情報の更新ハンドラ
	http.Handle(appconst.UserEditURL, baseHandlerFunc(userHandler.UserEditHandler, 1))
	// ユーザパスワード再発行申込ハンドラ
	http.Handle(appconst.UserPassWordOrderURL, baseHandlerFunc(userHandler.UserPassWordOrderHandler, 0))
	// ユーザパスワード再登録ハンドラ
	http.Handle(appconst.UserPassWordRegistURL, baseHandlerFunc(userHandler.UserPasswordRegist, 0))
	// 本一覧画面表示ハンドラ
	http.Handle(appconst.BookURL, baseHandlerFunc(bookhandler.BookListHandler, 1))
	// 本登録画面表示ハンドラ
	http.Handle(appconst.BookRegistURL, baseHandlerFunc(bookhandler.BookRegistHandler, 1))
	// 本登録処理ハンドラ
	http.Handle(appconst.BookRegistProcessURL, baseHandlerFunc(bookhandler.BookInsertHandler, 1))
	// 本登録結果画面表示ハンドラ
	http.Handle(appconst.BookRegistResultURL, baseHandlerFunc(bookhandler.BookInsertResultHandler, 1))
	// 本検索画面表示ハンドラ
	http.Handle(appconst.BookSearchURL, baseHandlerFunc(bookhandler.BookSearchHandler, 1))
	// 本詳細画面表示ハンドラ
	http.Handle(appconst.BookDetailLURL, baseHandlerFunc(bookhandler.BookDetailHandler, 1))
	// 本更新処理ハンドラ
	http.Handle(appconst.BookUpdatehURL, baseHandlerFunc(bookhandler.BookUpdateHandler, 1))
	// 本削除処理ハンドラ
	http.Handle(appconst.BookDeleteURL, baseHandlerFunc(bookhandler.BookDeleteHandler, 1))
	// ログイン処理ハンドラ
	http.Handle(appconst.LoginURL, baseHandlerFunc(loginHandler.LoginHandler, 0))
	// cssフレームワーク読み込み
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules/"))))
	// 画像フォルダ
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.ListenAndServe(":3000", nil)
}
