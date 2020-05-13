package netutil

import (
	"io"
	"net/http"
)

func NewRequest(method, url string) *http.Request {
	request, _ := http.NewRequest(method, url, nil)
	request.Header.Add("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 10; Pixel 2) com.chaoxing.mobile/ChaoXingStudy_3_4.3.7_android_phone_497_27 (@Kalimdor)_aed7e7f96119453a9c9727776a940d5e")
	return request
}

func BodyClose(body io.Closer) {
	_ = body.Close()
}
