package File

import (
	"log"
	"mime/multipart"
)

const fileUploadPath = "static/img/"

/***
	ファイルアップロード関数
	input:
		file mulipart.File	: upload　file
		file
		fileName			: saved file name
*/
func fileUpload(file multipart.File, fileName string){
	
}