package component

import (
	"github.com/lgdzz/vingo-utils/vingo"
	"github.com/lgdzz/vingo-utils/vingo/db/mysql"
)

type Shortcut struct {
	Id    string `gorm:"primaryKey;column:id" json:"id"`
	AccId uint   `gorm:"column:acc_id" json:"accId"` // 账户id
	Icon  string `gorm:"column:icon" json:"icon"`    // 图标
	Name  string `gorm:"column:name" json:"name"`    // 名称
	Link  string `gorm:"column:link" json:"link"`    // 地址
	Sort  uint   `gorm:"column:sort" json:"sort"`    // 排序
}

func (s *Shortcut) TableName() string {
	return "shortcut"
}

// 快捷方式列表
func ShortcutList(c *vingo.Context) {
	var rows []Shortcut
	mysql.Where("acc_id=?", c.GetAccId()).Order("sort asc").Find(&rows)
	c.ResponseBody(rows)
}

// 添加快捷方式
func ShortcutAdd(c *vingo.Context) {
	var body = vingo.GetRequestBody[Shortcut](c)
	body.Id = vingo.GetUUID()
	body.AccId = c.GetAccId()
	if !mysql.Exists(&Shortcut{}, "acc_id=? AND name=? AND link=?", body.AccId, body.Name, body.Link) {
		mysql.Model(&Shortcut{}).Where("acc_id=?", body.AccId).Select("sort").Order("sort desc").Scan(&body.Sort)
		body.Sort++
		mysql.Create(&body)
	}
	c.ResponseSuccess()
}

// 删除快捷方式
func ShortcutDel(c *vingo.Context) {
	var row = mysql.Get[Shortcut](c.Param("id"))
	mysql.Delete(&row)
	c.ResponseSuccess()
}

// 快捷方式排序
func ShortcutSort(c *vingo.Context) {
	var body = vingo.GetRequestBody[struct {
		Ids []string `json:"ids"`
	}](c)
	tx := mysql.Begin()
	defer mysql.AutoCommit(tx)
	for sort, id := range body.Ids {
		tx.Model(&Shortcut{}).Where("id=?", id).Update("sort", sort)
	}
	c.ResponseSuccess()
}
