package task

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	cxUrl, _ := url.Parse("https://mobilelearn.chaoxing.com/pptSign/stuSignajax")
	params := url.Values{}
	// 签到通用参数
	params.Set("activeId", task.Id)
	params.Set("uid", uid)
	params.Set("latitude", "-1")
	params.Set("longitude", "-1")
	params.Set("appType", "15")
	params.Set("clientip", "")
	params.Set("fid", "")
	params.Set("name", "")
	params.Set("useragent", "")

	// 针对特殊方式签到追加参数
	signType := task.getSignType(client)
	switch signType {
	case SignTypePhoto:
		params.Set("objectId", "")
	case SignTypeLocation:
		params.Set("address", "中国")
		params.Set("ifTiJiao", "1")
	}

	cxUrl.RawQuery = params.Encode()
	request := netutil.NewWebViewRequest(http.MethodGet, cxUrl.String())
	request.Header.Set("Referer", task.Referer)
	response, _ := client.Do(request)

	defer netutil.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(contentBytes))
}

/**
获取签到类型
*/
func (task *SignTask) getSignType(client *http.Client) (signType int) {
	// 模拟用户点击客户端签到任务打开网页
	request := netutil.NewWebViewRequest(http.MethodGet, task.Referer)
	response, _ := client.Do(request)

	defer netutil.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)

	// 通过签到网页中的字符串区分签到类型
	html := string(contentBytes)
	switch {
	default:
		signType = SignTypeGeneral
	case strings.Contains(html, "手势"):
		signType = SignTypeGesture
	case strings.Contains(html, "拍照"):
		signType = SignTypePhoto
	case strings.Contains(html, "位置"):
		signType = SignTypeLocation
	case strings.Contains(html, "二维码"):
		signType = SignTypeQrCode
	}
	return
}

/**
获取上传图片所需要的 Token
*/
func getToken(client *http.Client) string {
	request := netutil.NewClientRequest(http.MethodGet, "https://pan-yz.chaoxing.com/api/token/uservalid")
	response, _ := client.Do(request)

	defer netutil.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)
	jsonResp := make(map[string]string)
	_ = json.Unmarshal(contentBytes, &jsonResp)

	return jsonResp["_token"]
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
