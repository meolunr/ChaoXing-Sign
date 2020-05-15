package course

import (
	"encoding/json"
	"fmt"
	"global"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"task"
)

type Course struct {
	Id      string
	Name    string
	ClassId string
}

/**
获取所有未签到的任务
*/
func (course *Course) ObtainSignTasks() []*task.SignTask {
	cxUrl, _ := url.Parse("https://mobilelearn.chaoxing.com/ppt/activeAPI/taskactivelist")
	params := url.Values{}
	params.Set("classId", course.ClassId)
	params.Set("courseId", course.Id)
	params.Set("uid", global.Uid)

	cxUrl.RawQuery = params.Encode()
	request := global.NewClientRequest(http.MethodGet, cxUrl.String())
	response, _ := global.Client.Do(request)

	defer global.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)

	var jsonResp task.JsonResponse
	_ = json.Unmarshal(contentBytes, &jsonResp)

	return course.filterSignTask(&jsonResp)
}

/**
过滤非签到任务
*/
func (course *Course) filterSignTask(jsonResp *task.JsonResponse) []*task.SignTask {
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

			fmt.Printf("SignTask: %s, Course: %s\n", signTask.Name, course.Name)
		}
	}
	return signTasks
}

type JsonResponse struct {
	Result      int `json:"result"`
	ChannelList []struct {
		Content struct {
			Id     int `json:"id"`
			Course struct {
				Data []struct {
					Id   int    `json:"id"`
					Name string `json:"name"`
				} `json:"data"`
			} `json:"course"`
		} `json:"content"`
	} `json:"channelList"`
}
