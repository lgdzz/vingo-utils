package vingo

import (
	"database/sql/driver"
	"gorm.io/gorm"
	"strings"
	"time"
)

// 定位坐标
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// 文件信息
type FileInfo struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Mimetype  string `json:"mimetype"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	Realpath  string `json:"realpath"`
}

// 文件信息(简单)
type FileInfoSimple struct {
	Name     string `json:"name"`
	Realpath string `json:"realpath"`
}

// 身份证信息
// IdCardInfo{IdCard: ""}
type IdCardInfo struct {
	IdCard     string // 身份证号码
	RegionCode string // 6位行政区域编码
	Birthday   string // 2006-01-02 格式日期
	Age        int    // 年龄：精确到月份
	UniformAge int    // 年龄：按年份计算
	Gender     string // 性别
}

// 身份验证异常
type AuthException struct {
	Message string
}

// 数据库事务异常
type DbException struct {
	Message string
}

// 时间范围
type DateRange struct {
	Start time.Time
	End   time.Time
}

type QueryDateRange struct {
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
}

// 响应数据
type ResponseData struct {
	Status    int    // 状态
	Error     int8   // 0-无错误|1-有错误
	ErrorType string // 错误类型
	Message   string // 消息
	Data      any    // 返回数据内容
	NoLog     bool   // true时不记录日志
}

type DbModel struct {
	ID        uint           `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt *LocalTime     `gorm:"column:created_at;" json:"createdAt"`
	UpdatedAt *LocalTime     `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deletedAt"`
}

type UintIds []uint

func (s UintIds) Value() (driver.Value, error) {
	return strings.Join(SliceUintToString(s), ","), nil
}

func (s *UintIds) Scan(value interface{}) error {
	v := string(value.([]byte))
	if v == "" {
		s = &UintIds{}
	} else {
		CustomOutput(SliceStringToUint(strings.Split(v, ",")), s)
	}
	return nil
}

func (s *UintIds) Uints() (result []uint) {
	CustomOutput(s, &result)
	return
}

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	return strings.Join(s, ","), nil
}

func (s *StringSlice) Scan(value interface{}) error {
	v := string(value.([]byte))
	if v == "" {
		s = &StringSlice{}
	} else {
		*s = strings.Split(v, ",")
	}
	return nil
}

func (s *StringSlice) Strings() (result []string) {
	CustomOutput(s, &result)
	return
}

type IpInfo struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	Area    string `json:"area"`
}
