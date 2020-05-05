package appstructure

// 画面用のデータ構造
type ResponseData struct {
	ViewData  map[string]string
	Message   map[string][]string
	LoginFlag bool
}

func CreateResponseData(viewData map[string]string, message map[string][]string) ResponseData {
	responseData := ResponseData{
		ViewData:  viewData,
		Message:   message,
		LoginFlag: false,
	}

	return responseData
}
