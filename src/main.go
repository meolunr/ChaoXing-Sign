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
	"sort"
	"strconv"
	"task"
	"time"
)

var signedChan chan string
var signedIds = make([]string, 0)

func main() {
	global.Profile = loadProfile()
	global.Client = newHttpClient()
	signedChan = make(chan string)

	login()
	courses := obtainCourses()

	// 单个课程休眠时间 = 总休眠时间 / 课程数
	// 避免并发请求所有课程的任务列表
	delay := time.Second * time.Duration(16/len(courses))
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			for _, item := range courses {
				startSign(item)
				time.Sleep(delay)
			}
		}
	}()

	for id := range signedChan {
		signedIds = append(signedIds, id)
	}
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

func obtainCourses() (courses []*course.Course) {
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

func startSign(course *course.Course) {
	jsonResp := course.ObtainTasks()
	signTasks := filterSignTask(jsonResp)

	if len(signTasks) > 0 {
		fmt.Println("---------------------------------")
		fmt.Println(course.Name)
		for _, signTask := range signTasks {
			fmt.Printf("  * %s\n", signTask.Name)
		}

		for _, signTask := range signTasks {
			isSuccess := signTask.Sign()
			if isSuccess {
				signedChan <- signTask.Id
				fmt.Println("签到成功：", signTask.Name)
			} else {
				fmt.Println("签到失败：", signTask.Name)
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
		// 是否为签到任务
		if item.ActiveType == 2 {
			taskId := strconv.Itoa(item.Id)
			// 是否未过期
			if item.Status == 1 {
				// 是否未签到
				if !containInSlice(signedIds, taskId) {
					signTasks = append(signTasks, &task.SignTask{
						Id:      taskId,
						Name:    item.NameOne,
						Referer: item.Url,
					})
				}
			} else {
				// 签到任务已过期，从已签到切片中移除 taskId
				removeInSlice(signedIds, taskId)
			}
		}
	}
	fmt.Println(signedIds)
	return signTasks
}

func getUid(response *http.Response) string {
	for _, cookie := range response.Cookies() {
		if cookie.Name == "UID" {
			return cookie.Value
		}
	}
	return ""
}

/**
@return slice 内是否包含某个元素
*/
func containInSlice(haystack []string, needle string) bool {
	sort.Strings(haystack)

	index := sort.SearchStrings(haystack, needle)
	return index < len(haystack) && haystack[index] == needle
}

/**
从 slice 中删除某个元素
*/
func removeInSlice(haystack []string, needle string) []string {
	sort.Strings(haystack)

	index := sort.SearchStrings(haystack, needle)
	if index < len(haystack) && haystack[index] == needle {
		return append(haystack[:index], haystack[index+1:]...)
	}
	return haystack
}

type jsonResponse struct {
	Message string `json:"mes"`
	Status  bool   `json:"status"`
}
