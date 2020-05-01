package bookdao

import (
	"database/sql"
	"log"
	"strconv"
	"time"
)

/**
HTMLのフォームのnameに与える名前
*/
const (
	ID          string = "Id"
	TITLE       string = "Title"
	AUTHOR      string = "Author"
	LatestIssue string = "LatestIssue"
	IMGPATH     string = "FrontCoverImagePath"
	NewIMGPATH  string = "NewFrontCoverImagePath"
)

// Book D層とP層で本の情報を受け渡す構造体
type Book struct {
	ID                  int
	UserID				int
	Title               string
	Author              string
	LatestIssue         float64
	FrontCoverImagePath string
	active				string
	Created_at			time.Time
	Update_at			time.Time
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
	db, err := sql.Open("mysql", "racmy:racmy@tcp(db:3306)/book-management?parseTime=true")
	errCheck(err)
	return db
}

//GetBookByID ...BookテーブルのIDに紐つく情報を1件取得
/*
@param id string
@return book Book
@return err error
*/
func GetBookByID(id string) (Book, error) {
	db := dbSetUp()
	defer db.Close()
	rows, err := db.Query("SELECT id, user_id, title, author, latest_issue, front_cover_image_path,active,created_at, update_at FROM book WHERE Id = ?", id)
	defer rows.Close()

	// SELECT失敗時にbookがerrorでのリターン
	if err != nil {
		log.Print("【BookDao.GetBookByID】id = " + id + "not exist in book table.")
		return Book{}, err
	}

	//　本が検索できた場合は、本の情報を含めてリターン
	var book Book
	if rows.Next() {
		err = rows.Scan(&book.ID, &book.UserID , &book.Title, &book.Author, &book.LatestIssue, &book.FrontCoverImagePath, &book.active, &book.Created_at, &book.Update_at)
		errCheck(err)
		return book, err
	}

	// 検索したが「０件」の場合は、book・errが共に空
	return book, err
}
//GetBookByIDAndUserID ...BookテーブルのIDに紐つく情報を1件取得
/*
@param id string
@param userid int
@return book Book
@return err error
*/
func GetBookByIDAndUserID(id string, userID int) (Book, error) {
	db := dbSetUp()
	defer db.Close()
	rows, err := db.Query("SELECT id, user_id, title, author, latest_issue, front_cover_image_path,active,created_at, update_at FROM book WHERE Id = ? AND user_id = ?", id,userID)
	defer rows.Close()

	// SELECT失敗時にbookがerrorでのリターン
	if err != nil {
		log.Print("【BookDao.GetBookByID】id = " + id + "not exist in book table.")
		return Book{}, err
	}

	//　本が検索できた場合は、本の情報を含めてリターン
	var book Book
	if rows.Next() {
		err = rows.Scan(&book.ID, &book.UserID , &book.Title, &book.Author, &book.LatestIssue, &book.FrontCoverImagePath, &book.active, &book.Created_at, &book.Update_at)
		errCheck(err)
		return book, err
	}

	// 検索したが「０件」の場合は、book・errが共に空
	return book, err
}
//GetAllBooks ...DB内のすべての本を取得
/*
input:
output:[]Book
*/
func GetAllBooks() []Book {
	db := dbSetUp()
	defer db.Close() // 関数がリターンする直前に呼び出される
	rows, err := db.Query("SELECT id, user_id, title, author, latest_issue, front_cover_image_path,active,created_at, update_at FROM book")
	errCheck(err)
	// Bookを格納するArray作成
	var books = []Book{}
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.ID, &book.UserID , &book.Title, &book.Author, &book.LatestIssue, &book.FrontCoverImagePath, &book.active, &book.Created_at, &book.Update_at)
		errCheck(err)
		books = append(books, book)
	}
	return books
}

//GetAllBooks ...DB内のユーザの本をすべて取得
/*
input:
output:[]Book
*/
func GetAllBooksByUserID(user_id int) []Book {
	db := dbSetUp()
	defer db.Close() // 関数がリターンする直前に呼び出される
	rows, err := db.Query("SELECT id, user_id, title, author, latest_issue, front_cover_image_path,active,created_at, update_at FROM book WHERE user_id = ?",user_id)
	errCheck(err)
	// Bookを格納するArray作成
	var books = []Book{}
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.ID, &book.UserID , &book.Title, &book.Author, &book.LatestIssue, &book.FrontCoverImagePath, &book.active, &book.Created_at, &book.Update_at)
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
		ins, err := db.Prepare("INSERT INTO book (user_id,title,author,latest_issue) VALUES(?,?,?,?)")
		errCheck(err)
		// Bookを格納する
		result, err = ins.Exec(&book.UserID, &book.Title,&book.Author, &book.LatestIssue)
		errCheck(err)
	} else {
		ins, err := db.Prepare("INSERT INTO book (user_id,title,author,latest_issue,front_cover_image_path) VALUES(?,?,?,?,?)")
		errCheck(err)
		// Bookを格納する
		result, err = ins.Exec(&book.UserID, &book.Title,&book.Author, &book.LatestIssue, &book.FrontCoverImagePath)
		errCheck(err)
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
	rows, _ := db.Query("SELECT id, user_id, title, author, latest_issue, front_cover_image_path,active,created_at, update_at FROM book WHERE title LIKE ? OR author LIKE ?", keyword, keyword)
	// errCheck(err)
	// Bookを格納するArray作成
	var books = []Book{}
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.UserID , &book.Title, &book.Author, &book.LatestIssue, &book.FrontCoverImagePath, &book.active, &book.Created_at, &book.Update_at)
		errCheck(err)
		books = append(books, book)
	}
	return books
}

// GetSearchedBooks ...キーワードから本情報を取得する
/*
	本をキーワードで検索する
	input:keyword string
	input:keyword int
	output:[]Book
*/
func GetSearchedBooksByKeywordAndUserID(keyword string, userID int) []Book {
	keyword = "%" + keyword + "%"
	db := dbSetUp()
	defer db.Close()
	rows, _ := db.Query("SELECT id, user_id, title, author, latest_issue, front_cover_image_path,active,created_at, update_at FROM book WHERE user_id = ? AND (title LIKE ? OR author LIKE ?)", userID, keyword, keyword)
	// errCheck(err)
	// Bookを格納するArray作成
	var books = []Book{}
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.UserID , &book.Title, &book.Author, &book.LatestIssue, &book.FrontCoverImagePath, &book.active, &book.Created_at, &book.Update_at)
		errCheck(err)
		books = append(books, book)
	}
	return books
}

// UpdateBook ...本情報の更新処理を行う
/*
	本の更新
	input:book Book
	input:userID int
	output:bookid int, err error
*/
func UpdateBook(book Book,userID int) (int, error) {
	db := dbSetUp()
	defer db.Close()

	var err error = nil

	// DBに存在する確認する
	_, err = GetBookByIDAndUserID(strconv.Itoa(book.ID), userID)

	// 存在する場合は更新する
	if err == nil {
		upd, err := db.Prepare("UPDATE book SET title = ?, author = ?, latest_issue = ?, front_cover_image_path = ? WHERE id = ?")
		errCheck(err)
		_, err = upd.Exec(&book.Title, &book.Author, &book.LatestIssue, &book.FrontCoverImagePath, &book.ID)

		// 更新失敗時のエラー
		if err != nil {
			log.Print("【BookDao.UpdateBook】id = " + strconv.Itoa(book.ID) + "update error.")
			return -1, err
		}

		// 更新成功のため、IDとnilのerrを返す
		return book.ID, err
	}

	// IDを取得できなかったログ出力
	log.Print("【BookDao.UpdateBook】id = " + strconv.Itoa(book.ID) + "not exist in book table.")

	//　エラーが発生した場合は、存在しないIDである「-１」を返す
	return -1, err
}

// DeleteBookByID ... 本の削除を行う
/*
	本の削除機能
	@param id 本のID
	@param userId ユーザのID
	@return error 削除失敗時：nil
*/
func DeleteBookByIDAndUserID(id string, userID int) error {
	db := dbSetUp()
	defer db.Close()
	_, err := db.Query("DELETE FROM book WHERE id = ? and user_id = ?", id, userID)

	if err != nil {
		log.Print("【BookDao.DeleteBookByID】" + id + "error can't delete")
	}

	return nil
}
