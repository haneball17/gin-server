package models

import (
	"time"

	"gorm.io/gorm"
)

// UserBehavior 用户行为
type UserBehavior struct {
	gorm.Model
	UserID       int       `json:"user_id" gorm:"column:user_id;index;not null"`             // 用户ID
	BehaviorTime time.Time `json:"behavior_time" gorm:"column:behavior_time;index;not null"` // 行为开始时间
	BehaviorType int       `json:"behavior_type" gorm:"column:behavior_type;index;not null"` // 行为类型，1:发送，2:接收
	DataType     int       `json:"data_type" gorm:"column:data_type;not null"`               // 数据类型，1:文件，2:消息
	DataSize     int64     `json:"data_size" gorm:"column:data_size;not null"`               // 数据大小
}

// TableName 指定表名
func (UserBehavior) TableName() string {
	return "user_behaviors"
}

// UserInfo 用户信息（此结构体仅用于聚合数据，不映射到数据库表）
type UserInfo struct {
	UserID         int            `json:"user_id"`         // 用户id
	Status         int            `json:"status"`          // 用户状态
	OnlineDuration int            `json:"online_duration"` // 在线时长
	Behaviors      []UserBehavior `json:"behaviors"`       // 行为列表
}
