package global

import (
	"io"
	"net/http"
)

func NewClientRequest(method, url string) *http.Request {
	request, _ := http.NewRequest(method, url, nil)
	request.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 10; Pixel 2) com.chaoxing.mobile/ChaoXingStudy_3_4.3.7_android_phone_497_27 (@Kalimdor)_aed7e7f96119453a9c9727776a940d5e")
	return request
}

func NewWebViewRequest(method, url string) *http.Request {
	request, _ := http.NewRequest(method, url, nil)
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; Pixel 2) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/80.0.3987.99 Mobile Safari/537.36 com.chaoxing.mobile/ChaoXingStudy_3_4.3.7_android_phone_497_27 (@Kalimdor)_aed7e7f96119453a9c9727776a940d5e")
	return request
}

func NewFormRequest(url string, body io.Reader) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, url, body)
	request.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 10; Pixel 2) com.chaoxing.mobile/ChaoXingStudy_3_4.3.7_android_phone_497_27 (@Kalimdor)_aed7e7f96119453a9c9727776a940d5e")
	return request
}

func BodyClose(body io.Closer) {
	_ = body.Close()
}
