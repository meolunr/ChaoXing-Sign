package course

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"netutil"
)

var courses []*Course

func ObtainCourses(client *http.Client) {
	request := netutil.NewRequest(http.MethodGet, "https://mooc1-api.chaoxing.com/mycourse/backclazzdata")
	response, _ := client.Do(request)
	defer netutil.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)

	jsonResp := &CoursesResponse{}
	_ = json.Unmarshal(contentBytes, jsonResp)

	if jsonResp.Result == 1 {
		// Get courses success
		fmt.Println(jsonResp.ChannelList)
		courses = make([]*Course, len(jsonResp.ChannelList))

		for _, channel := range jsonResp.ChannelList {
			course := &Course{
				ClassId:    channel.Content.ClassId,
				CourseId:   channel.Content.Course.Data[0].CourseId,
				CourseName: channel.Content.Course.Data[0].CourseName,
			}
			courses = append(courses, course)
		}
	}

	for _, cours := range courses {
		fmt.Println(cours)
	}
}

type Course struct {
	ClassId    int `json:"id"`
	CourseId   int
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
