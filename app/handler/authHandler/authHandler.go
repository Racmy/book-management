package authHandler

import (
	"encoding/gob"
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/common/message"
	"github.com/docker_go_nginx/app/db/userdao"
	"github.com/docker_go_nginx/app/utility/ulogin"
	"log"
	"net/http"
	"text/template"
)

var Tpl *template.Template

var rootTemplatePath = "./template/"
var loginTemplatePath = rootTemplatePath + "login/"
var loginHTMLName = "testLogin.html"
var loginResultHTMLName = "testLoginResult.html"

/**
ログインチェックのハンドラ
*/
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// sessionにmap形式のデータを追加できるように設定
	gob.Register(map[string][]string{})
	gob.Register(map[string]string{})
	gob.Register(userdao.User{})
	// テンプレート前処理
	Tpl, _ := template.ParseGlob("./template/parts/*")
	Tpl.New(loginHTMLName).ParseFiles(loginTemplatePath + loginResultHTMLName)

	// クライアントからのデータ取得
	email := r.FormValue(appconst.EMAIL)
	password := r.FormValue(appconst.PASSWORD)
	log.Println("email/pass:" + email + "/" + password)

	// エラーメッセージ格納用の変数作成
	mailMsg := []string{}
	passWordMsg := []string{}

	// メールアドレスとパスワードの空文字チェック
	if email == "" {
		mailMsg = append(mailMsg, message.ErrMsgNoEmail)
	}
	if password == "" {
		passWordMsg = append(passWordMsg, message.ErrMsgNoPassword)
	}

	// セッションにつめる、エラーメッセージと画面データ格納用の変数作成
	errMsgMap := map[string][]string{}
	viewData := map[string]string{}

	// メールアドレスとパスワードにおいて、入力チェックで不正と判断された場合はリダイレクト
	if len(mailMsg) > 0 || len(passWordMsg) > 0 {
		// エラーメッセージを格納
		errMsgMap["mail"] = mailMsg
		errMsgMap["password"] = passWordMsg
		// 画面データを格納
		viewData["mail"] = email

		// セッションにエラーメッセージと画面データをつめる
		session, _ := ulogin.GetSession(r)
		session.AddFlash(errMsgMap, appconst.SessionMsg)
		session.AddFlash(viewData, appconst.SessionViewData)
		session.Save(r, w)

		// ログイン画面へ遷移
		http.Redirect(w, r, appconst.RootURL, http.StatusFound)
	}

	resUser, errMsg := ulogin.LoginCheck(email, password)
	// 一致するユーザIDとパスワードが存在するかチェック
	if errMsg != "" {
		// エラーメッセージを作成
		sokanMsg := []string{}
		sokanMsg = append(sokanMsg, "ログインIDとパスワードに誤りがあります。")
		errMsgMap["sokan"] = sokanMsg
		viewData["mail"] = email

		// セッションにエラーメッセージと画面データをつめる
		session, _ := ulogin.GetSession(r)
		session.AddFlash(errMsgMap, appconst.SessionMsg)
		session.AddFlash(viewData, appconst.SessionViewData)
		session.Save(r, w)

		// ホーム画面へリダイレクト
		http.Redirect(w, r, appconst.RootURL, http.StatusFound)
	}

	// ログインセッションに登録
	session, _ := ulogin.GetSession(r)
	// セッションにデータにデータをつめる
	session.Values[appconst.SessionLoginUser] = resUser
	session.Save(r, w)
	http.Redirect(w, r, appconst.BookURL, http.StatusFound)

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// sessionにmap形式のデータを追加できるように設定
	ulogin.DelSession(w, r)
	http.Redirect(w, r, appconst.RootURL, http.StatusFound)

}
