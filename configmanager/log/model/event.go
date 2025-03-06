package model

import (
	"time"
)

// EventType 事件类型
type EventType int

const (
	EventTypeSecurity EventType = 1 // 安全事件
	EventTypeFault    EventType = 2 // 故障事件
)

// Event 事件记录
type Event struct {
	EventID   int64     `json:"eventId" gorm:"column:eventId"`
	DeviceID  string    `json:"deviceId" gorm:"column:deviceId"`
	EventTime time.Time `json:"eventTime" gorm:"column:eventTime"`
	EventType EventType `json:"eventType" gorm:"column:eventType"`
	EventCode string    `json:"eventCode" gorm:"column:eventCode"`
	EventDesc string    `json:"eventDesc" gorm:"column:eventDesc"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt"`
}

// TableName 指定表名
func (Event) TableName() string {
	return "events"
}
