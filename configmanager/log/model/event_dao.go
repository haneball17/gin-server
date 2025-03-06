package model

import (
	"time"

	"gorm.io/gorm"
)

// EventDAO 事件数据访问对象
type EventDAO struct {
	db *gorm.DB
}

// NewEventDAO 创建事件DAO实例
func NewEventDAO(db *gorm.DB) *EventDAO {
	return &EventDAO{db: db}
}

// GetEventsByTimeRange 获取指定时间范围内的事件
func (dao *EventDAO) GetEventsByTimeRange(startTime, endTime time.Time) ([]Event, error) {
	var events []Event
	err := dao.db.Where("event_time BETWEEN ? AND ?", startTime, endTime).Find(&events).Error
	return events, err
}

// GetEventsByTypeAndTimeRange 获取指定类型和时间范围内的事件
func (dao *EventDAO) GetEventsByTypeAndTimeRange(eventType EventType, startTime, endTime time.Time) ([]Event, error) {
	var events []Event
	err := dao.db.Where("event_type = ? AND event_time BETWEEN ? AND ?",
		eventType, startTime, endTime).Find(&events).Error
	return events, err
}
