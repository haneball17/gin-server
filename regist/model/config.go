package model

import (
	"gin-server/config"
)

// LoadConfig 已废弃，请使用 config.GetConfig()
// 为了兼容现有代码，此函数返回与原函数相同的结构
func LoadConfig() *config.Config {
	return config.GetConfig()
}
