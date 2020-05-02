package userdao

import (
	"github.com/docker_go_nginx/app/utility/uDB"
	"log"
)

// UserData ...ユーザデータ
type UserData struct {
	ID        int
	Email     string
	Name      string
	ImagePath string
}

//GetBookByID ...BookテーブルのIDに紐つく情報を1件取得
/*
@param id string
@return book Book
@return err error
*/
func GetUserByEmailAndPass(mailAddress string, password string) (UserData, error) {
	db := uDB.DbSetUp()
	defer db.Close()
	rows, err := db.Query("SELECT id, email, name, user_image_path FROM user WHERE email = ? and password = ?", mailAddress, password)
	defer rows.Close()
	uDB.ErrCheck(err)

	// SELECT失敗時にuserがerrorでのリターン
	if err != nil {
		log.Print("【UserDao.GetUserByEmailAndPass】mailAddress = " + mailAddress + " password = " + password + " not exist in user table.")
		return UserData{}, err
	}

	//　本が検索できた場合は、本の情報を含めてリターン
	var responceUserData UserData
	if rows.Next() {
		err = rows.Scan(&responceUserData.ID, &responceUserData.Email, &responceUserData.Name, &responceUserData.ImagePath)
		uDB.ErrCheck(err)
		return responceUserData, err
	} else {
		log.Print("【UserDao.GetUserByEmailAndPass】mailAddress = " + mailAddress + " password = " + password + " not exist in user table.")
	}

	// 検索したが「０件」の場合は、book・errが共に空
	return responceUserData, err
}

//insertUser ...ユーザの登録
/*

 */
func InsertUser(book) {
	db := uDB.DbSetUp()
	defer db.Close() // 関数がリターンする直前に呼び出される
	var result sql.Result
	if book.FrontCoverImagePath == "" {
		ins, err := db.Prepare("INSERT INTO book (user_id,title,author,latest_issue) VALUES(?,?,?,?)")
		uDB.ErrCheck(err)
		// Bookを格納する
		result, err = ins.Exec(1, &book.Title, &book.Author, &book.LatestIssue)
		uDB.ErrCheck(err)
	} else {
		ins, err := db.Prepare("INSERT INTO book (user_id,title,author,latest_issue,front_cover_image_path) VALUES(?,?,?,?,?)")
		uDB.ErrCheck(err)
		// Bookを格納する
		result, err = ins.Exec(1, &book.Title, &book.Author, &book.LatestIssue, &book.FrontCoverImagePath)
		uDB.ErrCheck(err)
	}
	// Insertした結果を返す（id, error）
	return result.LastInsertId()
}
