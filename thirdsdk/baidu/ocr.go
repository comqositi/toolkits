package baidu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	tokenUrlBaiDu     = "https://aip.baidubce.com/oauth/2.0/token"
	transformUrlBaidu = "https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic?access_token=%s"
)

type BodyResultResponse struct {
	LogId          int         `json:"log_id"`
	WordsResultNum int         `json:"words_result_num"`
	WordsResult    []WordsList `json:"words_result"`
}

type WordsList struct {
	Words string `json:"words"`
}

type BaiDuTokenResponse struct {
	RefreshToken     string `json:"refresh_token,omitempty"`
	ExpiresIn        int64  `json:"expires_in,omitempty"`
	SessionKey       string `json:"session_key,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	Scope            string `json:"scope,omitempty"`
	SessionSecret    string `json:"session_secret,omitempty"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type Cache interface {
	Set(key string, value string, expires int) error
	Get(key string) (string, error)
}

type BaiduOcr struct {
	cache     Cache
	apiKey    string
	apiSecret string
}

func NewBaiduOcr(apiKey string, apiSecret string, cache Cache) (*BaiduOcr, error) {
	c := &BaiduOcr{apiKey: apiKey, apiSecret: apiSecret, cache: cache}
	_, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// 图片转文字
func (b *BaiduOcr) ImageToWord(filePath string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(filePath)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}
	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}

	encode := b.getFileContentAsBase64(filePath)
	contextLen := len(encode)
	if contextLen/1024/1024 > 8 {
		return "", "", 0, errors.New("文件大小不能大于8M！")
	}
	payload := strings.NewReader("image=" + url.QueryEscape(encode) + "&detect_direction=false&detect_language=false&paragraph=false&probability=false")
	str, err := b.commonFun(payload)
	if err != nil {
		return "", "", 0, errors.New("word文档解析失败！")
	}

	return str, suffix, size, nil
}

// 图片地址转文字
func (b *BaiduOcr) ImageUrlToWord(imageUrl string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(imageUrl)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}

	size, err := countImgSize(imageUrl)
	if err != nil {
		return "", "", 0, errors.New("获取图片大小失败！")
	}

	if len(imageUrl) > 1024 {
		return "", "", 0, errors.New("图片地址不能超过 1024 个字节")
	}

	encode := b.getFileContentAsBase64(imageUrl)
	contextLen := len(encode)
	if contextLen/1024/1024 > 8 {
		return "", "", 0, errors.New("文件大小不能大于8M")
	}
	payload := strings.NewReader("url=" + url.QueryEscape(imageUrl) + "&detect_direction=false&detect_language=false&paragraph=false&probability=false")
	str, err := b.commonFun(payload)
	if err != nil {
		return "", "", 0, errors.New("word文档解析失败！")
	}

	return str, suffix, size, nil
}

// pdf转文字
func (b *BaiduOcr) PdfToWord(filePath string) (word string, fileSuffix string, FileSize int, err error) {
	numPages, err := getPdfNum(filePath)
	if err != nil {
		return "", "", 0, err
	}
	if numPages > 21 {
		return "", "", 0, errors.New("pdf文件不能超过20页！")
	}

	suffix, err := getSuffix(filePath)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}
	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}

	defer os.Remove(filePath)

	encode := b.getFileContentAsBase64(filePath)
	contextLen := len(encode)
	if contextLen/1024/1024 > 5 {
		return "", "", 0, errors.New("文件大小不能大于5M")
	}

	var result string
	for i := 1; i < numPages+1; i++ {
		time.Sleep(time.Millisecond * 300)
		payload := strings.NewReader("pdf_file=" + url.QueryEscape(encode) + "&pdf_file_num=" + strconv.Itoa(i) + "&detect_direction=false&detect_language=false&paragraph=false&probability=false")

		str, err := b.commonFun(payload)
		if err != nil {
			return "", "", 0, errors.New("pdf解析失败！")
		}
		result += str
	}

	return result, suffix, size, nil
}

// pdf转文字
func (b *BaiduOcr) PdfUrlToWord(pdfUrl string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(pdfUrl)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}

	filePath, err := saveFile(pdfUrl, suffix)
	if err != nil {
		return "", "", 0, errors.New("文件保存在本地失败！")
	}

	numPages, err := getPdfNum(filePath)
	if err != nil {
		return "", "", 0, err
	}
	if numPages > 21 {
		return "", "", 0, errors.New("pdf文件不能超过20页！")
	}
	defer os.Remove(filePath)

	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}

	encode := b.getFileContentAsBase64(filePath)
	contextLen := len(encode)
	if contextLen/1024/1024 > 5 {
		return "", "", 0, errors.New("文件大小不能大于5M")
	}

	var result string
	for i := 1; i < numPages+1; i++ {
		time.Sleep(time.Millisecond * 300)
		payload := strings.NewReader("pdf_file=" + url.QueryEscape(encode) + "&pdf_file_num=" + strconv.Itoa(i) + "&detect_direction=false&detect_language=false&paragraph=false&probability=false")

		str, err := b.commonFun(payload)
		if err != nil {
			return "", "", 0, errors.New("pdf解析失败！")
		}
		result += str
	}

	return result, suffix, size, nil
}

// 获取token
func (b *BaiduOcr) getAccessToken() (token string, err error) {

	md5String, _ := md5ByString(b.apiKey)
	tokenKey := "kpai:baiduocr:" + md5String
	token, err = b.cache.Get(tokenKey)
	if err != nil {
		fmt.Printf("baidu gettoken redis token, err = %v \n", err)
	}
	if len(token) > 1 {
		fmt.Printf("baidu gettoken redis token is not empty, token 1%v1 \n", token)
		return token, nil
	}

	url := tokenUrlBaiDu + "?client_id=%s&client_secret=%s&grant_type=client_credentials"
	url = fmt.Sprintf(url, b.apiKey, b.apiSecret)
	payload := strings.NewReader(``)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Printf("baidu gettoken http.NewRequest, err %v\n", err)
		return token, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("baidu gettoken http.NewRequest Do, err %v\n", err)
		return token, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("baidu gettoken http.NewRequest ioutil.ReadAll, err %v\n", err)
		return token, err
	}

	var baiDuTokenResponse BaiDuTokenResponse
	err = json.Unmarshal(body, &baiDuTokenResponse)
	if err != nil {
		fmt.Printf("baidu gettoken http.NewRequest json.Unmarshal, body %v ;err %v\n", string(body), err)
		return token, err
	}
	if len(baiDuTokenResponse.Error) > 1 {
		fmt.Printf("baidu gettoken http.NewRequest err, ErrorMsg %v, ErrorCode %v \n", baiDuTokenResponse.Error, baiDuTokenResponse.ErrorDescription)
		return token, err
	}

	token = baiDuTokenResponse.AccessToken
	if len(token) > 0 {
		fmt.Printf("baidu translate save token, %v  \n", token)
		if err != nil {
			fmt.Printf("baidu gettoken save token, %v  \n", err.Error())
			return token, err
		}

		err = b.cache.Set(tokenKey, token, int(baiDuTokenResponse.ExpiresIn))
		if err != nil {
			fmt.Printf("baidu gettoken save token Expire,token = %v  err = %v \n", token, err)
			return token, err
		}
	}
	return token, nil
}
