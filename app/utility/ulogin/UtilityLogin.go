package ulogin

import (
	"log"
	"net/http"
	"errors"
	"github.com/gorilla/sessions"
	"github.com/docker_go_nginx/app/common/appstructure"
	"github.com/docker_go_nginx/app/common/message"
	"github.com/docker_go_nginx/app/db/UserDao"
	"github.com/docker_go_nginx/app/common/appconst"
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

func SessionCheck(w http.ResponseWriter , r *http.Request) (*sessions.Session, error){
	retSession, err := store.Get(r,CookieName)
	if err != nil {
		log.Println("SessionCheck err")
		return retSession, err
	} else{
		if retSession.Values[appconst.SessionUserID] == nil || retSession.Values[appconst.SessionUserID].(int) <= 0 {
			return nil, errors.New("no User ID")
		}
		return retSession, err
	}
}

