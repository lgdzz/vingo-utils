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

func (c *Context) UrlDecode() (decodeStr string) {
	var (
		err       error
		encodeStr = c.Request.RequestURI
	)
	decodeStr, err = url.QueryUnescape(encodeStr)
	if err != nil {
		decodeStr = encodeStr
	}
	return
}

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

type ResponseData struct {
	Status    int
	Error     int8
	ErrorType string
	Message   string
	Data      any
	NoLog     bool // true时不记录日志
}

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

func (c *Context) ResponseBody(data any) {
	c.Response(&ResponseData{Data: data})
}

func (c *Context) ResponseSuccess() {
	c.Response(&ResponseData{})
}

func (c *Context) RequestBody(body any) {
	_ = c.ShouldBindJSON(body)

	data, _ := json.Marshal(body)
	c.Set("requestBody", string(data))
}

func (c *Context) RequestQuery(query any) {
	_ = c.ShouldBindQuery(query)
}

// 获取账号ID
func (c *Context) GetUserId() uint {
	return c.GetUint("userId")
}

// 获取账户ID
func (c *Context) GetAccId() uint {
	return c.GetUint("accountId")
}

// 获取角色ID
func (c *Context) GetRoleId() uint {
	return c.GetUint("roleId")
}

// 获取组织ID
func (c *Context) GetOrgId() uint {
	return c.GetUint("orgId")
}
