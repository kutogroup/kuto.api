package models

import (
	"database/sql/driver"
	"time"
)

//ModelTables 所有model的集合
var ModelTables []interface{}

//Time 自定义时间类（主要用于定义时间格式）
type Time time.Time

const (
	timeFormart = "2006-01-02 15:04:05"
)

func (t *Time) Scan(value interface{}) error {
	*t = Time(value.(time.Time))
	return nil
}

func (t Time) Value() (driver.Value, error) {
	return time.Time(t), nil
}

//UnmarshalJSON 反序列化json
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

//MarshalJSON 序列化json
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormart)
}

//================================== time类型的json序列化 ======================================
