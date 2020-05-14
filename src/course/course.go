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

func (course *Course) ObtainSignTasks(uid string, client *http.Client) []*SignTask {
	cxUrl, _ := url.Parse("https://mobilelearn.chaoxing.com/ppt/activeAPI/taskactivelist")
	params := url.Values{}
	params.Set("classId", course.ClassId)
	params.Set("courseId", course.Id)
	params.Set("uid", uid)

	cxUrl.RawQuery = params.Encode()
	request := netutil.NewRequest(http.MethodGet, cxUrl.String())
	response, _ := client.Do(request)

	defer netutil.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)

	var jsonResp task.JsonResponse
	_ = json.Unmarshal(contentBytes, &jsonResp)

	return filterSignTask(course, &jsonResp)
}

/**
过滤非签到任务
*/
func filterSignTask(course *Course, jsonResp *task.JsonResponse) []*SignTask {
	signTasks := make([]*SignTask, 0)
	for _, item := range jsonResp.ActiveList {
		// 检查是否为未过期的签到任务
		//if item.ActiveType == 2 && item.Status == 1 {
		if item.ActiveType == 2 && item.Status == 2 { // 测试用
			signTask := &SignTask{
				Id:       strconv.Itoa(item.Id),
				Referer:  item.Url,
				SignType: item.NameOne,
			}
			signTasks = append(signTasks, signTask)

			fmt.Printf("SignTask: %s, Course : %s\n", item.NameOne, course.Name)
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

func sign(task *SignTask, uid string, client *http.Client) {
	cxUrl, _ := url.Parse("https://mobilelearn.chaoxing.com/pptSign/stuSignajax")
	params := url.Values{}
	params.Set("activeId", task.Id)
	params.Set("uid", uid)
	params.Set("latitude", "-1")
	params.Set("longitude", "-1")
	params.Set("appType", "15")
	params.Set("clientip", "")
	params.Set("fid", "")
	params.Set("name", "")

	cxUrl.RawQuery = params.Encode()
	request := netutil.NewRequest(http.MethodGet, cxUrl.String())
	request.Header.Set("Referer", task.Referer)
	response, _ := client.Do(request)

	defer netutil.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(contentBytes))
}

type SignTask struct {
	Id       string
	Referer  string
	SignType string
}

type jsonResponse1 struct {
	ActiveList []struct {
		Id         int    `json:"id"`
		Status     int    `json:"status"`
		ActiveType int    `json:"activeType"`
		NameOne    string `json:"nameOne"`
		Url        string `json:"url"`
	} `json:"activeList"`
}
