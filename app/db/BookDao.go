package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Id                     int
	Title                  string
	Author                 string
	Latest_Issue           float64
	Front_Cover_Image_Path string
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

/*
	DB内のすべての本を取得
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
		err = rows.Scan(&book.Id, &book.Title, &book.Author, &book.Latest_Issue, &book.Front_Cover_Image_Path)
		errCheck(err)
		books = append(books, book)
	}
	return books
}


/*
	本を1冊DBに挿入する
	input:Book
	output:error
*/
func InsertBook(book Book) error{
	db := dbSetUp()
	defer db.Close() // 関数がリターンする直前に呼び出される
	var err error

	if(book.Front_Cover_Image_Path == ""){
		ins, err := db.Prepare("INSERT INTO book (title,author,latest_issue) VALUES(?,?,?)")
		errCheck(err)
		// Bookを格納する
		_, err = ins.Exec(&book.Title, &book.Author, &book.Latest_Issue)
	}else{
		ins, err := db.Prepare("INSERT INTO book (title,author,latest_issue,front_cover_image_path) VALUES(?,?,?,?)")
		errCheck(err)
		// Bookを格納する
		_, err = ins.Exec(&book.Title, &book.Author, &book.Latest_Issue, &book.Front_Cover_Image_Path)
	}
	
	return err
}

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
		err = rows.Scan(&book.Id, &book.Title, &book.Author, &book.Latest_Issue, &book.Front_Cover_Image_Path)
		errCheck(err)
		books = append(books, book)
	}
	return books
}
