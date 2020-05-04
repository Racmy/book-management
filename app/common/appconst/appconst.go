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
	LoginCheckURL         string = LoginURL + "/check"
)

const (
	EMAIL    string = "MailAddress"
	PASSWORD string = "Password"
)

const (
	// 俺用
	SessionMsg       string = "Message"
	SessionViewData  string = "ViewData"
	SessionLoginUser string = "LoginUser"

	// 下、余田
	SessionErrFlg         string = "SessionErrFlg"
	SessionUserID         string = "SessionUserID"
	SessionUserName       string = "SessionUserName"
	SessionUserImagePath  string = "sessionUserImagePath"
	SessionErrMsgEmail    string = "emailNil"
	SessionErrMsgPassword string = "passwordlNil"
	SessionErrMsgNoUser   string = "noUser"
)
