package page

import (
	"fmt"
	"github.com/lgdzz/vingo-utils/vingo"
	"gorm.io/gorm"
	"strings"
)

// 分页查询处理
// page.New[flow.Approval](pool, page.Option{
// 		Limit: page.Limit{Page: query.Page, Size: query.Size},
// })

type Result struct {
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Total int64 `json:"total"` // 总的记录数
	Items any   `json:"items"` // 查询数据列表
}

type Limit struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

func (s *Limit) GetPage() int {
	if s.Page > 0 {
		return s.Page
	} else {
		return 1
	}
}

func (s *Limit) GetSize() int {
	if s.Size > 0 {
		return s.Size
	} else {
		return 10
	}
}

func (s *Limit) Offset() int64 {
	if s.Page > 0 {
		return int64((s.Page - 1) * s.Size)
	} else {
		return 0
	}
}

type Order struct {
	Column string `form:"sortField"`
	Sort   string `form:"sortOrder"`
}

func (s *Order) HandleColumn() string {
	var items = strings.Split(s.Column, ".")
	for index, item := range items {
		items[index] = "`" + item + "`"
	}
	return fmt.Sprintf("%v %v", strings.Join(items, "."), s.Sort)
}

type Option struct {
	Limit  Limit
	Order  []Order
	Handle func(item any, index int) any
}

// 创建一个新的分页查询
func New[T any](db *gorm.DB, option Option, handle func(T) any) (result Result) {
	var count int64
	var items = make([]T, 0)
	db.Count(&count)
	result.Total = count
	result.Page = option.Limit.GetPage()
	result.Size = option.Limit.GetSize()
	if count > 0 {
		if len(option.Order) == 0 {
			db = db.Order("`id` desc")
		} else {
			var orders = make([]string, 0)
			for _, order := range option.Order {
				orders = append(orders, order.HandleColumn())
			}
			db = db.Order(strings.Join(orders, ","))
		}

		db.Limit(option.Limit.GetSize()).Offset(int(option.Limit.Offset())).Scan(&items)

		if handle != nil {
			result.Items = vingo.ForEach(items, func(item T, index int) any {
				return handle(item)
			})
			return
		}
	}
	result.Items = items
	return
}