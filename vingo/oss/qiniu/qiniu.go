package qiniu

import (
	"fmt"
	"github.com/lgdzz/vingo-utils/vingo/db/redis"
	"github.com/lgdzz/vingo-utils/vingo/oss"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"time"
)

type Config struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Expires   uint64
	Cache     bool
}

type ClientApi struct {
	Config Config
	Client any
}

// 在主进程中只需要执行一次
func InitClient(config Config) (api ClientApi) {
	api.Config = config

	api.Client = ""

	return api
}

// 上传签名
func (s *ClientApi) Sign() (token string) {
	key := fmt.Sprintf("qiniu_%v", s.Config.AccessKey)
	putPolicy := storage.PutPolicy{
		Scope: s.Config.Bucket,
	}

	// 签名有效期
	if s.Config.Expires == 0 {
		putPolicy.Expires = 600 // 默认10分钟有效期
	} else {
		putPolicy.Expires = s.Config.Expires
	}

	if s.Config.Cache {
		tokenPoint := redis.Get[string](key)
		// 从缓存中读取凭证
		if tokenPoint != nil {
			token = *tokenPoint
			return
		}
	}

	mac := qbox.NewMac(s.Config.AccessKey, s.Config.SecretKey)
	token = putPolicy.UploadToken(mac)

	if s.Config.Cache {
		// 缓存提前100秒失效
		redis.Set(key, token, time.Duration(putPolicy.Expires-100)*time.Second)
	}
	return
}

// todo
func (s *ClientApi) Upload(object oss.Object, localFilePath string) *oss.UploadRes {
	return &oss.UploadRes{}
}

// todo
func (s *ClientApi) Delete(objectName string) error {
	return nil
}
