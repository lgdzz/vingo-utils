package vingo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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