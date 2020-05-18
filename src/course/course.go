package course

import (
	"encoding/json"
	"errors"
	"fmt"
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
func (course *Course) ObtainTasks() (jsonResp *task.JsonResponse) {
	global.Retry(func() error {
		cxUrl, _ := url.Parse("https://mobilelearn.chaoxing.com/ppt/activeAPI/taskactivelist")
		params := url.Values{}
		params.Set("classId", course.ClassId)
		params.Set("courseId", course.Id)
		params.Set("uid", global.Uid)

		jar := global.Client.Jar
		fmt.Println(&jar)
		cxUrl.RawQuery = params.Encode()
		request := global.NewClientRequest(http.MethodGet, cxUrl.String())
		response, err := global.Client.Do(request)
		if response == nil || response.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintln("obtain tasks failed.", err))
		}

		defer global.BodyClose(response.Body)
		contentBytes, _ := ioutil.ReadAll(response.Body)
		jsonErr := json.Unmarshal(contentBytes, &jsonResp)
		if jsonErr != nil {
			global.Login()
			return errors.New(fmt.Sprintln("attempting to login again to get cookies.", jsonErr))
		}

		return nil
	})
	return
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
