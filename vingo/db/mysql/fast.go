package mysql

import (
	"database/sql"
	"fmt"
	"github.com/lgdzz/vingo-utils/vingo"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"strings"
)

func NewDB() *gorm.DB {
	return Db
}

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
func Begin(opts ...*sql.TxOptions) *gorm.DB {
	return Db.Begin(opts...)
}

// Create 创建数据记录
func Create(value any) *gorm.DB {
	return Db.Create(value)
}

// FirstOrCreate 不存在则创建
func FirstOrCreate(dest any, conds ...any) *gorm.DB {
	return Db.FirstOrCreate(dest, conds...)
}

// Updates 更新指定模型字段
func Updates[T any](model *T, column string, columns ...any) *gorm.DB {
	return Db.Select(column, columns...).Updates(model)
}

// Delete 删除数据记录
func Delete[T any](model *T, conds ...any) *gorm.DB {
	return Db.Delete(model, conds...)
}

func Save(value any) *gorm.DB {
	return Db.Save(value)
}

func Debug() *gorm.DB {
	return Db.Debug()
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

func Order(value any) *gorm.DB {
	return Db.Order(value)
}

func Like(db *gorm.DB, keyword string) *gorm.DB {
	if keyword != "" {
		db = db.Where("name like @text OR description like @text", sql.Named("text", SqlLike(keyword)))
	}
	return db
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

// 关键词组装
func SqlLike(keyword string) string {
	return fmt.Sprintf("%%%v%%", strings.Trim(keyword, " "))
}

// like模糊查询
func LikeOr(db *gorm.DB, keyword string, column ...string) *gorm.DB {
	if keyword != "" {
		var s []string
		for _, item := range column {
			s = append(s, fmt.Sprintf("%v like @text", item))
		}
		db = db.Where(strings.Join(s, " OR "), sql.Named("text", SqlLike(keyword)))
	}
	return db
}

// 时间范围查询
func TimeBetween(db *gorm.DB, column string, dateAt vingo.DateAt) *gorm.DB {
	return db.Where(fmt.Sprintf("%v BETWEEN ? AND ?", column), dateAt.Start(), dateAt.End())
}

// 检查字段是否允许被修改
func CheckPatchWhite(field string, whites []string) {
	if !vingo.IsInSlice(field, whites) {
		panic(fmt.Sprintf("字段%v禁止修改", field))
	}
}

func QueryWhere(db *gorm.DB, query any, column string) *gorm.DB {
	valueOf := reflect.ValueOf(query)
	typeOf := valueOf.Type()
	if typeOf.Kind() == reflect.Ptr {
		if valueOf.IsNil() {
			//fmt.Println("空指针无条件")
			return db
		} else {
			query = valueOf.Elem().Interface()
		}
	} else {
		switch v := query.(type) {
		case string:
			if v == "" {
				//fmt.Println("string无条件")
				return db
			}
		}
		query = valueOf.Interface()
	}
	if query != nil {
		db = db.Where(fmt.Sprintf("%v=?", column), query)
	}
	return db
}

func QueryWhereDateAt(db *gorm.DB, query *vingo.DateAt, column string) *gorm.DB {
	if query != nil {
		db = TimeBetween(db, column, *query)
	}
	return db
}

func QueryWhereLike(db *gorm.DB, query string, column ...string) *gorm.DB {
	if query != "" {
		db = LikeOr(db, query, column...)
	}
	return db
}

func QueryWhereBetween(db *gorm.DB, query *[2]any, column string) *gorm.DB {
	if query != nil {
		db = db.Where(fmt.Sprintf("%v BETWEEN ? AND ?", column), query[0], query[1])
	}
	return db
}

func QueryWhereDeletedAt(db *gorm.DB, column string) *gorm.DB {
	db = db.Where(fmt.Sprintf("%v IS NULL", column))
	return db
}
