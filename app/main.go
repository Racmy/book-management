package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"os"
	"io"
	"mime/multipart"
	"github.com/docker_go_nginx/app/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	ROOT         string = "/"
	ID           string = "Id"
	TITLE        string = "Title"
	AUTHOR       string = "Author"
	LATEST_ISSUE string = "Latest_Issue"
	IMGPATH      string = "Front_Cover_Image_Path"
)

const (
	ERMSG_TITLE_NULL string = "タイトルを入力してください"
	ERMSG_AUTH_NULL  string = "著者を入力してください"
	ERMSG_LI_NULL    string = "最新所持巻数を数字で入力してください"
)

const (
	SUCMSG_UPDATE string = "更新が完了しました"
	FRONT_COVER_IMAGE string = "Front_Cover_Image"
	IMG_PATH string = "/static/img/"
)

type RegistResultValue struct {
	Title        string
	Author       string
	Latest_Issue float64
}

type RegistValue struct {
	Title        string
	Author       string
	Latest_Issue float64
	ErrString    []string
}

/*
	本詳細画面用のレスポンデータ
*/
type ResponseDataForDetail struct {
	Book   db.Book
	ErrMsg []string
	SucMsg []string
}

/*
	一覧画面用のレスポンスデータ
*/
type ResponseData struct {
	Keyword string
	Books   []db.Book
}

/*
	本を登録画面へのハンドラ
*/
func bookRegistHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("./template/bookRegist.html"))

	errString := []string{}
	tmpTitle := r.FormValue(TITLE)
	tmpAuthor := r.FormValue(AUTHOR)
	tmpLatest_Issue_String := r.FormValue(LATEST_ISSUE)
	tmpLatest_Issue, strConvErr := strconv.ParseFloat(tmpLatest_Issue_String, 64)
	tmpErrCheckFlag := r.FormValue("ErrCheckFlag")

	if tmpErrCheckFlag == "1" {
		if tmpTitle == "" {
			errString = append(errString, ERMSG_TITLE_NULL)
		}
		if tmpAuthor == "" {
			errString = append(errString, ERMSG_AUTH_NULL)
		}
		if strConvErr != nil {
			errString = append(errString, ERMSG_LI_NULL)
		}
	}

	if strConvErr != nil {
		tmpLatest_Issue = 1
	}

	tmp := RegistValue{
		Title:        tmpTitle,
		Author:       tmpAuthor,
		Latest_Issue: tmpLatest_Issue,
		ErrString:    errString,
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
	tmpFront_Cover_Image_Path := ""
	//表紙画像を格納する変数宣言
	var file multipart.File
	var fileHeader *multipart.FileHeader
	// POSTされたファイルデータをメモリに格納
	//33554432 約30MByte(8Kのping形式には耐えられない)
	err := r.ParseMultipartForm(32 << 20) 
    if err != nil {
        log.Println("not ParseMultipartForm")
        fileUploadFlag = false
    }else{
		file , fileHeader , err = r.FormFile (FRONT_COVER_IMAGE)
		if (err != nil) {
			log.Println("not file upload")
			fileUploadFlag = false
		}
	}
	//表紙画像がuploadされている時
	if(fileUploadFlag){
		tmpFront_Cover_Image_Path = IMG_PATH + fileHeader.Filename

		log.Println(tmpFront_Cover_Image_Path)
		// サーバー側に保存するために空ファイルを作成
		var saveImage *os.File
		saveImage, err = os.Create(tmpFront_Cover_Image_Path)
		if (err != nil) {
			log.Println(err)
		}
		defer saveImage.Close()
		defer file.Close()
		size, err := io.Copy(saveImage, file)
		if (err != nil) {
			log.Println(err)
		}
		log.Println(size)
		tmpFront_Cover_Image_Path = "/" + tmpFront_Cover_Image_Path
	}
	

	r.ParseForm()

	tmpTitle := r.Form[TITLE][0]
	tmpAuthor := r.Form[AUTHOR][0]
	tmpLatest_Issue_String := r.Form[LATEST_ISSUE][0]
	tmpLatest_Issue, strConvErr := strconv.ParseFloat(tmpLatest_Issue_String, 64)

	if (tmpTitle == "") || (tmpAuthor == "") || (strConvErr != nil) {
		var url = "/regist"
		url += "?Title=" + tmpTitle + "&Author=" + tmpAuthor + "&Latest_Issue=" + tmpLatest_Issue_String + "&ErrCheckFlag=1"
		http.Redirect(w, r, url, http.StatusFound)
	}

	insertBook := db.Book{
		Title: tmpTitle,
		Author: tmpAuthor,
		Latest_Issue: tmpLatest_Issue,
		Front_Cover_Image_Path: tmpFront_Cover_Image_Path,
	}

	dbErr := db.InsertBook(insertBook)
	if dbErr != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	
	// テンプレートにデータを埋め込む
	if err := tmpl.ExecuteTemplate(w, "bookRegistResult.html", insertBook); err != nil {
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
		responseData.Book = db.GetBookById(id)
		//更新成功時のメッセージを格納
		if query.Get("sucFlg") != "" {
			responseData.SucMsg = append(responseData.SucMsg, SUCMSG_UPDATE)
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
	id_int, strConvErr := strconv.Atoi(id)
	latest_Issue_String := r.FormValue(LATEST_ISSUE)
	latest_Issue, strConvErr := strconv.ParseFloat(latest_Issue_String, 64)

	var url = "/detail"
	var errMsg []string
	/*エラーチェック【相談】登録との共通化*/
	if title == "" {
		errMsg = append(errMsg, ERMSG_TITLE_NULL)
	}
	if author == "" {
		errMsg = append(errMsg, ERMSG_AUTH_NULL)
	}
	if strConvErr != nil {
		errMsg = append(errMsg, ERMSG_LI_NULL)
	}

	// 入力エラーがない場合は更新処理を実施
	if len(errMsg) == 0 {
		// 入力データで更新
		updateBook := db.Book{Id: id_int, Title: title, Author: author, Latest_Issue: latest_Issue}
		db.UpdateBook(updateBook)
		// 成功したことをDetailに伝えるためにsucFlgをつける
		url += "?Id=" + id + "&sucFlg=1"
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		var tmpl = template.Must(template.ParseFiles("./template/detail.html"))
		// エラー時は、画面から送られてきたデータを渡す
		inputBook := db.Book{Id: id_int, Title: title, Author: author, Latest_Issue: latest_Issue, Front_Cover_Image_Path: imgPath}

		// 画面に表示するデータを格納
		responseData := ResponseDataForDetail{Book: inputBook, ErrMsg: errMsg}

		if err := tmpl.ExecuteTemplate(w, "detail.html", responseData); err != nil {
			log.Fatal(err)
		}
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
