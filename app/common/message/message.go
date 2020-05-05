package message

const (
	ErrMsgTitleNull string = "タイトルを入力してください"
	ErrMsgAuthNull  string = "著者を入力してください"
	ErrMsgLiNull    string = "最新所持巻数を数字で入力してください"
	ErrMsgServerErr string = "現在不安定な状態です。再度、お試しください。"
	ErrMsgDelErr	string = "現在不安定な状態です。削除に失敗しました。再度、お試し下さい。"
	ErrMsgNoEmail	string = "メールアドレスを入力してください。"
	ErrMsgNoPassword	string = "パスワードを入力してください。"
	ErrMsgNoUserErr	string = "該当のユーザが見つかりませんでした。メールアドレス，パスワードを確かめてください。"
	ErrMsgNoSession	string = "ログインを行ってください。"
)

const (
	SucMsgUpdate string = "更新が完了しました。"
	SucMsgDel    string = "削除が完了しました。"
)
