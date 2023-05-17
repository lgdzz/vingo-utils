package db

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type DbHelper struct{}

func LikeKeyword(keyword string) string {
	return fmt.Sprintf("%%%v%%", keyword)
}

// 分页查询
func NewPage(db *gorm.DB, p *PageResult, args ...any) *PageResult {
	var count int64
	db.Count(&count)
	p.Total = count
	if count > 0 {
		if len(args) > 0 {
			order := "`id` desc"
			switch o := args[0].(type) {
			case string:
				order = o
			case *OrderBy:
				if o.SortField != "" {
					if o.SortOrder != "asc" && o.SortOrder != "desc" {
						panic("无效的排序")
					}
					order = fmt.Sprintf("`%v` %v", o.SortField, o.SortOrder)
				}
			}
			db = db.Order(order)
		}
		switch p.Items {
		case "map":
			var items = []map[string]any{}
			db.Limit(p.GetSize()).Offset(int(p.Offset())).Scan(&items)
			p.Items = &items
		default:
			db.Limit(p.GetSize()).Offset(int(p.Offset())).Find(&p.Items)
		}

	} else {
		if p.Items == "map" {
			p.Items = []map[string]any{}
		}
	}
	return p
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
	err := Pool.First(model, condition...).Error
	if err == gorm.ErrRecordNotFound {
		return false
	} else if err != nil {
		panic(err.Error())
	}
	return true
}

// 记录不存在时抛出错误
func NotExistsErr(model any, condition ...any) {
	err := Pool.First(model, condition...).Error
	if err == gorm.ErrRecordNotFound {
		panic(err.Error())
	} else if err != nil {
		panic(err.Error())
	}
}

// 记录不存在时抛出错误
func NotExistsErrMsg(msg string, model any, condition ...any) {
	err := Pool.First(model, condition...).Error
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
	err := Pool.First(model, "pid=?", id)
	if err.Error != gorm.ErrRecordNotFound {
		panic("记录有子项，删除失败")
	}
}

func (p *PageResult) GetPage() int {
	if p.Page > 0 {
		return p.Page
	} else {
		return 1
	}
}

func (p *PageResult) GetSize() int {
	if p.Size > 0 {
		return p.Size
	} else {
		return 10
	}
}

func (p *PageResult) Offset() int64 {
	if p.Page > 0 {
		return int64((p.Page - 1) * p.Size)
	} else {
		return 0
	}
}

type PageResult struct {
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Total int64 `json:"total"` // 总的记录数
	Items any   `json:"items"` // 查询数据列表
}

type OrderBy struct {
	SortField string `form:"sortField"`
	SortOrder string `form:"sortOrder"`
}

type PageQuery struct {
	Page int `form:"page"`
	Size int `form:"size"`
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
