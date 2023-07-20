package vingo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func Request(method, url string, headers map[string]string, body interface{}) []byte {
	var requestBody []byte
	if body != nil {
		requestBody, _ = json.Marshal(body)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
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

func RequestStream(method, url string, headers map[string]string, body interface{}, receive func(...byte)) {
	var requestBody []byte
	if body != nil {
		requestBody, _ = json.Marshal(body)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
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

// 发送文件到远程接口
func FilePostRequest(url string, filePath string, fieldName *string) []byte {
	if fieldName == nil {
		fieldName = StringPointer("file")
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
	filePart, err := writer.CreateFormFile(*fieldName, filePath)
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
		Timeout: 60 * time.Second,
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
