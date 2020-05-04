package loginHandler

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

// /**
// 	ログイン画面へのハンドラ
// */
// func LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	Tpl, _ := template.ParseGlob("./template/parts/*")
// 	Tpl.New(loginHTMLName).ParseFiles(loginTemplatePath + loginHTMLName)
// 	if err := Tpl.ExecuteTemplate(w, loginHTMLName, nil); err != nil {
// 		log.Fatal(err)
// 	}
// }

/**
ログインチェックのハンドラ
*/
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// sessionにmap形式のデータを追加できるように設定
	gob.Register(map[string][]string{})
	gob.Register(map[string]string{})
	gob.Register(userdao.User{})

	Tpl, _ := template.ParseGlob("./template/parts/*")
	Tpl.New(loginHTMLName).ParseFiles(loginTemplatePath + loginResultHTMLName)
	email := r.FormValue(appconst.EMAIL)
	password := r.FormValue(appconst.PASSWORD)
	log.Println("email/pass:" + email + "/" + password)
	mailMsg := []string{}
	passWordMsg := []string{}
	if email == "" {
		mailMsg = append(mailMsg, message.ErrMsgNoEmail)
	}
	if password == "" {
		passWordMsg = append(passWordMsg, message.ErrMsgNoPassword)
	}

	errMsgMap := map[string][]string{}
	viewData := map[string]string{}
	if len(mailMsg) > 0 || len(passWordMsg) > 0 {
		errMsgMap["mail"] = mailMsg
		errMsgMap["password"] = passWordMsg
		// sessionにmap形式のデータを追加できるように設定
		gob.Register(map[string][]string{})
		gob.Register(map[string]string{})

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
	if errMsg != "" {
		sokanMsg := []string{}
		sokanMsg = append(sokanMsg, "ログインIDとパスワードに誤りがあります。")

		errMsgMap["sokan"] = sokanMsg
		viewData["mail"] = email
		session, _ := ulogin.GetSession(r)
		session.AddFlash(errMsgMap, appconst.SessionMsg)
		session.AddFlash(viewData, appconst.SessionViewData)
		session.Save(r, w)
		http.Redirect(w, r, appconst.RootURL, http.StatusFound)
	}

	// ログインセッションに登録
	session, _ := ulogin.GetSession(r)
	// セッションにデータにデータをつめる
	session.Values[appconst.SessionLoginUser] = resUser
	session.Save(r, w)
	http.Redirect(w, r, appconst.BookURL, http.StatusFound)

}
