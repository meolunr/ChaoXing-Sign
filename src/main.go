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
	"strconv"
	"time"
)

var profile *Profile
var client *http.Client

var uid string
var courses []*course.Course

func main() {
	loadProfile()
	newHttpClient()

	login()
	obtainCourses()

	item := courses[2]
	tasks := item.ObtainSignTasks(uid, client)
	tasks[0].Sign(uid, client)
}

func loadProfile() {
	bytes, _ := ioutil.ReadFile("profile.json")
	profile = &Profile{}
	err := json.Unmarshal(bytes, profile)

	if err != nil {
		fmt.Println("用户配置文件读取失败")
		fmt.Println("请检查 profile.json")
		os.Exit(0)
	}
}

func newHttpClient() {
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
	}
}

func login() {
	cxUrl, _ := url.Parse("https://passport2-api.chaoxing.com/v11/loginregister")
	params := url.Values{}
	params.Set("uname", profile.Username)
	params.Set("code", profile.Password)

	cxUrl.RawQuery = params.Encode()
	request := netutil.NewRequest(http.MethodPost, cxUrl.String())
	response, err := client.Do(request)

	if err != nil || response.StatusCode != http.StatusOK {
		fmt.Println("超星 API 请求失败")
		fmt.Println("10 秒后自动重试...")

		time.Sleep(time.Second * 10)
		login()
		return
	}

	defer netutil.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)
	jsonResp := &jsonResponse{}
	_ = json.Unmarshal(contentBytes, jsonResp)

	if jsonResp.Status == true {
		uid = getUid(response)
		fmt.Println("登录成功")
	} else {
		fmt.Println("登录失败, message: ", jsonResp.Message)
	}
}

func obtainCourses() {
	request := netutil.NewRequest(http.MethodGet, "https://mooc1-api.chaoxing.com/mycourse/backclazzdata")
	response, _ := client.Do(request)
	defer netutil.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)

	jsonResp := &course.JsonResponse{}
	err := json.Unmarshal(contentBytes, jsonResp)

	if err != nil {
		fmt.Println("获取课程失败")
		os.Exit(0)
	}

	if jsonResp.Result == 1 {
		// 获取课程成功
		courses = make([]*course.Course, 0, len(jsonResp.ChannelList))

		for _, channel := range jsonResp.ChannelList {
			item := &course.Course{
				ClassId: strconv.Itoa(channel.Content.Id),
				Id:      strconv.Itoa(channel.Content.Course.Data[0].Id),
				Name:    channel.Content.Course.Data[0].Name,
			}
			courses = append(courses, item)

			fmt.Println("---------------------------------")
			fmt.Println("ClassId:    ", item.ClassId)
			fmt.Println("CourseId:   ", item.Id)
			fmt.Println("CourseName: ", item.Name)
		}
	}
	return
}

func getUid(response *http.Response) string {
	for _, cookie := range response.Cookies() {
		if cookie.Name == "UID" {
			return cookie.Value
		}
	}
	return ""
}

type Profile struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type jsonResponse struct {
	Message string `json:"mes"`
	Status  bool   `json:"status"`
}
