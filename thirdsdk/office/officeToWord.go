package office

import (
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/presentation"
	"errors"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

type excelRes struct {
	Question string `json:"question"` // 问题
	Answer   string `json:"answer"`   // 答案
}

type ExcelResult struct {
	Name    string     `json:"name"`    // 问题
	Content [][]string `json:"content"` // 答案
}

// word文件转文字
func WordToContent(filePath string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(filePath)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}
	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}

	doc, err := document.Open(filePath)
	if err != nil {
		return "", "", 0, errors.New("打开文件失败！")
	}

	// 遍历文档中的每个段落
	text := ""
	for _, p := range doc.Paragraphs() {
		// 遍历段落中的每个 Run
		for _, r := range p.Runs() {
			// 打印 Run 的文本
			text += r.Text()
		}
	}

	return text, suffix, size, nil
}

// word地址文件转文字
func WordUrlToContent(url string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(url)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}

	filePath, err := saveFile(url, suffix)
	if err != nil {
		return "", "", 0, errors.New("文件保存在本地失败！")
	}

	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}

	defer os.Remove(filePath)

	doc, err := document.Open(filePath)
	if err != nil {
		return "", "", 0, errors.New("打开文件失败！")
	}

	// 遍历文档中的每个段落
	text := ""
	for _, p := range doc.Paragraphs() {
		// 遍历段落中的每个 Run
		for _, r := range p.Runs() {
			// 打印 Run 的文本
			text += r.Text()
		}
	}
	return text, suffix, size, nil
}

// excel文件转文字
func ExcelToContent(filePath string) (word []excelRes, fileSuffix string, FileSize int, err error) {
	var list []excelRes
	suffix, err := getSuffix(filePath)
	if err != nil {
		return list, "", 0, errors.New("获取前缀失败！")
	}
	size, err := countSize(filePath)
	if err != nil {
		return list, "", 0, errors.New("计算文件大小失败！")
	}
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return list, "", 0, errors.New("打开文件失败！")
	}

	for _, sheet := range xlFile.Sheets {
		for k, row := range sheet.Rows {
			if k == 0 {
				continue
			}
			per := excelRes{}
			per.Question = row.Cells[0].String()
			per.Answer = row.Cells[1].String()
			list = append(list, per)
		}
	}
	return list, suffix, size, nil
}

// ExcelToContentTwo excel文件转文字（方法二）
func ExcelToContentTwo(filePath string) (word []ExcelResult, fileSuffix string, FileSize int, err error) {
	var excelResult []ExcelResult
	suffix, err := getSuffix(filePath)
	if err != nil {
		return excelResult, "", 0, errors.New("获取前缀失败！")
	}
	size, err := countSize(filePath)
	if err != nil {
		return excelResult, "", 0, errors.New("计算文件大小失败！")
	}
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return excelResult, "", 0, errors.New("打开文件失败！")
	}

	for _, sheet := range xlFile.Sheets {
		list := ExcelResult{}
		var table [][]string
		for _, row := range sheet.Rows {
			var rows []string
			for _, cell := range row.Cells {
				rows = append(rows, cell.String())
			}
			table = append(table, rows)
		}

		list.Content = table
		list.Name = sheet.Name
		excelResult = append(excelResult, list)
	}
	return excelResult, suffix, size, nil
}

// excel地址文件转文字
func ExcelUrlToContent(url string) (word []excelRes, fileSuffix string, FileSize int, err error) {
	var list []excelRes
	suffix, err := getSuffix(url)
	if err != nil {
		return list, "", 0, errors.New("获取前缀失败！")
	}

	filePath, err := saveFile(url, suffix)
	if err != nil {
		return list, "", 0, errors.New("文件保存在本地失败！")
	}

	size, err := countSize(filePath)
	if err != nil {
		return list, "", 0, errors.New("计算文件大小失败！")
	}

	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return list, "", 0, errors.New("打开文件失败！")
	}

	defer os.Remove(filePath)

	for _, sheet := range xlFile.Sheets {
		for k, row := range sheet.Rows {
			if k == 0 {
				continue
			}
			per := excelRes{}
			per.Question = row.Cells[0].String()
			per.Answer = row.Cells[1].String()
			list = append(list, per)
		}
	}

	return list, suffix, size, nil
}

// ExcelUrlToContentTwo  excel地址文件转文字（方法二）
func ExcelUrlToContentTwo(url string) (word []ExcelResult, fileSuffix string, FileSize int, err error) {
	var excelResult []ExcelResult
	suffix, err := getSuffix(url)
	if err != nil {
		return excelResult, "", 0, errors.New("获取前缀失败！")
	}

	filePath, err := saveFile(url, suffix)
	if err != nil {
		return excelResult, "", 0, errors.New("文件保存在本地失败！")
	}

	size, err := countSize(filePath)
	if err != nil {
		return excelResult, "", 0, errors.New("计算文件大小失败！")
	}

	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return excelResult, "", 0, errors.New("打开文件失败！")
	}

	defer os.Remove(filePath)

	for _, sheet := range xlFile.Sheets {
		list := ExcelResult{}
		var table [][]string
		for _, row := range sheet.Rows {
			var rows []string
			for _, cell := range row.Cells {
				rows = append(rows, cell.String())
			}
			table = append(table, rows)
		}

		list.Content = table
		list.Name = sheet.Name
		excelResult = append(excelResult, list)
	}

	return excelResult, suffix, size, nil
}

// ppt文件转文字
func PptToContent(filePath string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(filePath)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}
	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}

	ppt, err := presentation.Open(filePath)
	if err != nil {
		return "", "", 0, errors.New("打开文件失败！")
	}

	var text string
	for _, slide := range ppt.Slides() {
		//所有的控件
		for _, choice := range slide.X().CSld.SpTree.Choice {
			if choice.Sp == nil {
				continue
			}
			//一个文本框或一个控件
			for _, sp := range choice.Sp {
				if sp.TxBody == nil {
					continue
				}
				//数据
				for _, p := range sp.TxBody.P {
					textrun := p.EG_TextRun
					//所有的EG_TextRun中的数据组合起来是一段
					for _, run := range textrun {
						if run.R != nil {
							text += run.R.T
						}
					}
					if len(text) == 0 {
						continue
					}
				}
			}
		}
	}

	return text, suffix, size, nil
}

// ppt地址文件转文字
func PptUrlToContent(url string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(url)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}

	filePath, err := saveFile(url, suffix)
	if err != nil {
		return "", "", 0, errors.New("文件保存在本地失败！")
	}

	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}

	ppt, err := presentation.Open(filePath)
	if err != nil {
		return "", "", 0, errors.New("打开文件失败！")
	}

	var text string
	for _, slide := range ppt.Slides() {
		//所有的控件
		for _, choice := range slide.X().CSld.SpTree.Choice {
			if choice.Sp == nil {
				continue
			}
			//一个文本框或一个控件
			for _, sp := range choice.Sp {
				if sp.TxBody == nil {
					continue
				}
				//数据
				for _, p := range sp.TxBody.P {
					textrun := p.EG_TextRun
					//所有的EG_TextRun中的数据组合起来是一段
					for _, run := range textrun {
						if run.R != nil {
							text += run.R.T
						}
					}
					if len(text) == 0 {
						continue
					}
				}
			}
		}
	}

	return text, suffix, size, nil
}

// txt文件转文字
func TxtToContent(filePath string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(filePath)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}
	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	text := string(content)
	text = strings.Replace(text, "\n", " ", -1)

	if !utf8.ValidString(text) {
		return "", "", 0, errors.New("文件编码只能是UTF-8！")
	}

	return text, suffix, size, nil
}

// txt地址文件转文字
func TxtUrlToContent(url string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(url)
	if err != nil {
		return "", "", 0, err
	}

	filePath, err := saveFile(url, suffix)
	if err != nil {
		return "", "", 0, errors.New("文件保存在本地失败！")
	}

	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}

	defer os.Remove(filePath)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	text := string(content)

	if !utf8.ValidString(text) {
		return "", "", 0, errors.New("文件编码只能是UTF-8！")
	}

	text = strings.Replace(text, "\n", " ", -1)

	return text, suffix, size, nil
}
