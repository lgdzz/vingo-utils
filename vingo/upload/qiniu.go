package upload

import (
	"fmt"
	"github.com/lgdzz/vingo-utils/vingo"
	"github.com/lgdzz/vingo-utils/vingo/db"
	"github.com/lgdzz/vingo-utils/vingo/db/redis"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"time"
)

type QiniuOption struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Expires   uint64
	Cache     bool
}

// 七牛上传签名
func QiniuUploadSign(option QiniuOption) (token string) {
	key := fmt.Sprintf("qiniu_%v", option.AccessKey)
	putPolicy := storage.PutPolicy{
		Scope: option.Bucket,
	}

	// 签名有效期
	if option.Expires == 0 {
		putPolicy.Expires = 600 // 默认10分钟有效期
	} else {
		putPolicy.Expires = option.Expires
	}

	if option.Cache {
		tokenPoint := redis.Get[string](key)
		// 从缓存中读取凭证
		if tokenPoint != nil {
			token = *tokenPoint
			return
		}
	}

	mac := qbox.NewMac(option.AccessKey, option.SecretKey)
	token = putPolicy.UploadToken(mac)

	if option.Cache {
		// 缓存提前100秒失效
		redis.Set(key, token, time.Duration(putPolicy.Expires-100)*time.Second)
	}
	return
}

// 文件记录保存
func FileRecordSave(c *vingo.Context) {
	body := File{OrgID: c.GetOrgId(), Channel: "backstage"}
	c.RequestBody(&body)
	db.Pool.Create(&body)
	c.ResponseSuccess()
}

// 文件记录列表
func FileRecordList(c *vingo.Context) {
	query := FileQuery{}
	c.RequestQuery(&query)
	pool := db.Pool.Table("file").Where("org_id=?", c.GetOrgId())
	c.ResponseBody(db.NewPage(pool, &db.PageResult{Page: query.Page, Size: query.Size, Items: "map"}))
}
