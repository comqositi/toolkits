package baidu

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/unidoc/unipdf/v3/model"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func (b *BaiduOcr) commonFun(payload *strings.Reader) (word string, err error) {
	token, err := b.getAccessToken()
	if err != nil {
		return "", err
	}

	requestUrl := fmt.Sprintf(transformUrlBaidu, token)

	client := &http.Client{}
	req, err := http.NewRequest("POST", requestUrl, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	resBody1 := BodyResultResponse{}
	err = json.Unmarshal(body, &resBody1)
	if err != nil {
		return
	}

	var str string
	for _, val := range resBody1.WordsResult {
		str += val.Words + " "
	}
	str = strings.TrimRight(str, " ")

	return str, nil
}

func saveFile(url string, suffix string) (string, error) {
	byString, err := md5ByString(url)
	if err != nil {
		return "", err
	}
	targetName := "temporary" + byString + "." + suffix

	// 发起GET请求
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// 创建临时文件
	tmpFile, err := ioutil.TempFile("/tmp", "tempfile")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// 将HTTP响应的内容写入临时文件
	_, err = io.Copy(tmpFile, response.Body)
	if err != nil {
		return "", errors.New("写入临时文件时出错！")
	}

	// 重命名临时文件为目标文件名并保存到本地/tmp目录下
	err = os.Rename(tmpFile.Name(), "/tmp/"+targetName)
	if err != nil {
		return "", errors.New("重命名临时文件失败！")
	}
	return "/tmp/" + targetName, nil
}

// base64编码后进行urlEncode
func (b *BaiduOcr) getFileContentAsBase64(path string) string {
	srcByte, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(srcByte)
}

func getSuffix(url string) (string, error) {
	dotIndex := strings.LastIndex(url, ".")
	if dotIndex == -1 || dotIndex == len(url)-1 {
		return "", errors.New("没有文件后缀或者文件名以点号结尾！")
	}
	suffix := url[dotIndex+1:]
	return suffix, nil
}

func countSize(filePath string) (int, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, errors.New("计算文件大小失败！")
	}

	fileSize := int(fileInfo.Size())
	return fileSize, nil
}

func countImgSize(url string) (int, error) {
	// 获取响应
	resp, err := http.Get(url)
	if err != nil {
		return 0, errors.New("远程获取图片失败！")
	}
	defer resp.Body.Close()

	// 读取body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.New("后去图片大小失败！")
	}

	// 获取大小
	size := len(body)
	return size, nil
}

func md5ByString(str string) (string, error) {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		return "", err
	}
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr), nil
}

func getPdfNum(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, errors.New("无法打开 PDF 文件！")
	}
	defer file.Close()
	// 创建 PDF reader
	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		return 0, errors.New("无法创建 PDF reader！")
	}
	// 获取 PDF 文件总页数
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return 0, errors.New("无法获取 PDF 文件页数！")
	}

	return numPages, nil
}
