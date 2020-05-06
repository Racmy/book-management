package processTemplate

import (
	"log"
	"net/http"
	"text/template"

	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/common/appstructure"
	"github.com/docker_go_nginx/app/common/message"
	"github.com/docker_go_nginx/app/utility/ulogin"
)

///baseHandlerFunc...ハンドラに共通処理を付与する関数
/*
@param handler func(w http.ResponseWriter, r *http.Request) 固有処理関数
@param sessionFlg int セッションのログイン情報を用いる場合:1　それ以外:0
@return http.Handler 共通処理を付与したハンドラ
*/
func BaseHandlerFunc(handler func(w http.ResponseWriter, r *http.Request), sessionFlg int) http.Handler {
	return PreHandler(http.HandlerFunc(handler), sessionFlg)
}

///PreHandler...ハンドラの前処理
/*
@param handler http.Handler 固有処理関数
@param sessionFlg int セッションのログイン情報を用いる場合:1　それ以外:0
@return http.Handler 共通処理を付与したハンドラ
*/
func PreHandler(handler http.Handler, sessionFlg int) http.Handler {
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

///PostHandler...ハンドラの後処理
/*
@param w http.ResponseWriter
@param r *http.Request
@param templatePath string
@param htmlName string
@param responseData appstructure.ResponseData
@return
*/
func PostHandler(w http.ResponseWriter, r *http.Request, templatePath string, htmlName string, responseData appstructure.ResponseData) {
	var Tpl *template.Template
	Tpl, _ = template.ParseGlob("./template/parts/*")
	Tpl.New(htmlName).ParseFiles(templatePath + htmlName)
	responseData.LoginFlag = ulogin.IsLogined(r)
	if err := Tpl.ExecuteTemplate(w, htmlName, responseData); err != nil {
		log.Fatal(err)
	}
}
