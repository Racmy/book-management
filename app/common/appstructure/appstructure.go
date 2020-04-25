package appstructure

// UserData ...ユーザデータ
type UserData struct {
	ID			int
	Email       string
	Name		string
	ImagePath	string
}

// HomeErrorMessage ...エラーメッセージ格納
type HomeErrorMessage struct {
	EmailErr	string
	PasswordErr	string
	NoUserErr	string
}
