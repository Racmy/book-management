package main

import (
	"encoding/gob"
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/common/appstructure"
	"github.com/docker_go_nginx/app/handler/bookHandler"
	"github.com/docker_go_nginx/app/handler/loginHandler"
	"github.com/docker_go_nginx/app/utility/ulogin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"regexp"
	"text/template"
)

var rootTemplatePath = "./template/"
var homeTemplatePath = rootTemplatePath + "home/"
var homeHTMLName = "index.html"
var userTemplatePath = rootTemplatePath + "user/"
var userRegistHTMLName = "regist.html"
var userEditHTMLName = "edit.html"
var userPasswordOrderHTMLName = "password_order.html"
var userPasswordRegistHTMLName = "password_regist.html"

// アカウント登録画面用の画面データ構造
type UserRegistResponseData struct {
	ViewData map[string]string
	Message  map[string][]string
}

/*
	ホーム画面を表示するハンドラ
*/
func homeHandler(w http.ResponseWriter, r *http.Request) {
	Tpl, _ := template.ParseGlob("./template/parts/*")
	Tpl.New(homeHTMLName).ParseFiles(homeTemplatePath + homeHTMLName)
	var errMsg appstructure.HomeErrorMessage
	session, _ := ulogin.GetSession(r)

	if errFlg := session.Values[appconst.SessionErrFlg]; errFlg != nil && errFlg.(bool) == true {
		if flashErrMsg := session.Flashes(appconst.SessionErrMsgEmail); len(flashErrMsg) > 0 {
			errMsg.EmailErr = flashErrMsg[0].(string)
		}
		if flashErrMsg := session.Flashes(appconst.SessionErrMsgPassword); len(flashErrMsg) > 0 {
			errMsg.PasswordErr = flashErrMsg[0].(string)
		}
		if flashErrMsg := session.Flashes(appconst.SessionErrMsgNoUser); len(flashErrMsg) > 0 {
			errMsg.NoUserErr = flashErrMsg[0].(string)
		}
	}
	session.Save(r, w)

	if err := Tpl.ExecuteTemplate(w, homeHTMLName, errMsg); err != nil {
		log.Fatal(err)
	}
}

/*
	ユーザを新規登録するハンドラ
*/
func userRegistHandler(w http.ResponseWriter, r *http.Request) {
	Tpl, _ := template.ParseGlob("./template/parts/*")
	switch r.Method {
	case http.MethodPost:
		//ユーザ登録処理
		// メールアドレスのバリデーション
		mail := r.FormValue("mail")
		mailMsg := []string{}
		if noValueValidation(mail, "メールアドレス", &mailMsg) {
			reg := regexp.MustCompile(`^.+\@.+\..+$`)
			if !reg.MatchString(mail) {
				mailMsg = append(mailMsg, "メールアドレスの形式に誤りがあります。")
			}
		} else {
			mail = ""
		}

		// 名前のバリデーション
		name := r.FormValue("name")
		nameMsg := []string{}
		if !noValueValidation(name, "ユーザ名", &nameMsg) {
			name = ""
		}

		// パスワードのバリデーション
		password := r.FormValue("password")
		passwordMsg := []string{}
		isSetPassword := noValueValidation(password, "パスワード", &passwordMsg)

		// パスワード（再入力）のバリデーション
		rePassword := r.FormValue("re_password")
		rePasswordMsg := []string{}
		isSetRePassword := noValueValidation(rePassword, "再入力パスワード", &rePasswordMsg)

		// パスワードと再パスワードの一致チェック
		sokanCheckMsg := []string{}
		if isSetPassword && isSetRePassword {
			if password != rePassword {
				sokanCheckMsg = append(sokanCheckMsg, "パスワードが一致しません。")
			}
		}

		errMsgMap := map[string][]string{}
		viewData := map[string]string{}
		//　エラーがある場合は登録画面へリダイレクト
		if len(mailMsg) > 0 || len(nameMsg) > 0 || len(passwordMsg) > 0 || len(rePasswordMsg) > 0 || len(sokanCheckMsg) > 0 {
			session, _ := ulogin.GetSession(r)
			errMsgMap["mail"] = mailMsg
			errMsgMap["name"] = nameMsg
			errMsgMap["password"] = passwordMsg
			errMsgMap["repassword"] = rePasswordMsg
			errMsgMap["sokanCheck"] = sokanCheckMsg

			viewData["mail"] = mail
			viewData["name"] = name

			gob.Register(map[string][]string{})
			gob.Register(map[string]string{})

			session.AddFlash(errMsgMap, appconst.SessionMsg)
			session.AddFlash(viewData, appconst.SessionViewData)

			err := session.Save(r, w)

			if err != nil {
				log.Print(err)
			}
			http.Redirect(w, r, appconst.UserRegistURL, http.StatusFound)
		}
		// エラーがない場合はユーザテーブルに登録

		// ログインセッションに登録

		// 本一覧画面へ遷移

	default:
		// Get・PUT・PATCH・DELETEなどできた場合は登録画面を表示
		// ユーザ登録画面の表示
		session, _ := ulogin.GetSession(r)
		Tpl.New(userRegistHTMLName).Option("missingkey=zero").ParseFiles(userTemplatePath + userRegistHTMLName)
		if message := session.Flashes(appconst.SessionMsg); len(message) > 0 {
			castedMessage := message[0].(map[string][]string)
			viewData := session.Flashes(appconst.SessionViewData)[0].(map[string]string)
			session.Save(r, w)

			// 画面表示データ構造作成
			responseData := UserRegistResponseData{
				ViewData: viewData,
				Message:  castedMessage,
			}

			if err := Tpl.ExecuteTemplate(w, userRegistHTMLName, responseData); err != nil {
				log.Fatal(err)
			}
		} else {
			// 画面表示データ構造作成
			responseData := UserRegistResponseData{
				ViewData: map[string]string{},
				Message:  nil,
			}

			if err := Tpl.ExecuteTemplate(w, userRegistHTMLName, responseData); err != nil {
				log.Fatal(err)
			}
		}

	}
}

/*
	空文字が判定する。
	空文字の場合はメッセージに、バリデーションメッセージを追加する。
	@return true: 値有り、false：値無し
*/
func noValueValidation(value string, itemName string, msg *[]string) bool {
	if len(value) == 0 {
		*msg = append(*msg, itemName+"を入力してください。")
		return false
	}
	return true
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
	Tpl.New(userPasswordRegistHTMLName).ParseFiles(userTemplatePath + userPasswordRegistHTMLName)
	if err := Tpl.ExecuteTemplate(w, userPasswordRegistHTMLName, nil); err != nil {
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
	r.HandleFunc(appconst.UserRegistURL, userRegistHandler)
	// ユーザ登録情報の更新ハンドラ
	r.HandleFunc(appconst.UserEditURL, userEditHandler)
	// ユーザパスワード再発行申込ハンドラ
	r.HandleFunc(appconst.UserPassWordOrderURL, userPassWordOrderHandler)
	// ユーザパスワード再登録ハンドラ
	r.HandleFunc(appconst.UserPassWordRegistURL, userPasswordRegist)
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
