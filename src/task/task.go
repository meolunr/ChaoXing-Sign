package task

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"netutil"
	"strings"
)

type SignTask struct {
	Id       string
	Name     string
	Referer  string
	SignType int
}

func (task *SignTask) Sign(uid string, client *http.Client) {
	signType := task.getSignType(client)
	fmt.Println("-----------------> ", task.Name, ", signType = ", signType)
	/*cxUrl, _ := url.Parse("https://mobilelearn.chaoxing.com/pptSign/stuSignajax")
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
	request := netutil.NewWebViewRequest(http.MethodGet, cxUrl.String())
	request.Header.Set("Referer", task.Referer)
	response, _ := client.Do(request)

	defer netutil.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(contentBytes))*/
}

func (task *SignTask) getSignType(client *http.Client) (signType int) {
	request := netutil.NewWebViewRequest(http.MethodGet, task.Referer)
	response, _ := client.Do(request)

	defer netutil.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)

	html := string(contentBytes)
	switch {
	case strings.Contains(html, "手势"):
		signType = SignTypeGesture
	case strings.Contains(html, "拍照"):
		signType = SignTypePhoto
	case strings.Contains(html, "位置"):
		signType = SignTypeLocation
	case strings.Contains(html, "二维码"):
		signType = SignTypeQrCode
	default:
		signType = SignTypeGeneral
	}
	return
}

type JsonResponse struct {
	ActiveList []struct {
		Id         int    `json:"id"`
		Status     int    `json:"status"`
		ActiveType int    `json:"activeType"`
		NameOne    string `json:"nameOne"`
		Url        string `json:"url"`
	} `json:"activeList"`
}
