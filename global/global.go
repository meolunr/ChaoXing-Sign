package global

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

var Profile *ProfileStruct
var Client *http.Client
var Uid string

func LoadProfile() {
	bytes, _ := ioutil.ReadFile("profile.json")
	Profile = &ProfileStruct{}
	err := json.Unmarshal(bytes, Profile)

	if err != nil {
		fmt.Println("用户配置文件读取失败")
		fmt.Println("请检查 profile.json")
		os.Exit(0)
	}

	if Profile.Interval == 0 {
		// 默认刷新间隔时间为 60 秒
		Profile.Interval = 60
	}
}

func NewHttpClient() {
	jar, _ := cookiejar.New(nil)
	Client = &http.Client{
		Jar: jar,
	}
}

func Login() {
	Retry(func() error {
		cxUrl, _ := url.Parse("https://passport2-api.chaoxing.com/v11/loginregister")
		params := url.Values{}
		params.Set("uname", Profile.Username)
		params.Set("code", Profile.Password)

		cxUrl.RawQuery = params.Encode()
		request := NewClientRequest(http.MethodPost, cxUrl.String())
		response, err := Client.Do(request)
		if response == nil || response.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintln("login failed.", err))
		}

		defer BodyClose(response.Body)
		contentBytes, _ := ioutil.ReadAll(response.Body)
		jsonResp := &jsonResponse{}
		_ = json.Unmarshal(contentBytes, jsonResp)

		if jsonResp.Status == true {
			Uid = getUid(response)
			fmt.Println("登录成功")
		} else {
			fmt.Println("登录失败, message: ", jsonResp.Message)
		}
		return nil
	})
}

func getUid(response *http.Response) string {
	for _, cookie := range response.Cookies() {
		if cookie.Name == "UID" {
			return cookie.Value
		}
	}
	return ""
}

type ProfileStruct struct {
	Username      string   `json:"username"`
	Password      string   `json:"password"`
	Interval      int      `json:"interval"`
	StartTime     string   `json:"startTime"`
	EndTime       string   `json:"endTime"`
	ServerChan    string   `json:"serverChan"`
	ExcludeCourse []string `json:"excludeCourse"`
}

type jsonResponse struct {
	Message string `json:"mes"`
	Status  bool   `json:"status"`
}
