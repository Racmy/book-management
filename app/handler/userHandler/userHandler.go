package userHandler

import (
	"encoding/gob"
	"net/http"
	"regexp"

	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/common/message"
	"github.com/docker_go_nginx/app/common/processTemplate"
	"github.com/docker_go_nginx/app/db/userdao"
	"github.com/docker_go_nginx/app/utility/uDB"
	"github.com/docker_go_nginx/app/utility/ulogin"
	_ "github.com/go-sql-driver/mysql"
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

//
const (
	IMAGE string = "img"
)

/*
	ユーザを新規登録するハンドラ
*/
func UserRegistHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		//ユーザ登録処理
		// メールアドレスのバリデーション
		mail := r.FormValue("mail")
		mailMsg := []string{}
		if noValueValidation(mail, "メールアドレス", &mailMsg) {
			reg := regexp.MustCompile(`^.+\@.+\..+$`)
			if !reg.MatchString(mail) {
				mailMsg = append(mailMsg, message.ErrEmailStyle)
			} else {
				if userdao.IsSetEmail(mail) {
					mailMsg = append(mailMsg, message.RegisteredEmail)
				}
			}
		} else {
			mail = ""
		}

		// 名前のバリデーション
		name := r.FormValue("name")
		nameMsg := []string{}
		if !noValueValidation(name, "ユーザ名", &nameMsg) {
			name = ""
		} else {
			if userdao.IsSetName(name) {
				nameMsg = append(nameMsg, message.RegisteredUserName)
			}
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
		session, _ := ulogin.GetSession(r)
		if len(mailMsg) > 0 || len(nameMsg) > 0 || len(passwordMsg) > 0 || len(rePasswordMsg) > 0 || len(sokanCheckMsg) > 0 {
			// エラーメッセージをmapにつめる
			errMsgMap["mail"] = mailMsg
			errMsgMap["name"] = nameMsg
			errMsgMap["password"] = passwordMsg
			errMsgMap["repassword"] = rePasswordMsg
			errMsgMap["sokanCheck"] = sokanCheckMsg
			// １つ前に入力したデータをmapにつめる
			viewData["mail"] = mail
			viewData["name"] = name

			// sessionにmap形式のデータを追加できるように設定
			gob.Register(map[string][]string{})
			gob.Register(map[string]string{})

			// セッションにデータにデータをつめる
			session.AddFlash(errMsgMap, appconst.SessionMsg)
			session.AddFlash(viewData, appconst.SessionViewData)

			// セッションの保存
			err := session.Save(r, w)
			uDB.ErrCheck(err)
			http.Redirect(w, r, appconst.UserRegistURL, http.StatusFound)
		} else {
			user := userdao.GetUserInstance(0, mail, name, password, "")
			// エラーがない場合はユーザテーブルに登録
			registeredUser, err := userdao.InsertUser(user)

			// エラー場合は登録画面に戻す
			if err != nil {
				http.Redirect(w, r, appconst.UserRegistURL, http.StatusFound)
			}
			//　ホーム画面遷移する
			// ログインセッションに登録
			gob.Register(userdao.User{})

			// セッションにデータにデータをつめる
			session.Values[appconst.SessionLoginUser] = registeredUser
			session.Save(r, w)

			http.Redirect(w, r, appconst.BookURL, http.StatusFound)
		}

	default:
		// Get・PUT・PATCH・DELETEなどできた場合は登録画面を表示
		// ユーザ登録画面の表示
		responseData := ulogin.GetViewDataAndMessage(w, r)
		// 後処理
		processTemplate.PostHandler(w, r, userTemplatePath, userRegistHTMLName, responseData)
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
func UserEditHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := ulogin.GetLoginUser(r)
	switch r.Method {
	case http.MethodPost:
		//ユーザ登録処理
		// メールアドレスのバリデーション
		mail := r.FormValue("mail")
		mailMsg := []string{}
		if noValueValidation(mail, "メールアドレス", &mailMsg) {
			reg := regexp.MustCompile(`^.+\@.+\..+$`)
			if !reg.MatchString(mail) {
				mailMsg = append(mailMsg, message.ErrEmailStyle)
			} else {
				if userdao.IsSetEmailExceptID(mail, user.ID) {
					mailMsg = append(mailMsg, message.RegisteredEmail)
				}
			}
		} else {
			mail = ""
		}
		// 名前のバリデーション
		name := r.FormValue("name")
		nameMsg := []string{}
		if !noValueValidation(name, "ユーザ名", &nameMsg) {
			name = ""
		} else {
			if userdao.IsSetNameExceptID(name, user.ID) {
				nameMsg = append(nameMsg, message.RegisteredUserName)
			}
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
		session, _ := ulogin.GetSession(r)
		if len(mailMsg) > 0 || len(nameMsg) > 0 || len(passwordMsg) > 0 || len(rePasswordMsg) > 0 || len(sokanCheckMsg) > 0 {
			// エラーメッセージをmapにつめる
			errMsgMap["mail"] = mailMsg
			errMsgMap["name"] = nameMsg
			errMsgMap["password"] = passwordMsg
			errMsgMap["repassword"] = rePasswordMsg
			errMsgMap["sokanCheck"] = sokanCheckMsg
			// １つ前に入力したデータをmapにつめる
			viewData["mail"] = mail
			viewData["name"] = name
			// sessionにmap形式のデータを追加できるように設定
			gob.Register(map[string][]string{})
			gob.Register(map[string]string{})

			// セッションにデータにデータをつめる
			session.AddFlash(errMsgMap, appconst.SessionMsg)
			session.AddFlash(viewData, appconst.SessionViewData)

			// セッションの保存
			err := session.Save(r, w)
			uDB.ErrCheck(err)
			http.Redirect(w, r, appconst.UserEditURL, http.StatusFound)
		} else {
			pre_user, _ := ulogin.GetLoginUser(r)
			user := userdao.GetUserInstance(pre_user.ID, mail, name, password, "")
			// エラーがない場合はユーザテーブルに登録
			registeredUser, err := userdao.UpdateUser(user)

			// エラー場合は登録画面に戻す
			if err != nil {
				uDB.ErrCheck(err)
				http.Redirect(w, r, appconst.UserEditURL, http.StatusFound)
			}
			//　ホーム画面遷移する
			// ログインセッションに登録
			gob.Register(userdao.User{})
			// sessionにmap形式のデータを追加できるように設定
			gob.Register(map[string][]string{})
			sucMsg := []string{}
			sucMsg = append(sucMsg, message.SucMsgUpdate)
			errMsgMap["sucMsg"] = sucMsg
			session.AddFlash(errMsgMap, appconst.SessionMsg)
			// セッションにデータにデータをつめる
			session.Values[appconst.SessionLoginUser] = registeredUser
			session.Save(r, w)

			http.Redirect(w, r, appconst.UserEditURL, http.StatusFound)
		}

	default:
		// セッションに残っっている画面データ等を含めたレスポンスデータの取得
		responseData := ulogin.GetViewDataAndMessage(w, r)
		// 後処理
		processTemplate.PostHandler(w, r, userTemplatePath, userEditHTMLName, responseData)
	}
}

/*
	ユーザのログインパスワード再発行
*/
func UserPassWordOrderHandler(w http.ResponseWriter, r *http.Request) {
	// セッションに残っっている画面データ等を含めたレスポンスデータの取得
	responseData := ulogin.GetViewDataAndMessage(w, r)
	// 後処理
	processTemplate.PostHandler(w, r, userTemplatePath, userPasswordOrderHTMLName, responseData)
}

/*
	ユーザのパスワード再登録画面
*/
func UserPasswordRegist(w http.ResponseWriter, r *http.Request) {
	// セッションに残っっている画面データ等を含めたレスポンスデータの取得
	responseData := ulogin.GetViewDataAndMessage(w, r)
	// 後処理
	processTemplate.PostHandler(w, r, userTemplatePath, userPasswordRegistHTMLName, responseData)
}
