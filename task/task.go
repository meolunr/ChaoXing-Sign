package task

import (
	"bytes"
	"chaoxing-sign/global"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type SignTask struct {
	Id       string
	Name     string
	Referer  string
	SignType int
	Course   string
}

/**
@return 签到是否成功
*/
func (task *SignTask) Sign() bool {
	cxUrl, _ := url.Parse("https://mobilelearn.chaoxing.com/pptSign/stuSignajax")
	params := url.Values{}
	// 签到通用参数
	params.Set("activeId", task.Id)
	params.Set("uid", global.Uid)
	params.Set("latitude", "-1")
	params.Set("longitude", "-1")
	params.Set("appType", "15")
	params.Set("clientip", "")
	params.Set("fid", "")
	params.Set("name", "")
	params.Set("useragent", "")

	// 针对特殊方式签到追加参数
	signType := task.getSignType()
	switch signType {
	case SignTypePhoto:
		params.Set("objectId", uploadPhoto())
	case SignTypeLocation:
		params.Set("address", "中国")
		params.Set("ifTiJiao", "1")
	}

	cxUrl.RawQuery = params.Encode()
	request := global.NewWebViewRequest(http.MethodGet, cxUrl.String())
	request.Header.Set("Referer", task.Referer)
	response, _ := global.Client.Do(request)
	if response == nil || response.StatusCode != http.StatusOK {
		return false
	}

	defer global.BodyClose(response.Body)
	contentBytes, _ := ioutil.ReadAll(response.Body)
	contentStr := string(contentBytes)

	return strings.Contains(contentStr, "success") ||
		strings.Contains(contentStr, "已签到")
}

/**
@return 签到类型
*/
func (task *SignTask) getSignType() (signType int) {
	// 模拟用户点击客户端签到任务打开网页
	request := global.NewWebViewRequest(http.MethodGet, task.Referer)
	response, _ := global.Client.Do(request)
	if response == nil || response.StatusCode != http.StatusOK {
		return
	}

	defer global.BodyClose(response.Body)
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
上传照片
@return 用于拍照签到的 ObjectId
*/
func uploadPhoto() (objectId string) {
	_, err := os.Stat("photo.jpg")
	if os.IsNotExist(err) {
		// 没有自定义签到照片时，返回默认 ObjectId
		// 是一张 654x872 尺寸的纯黑色背景
		return "92bef3565bcd2a950b10e5cd98190e39"
	}

	cxUrl, _ := url.Parse("https://pan-yz.chaoxing.com/upload")
	params := url.Values{}
	params.Set("_token", getToken())
	cxUrl.RawQuery = params.Encode()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	form, formErr := writer.CreateFormFile("file", "photo.jpg")
	file, openErr := os.Open("photo.jpg")
	defer func() { _ = file.Close() }()
	if formErr != nil || openErr != nil {
		fmt.Println("photo.jpg 打开失败")
		return
	}
	_, writeErr := io.Copy(form, file)
	if writeErr != nil {
		fmt.Println("photo.jpg 读取失败")
		return
	}

	_ = writer.WriteField("puid", global.Uid)
	_ = writer.Close()

	request := global.NewFormRequest(cxUrl.String(), body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, _ := global.Client.Do(request)
	if response == nil || response.StatusCode != http.StatusOK {
		return
	}
	defer global.BodyClose(response.Body)

	contentBytes, _ := ioutil.ReadAll(response.Body)
	jsonResp := make(map[string]string)
	_ = json.Unmarshal(contentBytes, &jsonResp)

	return jsonResp["objectId"]
}

/**
@return 上传图片需要的 Token
*/
func getToken() string {
	request := global.NewClientRequest(http.MethodGet, "https://pan-yz.chaoxing.com/api/token/uservalid")
	response, _ := global.Client.Do(request)
	if response == nil || response.StatusCode != http.StatusOK {
		return ""
	}

	defer global.BodyClose(response.Body)
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
