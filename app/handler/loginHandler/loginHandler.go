package loginHandler

import (
	"log"
	"net/http"
	"text/template"
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/common/message"
	"github.com/docker_go_nginx/app/utility/ulogin"
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
	session, err := ulogin.GetSession(r)
	if err != nil {
		log.Println("session err")
	}
	session.Values[appconst.SessionErrFlg] = false
	Tpl, _ := template.ParseGlob("./template/parts/*")
	Tpl.New(loginHTMLName).ParseFiles(loginTemplatePath + loginResultHTMLName)
	email := r.FormValue(appconst.EMAIL)
	password := r.FormValue(appconst.PASSWORD)
	log.Println("email/pass:" + email + "/" + password)
	if email == "" || password == ""{
		session.Values[appconst.SessionErrFlg] = true
		log.Println("test")
		if (email == ""){
			session.AddFlash(message.ErrMsgNoEmail,appconst.SessionErrMsgEmail)
		}
		if (password == ""){
			session.AddFlash(message.ErrMsgNoPassword,appconst.SessionErrMsgPassword)
		}
		session.Save(r, w)
		http.Redirect(w, r, appconst.RootURL, http.StatusFound)
	} else{
		resUser,ErrMsg := ulogin.LoginCheck(email,password)
		if ErrMsg != ""{
			session.Values[appconst.SessionErrFlg] = true
			session.AddFlash(ErrMsg,appconst.SessionErrMsgNoUser)
			session.Save(r, w)
			http.Redirect(w, r, appconst.RootURL, http.StatusFound)
		} else{
			session.Values[appconst.SessionUserID] = resUser.ID
			session.Values[appconst.SessionUserName] = resUser.Name
			session.Values[appconst.SessionUserImagePath] = resUser.ImagePath
			session.Save(r, w)
			http.Redirect(w, r, appconst.BookURL, http.StatusFound)
		}
	}
	
	
}


