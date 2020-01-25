package main

import (
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/docker_go_nginx/app/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

/**
【Book構造体用の定数】
ROOT ・・・ドキュメントルート
ID・・・本のID
TITLE・・・本のタイトル
AUTHOR・・・本の著者
LatestIssue・・・所持巻数
IMGPATH・・・画像へのパス
*/
const (
	ID          string = "Id"
	TITLE       string = "Title"
	AUTHOR      string = "Author"
	LatestIssue string = "LatestIssue"
	IMGPATH     string = "FrontCoverImagePath"
)

/*
【エラーメッセージの文言】
*/
const (
	ErrMsgTitleNull string = "タイトルを入力してください"
	ErrMsgAuthNull  string = "著者を入力してください"
	ErrMsgLiNull    string = "最新所持巻数を数字で入力してください"
	ErrMsgServerErr string = "現在不安定な状態です。再度、お試しください。"
)

/*
【成功時のメッセージ文言】
*/
const (
	SucMsgUpdate        string = "更新が完了しました"
	FrontCoverImageName string = "FrontCoverImageName"
)

/*
【パス関係の定数】
*/
const (
	ROOT    string = "/"
	ImgPath string = "/static/img/"
)

//ResponseDataForDetail ...　本詳細画面用の構造体
type ResponseDataForDetail struct {
	Book   db.Book
	ErrMsg []string
	SucMsg []string
}

// ResponseData ...　一覧画面用のレスポンスデータ
type ResponseData struct {
	Keyword string
	Books   []db.Book
}

// RegistValue ...登録用の構造体
type RegistValue struct {
	Title       string
	Author      string
	LatestIssue float64
	ErrString   []string
}

/*
	本を登録画面へのハンドラ
*/
func bookRegistHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("./template/bookRegist.html"))

	errString := []string{}
	title := r.FormValue(TITLE)
	author := r.FormValue(AUTHOR)
	latestIssueString := r.FormValue(LatestIssue)
	latestIssue, strConvErr := strconv.ParseFloat(latestIssueString, 64)
	tmpErrCheckFlag := r.FormValue("ErrCheckFlag")

	if tmpErrCheckFlag == "1" {
		if title == "" {
			errString = append(errString, ErrMsgTitleNull)
		}
		if author == "" {
			errString = append(errString, ErrMsgAuthNull)
		}
		if strConvErr != nil {
			errString = append(errString, ErrMsgLiNull)
		}
	}

	if strConvErr != nil {
		latestIssue = 1
	}

	tmp := RegistValue{
		Title:       title,
		Author:      author,
		LatestIssue: latestIssue,
		ErrString:   errString,
	}

	if err := tmpl.ExecuteTemplate(w, "bookRegist.html", tmp); err != nil {
		log.Fatal(err)
	}
}

/*
	本を登録完了画面へのハンドラ
*/
func bookInsertHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("./template/bookRegistResult.html"))

	//表紙画像がuploadされたかどうかを判定するフラグの初期化
	fileUploadFlag := true
	frontCoverImagePath := ""
	//表紙画像を格納する変数宣言
	var file multipart.File
	var fileHeader *multipart.FileHeader
	// POSTされたファイルデータをメモリに格納
	//33554432 約30MByte(8Kのping形式には耐えられない)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println("【main.go bookInsertHandler】not ParseMultipartForm")
		fileUploadFlag = false
	} else {
		file, fileHeader, err = r.FormFile(FrontCoverImageName)
		if err != nil {
			log.Println("【main.go bookInsertHandler】not file upload")
			fileUploadFlag = false
		}
	}
	//表紙画像がuploadされている時
	if fileUploadFlag {
		//【TODO】余田へ
		// Filerクラスを作成する際に、static/img/（相対パス）を絶対パスでできるようにしてください。
		frontCoverImagePath = "static/img/" + fileHeader.Filename

		log.Println(frontCoverImagePath)
		// サーバー側に保存するために空ファイルを作成
		var saveImage *os.File
		saveImage, err = os.Create(frontCoverImagePath)
		if err != nil {
			log.Println("【main.go bookInsertHandler】os.Create Error")
			log.Println(err)
			return
		}
		defer saveImage.Close()
		defer file.Close()
		size, err := io.Copy(saveImage, file)
		if err != nil {
			log.Println("【main.go bookInsertHandler】io.Copy Error")
			log.Println(err)
		}
<<<<<<< HEAD
		log.Print("File Upload データサイズ")
		log.Println(size)
=======
		log.Printf("File Upload データサイズ" + strconv.FormatInt(size, 10))
>>>>>>> develop
		frontCoverImagePath = "/" + frontCoverImagePath
	}

	r.ParseForm()

	title := r.Form[TITLE][0]
	author := r.Form[AUTHOR][0]
	latestIssueString := r.Form[LatestIssue][0]
	latestIssue, strConvErr := strconv.ParseFloat(latestIssueString, 64)

	if (title == "") || (author == "") || (strConvErr != nil) {
		var url = "/regist"
		url += "?Title=" + title + "&Author=" + author + "&LatestIssue=" + latestIssueString + "&ErrCheckFlag=1"
		http.Redirect(w, r, url, http.StatusFound)
	}

	insertBook := db.Book{
		Title:               title,
		Author:              author,
		LatestIssue:         latestIssue,
		FrontCoverImagePath: frontCoverImagePath,
	}

	// 入力チェック後の本データを登録
	id, insErr := db.InsertBook(insertBook)
	//　本の登録失敗時には、ホーム画面へ遷移
	if insErr != nil {
		http.Redirect(w, r, ROOT, http.StatusFound)
	}

	// 登録した際に発行されるIDで本情報をDBから取得
	book, canGet := db.GetBookByID(strconv.FormatInt(id, 10))
	// 取得失敗時は、ホーム画面へ遷移
	if canGet == false {
		http.Redirect(w, r, ROOT, http.StatusFound)
	}

	// テンプレートにデータを埋め込む
	if err := tmpl.ExecuteTemplate(w, "bookRegistResult.html", book); err != nil {
		log.Fatal(err)
	}

}

/*
	ホーム画面へのハンドラ
*/
func homeHandler(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("./template/list.html"))
	var responseData ResponseData
	responseData.Books = db.GetAllBooks()
	responseData.Keyword = ""
	if err := tpl.ExecuteTemplate(w, "list.html", responseData); err != nil {
		log.Fatal(err)
	}
}

/*
	本詳細画面へのハンドラ
*/
func bookDetailHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	if id := query.Get("Id"); query.Get("Id") != "" {
		// 画面からIdを取得し、DBから紐つくデータを取得
		var responseData ResponseDataForDetail
		var canGet bool
		responseData.Book, canGet = db.GetBookByID(id)
		// データ取得失敗時はホームへ戻す
		if canGet == false {
			http.Redirect(w, r, ROOT, http.StatusFound)
		}

		//更新成功時のメッセージを格納
		if query.Get("sucFlg") != "" {
			responseData.SucMsg = append(responseData.SucMsg, SucMsgUpdate)
		}
		var tpl = template.Must(template.ParseFiles("./template/detail.html"))
		if err := tpl.ExecuteTemplate(w, "detail.html", responseData); err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}

}

/*
	本の検索のためのハンドラ
*/
func bookSearchHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	if keyword := query.Get("keyword"); query.Get("keyword") != "" {
		// keywordがnullの場合は、HOMEへリダイレクト
		if keyword == "" {
			http.Redirect(w, r, ROOT, http.StatusFound)
		}

		var tpl = template.Must(template.ParseFiles("./template/list.html"))

		// ResponseDataの作成
		var responseData ResponseData
		responseData.Keyword = keyword
		responseData.Books = db.GetSearchedBooks(keyword)

		if err := tpl.ExecuteTemplate(w, "list.html", responseData); err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(w, r, ROOT, http.StatusFound)
	}
}

/*
	本情報を更新するHandler
*/
func bookUpdateHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue(TITLE)
	author := r.FormValue(AUTHOR)
	imgPath := r.FormValue(IMGPATH)
	id := r.FormValue(ID)
	idInt, strConvErr := strconv.Atoi(id)
	latestIssueString := r.FormValue(LatestIssue)
	latestIssue, strConvErr := strconv.ParseFloat(latestIssueString, 64)

	var errMsg []string
	/*エラーチェック【相談】登録との共通化*/
	if title == "" {
		errMsg = append(errMsg, ErrMsgTitleNull)
	}
	if author == "" {
		errMsg = append(errMsg, ErrMsgAuthNull)
	}
	if strConvErr != nil {
		errMsg = append(errMsg, ErrMsgLiNull)
	}

	// 入力エラーがない場合は更新処理を実施
	if len(errMsg) == 0 {
		// 入力データで更新
		updateBook := db.Book{ID: idInt, Title: title, Author: author, LatestIssue: latestIssue}
		id := db.UpdateBook(updateBook)

		idString := strconv.Itoa(id)

		// 更新処理が失敗していない場合は、詳細画面へ遷移（detail.html）
		var url string
		if id != -1 {
			log.Print("【main.go　UpdateBookHander】success update")
			// 成功したことをDetailに伝えるためにsucFlgをつける
			url = "/detail" + "?Id=" + idString + "&sucFlg=1"
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			// 更新に失敗したことを、エラーメッセージにつめる
			errMsg = append(errMsg, ErrMsgServerErr)
		}
	}

	// 以下、入力ミス・更新失敗時の処理
	log.Print("【main.go　UpdateBookHander】invalid input value or fail update")

	var tmpl = template.Must(template.ParseFiles("./template/detail.html"))
	// エラー時は、画面から送られてきたデータを渡す
	inputBook := db.Book{ID: idInt, Title: title, Author: author, LatestIssue: latestIssue, FrontCoverImagePath: imgPath}

	// 画面に表示するデータを格納
	responseData := ResponseDataForDetail{Book: inputBook, ErrMsg: errMsg}

	if err := tmpl.ExecuteTemplate(w, "detail.html", responseData); err != nil {
		log.Fatal(err)
	}

}

/*
	ルーティング
*/
func main() {

	r := mux.NewRouter()
	r.HandleFunc(ROOT+"regist", bookRegistHandler)
	r.HandleFunc(ROOT+"regist/success", bookInsertHandler)
	r.HandleFunc(ROOT, homeHandler)
	r.HandleFunc(ROOT+"search", bookSearchHandler)
	r.HandleFunc(ROOT+"detail", bookDetailHandler)
	r.HandleFunc(ROOT+"update", bookUpdateHandler)
	// cssフレームワーク読み込み
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules/"))))
	// 画像フォルダ
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.Handle(ROOT, r)
	http.ListenAndServe(":3000", nil)
}
