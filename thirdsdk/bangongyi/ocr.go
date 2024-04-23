package bangongyi

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type InfoRequest struct {
	Url     string `json:"url"`
	PageNum int64  `json:"page_num"`
}

type InfoResponse struct {
	Success bool     `json:"success"`
	Msg     string   `json:"msg"`
	Data    []string `json:"data"`
}

func ImageToContent(url string, imageUrl string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(imageUrl)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}

	filePath, err := saveFile(imageUrl, suffix)
	if err != nil {
		return "", "", 0, errors.New("文件保存在本地失败！")
	}

	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}

	if size > 1048576*10 {
		return "", "", 0, errors.New("图片不能大于 10 MB！")
	}

	defer os.Remove(filePath)

	body, _ := PostRequest(url, &InfoRequest{
		Url: imageUrl,
	})

	resBody := InfoResponse{}
	err = json.Unmarshal(body, &resBody)
	if err != nil {
		return "", "", 0, errors.New("解析数据失败！")
	}
	if resBody.Success != true {
		return "", "", 0, errors.New("解析数据失败！")
	}
	for _, v := range resBody.Data {
		word += v + ","
	}
	word = strings.Trim(word, ",")

	return word, suffix, size, nil
}

func PdfToContent(url string, pdfUrl string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(pdfUrl)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}

	filePath, err := saveFile(pdfUrl, suffix)
	if err != nil {
		return "", "", 0, errors.New("文件保存在本地失败！")
	}

	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}

	if size > 1048576*20 {
		return "", "", 0, errors.New("图片不能大于 20 MB！")
	}

	defer os.Remove(filePath)

	body, _ := PostRequest(url, &InfoRequest{
		Url:     pdfUrl,
		PageNum: 30,
	})

	resBody := InfoResponse{}
	err = json.Unmarshal(body, &resBody)
	if err != nil {
		return "", "", 0, errors.New("解析数据失败！")
	}
	if resBody.Success != true {
		return "", "", 0, errors.New("解析数据失败！")
	}
	for _, v := range resBody.Data {
		word += v + ","
	}
	word = strings.Trim(word, ",")

	return word, suffix, size, nil

}
