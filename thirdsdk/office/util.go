package office

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/nguyenthenguyen/docx"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

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

func md5ByString(str string) (string, error) {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		return "", err
	}
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr), nil
}

func wordToData(local string) (string, error) {
	r, err := docx.ReadDocxFile(local)
	if err != nil {
		return "", errors.New("读取pdf文件失败！")
	}

	docx := r.Editable()
	docx.Replace("old_2_1", "new_2_1", -1)
	docx.Replace("old_2_2", "new_2_2", -1)
	res := docx.GetContent()
	regex := regexp.MustCompile(`<w:t>(.*?)</w:t>`)
	// 查找所有匹配项
	matches := regex.FindAllStringSubmatch(res, -1)

	resStr := ""
	for _, match := range matches {
		resStr += match[1]
	}

	defer r.Close()

	return resStr, nil
}
