package vingo

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type LocalTime time.Time

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t).Local()
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format(DatetimeFormat))), nil
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	var err error
	var parsedTime time.Time
	if string(data) == "null" {
		*t = LocalTime{}
		return nil
	}

	parsedTime, err = time.ParseInLocation(`"`+DatetimeFormat+`"`, string(data), time.Local)
	if err != nil {
		return err
	}

	*t = LocalTime(parsedTime.Local())
	return nil
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t).Local()
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value.Local())
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t LocalTime) Now() LocalTime {
	return LocalTime(time.Now().Local())
}

func (t *LocalTime) SetNow() {
	*t = LocalTime(time.Now().Local())
}

func (t *LocalTime) To(value time.Time) {
	*t = LocalTime(value.Local())
}

func (t LocalTime) String() string {
	return time.Time(t).Format(DatetimeFormat)
}

func (t LocalTime) Time() time.Time {
	return time.Time(t)
}

func (t LocalTime) Format(layout string) string {
	return t.Time().Format(layout)
}

func (t *LocalTime) ScanFromRow(rows *sql.Rows, columnName string) error {
	var tmp time.Time
	err := rows.Scan(&tmp)
	if err != nil {
		return err
	}
	*t = LocalTime(tmp.Local())
	return nil
}

func (t LocalTime) ValueFromRow(rows *sql.Rows, columnName string) (interface{}, error) {
	return t, nil
}

func TimeAddDays(t time.Time, days int, hour int, min int, sec int) time.Time {
	// Add the specified number of days
	t = t.AddDate(0, 0, days)

	// Set the time to midnight
	year, month, day := t.Date()

	midnight := time.Date(year, month, day, hour, min, sec, 0, t.Location())

	return midnight
}

// 判断当前时间是否大于指定的时间
func TimeIsAfterNow(t time.Time) bool {
	// 获取当前时间
	now := time.Now()

	// 判断当前时间是否大于指定时间
	return now.After(t)
}

// 获取昨日开始时间
func GetYesterdayStartTime() time.Time {
	now := time.Now().Local()
	yesterday := now.AddDate(0, 0, -1)
	return time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())
}

// 获取日期范围
func GetDateRange(dateStr string) (string, string) {
	now := time.Now().Local()
	var start, end time.Time
	var err error

	switch dateStr {
	case "yesterday":
		start = GetYesterdayStartTime()
	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	default:
		start, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			panic(err.Error())
		}
	}
	end = time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())

	return start.Format(DatetimeFormat), end.Format(DatetimeFormat)
}

// 是否超过时间
func IsTimeExceeded(t time.Time, days int) bool {
	duration := time.Since(t)
	return duration > time.Duration(days)*24*time.Hour
}

// 获取最近一周日期
func GetLastWeekDates(startDay ...time.Time) []string {
	var t time.Time
	if len(startDay) == 0 {
		t = time.Now()
	} else {
		t = startDay[0]
	}
	t = t.AddDate(0, 0, -6)

	var lastWeek []string
	for i := 0; i < 7; i++ {
		day := t.AddDate(0, 0, i).Format("2006-01-02")
		lastWeek = append(lastWeek, day)
	}
	return lastWeek
}

// 获取最近一月日期
func GetLastMonthDates(startDay ...time.Time) []string {
	var t time.Time
	if len(startDay) == 0 {
		t = time.Now()
	} else {
		t = startDay[0]
	}
	t = t.AddDate(0, 0, -29)

	var lastWeek []string
	for i := 0; i < 30; i++ {
		day := t.AddDate(0, 0, i).Format("2006-01-02")
		lastWeek = append(lastWeek, day)
	}
	return lastWeek
}

// 获取当前时间值指针
func GetNowTime() *LocalTime {
	t := LocalTime{}.Now()
	return &t
}

// 获取当前时间值
func GetNowTimeValue() LocalTime {
	return LocalTime{}.Now()
}

// 获取今年开始时间
func GetThisYearStartTime() time.Time {
	now := time.Now()        // 获取当前时间
	year, _, _ := now.Date() // 获取当前年份
	return time.Date(year, 1, 1, 0, 0, 0, 0, now.Location())
}

// 获取今年结束时间
func GetThisYearEndTime() time.Time {
	now := time.Now()        // 获取当前时间
	year, _, _ := now.Date() // 获取当前年份
	return time.Date(year, 12, 31, 23, 59, 59, 0, now.Location())
}
