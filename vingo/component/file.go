package component

import (
	"github.com/lgdzz/vingo-utils/vingo"
	"github.com/lgdzz/vingo-utils/vingo/db/mysql"
	"github.com/lgdzz/vingo-utils/vingo/db/page"
)

// File 文件
type File struct {
	ID        uint             `gorm:"primaryKey;column:id" json:"id"`
	CID       uint             `gorm:"column:c_id" json:"cId"`          // 分类ID，预留
	Channel   string           `gorm:"column:channel" json:"channel"`   // 渠道
	OrgID     uint             `gorm:"column:org_id" json:"orgId"`      // 关联组织ID
	FromID    string           `gorm:"column:from_id" json:"fromId"`    // 来源ID
	Type      string           `gorm:"column:type" json:"type"`         // 文件类型
	Filename  string           `gorm:"column:filename" json:"filename"` // 资源名称
	Filepath  string           `gorm:"column:filepath" json:"filepath"` // 资源路径
	Filesize  uint             `gorm:"column:filesize" json:"filesize"` // 文件大小
	Mimetype  string           `gorm:"column:mimetype" json:"mimetype"`
	Extension string           `gorm:"column:extension" json:"extension"` // 文件后缀
	Extra     string           `gorm:"column:extra" json:"extra"`         // 附件属性
	CreatedAt *vingo.LocalTime `gorm:"column:created_at;" json:"createdAt"`
	UpdatedAt *vingo.LocalTime `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName get sql table name.获取数据库表名
func (m *File) TableName() string {
	return "file"
}

// 文件记录保存
func FileRecordSave(c *vingo.Context) {
	body := vingo.GetRequestBody[File](c)
	body.OrgID = c.GetOrgId()
	body.Channel = "backstage"
	mysql.Create(&body)
	c.ResponseSuccess()
}

// 文件记录列表
func FileRecordList(c *vingo.Context) {
	query := vingo.GetRequestQuery[page.Limit](c)
	db := mysql.Table("file").Where("org_id=?", c.GetOrgId())
	result := page.New[File](db, page.Option{
		Limit: page.Limit{Page: query.Page, Size: query.Size},
		Order: page.OrderDefault(nil),
	}, nil)
	c.ResponseBody(result)
}
