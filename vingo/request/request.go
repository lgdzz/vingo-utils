package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lgdzz/vingo-utils/vingo"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type Option struct {
	Headers       *map[string]string
	Timeout       *int    // 请求超时时间，单位：秒，默认30秒
	FileFieldName *string // 发送文件的字段名，默认：file
}

// 发送get请求
func Get(url string, option Option) []byte {
	if option.Timeout == nil {
		option.Timeout = vingo.IntPointer(30)
	}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err.Error())
	}
	if option.Headers != nil {
		for key, value := range *option.Headers {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{
		Timeout: time.Duration(*option.Timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return responseBody
}

// 发送json格式的post请求
func PostByJson(url string, body interface{}, option Option) []byte {
	if option.Timeout == nil {
		option.Timeout = vingo.IntPointer(30)
	}
	var requestBody []byte
	if body != nil {
		requestBody, _ = json.Marshal(body)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	if option.Headers != nil {
		for key, value := range *option.Headers {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{
		Timeout: time.Duration(*option.Timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return responseBody
}

// 发送form-data格式的post请求
func PostByFormData(url string, formData map[string]string, option Option) []byte {
	if option.Timeout == nil {
		option.Timeout = vingo.IntPointer(30)
	}

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	for key, value := range formData {
		writer.WriteField(key, value)
	}

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		panic(err.Error())
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	if option.Headers != nil {
		for key, value := range *option.Headers {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{
		Timeout: time.Duration(*option.Timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return responseBody
}

// 发送文件请求
func PostFile(url string, option Option, filePath string) []byte {
	if option.FileFieldName == nil {
		option.FileFieldName = vingo.StringPointer("file")
	}
	if option.Timeout == nil {
		option.Timeout = vingo.IntPointer(30)
	}
	// 打开文件
	fileHandle, err := os.Open(filePath)
	if err != nil {
		panic(err.Error())
	}
	defer fileHandle.Close()

	// 创建一个 buffer 用于存储文件内容
	body := &bytes.Buffer{}

	// 创建一个新的 multipart writer
	writer := multipart.NewWriter(body)

	// 创建一个文件表单字段
	filePart, err := writer.CreateFormFile(*option.FileFieldName, filePath)
	if err != nil {
		panic(err.Error())
	}

	// 将文件内容复制到文件表单字段中
	_, err = io.Copy(filePart, fileHandle)
	if err != nil {
		panic(err.Error())
	}

	// 完成写入
	writer.Close()

	// 创建一个 HTTP POST 请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		panic(err.Error())
	}

	// 设置请求头，包括 Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求并获取响应
	client := &http.Client{
		Timeout: time.Duration(*option.Timeout) * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	// 处理响应
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("server returned non-200 status: %v", resp.Status))
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return responseBody
}

// 发送post请求，实时返回数据（未测试）
func PostByJsonStream(url string, body interface{}, option Option, receive func(...byte)) {
	if option.Timeout == nil {
		option.Timeout = vingo.IntPointer(30)
	}
	var requestBody []byte
	if body != nil {
		requestBody, _ = json.Marshal(body)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	if option.Headers != nil {
		for key, value := range *option.Headers {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{
		Timeout: time.Duration(*option.Timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	//var responseBody []byte
	buf := make([]byte, 1024) // 设置缓冲区大小
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			panic(err.Error())
		}
		if n == 0 {
			break
		}

		receive(buf[:n]...)
	}
}
