package wechat

import (
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/config"
)

type MiniProgramConfig struct {
	AppID     string `json:"app_id"`     // appid
	AppSecret string `json:"app_secret"` // appSecret
}

// 微信小程序
func MiniProgram(miniProgramConfig *MiniProgramConfig) *miniprogram.MiniProgram {
	wc := wechat.NewWechat()
	cfg := &config.Config{
		AppID:     miniProgramConfig.AppID,
		AppSecret: miniProgramConfig.AppSecret,
		Cache:     CacheApi,
	}
	return wc.GetMiniProgram(cfg)
}
