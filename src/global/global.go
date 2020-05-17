package global

import "net/http"

var Profile *ProfileStruct
var Client *http.Client

var Uid string

type ProfileStruct struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Interval   int    `json:"interval"`
	ServerChan string `json:"serverChan"`
}
