package userdao

import (
	"database/sql"
	"errors"
	"github.com/docker_go_nginx/app/utility/uDB"
	"log"
)

// User ...ユーザデータ
type User struct {
	ID        int64
	Email     string
	Name      string
	Password  string
	ImagePath string
}

func GetUserInstance(email string, name string, password string, imagePath string) User {
	return User{Email: email, Name: name, Password: password, ImagePath: imagePath}
}

//GetBookByID ...BookテーブルのIDに紐つく情報を1件取得
/*
@param id string
@return book Book
@return err error
*/
func GetUserByEmailAndPass(mailAddress string, password string) (User, error) {
	db := uDB.DbSetUp()
	defer db.Close()
	rows, err := db.Query("SELECT id, email, name, user_image_path FROM user WHERE email = ? and password = ?", mailAddress, password)
	defer rows.Close()
	uDB.ErrCheck(err)

	// SELECT失敗時にuserがerrorでのリターン
	if err != nil {
		log.Print("【UserDao.GetUserByEmailAndPass】mailAddress = " + mailAddress + " password = " + password + " not exist in user table.")
		return User{}, err
	}

	//　本が検索できた場合は、本の情報を含めてリターン
	var responceUserData User
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

//IsSetEmail...Userテーブルに登録されているメールアドレスか判定
/*
@param mailAddress string メールアドレス
@return bool true：存在する、false：存在しない
*/
func IsSetEmail(mailAddress string) bool {
	db := uDB.DbSetUp()
	defer db.Close()
	rows, err := db.Query("SELECT id, email, name, user_image_path FROM user WHERE email = ?", mailAddress)
	defer rows.Close()
	uDB.ErrCheck(err)
	return rows.Next()
}

//IsSetName...Userテーブルに登録されている名前か判定
/*
@param name string 名前
@return bool true：存在する、false：存在しない
*/
func IsSetName(name string) bool {
	db := uDB.DbSetUp()
	defer db.Close()
	rows, err := db.Query("SELECT id, email, name, user_image_path FROM user WHERE name = ?", name)
	defer rows.Close()
	uDB.ErrCheck(err)
	return rows.Next()
}

//insertUser ...ユーザの登録
/*
@param User　ユーザID未セット
@return User ユーザIDセット済
*/
func InsertUser(user User) (User, error) {
	db := uDB.DbSetUp()
	defer db.Close() // 関数がリターンする直前に呼び出される
	var result sql.Result
	var err error

	// 登録済メールアドレスか判定
	// 登録済ユーザ名か判定
	if IsSetEmail(user.Email) || IsSetName(user.Name) {
		err = errors.New("登録済みユーザ名・またはメールアドレスです。")
		return user, err
	}

	// 画像ない版
	if len(user.ImagePath) == 0 {
		ins, err := db.Prepare("INSERT INTO user (email, password, name) VALUES(?,?,?)")
		uDB.ErrCheck(err)
		// Userを格納する
		result, err = ins.Exec(&user.Email, &user.Password, &user.Name)
		uDB.ErrCheck(err)
	} else {
		ins, err := db.Prepare("INSERT INTO user (email, password, name, user_image_path) VALUES(?,?,?,?)")
		uDB.ErrCheck(err)
		// Userを格納する
		result, err = ins.Exec(1, &user.Email, &user.Password, &user.Name, &user.ImagePath)
		uDB.ErrCheck(err)
	}
	// Insertした結果を返す（id, error）
	user.ID, _ = result.LastInsertId()
	return user, err
}
