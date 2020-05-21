package main

import (
	"chaoxing-sign/course"
	"chaoxing-sign/global"
	"chaoxing-sign/task"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"
)

var signedChan chan *task.SignTask
var signedIds = make([]string, 0)

func main() {
	global.LoadProfile()
	global.NewHttpClient()
	global.Login()

	signedChan = make(chan *task.SignTask)
	courses := obtainCourses()

	startTime, endTime, rangeErr := getRefreshTimeRange()
	// 单个课程休眠时间 = 总休眠时间 / 课程数
	// 避免并发请求所有课程的任务列表
	delay := time.Second * time.Duration(global.Profile.Interval/len(courses))

	ticker := time.NewTicker(time.Second * time.Duration(global.Profile.Interval))
	defer ticker.Stop()
	go func() {
		for now := range ticker.C {
			if rangeErr == nil {
				// 当前时间是否不在可签到时间段内
				if !(now.After(startTime) && now.Before(endTime)) {
					continue
				}
			}
			for _, item := range courses {
				startSign(item)
				time.Sleep(delay)
			}
		}
	}()

	for signTask := range signedChan {
		signedIds = append(signedIds, signTask.Id)
		pushServerChan(signTask)
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
		fmt.Printf("共获取到 %d 个课程，以下课程将会自动签到：\n", len(jsonResp.ChannelList))

		for _, channel := range jsonResp.ChannelList {
			// 排除不是学生的课程和未开课的课程
			if !channel.Content.IsStart {
				continue
			}
			item := &course.Course{
				ClassId: strconv.Itoa(channel.Content.Id),
				Id:      strconv.Itoa(channel.Content.Course.Data[0].Id),
				Name:    channel.Content.Course.Data[0].Name,
			}
			// 排除不需要签到的课程
			if containInSlice(global.Profile.ExcludeCourse, item.Id) {
				continue
			}

			courses = append(courses, item)
			fmt.Printf("[ %9s ] %s\n", item.Id, item.Name)
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
		fmt.Println()

		for _, signTask := range signTasks {
			isSuccess := signTask.Sign()
			if isSuccess {
				signTask.Course = course.Name
				signedChan <- signTask
				fmt.Printf("签到成功：%s (%s)\n", signTask.Name, time.Now().Format("2006-01-02 15:04"))
			} else {
				fmt.Printf("签到失败：%s (%s)\n", signTask.Name, time.Now().Format("2006-01-02 15:04"))
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
			// 是否未过期
			if taskId := strconv.Itoa(item.Id); item.Status == 1 {
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
				signedIds = removeInSlice(signedIds, taskId)
			}
		}
	}
	return signTasks
}

func pushServerChan(signTask *task.SignTask) {
	if global.Profile.ServerChan == "" {
		return
	}

	serverChanUrl, _ := url.Parse(global.Profile.ServerChan)
	params := url.Values{}
	params.Set("text", signTask.Course+" 签到成功")
	params.Set("desp", signTask.Name)

	serverChanUrl.RawQuery = params.Encode()
	request := global.NewWebViewRequest(http.MethodGet, serverChanUrl.String())
	_, _ = global.Client.Do(request)
}

/**
获取签到周期，在此时间范围内可刷新签到任务
*/
func getRefreshTimeRange() (startTime time.Time, endTime time.Time, err error) {
	currentDate := time.Now().Format("2006-01-02")
	var startErr, endErr error
	startTime, startErr = time.ParseInLocation("2006-01-02 15:04",
		fmt.Sprintf("%s %s", currentDate, global.Profile.StartTime), time.Local)
	endTime, endErr = time.ParseInLocation("2006-01-02 15:04",
		fmt.Sprintf("%s %s", currentDate, global.Profile.EndTime), time.Local)

	if startErr != nil || endErr != nil {
		err = errors.New("start time or end time format is incorrect")
	}
	return
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
