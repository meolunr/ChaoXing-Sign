package course

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"netutil"
	"os"
	"strconv"
)

func ObtainCourses(client *http.Client) (courses []*Course) {
	request := netutil.NewRequest(http.MethodGet, "https://mooc1-api.chaoxing.com/mycourse/backclazzdata")
	response, _ := client.Do(request)
	defer netutil.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)

	jsonResp := &CoursesResponse{}
	err := json.Unmarshal(contentBytes, jsonResp)

	if err != nil {
		fmt.Println("Obtain course failed")
		os.Exit(0)
	}

	if jsonResp.Result == 1 {
		// Get courses success
		courses = make([]*Course, 0, len(jsonResp.ChannelList))

		for _, channel := range jsonResp.ChannelList {
			course := &Course{
				ClassId:    strconv.Itoa(channel.Content.ClassId),
				CourseId:   strconv.Itoa(channel.Content.Course.Data[0].CourseId),
				CourseName: channel.Content.Course.Data[0].CourseName,
			}
			courses = append(courses, course)
		}
	}
	return
}

func ObtainTaskList(course *Course, uid string, client *http.Client) {
	taskListUrl, _ := url.Parse("https://mobilelearn.chaoxing.com/ppt/activeAPI/taskactivelist")
	params := url.Values{}
	params.Set("classId", course.ClassId)
	params.Set("courseId", course.CourseId)
	params.Set("uid", uid)
	taskListUrl.RawQuery = params.Encode()

	request := netutil.NewRequest(http.MethodGet, taskListUrl.String())
	response, _ := client.Do(request)

	defer netutil.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)

	var jsonResp jsonResponse
	_ = json.Unmarshal(contentBytes, &jsonResp)
	filterSignTask(&jsonResp)
}

func filterSignTask(jsonResp *jsonResponse) []*SignTask {
	signTasks := make([]*SignTask, 0)
	for _, task := range jsonResp.ActiveList {
		// It's a sign task that has not expired
		if task.ActiveType == 2 && task.Status == 1 {
			signTask := &SignTask{
				Id:       task.Id,
				Referer:  task.Url,
				SignType: task.NameOne,
			}
			signTasks = append(signTasks, signTask)
		}
	}
	return signTasks
}

type SignTask struct {
	Id       int
	Referer  string
	SignType string
}

type jsonResponse struct {
	ActiveList []struct {
		Id         int    `json:"id"`
		Status     int    `json:"status"`
		ActiveType int    `json:"activeType"`
		NameOne    string `json:"nameOne"`
		Url        string `json:"url"`
	} `json:"activeList"`
}

type Course struct {
	ClassId    string
	CourseId   string
	CourseName string
}

type CoursesResponse struct {
	Result      int `json:"result"`
	ChannelList []struct {
		Content struct {
			ClassId int `json:"id"`
			Course  struct {
				Data []struct {
					CourseId   int    `json:"id"`
					CourseName string `json:"name"`
				} `json:"data"`
			} `json:"course"`
		} `json:"content"`
	} `json:"channelList"`
}
