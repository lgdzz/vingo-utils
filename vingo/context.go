package vingo

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"strings"
	"time"
)

type Context struct {
	*gin.Context
}

// url解码
func (c *Context) UrlDecode() (dStr string) {
	var (
		err  error
		eStr = c.Request.RequestURI
	)
	dStr, err = url.QueryUnescape(eStr)
	if err != nil {
		dStr = eStr
	}
	return
}

// 获取客户端真实IP
func (c *Context) GetRealClientIP() string {
	var ip string
	if ip = c.Request.Header.Get("X-Forwarded-For"); ip == "" {
		ip = c.Request.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = c.Request.RemoteAddr
	} else {
		ips := strings.Split(ip, ", ")
		ip = ips[len(ips)-1]
	}
	return ip
}

// 请求成功
func (c *Context) Response(d *ResponseData) {
	if d.Message == "" {
		d.Message = "Success"
	}
	if d.Status == 0 {
		d.Status = 200
	}
	uuid := GetUUID()
	c.Set("requestUUID", uuid)
	c.JSON(d.Status, gin.H{
		"uuid":      uuid,
		"error":     d.Error,
		"message":   d.Message,
		"data":      d.Data,
		"timestamp": time.Now().Unix(),
	})

	if !d.NoLog {
		// 记录请求日志
		go func(context *Context, uuid string, d *ResponseData) {
			startTime := context.GetTime("requestStart")
			endTime := time.Now()
			latency := endTime.Sub(startTime)
			duration := fmt.Sprintf("%.3fms", float64(latency.Nanoseconds())/float64(time.Millisecond))

			var err string
			if d.Error == 1 {
				err = d.Message
			}

			if context.Request.Method == "GET" {
				LogRequest(duration, fmt.Sprintf("{\"uuid\":\"%v\",\"method\":\"%v\",\"url\":\"%v\",\"err\":\"%v\",\"errType\":\"%v\",\"userAgent\":\"%v\",\"clientIP\":\"%v\"}", uuid, context.Request.Method, context.UrlDecode(), err, d.ErrorType, c.GetHeader("User-Agent"), c.GetRealClientIP()))
			} else {
				body := context.GetString("requestBody")
				if body == "" {
					body = "\"\""
				}
				LogRequest(duration, fmt.Sprintf("{\"uuid\":\"%v\",\"method\":\"%v\",\"url\":\"%v\",\"body\":%v,\"err\":\"%v\",\"errType\":\"%v\",\"userAgent\":\"%v\",\"clientIP\":\"%v\"}", uuid, context.Request.Method, context.Request.RequestURI, body, err, d.ErrorType, c.GetHeader("User-Agent"), c.GetRealClientIP()))
			}
		}(c, uuid, d)
	}
}

// 请求成功，带data数据
func (c *Context) ResponseBody(data any) {
	c.Response(&ResponseData{Data: data})
}

// 请求成功，默认
func (c *Context) ResponseSuccess() {
	c.Response(&ResponseData{})
}

// 绑定body参数结构体
func (c *Context) RequestBody(body any) {
	if err := c.ShouldBindJSON(body); err != nil {
		panic(err.Error())
	}

	if err := Valid.Struct(body); err != nil {
		// handle validation error
		panic(err)
	}

	if data, err := json.Marshal(body); err != nil {
		panic(err.Error())
	} else {
		c.Set("requestBody", string(data))
	}
}

// 绑定get参数结构体
func (c *Context) RequestQuery(query any) {
	if err := c.ShouldBindQuery(query); err != nil {
		panic(err.Error())
	}
}

func (c *Context) GetUserId() uint {
	return c.GetUint("userId")
}

func (c *Context) GetAccId() uint {
	return c.GetUint("accountId")
}

func (c *Context) GetOrgId() uint {
	return c.GetUint("orgId")
}

func (c *Context) GetRoleId() uint {
	return c.GetUint("roleId")
}

func (c *Context) GetMemberId() uint {
	return c.GetUint("memberId")
}

func (c *Context) GetCustomerId() uint {
	return c.GetUint("customerId")
}

func (c *Context) GetRealName() string {
	return c.GetString("realName")
}
