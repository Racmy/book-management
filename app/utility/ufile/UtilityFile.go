package ufile

import (
	"github.com/docker_go_nginx/app/utility/uDB"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultFileUploadPath = "static/img/"

const TimeFormat = "2006-01-02_15-04-05"

/*
	ファイルをサーバへアップロードする関数

	@param 	file 				mulipart.File
	@param 	filePath 			string
	@param 	fileName 			string
	@return fileName(with path) string
	@return err					error
*/
func FileUpload(file multipart.File, filePath string, fileName string) (string, error) {
	savedFilePath := filePath
	//日本時間用*Location
	Jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowTime := time.Now().In(Jst).Format(TimeFormat)
	log.Println(nowTime)
	savedFileName := savedFilePath + nowTime + "_" + fileName

	err := fileSaved(file, savedFileName)
	if err != nil {
		return "", err
	}
	return "/" + savedFileName, err
}

/*
	"static/img/"配下にファイルをサーバへアップロードする関数
	@param file 	mulipart.File
	@param fileName string
	@param err		err
*/
func DefaultFileUpload(file multipart.File, fileName string) (string, error) {
	return FileUpload(file, defaultFileUploadPath, fileName)
}

/*
	"static/img/{userid}"配下にファイルをサーバへアップロードする関数
	@param file 	mulipart.File
	@param fileName string
	@param err		err
*/
func UserFileUpload(file multipart.File, fileName string, userId int) (string, error) {
	path := defaultFileUploadPath + "user/" + strconv.Itoa(userId) + "/"

	// フォルダが存在していない場合はフォルダを作成する
	if !Exists(path) {
		MakeDir(path)
	}
	return FileUpload(file, path, fileName)
}

/*
	ファイルをサーバへ保存する関数
	@param file 	mulipart.File
	@param fileName string
	@return err		err
*/
func fileSaved(file multipart.File, fileName string) error {
	log.Println("saved file name : " + fileName)

	// サーバー側に保存するために空ファイルを作成
	saveImage, err := os.Create(fileName)
	if err != nil {
		log.Println("【File.go fileSaved】os.Create Error")
		log.Println(err)
		return err
	}
	//ファイルクローズ予約
	defer saveImage.Close()
	defer file.Close()
	//目的ファイル保存
	size, err := io.Copy(saveImage, file)
	if err != nil {
		log.Println("【File.go fileSaved】io.Copy Error")
		log.Println(err)
		return err
	}
	log.Printf("File Saved data size:" + strconv.FormatInt(size, 10))
	return nil
}

/*
	画像がRequestに格納されているか判定する
	メモリに画像を格納する
	@param 	request http.Request
	@param  name    string
	@return file    multipart.File
	@return fileHeader multipart.FileHeader
	@return err		error
*/
func ParseFile(request *http.Request, name string) (multipart.File, *multipart.FileHeader, error) {
	var file multipart.File
	var fileHeader *multipart.FileHeader
	file = nil
	fileHeader = nil
	// POSTされたファイルデータをメモリに格納
	//33554432 約30MByte(8Kのping形式には耐えられない)
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println("not ParseMultipartForm")
	} else {
		file, fileHeader, err = request.FormFile(name)
		if err != nil {
			log.Println("not file upload")
		}
	}

	return file, fileHeader, err
}

// ファイル・フォルダの存在チェック
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// フォルダの作成
func MakeDir(path string) bool {
	slice := strings.Split(path, "/")
	_path := ""
	for i := range slice {
		_path += slice[i] + "/"
		if !Exists(_path) {
			err := os.Mkdir(_path, 0775)
			uDB.ErrCheck(err)
		}
	}
	return true
}
