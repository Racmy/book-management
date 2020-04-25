package loginHandler

import (
	"log"
	"net/http"
	"text/template"
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/utility/ulogin"
)

var Tpl *template.Template

var rootTemplatePath = "./template/"
var loginTemplatePath = rootTemplatePath + "login/"
var loginHTMLName = "testLogin.html"
var loginResultHTMLName = "testLoginResult.html"

/**
	ログイン画面へのハンドラ
*/
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	Tpl, _ := template.ParseGlob("./template/parts/*")
	Tpl.New(loginHTMLName).ParseFiles(loginTemplatePath + loginHTMLName)
	if err := Tpl.ExecuteTemplate(w, loginHTMLName, nil); err != nil {
		log.Fatal(err)
	}
}

/**
	ログインチェックのハンドラ
*/
func LoginCheckHandler(w http.ResponseWriter, r *http.Request) {
	session, err := ulogin.GetSession(r)
	if err != nil {
		log.Println("session err")
	}
	Tpl, _ := template.ParseGlob("./template/parts/*")
	Tpl.New(loginHTMLName).ParseFiles(loginTemplatePath + loginResultHTMLName)
	email := r.FormValue(appconst.EMAIL)
	password := r.FormValue(appconst.PASSWORD)
	
	resUser,ErrMsg := ulogin.LoginCheck(email,password)
	if ErrMsg != ""{
		log.Println(ErrMsg)
	}

	session.Values[appconst.SessionUserID] = resUser.ID
	session.Values[appconst.SessionUserName] = resUser.Name
	session.Values[appconst.SessionUserImagePath] = resUser.ImagePath
	session.Save(r, w)
	if err := Tpl.ExecuteTemplate(w, loginResultHTMLName, resUser); err != nil {
		log.Fatal(err)
	}
}


