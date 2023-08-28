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

// Updates 更新指定模型字段
func Updates(model any, column string, columns ...any) {
	Db.Select(column, columns...).Updates(model)
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
