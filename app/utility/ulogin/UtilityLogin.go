package ulogin

import (
	"errors"
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/common/appstructure"
	"github.com/docker_go_nginx/app/common/message"
	"github.com/docker_go_nginx/app/db/userdao"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

const (
	CookieName string = "cokkieName1"
)

var (
	// キーの長さは 16, 24, 32 バイトのいずれかでなければならない。
	// (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

/**
ログインチェック
*/
func LoginCheck(mailAddress string, password string) (userdao.User, string) {
	var responceUserData userdao.User
	var errMsg string
	responceUserData, err := userdao.GetUserByEmailAndPass(mailAddress, password)
	if err != nil {
		errMsg = message.ErrMsgServerErr
	}
	if responceUserData.ID == 0 {
		errMsg = message.ErrMsgNoUserErr
	}
	return responceUserData, errMsg

}

/**
ログイン済みかチェック
@param http.Request
@return bool true：ログイン済、false：未ログイン
*/
func IsLogined(r *http.Request) bool {
	session, _ := GetSession(r)
	return session.Values[appconst.SessionLoginUser] != nil
}

/*
セションの取得
*/
func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, CookieName)
}

func SessionCheck(w http.ResponseWriter, r *http.Request) (*sessions.Session, error) {
	retSession, err := store.Get(r, CookieName)
	if err != nil {
		log.Println("SessionCheck err")
		return retSession, err
	} else {
		if retSession.Values[appconst.SessionLoginUser] == nil {
			return nil, errors.New("no User ID")
		}
		return retSession, err
	}
}

/*
ログインユーザの取得
*/
func GetLoginUser(r *http.Request) (userdao.User, error) {
	session, _ := GetSession(r)
	if session.Values[appconst.SessionLoginUser] == nil {
		err := errors.New("ログインセッションが切れています。")
		return userdao.User{}, err
	}
	user := session.Values[appconst.SessionLoginUser].(userdao.User)
	return user, nil
}

/*
ログインユーザIDの取得
*/
func GetLoginUserId(r *http.Request) (int, error) {
	session, _ := GetSession(r)
	if session.Values[appconst.SessionLoginUser] == nil {
		err := errors.New("ログインセッションが切れています。")
		return -1, err
	}
	user := session.Values[appconst.SessionLoginUser].(userdao.User)
	return int(user.ID), nil
}

/*
登録・更新・削除成功フラグがセッションに格納されているかの判定
*/
func GetSessionFlg(w http.ResponseWriter, r *http.Request) bool {
	session, _ := GetSession(r)
	sessionFlg := session.Flashes(appconst.SessionFlg)
	session.Save(r, w)
	if len(sessionFlg) > 0 {
		log.Print(sessionFlg)
		return sessionFlg[0].(bool)
	}
	return false
}

/*
登録・更新・削除成功フラグをtrueにする
*/
func SetSessionFlg(w http.ResponseWriter, r *http.Request) {
	session, _ := GetSession(r)
	session.AddFlash(true, appconst.SessionFlg)
	session.Save(r, w)
}

/*
セッションから画面データとメッセージを取り出し、ResponseDataにつめた状態で返す
*/
func GetViewDataAndMessage(w http.ResponseWriter, r *http.Request) appstructure.ResponseData {
	// セッションの取得
	session, _ := GetSession(r)

	var castedMessage map[string][]string
	var castedViewData map[string]string
	if message := session.Flashes(appconst.SessionMsg); len(message) > 0 {
		// 画面データとメッセージの取得
		castedMessage = message[0].(map[string][]string)
	}
	if viewData := session.Flashes(appconst.SessionViewData); len(viewData) > 0 {
		// 画面データとメッセージの取得
		castedViewData = viewData[0].(map[string]string)
	}
	session.Save(r, w)

	return appstructure.CreateResponseData(castedViewData, castedMessage)

}
