package appstructure

import (
	"github.com/docker_go_nginx/app/db/bookdao"
)

// 画面用のデータ構造
type ResponseData struct {
	Books     []bookdao.Book
	ViewData  map[string]string
	Message   map[string][]string
	LoginFlag bool
}

/**
レスポンスデータの構造体をセットした状態で返す
*/
func CreateResponseData(viewData map[string]string, message map[string][]string) ResponseData {
	responseData := ResponseData{
		Books:     []bookdao.Book{},
		ViewData:  viewData,
		Message:   message,
		LoginFlag: false,
	}
	return responseData
}

/**
レスポンスデータの構造体をセットした状態で返す
*/
func CreateResponseDataSetBook(books []bookdao.Book) ResponseData {
	responseData := ResponseData{
		Books:     books,
		ViewData:  map[string]string{},
		Message:   map[string][]string{},
		LoginFlag: false,
	}
	return responseData
}
