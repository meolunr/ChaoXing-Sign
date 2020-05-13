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

var profile *Profile
var client *http.Client
var courses []*Course

func main() {
	loadProfile()
	newHttpClient()
	login(profile.Username, profile.Password)
	getCourses()

	for _, cours := range courses {
		fmt.Println(cours)
	}
}

func loadProfile() {
	bytes, _ := ioutil.ReadFile("profile.json")
	profile = &Profile{}
	_ = json.Unmarshal(bytes, profile)
}

func newHttpClient() {
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
	}
}

func login(username string, password string) {
	loginUrl, _ := url.Parse("https://passport2-api.chaoxing.com/v11/loginregister")
	params := url.Values{}
	params.Set("uname", username)
	params.Set("code", password)
	loginUrl.RawQuery = params.Encode()

	request, _ := http.NewRequest(http.MethodPost, loginUrl.String(), nil)
	request.Header.Add("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 10; Pixel 2) com.chaoxing.mobile/ChaoXingStudy_3_4.3.7_android_phone_497_27 (@Kalimdor)_aed7e7f96119453a9c9727776a940d5e")

	response, _ := client.Do(request)
	defer bodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)

	jsonResp := &Response{}
	_ = json.Unmarshal(contentBytes, jsonResp)

	if jsonResp.Status == true {
		fmt.Println("User login successfully")
	}
}

func getCourses() {
	request, _ := http.NewRequest(http.MethodGet, "https://mooc1-api.chaoxing.com/mycourse/backclazzdata", nil)
	request.Header.Add("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 10; Pixel 2) com.chaoxing.mobile/ChaoXingStudy_3_4.3.7_android_phone_497_27 (@Kalimdor)_aed7e7f96119453a9c9727776a940d5e")

	response, _ := client.Do(request)
	defer bodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)

	jsonResp := &CoursesResponse{}
	_ = json.Unmarshal(contentBytes, jsonResp)

	if jsonResp.Result == 1 {
		// Get courses success
		fmt.Println(jsonResp.ChannelList)
		courses = make([]*Course, len(jsonResp.ChannelList))

		for _, channel := range jsonResp.ChannelList {
			course := &Course{
				ClassId:    channel.Content.ClassId,
				CourseId:   channel.Content.Course.Data[0].CourseId,
				CourseName: channel.Content.Course.Data[0].CourseName,
			}
			courses = append(courses, course)
		}
	}
}

func bodyClose(body io.Closer) {
	_ = body.Close()
}

type Course struct {
	ClassId    int `json:"id"`
	CourseId   int
	CourseName string
}

type CoursesResponse struct {
	Result      int `json:"result"`
	ChannelList []struct {
		Content struct {
			ClassId int `json:"id"`
			Course  struct {
				Data []struct {
					CourseId   int    `json:"id"`
					CourseName string `json:"name"`
				} `json:"data"`
			} `json:"course"`
		} `json:"content"`
	} `json:"channelList"`
}

type Response struct {
	Message string `json:"mes"`
	Status  bool   `json:"status"`
}

type Profile struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
