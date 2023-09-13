package aliyun

import (
	"github.com/lgdzz/vingo-utils/vingo/oss"
)

type Config struct {
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

// todo
func (s *ClientApi) Upload(object oss.Object, localFilePath string) *oss.UploadRes {
	return &oss.UploadRes{}
}

// todo
func (s *ClientApi) Delete(objectName string) error {
	return nil
}
