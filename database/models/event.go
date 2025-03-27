package models

import (
	"time"

	"gorm.io/gorm"
)

// EventType 事件类型
type EventType int

const (
	EventTypeSecurity EventType = 1 // 安全事件
	EventTypeFault    EventType = 2 // 故障事件
)

// Event 事件记录
type Event struct {
	gorm.Model
	EventID   int64     `json:"event_id" gorm:"column:event_id;uniqueIndex"`
	DeviceID  int       `json:"device_id" gorm:"column:device_id;index"`
	EventTime time.Time `json:"event_time" gorm:"column:event_time;index"`
	EventType EventType `json:"event_type" gorm:"column:event_type;index"`
	EventCode string    `json:"event_code" gorm:"column:event_code;type:varchar(64)"`
	EventDesc string    `json:"event_desc" gorm:"column:event_desc;type:varchar(255)"`
}

// TableName 指定表名
func (Event) TableName() string {
	return "events"
}
