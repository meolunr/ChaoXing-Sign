package course

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
