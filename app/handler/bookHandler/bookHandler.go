package bookhandler

import (
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"text/template"
	"github.com/docker_go_nginx/app/common/message"
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/db"
	"github.com/docker_go_nginx/app/utility/ufile"
	"github.com/docker_go_nginx/app/utility/ulogin"
)

var Tpl *template.Template

var rootTemplatePath = "./template/"
var bookTemplatePath = rootTemplatePath + "book/"
var bookListHTMLName = "bookList.html"
var bookDetailHTMLName = "bookDetail.html"
var bookRegistHTMLName = "bookRegist.html"
var bookRegistResultHTMLName = "bookRegistResult.html"

//BookDetailResponseData ...　本詳細画面用のレスポンスデータ
type BookDetailResponseData struct {
	Book   bookdao.Book
	NextURL string
	ErrMsg []string
	SucMsg []string
}

// BookListResponseData ...　本一覧画面用のレスポンスデータ
type BookListResponseData struct {
	Keyword string
	Books   []bookdao.Book
	SucMsg  []string
}

// BookRegistResponseData ...本登録画面用のレスポンスデータ
type BookRegistResponseData struct {
	Title       string
	Author      string
	LatestIssue float64
	ErrString   []string
}

/*
	本を登録画面へのハンドラ
*/
func BookRegistHandler(w http.ResponseWriter, r *http.Request) {
	session, err := ulogin.GetSession(r)
	if err != nil {
		log.Println("session err")
	}
	Tpl.New(bookRegistHTMLName).ParseFiles(bookTemplatePath + bookRegistHTMLName)

	errString := []string{}
	title := r.FormValue(bookdao.TITLE)
	author := r.FormValue(bookdao.AUTHOR)
	latestIssueString := r.FormValue(bookdao.LatestIssue)
	latestIssue, strConvErr := strconv.ParseFloat(latestIssueString, 64)
	tmpErrCheckFlag := r.FormValue("ErrCheckFlag")

	if tmpErrCheckFlag == "1" {
		if title == "" {
			errString = append(errString, message.ErrMsgTitleNull)
		}
		if author == "" {
			errString = append(errString, message.ErrMsgAuthNull)
		}
		if strConvErr != nil {
			errString = append(errString, message.ErrMsgLiNull)
		}
	}

	if strConvErr != nil {
		latestIssue = 1
	}
	userIP := session.Values[appconst.SessionUserImagePath].(string)
	userName := session.Values[appconst.SessionUserName].(string)

	responseData := BookRegistResponseData{
		Title:       userIP,
		Author:      userName,
		LatestIssue: latestIssue,
		ErrString:   errString,
	}

	if err := Tpl.ExecuteTemplate(w, bookRegistHTMLName, responseData); err != nil {
		log.Fatal(err)
	}
}

/*
	本を登録処理のハンドラ
*/
func BookInsertHandler(w http.ResponseWriter, r *http.Request) {
	Tpl.New(bookRegistResultHTMLName).ParseFiles(bookTemplatePath + bookRegistResultHTMLName)

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
		file, fileHeader, err = r.FormFile(bookdao.IMGPATH)
		if err != nil {
			log.Println("【main.go bookInsertHandler】not file upload")
			fileUploadFlag = false
		}
	}
	//表紙画像がuploadされている時
	if fileUploadFlag {
		frontCoverImagePath, err = ufile.DefaultFileUpload(file, fileHeader.Filename)
		if err != nil {
			//ファイルアップロード失敗
			fileUploadFlag = false
		}
	}

	r.ParseForm()

	title := r.Form[bookdao.TITLE][0]
	author := r.Form[bookdao.AUTHOR][0]
	latestIssueString := r.Form[bookdao.LatestIssue][0]
	latestIssue, strConvErr := strconv.ParseFloat(latestIssueString, 64)

	if (title == "") || (author == "") || (strConvErr != nil) {
		var url = appconst.BookRegistURL
		url += "?Title=" + title + "&Author=" + author + "&bookdao.LatestIssue=" + latestIssueString + "&ErrCheckFlag=1"
		http.Redirect(w, r, url, http.StatusFound)
	}

	insertBook := bookdao.Book{
		Title:               title,
		Author:              author,
		LatestIssue:         latestIssue,
		FrontCoverImagePath: frontCoverImagePath,
	}

	// 入力チェック後の本データを登録
	id, insErr := bookdao.InsertBook(insertBook)
	//　本の登録失敗時には、ホーム画面へ遷移
	if insErr != nil {
		http.Redirect(w, r, appconst.BookURL, http.StatusFound)
	}

	idString := strconv.FormatInt(id, 10)
	

	http.Redirect(w, r, appconst.BookRegistResultURL + "?Id=" + idString, http.StatusFound)

}

/*
	本を登録完了画面へのハンドラ
*/
func BookInsertResultHandler(w http.ResponseWriter, r *http.Request) {
	Tpl.New(bookRegistResultHTMLName).ParseFiles(bookTemplatePath + bookRegistResultHTMLName)
	r.ParseForm()
	id := r.FormValue(bookdao.ID)

	if book, err := bookdao.GetBookByID(id); err == nil {
		// テンプレートにデータを埋め込む
		if err := Tpl.ExecuteTemplate(w, bookRegistResultHTMLName, book); err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(w, r, appconst.BookURL, http.StatusFound)
	}
}

/*
	ホーム画面へのハンドラ
*/
func BookListHandler(w http.ResponseWriter, r *http.Request) {
	
	Tpl.New(bookListHTMLName).ParseFiles(bookTemplatePath + bookListHTMLName)
	var responseData BookListResponseData
	responseData.Books = bookdao.GetAllBooks()

	query := r.URL.Query()
	if query.Get("sucDelFlg") != "" {
		responseData.SucMsg = append(responseData.SucMsg, message.SucMsgDel)
	}

	responseData.Keyword = ""
	if err := Tpl.ExecuteTemplate(w, bookListHTMLName, responseData); err != nil {
		log.Fatal(err)
	}
}

/*
	本詳細画面へのハンドラ
*/
func BookDetailHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	if id := query.Get("Id"); query.Get("Id") != "" {
		// 画面からIdを取得し、DBから紐つくデータを取得
		var responseData BookDetailResponseData
		var err error
		responseData.Book, err = bookdao.GetBookByID(id)
		// データ取得失敗時はホームへ戻す
		if err != nil {
			http.Redirect(w, r, appconst.BookURL, http.StatusFound)
		}

		//更新成功時のメッセージを格納
		if query.Get("sucFlg") != "" {
			responseData.SucMsg = append(responseData.SucMsg, message.SucMsgUpdate)
		}
		Tpl.New(bookDetailHTMLName).ParseFiles(bookTemplatePath + bookDetailHTMLName)
		if err := Tpl.ExecuteTemplate(w, bookDetailHTMLName, responseData); err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(w, r, appconst.BookURL, http.StatusFound)
	}

}

/*
	本の検索のためのハンドラ
*/
func BookSearchHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	if keyword := query.Get("keyword"); query.Get("keyword") != "" {
		// keywordがnullの場合は、HOMEへリダイレクト
		if keyword == "" {
			http.Redirect(w, r, appconst.BookURL, http.StatusFound)
		}

		Tpl.New(bookListHTMLName).ParseFiles(bookTemplatePath + bookListHTMLName)

		// ResponseDataの作成
		var responseData BookListResponseData
		responseData.Keyword = keyword
		responseData.Books = bookdao.GetSearchedBooks(keyword)

		if err := Tpl.ExecuteTemplate(w, bookListHTMLName, responseData); err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(w, r, appconst.BookURL, http.StatusFound)
	}
}

/*
	本情報を更新するHandler
*/
func BookUpdateHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue(bookdao.TITLE)
	author := r.FormValue(bookdao.AUTHOR)
	imgPath := r.FormValue(bookdao.IMGPATH)
	id := r.FormValue(bookdao.ID)
	idInt, strConvErr := strconv.Atoi(id)
	latestIssueString := r.FormValue(bookdao.LatestIssue)
	latestIssue, strConvErr := strconv.ParseFloat(latestIssueString, 64)

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
		log.Println("【main.go bookUpdateHandler】not ParseMultipartForm")
		fileUploadFlag = false
	} else {
		file, fileHeader, err = r.FormFile(bookdao.NewIMGPATH)
		if err != nil {
			log.Println("【main.go bookUpdateHandler】not file upload")
			fileUploadFlag = false
		}
	}

	var errMsg []string
	/*エラーチェック【相談】登録との共通化*/
	if title == "" {
		errMsg = append(errMsg, message.ErrMsgTitleNull)
	}
	if author == "" {
		errMsg = append(errMsg, message.ErrMsgAuthNull)
	}
	if strConvErr != nil {
		errMsg = append(errMsg, message.ErrMsgLiNull)
	}

	// 入力エラーがない場合は更新処理を実施
	if len(errMsg) == 0 {
		// 入力データで更新
		updateBook := bookdao.Book{ID: idInt, Title: title, Author: author, LatestIssue: latestIssue, FrontCoverImagePath: imgPath}
		// 本画像の登録
		if fileUploadFlag {
			frontCoverImagePath, err = ufile.DefaultFileUpload(file, fileHeader.Filename)
			if err != nil {
				//ファイルアップロード失敗
				fileUploadFlag = false
				log.Println("【main.go　UpdateBookHander】fail file image upload")
			} else {
				updateBook.FrontCoverImagePath = frontCoverImagePath
			}
		}

		id, err := bookdao.UpdateBook(updateBook)
		idString := strconv.Itoa(id)

		// 更新処理が失敗していない場合は、詳細画面へ遷移（bookDetail.html）
		var url string
		if err == nil {
			log.Print("【main.go　UpdateBookHander】success update")
			// 成功したことをDetailに伝えるためにsucFlgをつける
			url = appconst.BookDetailLURL + "?Id=" + idString + "&sucFlg=1"
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			// 更新に失敗したことを、エラーメッセージにつめる
			errMsg = append(errMsg, message.ErrMsgServerErr)
		}
	}

	// 以下、入力ミス・更新失敗時の処理
	log.Print("【main.go　UpdateBookHander】invalid input value or fail update")
	Tpl.New(bookDetailHTMLName).ParseFiles(bookTemplatePath + bookDetailHTMLName)
	// エラー時は、画面から送られてきたデータを渡す
	inputBook := bookdao.Book{ID: idInt, Title: title, Author: author, LatestIssue: latestIssue, FrontCoverImagePath: imgPath}

	// 画面に表示するデータを格納
	responseData := BookDetailResponseData{Book: inputBook, ErrMsg: errMsg}

	if err := Tpl.ExecuteTemplate(w, bookDetailHTMLName, responseData); err != nil {
		log.Fatal(err)
	}
}

/*
	本の削除ハンドラ
	/detail →　/book
*/
func BookDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue(bookdao.ID)
	err := bookdao.DeleteBookByID(id)

	var errMsg []string
	// 削除失敗時は詳細画面へ遷移してエラーメッセージを表示
	if err != nil {
		// エラーの表示
		log.Fatal("【main.go bookDeleteHandler】本の削除に失敗しました。")
		// 表示用のデータ準備
		Tpl.New(bookDetailHTMLName).ParseFiles(bookTemplatePath + bookDetailHTMLName)
		errMsg = append(errMsg, message.ErrMsgDelErr)
		// 画面に表示するデータを格納
		targetBook, _ := bookdao.GetBookByID(id)
		responseData := BookDetailResponseData{Book: targetBook, ErrMsg: errMsg}
		err := Tpl.ExecuteTemplate(w, bookDetailHTMLName, responseData)
		if err != nil {
			log.Fatal("【main.go bookDeleteHandler】画面の描画中にエラーが発生しました。")
		}
		// 削除失敗時は本詳細画面へ遷移
	} else {
		url := appconst.BookURL + "?sucDelFlg=1"
		http.Redirect(w, r, url, http.StatusFound)
	}
}