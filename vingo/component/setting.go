package component

import (
	"fmt"
	"github.com/lgdzz/vingo-utils/vingo"
	"github.com/lgdzz/vingo-utils/vingo/cache"
	"github.com/lgdzz/vingo-utils/vingo/db/mysql"
	"gorm.io/gorm"
)

type Setting[T any] struct {
	OrgId     uint             `gorm:"primaryKey;column:org_id" json:"orgId"`
	AccId     uint             `gorm:"primaryKey;column:acc_id" json:"accId"`
	Name      string           `gorm:"primaryKey;column:name" json:"name"`
	Value     T                `gorm:"column:value;serializer:json" json:"value"`
	CreatedAt *vingo.LocalTime `gorm:"column:created_at;" json:"createdAt"`
	UpdatedAt *vingo.LocalTime `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt   `gorm:"column:deleted_at" json:"deletedAt"`
}

func (s *Setting[T]) TableName() string {
	return "setting"
}

// 获取设置值，如果不存在则创建
// Example component.SettingRead[Test]("test", 0, 0)
func SettingRead[T any](name string, orgId uint, accId uint) T {
	return *cache.Fast(settingCacheKey(name, orgId, accId), 0, func() *T {
		var setting = settingDefault[T](name, orgId, accId)
		return &setting.Value
	})
}

// 保存设置值
// Example component.SettingSave("test", 0, 0)
func SettingSave[T any](name string, orgId uint, accId uint, value T) {
	cache.FastRefresh(settingCacheKey(name, orgId, accId), 0, func() *T {
		var setting = settingDefault[T](name, orgId, accId)
		setting.Value = value
		mysql.Updates(&setting, "value")
		return &setting.Value
	}, true)
}

func settingDefault[T any](name string, orgId uint, accId uint) (setting Setting[T]) {
	if !mysql.Exists(&setting, "name=? AND org_id=? AND acc_id=?", name, orgId, accId) {
		var value T
		setting.Name = name
		setting.OrgId = orgId
		setting.AccId = accId
		setting.Value = value
		mysql.Create(&setting)
	}
	return setting
}

func settingCacheKey(name string, orgId uint, accId uint) string {
	return fmt.Sprintf("setting:%v.%d.%d", name, orgId, accId)
}
