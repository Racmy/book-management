package userdao

import (

	"log"
	"github.com/docker_go_nginx/app/common/appstructure"
	"github.com/docker_go_nginx/app/utility/uDB"
)

//GetBookByID ...BookテーブルのIDに紐つく情報を1件取得
/*
@param id string
@return book Book
@return err error
*/
func GetUserByEmailAndPass(mailAddress string, password string) (appstructure.UserData, error) {
	db := uDB.DbSetUp()
	defer db.Close()
	rows, err := db.Query("SELECT id, email, name, user_image_path FROM user WHERE email = ? and password = ?",mailAddress,password)
	defer rows.Close()
	uDB.ErrCheck(err)

	// SELECT失敗時にuserがerrorでのリターン
	if err != nil {
		log.Print("【UserDao.GetUserByEmailAndPass】mailAddress = " + mailAddress + " password = " + password + " not exist in user table.")
		return appstructure.UserData{}, err
	}

	//　本が検索できた場合は、本の情報を含めてリターン
	var responceUserData appstructure.UserData
	if rows.Next() {
		err = rows.Scan(&responceUserData.ID, &responceUserData.Email , &responceUserData.Name, &responceUserData.ImagePath)
		uDB.ErrCheck(err)
		return responceUserData, err
	}else{
		log.Print("【UserDao.GetUserByEmailAndPass】mailAddress = " + mailAddress + " password = " + password + " not exist in user table.")
	}

	// 検索したが「０件」の場合は、book・errが共に空
	return responceUserData, err
}
