package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func main() {
	bytes, _ := ioutil.ReadFile("profile.json")
	profile := &Profile{}
	_ = json.Unmarshal(bytes, profile)

	login(profile.Username, profile.Password)
}

func login(username string, password string) {
	loginUrl, _ := url.Parse("https://passport2-api.chaoxing.com/v11/loginregister")
	params := url.Values{}
	params.Set("uname", username)
	params.Set("code", password)
	loginUrl.RawQuery = params.Encode()

	request, _ := http.NewRequest(http.MethodPost, loginUrl.String(), nil)
	request.Header.Add("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 10; Pixel 2) com.chaoxing.mobile/ChaoXingStudy_3_4.3.7_android_phone_497_27 (@Kalimdor)_aed7e7f96119453a9c9727776a940d5e")

	cookieJar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: cookieJar,
	}
	response, _ := client.Do(request)
	defer bodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)

	jsonResp := &Response{}
	_ = json.Unmarshal(contentBytes, jsonResp)

	fmt.Println(cookieJar)
}

func bodyClose(body io.Closer) {
	_ = body.Close()
}

type Response struct {
	Message string `json:"mes"`
	Status  bool   `json:"status"`
}

type Profile struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Cookies  string `json:"cookies"`
}
