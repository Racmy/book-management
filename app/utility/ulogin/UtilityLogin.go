package ulogin

import (
	"net/http"
	"github.com/gorilla/sessions"
	"github.com/docker_go_nginx/app/common/appstructure"
	"github.com/docker_go_nginx/app/common/message"
	"github.com/docker_go_nginx/app/db/UserDao"
)

const (
	CookieName string = "cokkieName1"
)
var(
	// キーの長さは 16, 24, 32 バイトのいずれかでなければならない。
    // (AES-128, AES-192 or AES-256)
    key = []byte("super-secret-key")
    store = sessions.NewCookieStore(key)
)
/**
	ログインチェック
*/
func LoginCheck(mailAddress string, password string) (appstructure.UserData ,string){
	var responceUserData appstructure.UserData
	var errMsg string
	responceUserData, err := userdao.GetUserByEmailAndPass(mailAddress,password)
	if err != nil{
		errMsg = message.ErrMsgServerErr
	}
	if responceUserData.ID == 0 {
		errMsg = message.ErrMsgNoUserErr
	}
	return responceUserData, errMsg

}

func GetSession(r *http.Request) (*sessions.Session, error){
	return store.Get(r,CookieName)
}

// func SessionCheck(){

// }

