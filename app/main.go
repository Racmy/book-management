package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/docker_go_nginx/app/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)


const(
	ROOT string = "/"
	TITLE string = "Title"
	AUTHOR string = "Author"
	LATEST_ISSUE string = "Latest_Issue"
)

type RegistResultValue struct {
	Title string
	Author string
	Latest_Issue float64
}

type RegistValue struct {
	Title        string
	Author       string
	Latest_Issue float64
	ErrString []string
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
	tmpLatest_Issue , strConvErr := strconv.ParseFloat(tmpLatest_Issue_String,64)
	tmpErrCheckFlag := r.FormValue("ErrCheckFlag")

	//登録ボタン押下時エラーチェックに引っかかった時のメッセージ作成
	if(tmpErrCheckFlag == "1"){
		if(tmpTitle == ""){
			errString = append(errString,"タイトルを入力してください")	
		}
		if(tmpAuthor == ""){
			errString = append(errString,"著者を入力してください")
		}
		if(strConvErr != nil){
			errString = append(errString,"最新所持巻数を数字で入力してください")
		}
	}

	if strConvErr != nil {
		tmpLatest_Issue = 1
	}

	tmp := RegistValue{
		Title:        tmpTitle,
		Author:       tmpAuthor,
		Latest_Issue: tmpLatest_Issue,
		ErrString: errString,
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
	r.ParseForm()

	tmpTitle := r.Form[TITLE][0]
	tmpAuthor := r.Form[AUTHOR][0]
	tmpLatest_Issue_String := r.Form[LATEST_ISSUE][0]
	tmpLatest_Issue , strConvErr := strconv.ParseFloat(tmpLatest_Issue_String,64)

	if (tmpTitle == "") || (tmpAuthor == "") || (strConvErr != nil) {
		var url = "/regist"
		url += "?Title=" + tmpTitle + "&Author=" + tmpAuthor + "&Latest_Issue=" + tmpLatest_Issue_String + "&ErrCheckFlag=1"
		http.Redirect(w,r,url,http.StatusFound)
	}

	insertBook := db.Book{Title: tmpTitle, Author: tmpAuthor, Latest_Issue: tmpLatest_Issue}

	dbErr := db.InsertBook(insertBook)
	if dbErr != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	// テンプレートに埋め込むデータ作成
	dat := RegistResultValue{
		Title: tmpTitle,
		Author: tmpAuthor,
		Latest_Issue: tmpLatest_Issue,
	}

	// テンプレートにデータを埋め込む
	if err := tmpl.ExecuteTemplate(w, "bookRegistResult.html", dat); err != nil {
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
	http.Handle(ROOT+"node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules/"))))
	http.Handle(ROOT, r)
	http.ListenAndServe(":3000", nil)
}
