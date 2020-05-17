package main

import (
	"course"
	"encoding/json"
	"fmt"
	"global"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"task"
	"time"
)

var courses []*course.Course

var courseChan chan *course.Course
var signedChan chan string

func main() {
	global.Profile = loadProfile()
	global.Client = newHttpClient()
	courseChan = make(chan *course.Course)
	signedChan = make(chan string)

	login()
	obtainCourses()

	for _, item := range courses {
		go startSign()
		courseChan <- item
	}

	signedIds := make([]string, 0)
	for id := range signedChan {
		fmt.Println(id)
		signedIds = append(signedIds, id)
	}

	fmt.Println("main func stop...")
}

func startSign() {
	courseItem := <-courseChan
	jsonResp := courseItem.ObtainTasks()
	signTasks := filterSignTask(jsonResp)

	if len(signTasks) > 0 {
		fmt.Println("---------------------------------")
		fmt.Println(courseItem.Name)
		for _, signTask := range signTasks {
			fmt.Printf("  * %s\n", signTask.Name)
		}

		for _, signTask := range signTasks {
			//isSuccess := signTask.Sign()
			isSuccess := true
			if isSuccess {
				signedChan <- signTask.Id
				fmt.Println("签到成功：", signTask.Name)
			}
		}
	}
}

/**
过滤非签到任务
*/
func filterSignTask(jsonResp *task.JsonResponse) []*task.SignTask {
	signTasks := make([]*task.SignTask, 0)
	for _, item := range jsonResp.ActiveList {
		// 检查是否为未过期的签到任务
		if item.ActiveType == 2 && item.Status == 1 {
			signTask := &task.SignTask{
				Id:      strconv.Itoa(item.Id),
				Name:    item.NameOne,
				Referer: item.Url,
			}
			signTasks = append(signTasks, signTask)
		}
	}
	return signTasks
}

func loadProfile() *global.ProfileStruct {
	bytes, _ := ioutil.ReadFile("profile.json")
	profile := &global.ProfileStruct{}
	err := json.Unmarshal(bytes, profile)

	if err != nil {
		fmt.Println("用户配置文件读取失败")
		fmt.Println("请检查 profile.json")
		os.Exit(0)
	}
	return profile
}

func newHttpClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	return &http.Client{
		Jar: jar,
	}
}

func login() {
	cxUrl, _ := url.Parse("https://passport2-api.chaoxing.com/v11/loginregister")
	params := url.Values{}
	params.Set("uname", global.Profile.Username)
	params.Set("code", global.Profile.Password)

	cxUrl.RawQuery = params.Encode()
	request := global.NewClientRequest(http.MethodPost, cxUrl.String())
	response, err := global.Client.Do(request)

	if err != nil || response.StatusCode != http.StatusOK {
		fmt.Println("超星 API 请求失败")
		fmt.Println("10 秒后自动重试...")

		time.Sleep(time.Second * 10)
		login()
		return
	}

	defer global.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)
	jsonResp := &jsonResponse{}
	_ = json.Unmarshal(contentBytes, jsonResp)

	if jsonResp.Status == true {
		global.Uid = getUid(response)
		fmt.Println("登录成功")
	} else {
		fmt.Println("登录失败, message: ", jsonResp.Message)
	}
}

func obtainCourses() {
	request := global.NewClientRequest(http.MethodGet, "https://mooc1-api.chaoxing.com/mycourse/backclazzdata")
	response, _ := global.Client.Do(request)
	defer global.BodyClose(response.Body)
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
		fmt.Printf("共获取到 %d 个课程\n", len(jsonResp.ChannelList))

		for _, channel := range jsonResp.ChannelList {
			item := &course.Course{
				ClassId: strconv.Itoa(channel.Content.Id),
				Id:      strconv.Itoa(channel.Content.Course.Data[0].Id),
				Name:    channel.Content.Course.Data[0].Name,
			}
			courses = append(courses, item)
			fmt.Println("  * ", item.Name)
		}
		fmt.Println("---------------------------------")
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

type jsonResponse struct {
	Message string `json:"mes"`
	Status  bool   `json:"status"`
}
