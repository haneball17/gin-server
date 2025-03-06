package model

import "time"

// LogFile 日志文件结构
type LogFile struct {
	// 统计时间区间
	TimeRange struct {
		StartTime time.Time `json:"startTime"` // 统计起始时间
		Duration  int64     `json:"duration"`  // 统计时长（秒）
	} `json:"timeRange"`

	// 安全事件
	SecurityEvents struct {
		Events []Event `json:"events"` // 安全事件列表
	} `json:"securityEvents"`

	// 性能事件
	PerformanceEvents struct {
		SecurityDevices []SecurityDevice `json:"securityDevices"` // 安全接入管理设备列表
	} `json:"performanceEvents"`

	// 故障事件
	FaultEvents struct {
		Events []Event `json:"events"` // 故障事件列表
	} `json:"faultEvents"`
}

// SecurityDevice 安全接入管理设备
type SecurityDevice struct {
	DeviceID       string          `json:"deviceId"`       // 安全接入管理设备id
	CPUUsage       int             `json:"cpuUsage"`       // 峰值CPU占用率
	MemoryUsage    int             `json:"memoryUsage"`    // 峰值内存使用率
	OnlineDuration int             `json:"onlineDuration"` // 设备在线时间
	Status         int             `json:"status"`         // 设备状态
	GatewayDevices []GatewayDevice `json:"gatewayDevices"` // 网关设备列表
}

// GatewayDevice 网关设备
type GatewayDevice struct {
	DeviceID       string `json:"deviceId"`       // 网关设备id
	CPUUsage       int    `json:"cpuUsage"`       // 峰值CPU占用率
	MemoryUsage    int    `json:"memoryUsage"`    // 峰值内存使用率
	OnlineDuration int    `json:"onlineDuration"` // 设备在线时间
	Status         int    `json:"status"`         // 设备状态
	Users          []User `json:"users"`          // 用户列表
}

// User 用户信息
type User struct {
	UserID         int        `json:"userId"`         // 用户id
	Status         int        `json:"status"`         // 用户状态
	OnlineDuration int        `json:"onlineDuration"` // 在线时长
	Behaviors      []Behavior `json:"behaviors"`      // 行为列表
}

// Behavior 用户行为
type Behavior struct {
	Time     time.Time `json:"time"`     // 发生时间
	Type     int       `json:"type"`     // 行为类型
	DataType int       `json:"dataType"` // 数据类型
	DataSize int64     `json:"dataSize"` // 数据大小
}
