package course

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"netutil"
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
func (course *Course) ObtainSignTasks(uid string, client *http.Client) []*task.SignTask {
	cxUrl, _ := url.Parse("https://mobilelearn.chaoxing.com/ppt/activeAPI/taskactivelist")
	params := url.Values{}
	params.Set("classId", course.ClassId)
	params.Set("courseId", course.Id)
	params.Set("uid", uid)

	cxUrl.RawQuery = params.Encode()
	request := netutil.NewClientRequest(http.MethodGet, cxUrl.String())
	response, _ := client.Do(request)

	defer netutil.BodyClose(response.Body)
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
				Referer: item.Url,
			}
			signTasks = append(signTasks, signTask)

			fmt.Printf("SignTask: %s, Course: %s\n", item.NameOne, course.Name)
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
