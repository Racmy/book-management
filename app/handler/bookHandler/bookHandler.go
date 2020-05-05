package bookhandler

import (
	"github.com/docker_go_nginx/app/common/appconst"
	"github.com/docker_go_nginx/app/common/appstructure"
	"github.com/docker_go_nginx/app/common/message"
	"github.com/docker_go_nginx/app/db/bookdao"
	"github.com/docker_go_nginx/app/utility/uDB"
	"github.com/docker_go_nginx/app/utility/ufile"
	"github.com/docker_go_nginx/app/utility/ulogin"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"text/template"
)

var Tpl *template.Template

var rootTemplatePath = "./template/"
var bookTemplatePath = rootTemplatePath + "book/"
var bookListHTMLName = "bookList.html"
var bookDetailHTMLName = "bookDetail.html"
var bookRegistHTMLName = "bookRegist.html"
var bookRegistResultHTMLName = "bookRegistResult.html"

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
	Tpl.New(bookRegistHTMLName).Option("missingkey=zero").ParseFiles(bookTemplatePath + bookRegistHTMLName)

	//　セッションから画面表示データを取得
	responseData := ulogin.GetViewDataAndMessage(w, r)

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
	frontCoverImagePath := ""

	//表紙画像がuploadされている時
	userId, err := ulogin.GetLoginUserId(r)
	if err != nil {
		http.Redirect(w, r, appconst.RootURL, http.StatusFound)
	}

	// ファイルオブジェクトを取得
	file, fileHeader, err := ufile.ParseFile(r, bookdao.IMGPATH)

	if file != nil {
		frontCoverImagePath, err = ufile.UserFileUpload(file, fileHeader.Filename, userId)
	}

	r.ParseForm()

	title := r.Form[bookdao.TITLE][0]
	author := r.Form[bookdao.AUTHOR][0]
	latestIssueString := r.Form[bookdao.LatestIssue][0]
	latestIssue, strConvErr := strconv.ParseFloat(latestIssueString, 64)

	titleMsg := []string{}
	authorMsg := []string{}
	latestIssueMsg := []string{}

	// 入力チェック
	if title == "" {
		titleMsg = append(titleMsg, "タイトルを入力してください。")
	}
	if author == "" {
		authorMsg = append(authorMsg, "著者を入力してください。")
	}
	if strConvErr != nil {
		latestIssueMsg = append(latestIssueMsg, "巻数の形式に誤りがあります。")
	}

	errMsgMap := map[string][]string{}
	viewData := map[string]string{}
	if len(titleMsg) > 0 || len(authorMsg) > 0 || len(latestIssueMsg) > 0 {
		errMsgMap["title"] = titleMsg
		errMsgMap["author"] = authorMsg
		errMsgMap["latestIssue"] = latestIssueMsg

		// 画面データ格納
		viewData["title"] = title
		viewData["author"] = author
		viewData["latestIssue"] = latestIssueString

		// セッションにエラーメッセージと画面データをつめる
		session, _ := ulogin.GetSession(r)
		session.AddFlash(errMsgMap, appconst.SessionMsg)
		session.AddFlash(viewData, appconst.SessionViewData)
		session.Save(r, w)

		// 本登録画面へ遷移
		http.Redirect(w, r, appconst.BookRegistURL, http.StatusFound)
	}

	userId, getLoginUserErr := ulogin.GetLoginUserId(r)
	if getLoginUserErr != nil {
		http.Redirect(w, r, appconst.RootURL, http.StatusFound)
	}
	insertBook := bookdao.Book{
		UserID:              userId,
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

	http.Redirect(w, r, appconst.BookRegistResultURL+"?Id="+idString, http.StatusFound)
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
	本棚画面へのハンドラ
*/
func BookListHandler(w http.ResponseWriter, r *http.Request) {
	// 未ログインの場合はホームへリダイレクト
	if !ulogin.IsLogined(r) {
		http.Redirect(w, r, appconst.RootURL, http.StatusFound)
	}

	Tpl.New(bookListHTMLName).ParseFiles(bookTemplatePath + bookListHTMLName)

	userId, err := ulogin.GetLoginUserId(r)
	if err != nil {
		http.Redirect(w, r, appconst.RootURL, http.StatusFound)
	}
	var responseData BookListResponseData
	responseData.Books = bookdao.GetAllBooksByUserID(userId)

	if ulogin.GetSessionFlg(w, r) {
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
	// 更新失敗のハンドリング後の場合はそのデータを利用する
	responseData := ulogin.GetViewDataAndMessage(w, r)

	// 初回ではない場合はセッションのデータを表示
	// 更新失敗→本詳細の遷移
	if len(responseData.ViewData) != 0 {
		log.Println("piyo")
		log.Print(responseData)
		Tpl.New(bookDetailHTMLName).ParseFiles(bookTemplatePath + bookDetailHTMLName)
		if err := Tpl.ExecuteTemplate(w, bookDetailHTMLName, responseData); err != nil {
			log.Fatal(err)
		}
	} else {
		// クエリの取得
		query := r.URL.Query()
		//　初回時
		if id := query.Get("Id"); query.Get("Id") != "" {
			log.Print("fuga")
			// 画面からIdを取得し、DBから紐つくデータを取得
			book, err := bookdao.GetBookByID(id)
			// データ取得失敗時はホームへ戻す
			if err != nil {
				http.Redirect(w, r, appconst.BookURL, http.StatusFound)
			}

			viewData := map[string]string{}
			viewData["Id"] = strconv.Itoa(book.ID)
			viewData["title"] = book.Title
			viewData["author"] = book.Author
			viewData["latestIssue"] = strconv.FormatFloat(book.LatestIssue, 'f', 1, 64)
			viewData["imagePath"] = book.FrontCoverImagePath

			/**
			更新完了のメッセージがある場合は取得する
			*/
			var message map[string][]string
			if len(responseData.Message["success"]) > 0 {
				message = responseData.Message
			}

			responseData = appstructure.CreateResponseData(viewData, message)
			log.Println(responseData)
			Tpl.New(bookDetailHTMLName).ParseFiles(bookTemplatePath + bookDetailHTMLName)
			if err := Tpl.ExecuteTemplate(w, bookDetailHTMLName, responseData); err != nil {
				log.Fatal(err)
			}
		}

		// クエリにIDがない場合は本一覧にリダイレクト
		http.Redirect(w, r, appconst.BookURL, http.StatusFound)
	}
}

/*
	本の検索のためのハンドラ
*/
func BookSearchHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := ulogin.GetLoginUserId(r)
	if err != nil {
		http.Redirect(w, r, appconst.RootURL, http.StatusFound)
	}

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
		responseData.Books = bookdao.GetSearchedBooksByKeywordAndUserID(keyword, userId)

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

	titleMsg := []string{}
	authorMsg := []string{}
	latestIssueMsg := []string{}

	// 入力チェック
	if title == "" {
		titleMsg = append(titleMsg, "タイトルを入力してください。")
	}
	if author == "" {
		authorMsg = append(authorMsg, "著者を入力してください。")
	}
	if strConvErr != nil {
		latestIssueMsg = append(latestIssueMsg, "巻数の形式に誤りがあります。")
	}

	errMsgMap := map[string][]string{}
	viewData := map[string]string{}
	if len(titleMsg) > 0 || len(authorMsg) > 0 || len(latestIssueMsg) > 0 {
		errMsgMap["title"] = titleMsg
		errMsgMap["author"] = authorMsg
		errMsgMap["latestIssue"] = latestIssueMsg

		// 画面データ格納
		log.Print("hohoho")
		log.Println(id)
		viewData["Id"] = id
		viewData["title"] = title
		viewData["author"] = author
		viewData["latestIssue"] = latestIssueString
		viewData["imagePath"] = imgPath

		// セッションにエラーメッセージと画面データをつめる
		session, _ := ulogin.GetSession(r)
		session.AddFlash(errMsgMap, appconst.SessionMsg)
		session.AddFlash(viewData, appconst.SessionViewData)
		session.Save(r, w)

		// 本登録画面へ遷移
		url := appconst.BookDetailLURL + "?Id=" + id
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		// 入力エラーがない場合は更新処理を実施
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

		// ログインユーザの取得
		userId, err := ulogin.GetLoginUserId(r)
		if err != nil {
			http.Redirect(w, r, appconst.RootURL, http.StatusFound)
		}

		log.Println("hoge")
		log.Println(updateBook.ID)

		// 本の更新
		updatedBookId, errUpd := bookdao.UpdateBook(updateBook, userId)
		log.Println(updatedBookId)
		// 本の更新に失敗チェック
		uDB.ErrCheck(errUpd)
		idString := strconv.Itoa(updatedBookId)
		log.Print(idString)
		// 更新処理が失敗していない場合は、詳細画面へ遷移（bookDetail.html）
		var url string
		log.Print("【main.go　UpdateBookHander】success update")
		url = appconst.BookDetailLURL + "?Id=" + idString

		message := map[string][]string{}
		sucMessage := []string{}
		sucMessage = append(sucMessage, "本の更新が完了しました。")
		// 成功したセッションに成功フラグを立てる
		message["success"] = sucMessage
		// セッションにエラーメッセージと画面データをつめる
		session, _ := ulogin.GetSession(r)
		session.AddFlash(message, appconst.SessionMsg)
		session.Save(r, w)

		http.Redirect(w, r, url, http.StatusFound)
	}

}

/*
	本の削除ハンドラ
	/detail →　/book
*/
func BookDeleteHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := ulogin.GetLoginUserId(r)
	if err != nil {
		http.Redirect(w, r, appconst.RootURL, http.StatusFound)
	}

	r.ParseForm()
	id := r.FormValue(bookdao.ID)
	err = bookdao.DeleteBookByIdAndUserId(id, userId)

	// 削除失敗時は詳細画面へ遷移してエラーメッセージを表示
	if err != nil {
		// エラーの表示
		log.Fatal("【main.go bookDeleteHandler】本の削除に失敗しました。")
		//TODO削除失敗メッセージを出す
		http.Redirect(w, r, appconst.BookURL, http.StatusFound)
	} else {
		// セッションフラグをONにする
		ulogin.SetSessionFlg(w, r)
		http.Redirect(w, r, appconst.BookURL, http.StatusFound)
	}
}
