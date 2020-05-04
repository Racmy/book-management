package ulogin

import (
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/common/message"
	"github.com/docker_go_nginx/app/db/userdao"
	"github.com/gorilla/sessions"
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

/*
ログインユーザの取得
*/
func GetLoginUser(r *http.Request) userdao.User {
	session, _ := GetSession(r)
	user := session.Values[appconst.SessionLoginUser].(userdao.User)
	return user
}
