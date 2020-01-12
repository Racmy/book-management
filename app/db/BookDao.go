package db

import (
	"database/sql"
	"strconv"
)

// Book D層とP層で本の情報を受け渡す構造体
type Book struct {
	ID                  int
	Title               string
	Author              string
	LatestIssue         float64
	FrontCoverImagePath string
}

func errCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

/*
	DBの初期化
	input:
	output:*sql.DB
*/
func dbSetUp() *sql.DB {
	db, err := sql.Open("mysql", "racmy:racmy@tcp(db:3306)/book-management")
	errCheck(err)
	return db
}

//GetBookByID ...BookテーブルのIDに紐つく情報を1件取得
/*
@param id string
@return book Book
@return canGet bool
*/
func GetBookByID(id string) (Book, bool) {
	db := dbSetUp()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM book WHERE Id = ?", id)
	defer rows.Close()

	// SELECT失敗時にbookがerrorでのリターン
	if err != nil {
		return Book{}, false
	}

	//　本が検索できた場合は、本の情報を含めてリターン
	var book Book
	if rows.Next() {
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.LatestIssue, &book.FrontCoverImagePath)
		errCheck(err)
		return book, true
	}
	// 検索したが「０件」の場合は、book・errが共に空
	return book, false
}

//GetAllBooks ...DB内のすべての本を取得
/*
input:
output:[]Book
*/
func GetAllBooks() []Book {
	db := dbSetUp()
	defer db.Close() // 関数がリターンする直前に呼び出される
	rows, err := db.Query("SELECT * FROM book")
	errCheck(err)
	// Bookを格納するArray作成
	var books = []Book{}
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.LatestIssue, &book.FrontCoverImagePath)
		errCheck(err)
		books = append(books, book)
	}
	return books
}

//InsertBook ...
/*
	本を1冊DBに挿入する
	input:Book
	output:error
*/
func InsertBook(book Book) (int64, error) {
	db := dbSetUp()
	defer db.Close() // 関数がリターンする直前に呼び出される
	var result sql.Result
	if book.FrontCoverImagePath == "" {
		ins, err := db.Prepare("INSERT INTO book (title,author,latest_issue) VALUES(?,?,?)")
		errCheck(err)
		// Bookを格納する
		result, err = ins.Exec(&book.Title, &book.Author, &book.LatestIssue)
	} else {
		ins, err := db.Prepare("INSERT INTO book (title,author,latest_issue,front_cover_image_path) VALUES(?,?,?,?)")
		errCheck(err)
		// Bookを格納する
		result, err = ins.Exec(&book.Title, &book.Author, &book.LatestIssue, &book.FrontCoverImagePath)

	}
	// Insertした結果を返す（id, error）
	return result.LastInsertId()
}

// GetSearchedBooks ...キーワードから本情報を取得する
/*
	本をキーワードで検索する
	input:keyword string
	output:[]Book
*/
func GetSearchedBooks(keyword string) []Book {
	keyword = "%" + keyword + "%"
	db := dbSetUp()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM book WHERE title LIKE ? OR author LIKE ?", keyword, keyword)
	errCheck(err)
	// Bookを格納するArray作成
	var books = []Book{}
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.LatestIssue, &book.FrontCoverImagePath)
		errCheck(err)
		books = append(books, book)
	}
	return books
}

// UpdateBook ...本情報の更新処理を行う
/*
	本の更新
*/
func UpdateBook(book Book) int {
	db := dbSetUp()
	defer db.Close()
	_, canGet := GetBookByID(strconv.Itoa(book.ID))
	GetBookByID(strconv.Itoa(book.ID))
	var err error

	if canGet {
		upd, err := db.Prepare("UPDATE book SET title = ?, author = ?, latest_issue = ? WHERE id = ?")
		errCheck(err)
		_, err = upd.Exec(&book.Title, &book.Author, &book.LatestIssue, &book.ID)
	}

	// 更新に成功した場合は、更新対象のIDを返す
	if err == nil {
		return book.ID
	}
	//　エラーが発生した場合は、存在しないIDである「-１」を返す
	return -1
}
