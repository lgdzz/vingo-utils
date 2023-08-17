package mysql

import (
	"fmt"
	"github.com/lgdzz/vingo-utils/vingo"
	"gorm.io/gorm"
	"strings"
)

func LikeKeyword(keyword string) string {
	return fmt.Sprintf("%%%v%%", keyword)
}

// 组织like语句，返回值直接放在Where中使用
func NewLike(column string, value string) string {
	return fmt.Sprintf("%v like '%v'", column, LikeKeyword(value))
}

// 组织多个like语句，OR条件，返回值直接放在Where中使用
func NewLikes(columns *[]string, value string) string {
	strArray := []string{}
	value = LikeKeyword(value)
	for _, c := range *columns {
		strArray = append(strArray, fmt.Sprintf("%v like '%v'", c, value))
	}
	//fmt.Println(strings.Join(sql, " OR "))
	return strings.Join(strArray, " OR ")
}

// 指定字段第一个汉字按A-Z排序
func NewChineseSort(column string) string {
	return fmt.Sprintf("CONVERT(SUBSTR(%v, 1, 1) USING gbk)", column)
}

func Exists(model any, condition ...any) bool {
	err := Db.First(model, condition...).Error
	if err == gorm.ErrRecordNotFound {
		return false
	} else if err != nil {
		panic(err.Error())
	}
	return true
}

// 记录不存在时抛出错误
func NotExistsErr(model any, condition ...any) {
	err := Db.First(model, condition...).Error
	if err == gorm.ErrRecordNotFound {
		panic(err.Error())
	} else if err != nil {
		panic(err.Error())
	}
}

// 记录不存在时抛出错误
func NotExistsErrMsg(msg string, model any, condition ...any) {
	err := Db.First(model, condition...).Error
	if err == gorm.ErrRecordNotFound {
		panic(msg)
	} else if err != nil {
		panic(err.Error())
	}
}

// 记录不存在时抛出错误(事务内)
func TXNotExistsErr(tx *gorm.DB, model any, condition ...any) {
	err := tx.First(model, condition...).Error
	if err == gorm.ErrRecordNotFound {
		panic(err.Error())
	} else if err != nil {
		panic(err.Error())
	}
}

func CheckHasChild(model any, id uint) {
	err := Db.First(model, "pid=?", id)
	if err.Error != gorm.ErrRecordNotFound {
		panic("记录有子项，删除失败")
	}
}

// 数据库事务自动提交
func AutoCommit(tx *gorm.DB, callback ...func()) {
	if r := recover(); r != nil {
		//fmt.Printf("%T\n%v\n", r, r)
		tx.Rollback()
		if len(callback) > 0 && callback[0] != nil {
			callback[0]()
		}
		panic(r)
	} else if err := tx.Statement.Error; err != nil {
		//fmt.Println("数据库异常事务回滚")
		tx.Rollback()
		if len(callback) > 0 && callback[0] != nil {
			callback[0]()
		}
		panic(err.Error())
	} else {
		//fmt.Println("事务提交")
		tx.Commit()
		if len(callback) > 1 && callback[1] != nil {
			callback[1]()
		}
	}
}

type TableColumn struct {
	Column  string `gorm:"column:Field" json:"column"`
	Type    string `gorm:"column:Type" json:"type"`
	Comment string `gorm:"column:Comment" json:"comment"`
}

// 获取表字段
func GetTableColumn(tableName string) []TableColumn {
	var columns []TableColumn
	Db.Raw("SHOW FULL COLUMNS FROM " + tableName).Select("Field,Type,Comment").Scan(&columns)
	for index, item := range columns {
		if vingo.StringContainsOr(item.Type, []string{"int", "tinyint", "bigint", "float", "decimal"}) {
			columns[index].Type = "number"
		} else if vingo.StringContainsOr(item.Type, []string{"char", "varchar", "text", "longtext"}) {
			columns[index].Type = "string"
		} else if vingo.StringContainsOr(item.Type, []string{"datetime"}) {
			columns[index].Type = "datetime"
		}
	}
	return columns
}
