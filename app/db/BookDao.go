package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Id                     int
	Title                  string
	Author                 string
	Latest_Issue           float32
	Front_Cover_Image_Path string
}

func errCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func dbSetUp() *sql.DB {
	db, err := sql.Open("mysql", "racmy:racmy@tcp(db:3306)/book-management")
	errCheck(err)
	return db
}

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
