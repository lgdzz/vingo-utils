package sms

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/lgdzz/vingo-utils/vingo"
)

// 阿里云短信配置参数
type AliyunOption struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SignName        string `json:"signName"`
}

type AliyunBody struct {
	UUID          string         // 短信发送批次号
	Phone         string         // 接收短信的手机号
	TemplateCode  string         // 模板ID
	TemplateParam map[string]any // 模板参数
	Option        AliyunOption
}

type AliyunResponse struct {
	RequestId string `json:"requestId"`
	Message   string `json:"message"`
	Code      string `json:"code"`
}

func (s *AliyunBody) Log(msg string) {
	vingo.LogInfo(fmt.Sprintf("[发送短信] UUID: %v AccessKeyId: %v 接收手机号: %v %v", s.UUID, s.Option.AccessKeyId, s.Phone, msg))
}

func (s *AliyunBody) AliyunSendSMS() bool {
	s.UUID = vingo.GetUUID()

	s.Log(fmt.Sprintf("内容：%v", vingo.JsonToString(s.TemplateParam)))

	// 初始化SDK客户端
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", s.Option.AccessKeyId, s.Option.AccessKeySecret)
	if err != nil {
		s.Log(fmt.Sprintf("初始化SDK失败：%v", err.Error()))
		return false
	}

	// 创建请求对象
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"                 // 支持HTTP和HTTPS协议
	request.Domain = "dysmsapi.aliyuncs.com" // API服务地址
	request.Version = "2017-05-25"           // API版本号
	request.ApiName = "SendSms"              // API名称

	// 设置请求参数
	request.QueryParams["PhoneNumbers"] = s.Phone
	request.QueryParams["SignName"] = s.Option.SignName
	request.QueryParams["TemplateCode"] = s.TemplateCode
	request.QueryParams["TemplateParam"] = vingo.JsonToString(s.TemplateParam)

	// 发送请求，获取响应
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		s.Log(fmt.Sprintf("发送请求失败：%v", err.Error()))
		return false
	}

	// 解析响应结果
	if response.GetHttpStatus() == 200 {
		var resp AliyunResponse
		vingo.StringToJson(response.GetHttpContentString(), resp)
		if resp.Code == "OK" {
			s.Log(fmt.Sprintf("发送短信成功：%v", response.GetHttpContentString()))
			return true
		} else {
			s.Log(fmt.Sprintf("发送短信失败：%v", response.GetHttpContentString()))
			return false
		}
	} else {
		s.Log(fmt.Sprintf("发送短信失败：%v", response.GetHttpContentString()))
		return false
	}
}
