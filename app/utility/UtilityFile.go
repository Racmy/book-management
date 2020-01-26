package utilityFile

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"
	"mime/multipart"
)

const defaultFileUploadPath = "static/img/"

const TimeFormat = "2006-01-02_15-04-05"
/*
	ファイルをサーバへアップロードする関数

	input:
		file mulipart.File		: saved file
		filePath string			: saved file path
		fileName string			: saved file path
	output:
		string				: fileName(with path)
		err					: err
*/
func FileUpload(file multipart.File, filePath string, fileName string) (string ,error){
	savedFilePath := filePath
	//日本時間用*Location
	Jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowTime := time.Now().In(Jst).Format(TimeFormat)
	log.Println(nowTime)
	savedFileName := savedFilePath + nowTime +"_" + fileName

	err := fileSaved(file, savedFileName)
	if err != nil{
		return "", err
	}
	return "/" + savedFileName, err
}

/*
	ファイルをサーバへアップロードする関数

	input:
		file mulipart.File		: saved file
		fileName string
		fileType int				: saved file type
			1:表示画像用
	output:
		err					: err
*/
func DefaultFileUpload(file multipart.File, fileName string) (string ,error){
	return FileUpload(file,defaultFileUploadPath,fileName)
}

/*
	ファイルをサーバへ保存する関数

	input:
		file mulipart.File		: saved file
		fileName string			: saved file name(with path)
	output:
		err					: err
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
