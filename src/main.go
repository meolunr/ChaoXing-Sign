package main

import (
	"course"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"netutil"
	"os"
	"time"
)

var profile *Profile
var client *http.Client

func main() {
	loadProfile()
	newHttpClient()
	login(profile.Username, profile.Password)
	course.ObtainCourses(client)
}

func loadProfile() {
	bytes, _ := ioutil.ReadFile("profile.json")
	profile = &Profile{}
	err := json.Unmarshal(bytes, profile)

	if err != nil {
		fmt.Println("Profile read error")
		fmt.Println("Please check the \"profile.json\" file")
		os.Exit(0)
	}
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

	request := netutil.NewRequest(http.MethodPost, loginUrl.String())
	response, err := client.Do(request)

	if err != nil || response.StatusCode != http.StatusOK {
		fmt.Println("Request for ChaoXing API failed")
		fmt.Println("Try again in 10 seconds...")

		time.Sleep(time.Second * 10)
		login(username, password)
		return
	}

	defer netutil.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)
	jsonResp := &Response{}
	_ = json.Unmarshal(contentBytes, jsonResp)

	if jsonResp.Status == true {
		fmt.Println("User login successfully")
	} else {
		fmt.Println("User login failed, message: ", jsonResp.Message)
	}
}

type Profile struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Message string `json:"mes"`
	Status  bool   `json:"status"`
}
