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
	TITLE        string = "Title"
	AUTHOR       string = "Author"
	LATEST_ISSUE string = "Latest_Issue"
	FRONT_COVER_IMAGE string = "Front_Cover_Image"
	IMG_PATH string = "static/img/"
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
	レスポンスデータ
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

	//登録ボタン押下時エラーチェックに引っかかった時のメッセージ作成
	if tmpErrCheckFlag == "1" {
		if tmpTitle == "" {
			errString = append(errString, "タイトルを入力してください")
		}
		if tmpAuthor == "" {
			errString = append(errString, "著者を入力してください")
		}
		if strConvErr != nil {
			errString = append(errString, "最新所持巻数を数字で入力してください")
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

	if(tmpFront_Cover_Image_Path != ""){
		tmpFront_Cover_Image_Path = "/" + tmpFront_Cover_Image_Path
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
	http.Redirect(w, r, "/", http.StatusFound)
}

/*
	本の検索のためのハンドラ
*/
func bookSearchHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	if keyword := query.Get("keyword"); query.Get("keyword") != "" {
		// keywordがnullの場合は、HOMEへリダイレクト
		if keyword == "" {
			http.Redirect(w, r, "/", http.StatusFound)
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
		http.Redirect(w, r, "/", http.StatusFound)
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
	// cssフレームワーク読み込み
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules/"))))
	// 画像フォルダ
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.Handle(ROOT, r)
	http.ListenAndServe(":3000", nil)
}
