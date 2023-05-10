package vingo

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// JsonToString 结构体转字符串
func JsonToString(v any) string {
	o, _ := json.Marshal(v)
	return string(o)
}

// StringToJson 字符串转结构体
func StringToJson(v string, r any) {
	_ = json.Unmarshal([]byte(v), &r)
}

func MD5(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

func PasswordToCipher(text string, salt string) string {
	return MD5(MD5(text) + salt)
}

// 自定义输出格式
func CustomOutput(input any, output any) {
	treeByte, _ := json.Marshal(input)
	_ = json.Unmarshal(treeByte, output)
}

// 转金额保留两位小数
func ToMoney(value float64) float64 {
	return math.Round(value*100) / 100
}

func ToUint(v any) uint {
	switch t := v.(type) {
	case uint:
		return t
	case int32:
		return uint(t)
	case int64:
		return uint(t)
	case uint32:
		return uint(t)
	case float32:
		return uint(t)
	case float64:
		return uint(t)
	case string:
		v, _ := strconv.Atoi(t)
		return uint(v)
	default:
		panic(fmt.Sprintf("Cannot convert to uint: %v", v))
	}
	return 0
}

func ToString(v any) string {
	switch value := v.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.FormatInt(int64(value), 10)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	}
	return ""
}

func ToFloat64(v string) float64 {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic(fmt.Sprintf("字符串：%v转换float64失败，错误：%v", v, err.Error()))
	}
	return f
}
