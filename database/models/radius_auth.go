package models

import (
	"time"
)

// RadPostAuth Radius认证记录表
// 对应radius库中的radpostauth表
type RadPostAuth struct {
	ID       int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Username string    `json:"username" gorm:"column:username;type:varchar(64);not null"`
	Pass     string    `json:"pass" gorm:"column:pass;type:varchar(64);not null"`
	Reply    string    `json:"reply" gorm:"column:reply;type:varchar(32);not null"`
	AuthDate time.Time `json:"authdate" gorm:"column:authdate;type:timestamp(6);not null;default:CURRENT_TIMESTAMP(6)"`
	Class    string    `json:"class" gorm:"column:class;type:varchar(64);null"`
}

// TableName 指定表名
// 这里明确指定表名，不使用GORM的默认命名规则（蛇形命名法）
func (RadPostAuth) TableName() string {
	return "radpostauth"
}

// RadPostAuthQuery 查询条件结构体
type RadPostAuthQuery struct {
	Username  string `form:"username"`
	Reply     string `form:"reply"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
	Class     string `form:"class"`
}
