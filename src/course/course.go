package course

import (
	"encoding/json"
	"global"
	"io/ioutil"
	"net/http"
	"net/url"
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
func (course *Course) ObtainTasks() *task.JsonResponse {
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

	return &jsonResp
}

type JsonResponse struct {
	Result      int `json:"result"`
	ChannelList []struct {
		Content struct {
			Id      int  `json:"id"`
			IsStart bool `json:"isstart"`
			Course  struct {
				Data []struct {
					Id   int    `json:"id"`
					Name string `json:"name"`
				} `json:"data"`
			} `json:"course"`
		} `json:"content"`
	} `json:"channelList"`
}
