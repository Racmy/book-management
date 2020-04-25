package uDB

import (
	"database/sql"
	"log"
)
/**
	エラーチェック
*/
func ErrCheck(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

/*
	DBの初期化
	input:
	output:*sql.DB
*/
func DbSetUp() *sql.DB {
	db, err := sql.Open("mysql", "racmy:racmy@tcp(db:3306)/book-management?parseTime=true")
	ErrCheck(err)
	return db
}
