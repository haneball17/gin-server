package models

import (
	"time"

	"gorm.io/gorm"
)

// UserBehavior 用户行为
type UserBehavior struct {
	// 行为ID，设置为主键并自增
	BehaviorID int `json:"behavior_id" gorm:"column:behavior_id;primaryKey;autoIncrement;not null"`
	// 其他字段
	UserID       int       `json:"user_id" gorm:"column:user_id;index;not null"`             // 用户ID
	BehaviorTime time.Time `json:"behavior_time" gorm:"column:behavior_time;index;not null"` // 行为开始时间
	BehaviorType int       `json:"behavior_type" gorm:"column:behavior_type;index;not null"` // 行为类型，1:发送，2:接收
	DataType     int       `json:"data_type" gorm:"column:data_type;not null"`               // 数据类型，1:文件，2:消息
	DataSize     int64     `json:"data_size" gorm:"column:data_size;not null"`               // 数据大小

	// GORM的标准字段
	CreatedAt time.Time      `json:"created_at"`              // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`              // 更新时间
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"` // 删除时间（软删除）
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

// BehaviorLog 行为数据在日志中的格式（仅用于日志输出，不映射到数据库表）
type BehaviorLog struct {
	Time     string `json:"time"`      // ISO8601格式的时间
	Type     int    `json:"type"`      // 行为类型，1:发送，2:接收
	DataType int    `json:"data_type"` // 数据类型，1:文件，2:消息
	DataSize int64  `json:"data_size"` // 数据大小（字节）
}

// UserInfoLog 用户信息在日志中的格式（仅用于日志输出，不映射到数据库表）
type UserInfoLog struct {
	UserID         int           `json:"user_id"`         // 用户id
	Status         int           `json:"status"`          // 用户状态
	OnlineDuration int           `json:"online_duration"` // 在线时长
	Behaviors      []BehaviorLog `json:"behaviors"`       // 行为列表（使用简化的BehaviorLog格式）
}

// ToUserInfoLog 将UserInfo转换为UserInfoLog
func (ui *UserInfo) ToUserInfoLog() UserInfoLog {
	userInfoLog := UserInfoLog{
		UserID:         ui.UserID,
		Status:         ui.Status,
		OnlineDuration: ui.OnlineDuration,
		Behaviors:      make([]BehaviorLog, 0, len(ui.Behaviors)),
	}

	// 转换每个行为记录为简化格式
	for _, behavior := range ui.Behaviors {
		userInfoLog.Behaviors = append(userInfoLog.Behaviors, behavior.ToBehaviorLog())
	}

	return userInfoLog
}

// ToBehaviorLog 将UserBehavior转换为BehaviorLog
func (ub *UserBehavior) ToBehaviorLog() BehaviorLog {
	return BehaviorLog{
		Time:     ub.BehaviorTime.Format(time.RFC3339),
		Type:     ub.BehaviorType,
		DataType: ub.DataType,
		DataSize: ub.DataSize,
	}
}
