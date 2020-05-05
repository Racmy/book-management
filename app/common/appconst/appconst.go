package appconst

// URL
const (
	RootURL               string = "/"
	UserURL               string = RootURL + "user"
	UserRegistURL         string = UserURL + "/regist"
	UserEditURL           string = UserURL + "/edit"
	UserPassWordOrderURL  string = UserURL + "/password_order"
	UserPassWordRegistURL string = UserURL + "/password_regist"
	BookURL               string = RootURL + "book"
	BookDetailLURL        string = BookURL + "/detail"
	BookRegistURL         string = BookURL + "/regist"
	BookRegistProcessURL  string = BookRegistURL + "/process"
	BookRegistResultURL   string = BookRegistURL + "-result"
	BookSearchURL         string = BookURL + "/search"
	BookUpdatehURL        string = BookURL + "/update"
	BookDeleteURL         string = BookURL + "/delete"
	LoginURL              string = RootURL + "login"
	LogoutURL              string = RootURL + "logout"
	LoginCheckURL         string = LoginURL + "/check"
)

const (
	EMAIL    string = "MailAddress"
	PASSWORD string = "Password"
)

const (
	SessionMsg       string = "Message"   //メッセージ格納用のキー
	SessionViewData  string = "ViewData"  //画面データ格納用のキー
	SessionLoginUser string = "LoginUser" //ログインユーザ格納用のキー
	SessionFlg       string = "Flg"       //登録・更新処理の成功フラグ格納用のキー
)
