package ip

import (
	"fmt"
	"github.com/lgdzz/vingo-utils/vingo"
	"github.com/lgdzz/vingo-utils/vingo/request"
)

type Info struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	Area    string `json:"area"`
}

// 获取ip信息
func GetInfo(serverUrl string, ip string) (info Info) {
	var result map[string]any
	vingo.StringToJson(string(request.Get(fmt.Sprintf("%v/?ip=%v", serverUrl, ip), request.Option{})), &result)
	vingo.CustomOutput(result[ip], &info)
	if info.Country == "" {
		info.Country = "未知区域"
	}
	return
}
