package mysql

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strconv"
)

// Get 通过主键id获取记录
func Get[T any, I string | int | uint](id I) (data T) {
	NotExistsErr(&data, "id=?", id)
	return
}

// GetByColumn 通过条件获取记录
func GetByColumn[T any](condition ...any) (data T) {
	NotExistsErr(&data, condition...)
	return
}

// Begin 开始事务
func Begin() *gorm.DB {
	return Db.Begin()
}

// Create 创建数据记录
func Create(value any) *gorm.DB {
	return Db.Create(value)
}

// Updates 更新指定模型字段
func Updates[T any](model *T, column string, columns ...any) *gorm.DB {
	return Db.Select(column, columns...).Updates(model)
}

// Delete 删除数据记录
func Delete[T any](model *T) *gorm.DB {
	return Db.Delete(model)
}

func Table(name string, args ...any) *gorm.DB {
	return Db.Table(name, args...)
}

func Model(value any) *gorm.DB {
	return Db.Model(value)
}

func Select(query any, args ...any) *gorm.DB {
	return Db.Select(query, args...)
}

func Where(query any, args ...any) *gorm.DB {
	return Db.Where(query, args...)
}

// 设置数据路径，上下级数据结构包含（path、len）字段使用
// model传入必须是指针类型
func SetPath[T any](tx *gorm.DB, model T) {
	s := reflect.ValueOf(model).Elem()
	pid := s.FieldByName("Pid").Uint()
	if pid > 0 {
		var parent T
		TXNotExistsErr(tx, &parent, pid)
		parentValue := reflect.ValueOf(parent).Elem()
		s.FieldByName("Path").SetString(fmt.Sprintf("%v,%d", parentValue.FieldByName("Path").String(), s.FieldByName("Id").Uint()))
		s.FieldByName("Len").SetUint(parentValue.FieldByName("Len").Uint() + 1)
	} else {
		s.FieldByName("Path").SetString(strconv.Itoa(int(s.FieldByName("Id").Uint())))
		s.FieldByName("Len").SetUint(1)
	}
	tx.Model(model).Select("path", "len").Updates(s.Interface())
}

// 设置所有子级路径，一般在更新pid时使用
func SetPathChild[T any](tx *gorm.DB, model T) {
	s := reflect.ValueOf(model).Elem()
	var rows []T
	tx.Find(&rows, "pid=?", s.FieldByName("Id").Uint())
	for _, row := range rows {
		SetPath(tx, row)
		SetPathChild(tx, row)
	}
}

// 设置自身path和所有子级path
func SetPathAndChildPath[T any](tx *gorm.DB, model T) {
	SetPath(tx, model)
	SetPathChild(tx, model)
}
