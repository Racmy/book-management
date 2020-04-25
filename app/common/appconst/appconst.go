package appconst

// URL
const (
	RootURL    string = "/"
	UserURL    string = "/user/regist"
	BookURL	   string = RootURL + "book"
	BookDetailLURL string = BookURL + "/detail"
	BookRegistURL string = BookURL + "/regist"
	BookRegistProcessURL string = BookRegistURL + "/process"
	BookRegistResultURL string = BookRegistURL + "-result"
	BookSearchURL string = BookURL + "/search"
	BookUpdatehURL string = BookURL + "/update"
	BookDeleteURL string = BookURL + "/delete"
)