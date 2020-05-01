package main

import (
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/common/appstructure"
	"github.com/docker_go_nginx/app/utility/ulogin"
	"github.com/docker_go_nginx/app/handler/bookHandler"
	"github.com/docker_go_nginx/app/handler/loginHandler"
	"text/template"
	"net/http"
)

var rootTemplatePath = "./template/"
var homeTemplatePath = rootTemplatePath + "home/"
var homeHTMLName = "index.html"
var userTemplatePath = rootTemplatePath + "user/"
var userRegistHTMLName = "regist.html"
var userEditHTMLName = "edit.html"
var userPasswordOrderHTMLName = "password_order.html"
var userPasswordRegistHTMLName = "password_regist.html"


/*
	ホーム画面を表示するハンドラ
*/
func homeHandler(w http.ResponseWriter, r *http.Request) {
	Tpl, _ := template.ParseGlob("./template/parts/*")
	Tpl.New(homeHTMLName).ParseFiles(homeTemplatePath + homeHTMLName)
	var errMsg appstructure.HomeErrorMessage
	session, _ := ulogin.GetSession(r)

	if errFlg := session.Values[appconst.SessionErrFlg]; errFlg!= nil && errFlg.(bool) == true{
		if flashErrMsg := session.Flashes(appconst.SessionErrMsgEmail); len(flashErrMsg) > 0{
			errMsg.EmailErr = flashErrMsg[0].(string)
		}
		if flashErrMsg := session.Flashes(appconst.SessionErrMsgPassword); len(flashErrMsg) > 0{
			errMsg.PasswordErr = flashErrMsg[0].(string)
		}
		if flashErrMsg := session.Flashes(appconst.SessionErrMsgNoUser); len(flashErrMsg) > 0{
			errMsg.NoUserErr = flashErrMsg[0].(string)
		}
	}
	session.Save(r,w)
	
	if err := Tpl.ExecuteTemplate(w, homeHTMLName, errMsg); err != nil {
		log.Fatal(err)
	}
}

/*
	ユーザを新規登録するハンドラ
*/
func userRegistHandler(w http.ResponseWriter, r *http.Request) {
	Tpl, _ := template.ParseGlob("./template/parts/*")
	Tpl.New(userRegistHTMLName).ParseFiles(userTemplatePath + userRegistHTMLName)
	if err := Tpl.ExecuteTemplate(w, userRegistHTMLName, nil); err != nil {
		log.Fatal(err)
	}
}

/*
	ユーザ情報を更新するハンドラ
*/
func userEditHandler(w http.ResponseWriter, r *http.Request) {
	Tpl, _ := template.ParseGlob("./template/parts/*")
	Tpl.New(userEditHTMLName).ParseFiles(userTemplatePath + userEditHTMLName)
	if err := Tpl.ExecuteTemplate(w, userEditHTMLName, nil); err != nil {
		log.Fatal(err)
	}
}

/*
	ユーザのログインパスワード再発行
*/
func userPassWordOrderHandler(w http.ResponseWriter, r *http.Request) {
	Tpl, _ := template.ParseGlob("./template/parts/*")
	Tpl.New(userPasswordOrderHTMLName).ParseFiles(userTemplatePath + userPasswordOrderHTMLName)
	if err := Tpl.ExecuteTemplate(w, userPasswordOrderHTMLName, nil); err != nil {
		log.Fatal(err)
	}
}

/*
	ユーザのパスワード再登録画面
*/
func userPasswordRegist(w http.ResponseWriter, r *http.Request) {
	Tpl, _ := template.ParseGlob("./template/parts/*")
	log.Print("hoge")
	Tpl.New(userPasswordRegistHTMLName).ParseFiles(userTemplatePath + userPasswordRegistHTMLName)
	if err := Tpl.ExecuteTemplate(w, userPasswordRegistHTMLName, nil); err != nil {
		log.Fatal(err)
	}
}


func baseHandlerFunc(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
    return baseHandler(http.HandlerFunc(handler))
}

func baseHandler(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// common
		_,err := ulogin.SessionCheck(w, r)
		if err != nil {
			//エラーの時
			log.Println("session err")
			http.Redirect(w,r,appconst.RootURL,http.StatusFound)
		}else{
			log.Println(r.URL, r.Method)
			handler.ServeHTTP(w, r)
		}
        
    })
}


// ルーティング
func main() {
	bookhandler.Tpl, _ = template.ParseGlob("./template/parts/*")

	// ホーム画面のハンドラ
	http.Handle(appconst.RootURL, baseHandlerFunc(homeHandler))
	// ユーザ登録のハンドラ
	http.Handle(appconst.UserRegistURL, baseHandlerFunc(userRegistHandler))
	// ユーザ登録情報の更新ハンドラ
	http.Handle(appconst.UserEditURL, baseHandlerFunc(userEditHandler))
	// ユーザパスワード再発行申込ハンドラ
	http.Handle(appconst.UserPassWordOrderURL, baseHandlerFunc(userPassWordOrderHandler))
	// ユーザパスワード再登録ハンドラ
	http.Handle(appconst.UserPassWordRegistURL, baseHandlerFunc(userPasswordRegist))
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
	http.Handle(appconst.LoginURL,baseHandlerFunc(loginHandler.LoginHandler))
	// cssフレームワーク読み込み
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules/"))))
	// 画像フォルダ
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.ListenAndServe(":3000", nil)
}