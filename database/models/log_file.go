package models

import (
	"time"

	"gorm.io/gorm"
)

// LogFile 日志文件
type LogFile struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index:idx_log_files_deleted_at"`
	FileName     string         `json:"file_name" gorm:"column:file_name;uniqueIndex:idx_log_files_file_name;type:varchar(255)"`
	FilePath     string         `json:"file_path" gorm:"column:file_path;not null;type:varchar(255)"`
	FileSize     int64          `json:"file_size" gorm:"column:file_size;not null"`
	StartTime    time.Time      `json:"start_time" gorm:"column:start_time;index:idx_log_files_start_time"`
	EndTime      time.Time      `json:"end_time" gorm:"column:end_time;index:idx_log_files_end_time"`
	IsEncrypted  bool           `json:"is_encrypted" gorm:"column:is_encrypted;default:false"`
	IsUploaded   bool           `json:"is_uploaded" gorm:"column:is_uploaded;default:false"`
	RemotePath   string         `json:"remote_path" gorm:"column:remote_path;default:'';type:varchar(255)"`
	UploadedTime *time.Time     `json:"uploaded_time" gorm:"column:uploaded_time"`
}

// TableName 指定表名
func (LogFile) TableName() string {
	return "log_files"
}

// LogTimeRange 日志时间范围
type LogTimeRange struct {
	StartTime time.Time `json:"start_time"` // 统计起始时间
	Duration  int64     `json:"duration"`   // 统计时长（秒）
}

// LogContentLog 日志内容在标准格式中的结构
type LogContentLog struct {
	// 统计时间区间
	TimeRange struct {
		StartTime string `json:"start_time"` // ISO8601格式
		Duration  int64  `json:"duration"`   // 统计时长（秒）
	} `json:"time_range"`

	// 安全事件
	SecurityEvents struct {
		Events []Event `json:"events"` // 安全事件列表
	} `json:"security_events"`

	// 性能事件
	PerformanceEvents struct {
		SecurityDevices []SecurityDeviceLog `json:"security_devices"` // 安全接入管理设备列表
	} `json:"performance_events"`

	// 故障事件
	FaultEvents struct {
		Events []Event `json:"events"` // 故障事件列表
	} `json:"fault_events"`
}

// LogContent 日志内容结构（不映射到数据库表）
type LogContent struct {
	// 统计时间区间
	TimeRange LogTimeRange `json:"time_range"`

	// 安全事件
	SecurityEvents struct {
		Events []Event `json:"events"` // 安全事件列表
	} `json:"security_events"`

	// 性能事件
	PerformanceEvents struct {
		SecurityDevices []SecurityDevice `json:"security_devices"` // 安全接入管理设备列表
	} `json:"performance_events"`

	// 故障事件
	FaultEvents struct {
		Events []Event `json:"events"` // 故障事件列表
	} `json:"fault_events"`
}

// ToLogContentLog 将LogContent转换为LogContentLog
func (lc *LogContent) ToLogContentLog() LogContentLog {
	// 创建标准格式日志内容
	logContentLog := LogContentLog{}

	// 转换时间范围
	logContentLog.TimeRange.StartTime = lc.TimeRange.StartTime.Format(time.RFC3339)
	logContentLog.TimeRange.Duration = lc.TimeRange.Duration

	// 复制安全事件
	logContentLog.SecurityEvents.Events = lc.SecurityEvents.Events

	// 转换安全设备
	logContentLog.PerformanceEvents.SecurityDevices = make([]SecurityDeviceLog, 0, len(lc.PerformanceEvents.SecurityDevices))
	for _, device := range lc.PerformanceEvents.SecurityDevices {
		logContentLog.PerformanceEvents.SecurityDevices = append(
			logContentLog.PerformanceEvents.SecurityDevices,
			device.ToSecurityDeviceLog(),
		)
	}

	// 复制故障事件
	logContentLog.FaultEvents.Events = lc.FaultEvents.Events

	return logContentLog
}

// SecurityDeviceLog 安全接入管理设备在日志中的格式
type SecurityDeviceLog struct {
	DeviceID       int                `json:"device_id"`       // 安全接入管理设备id
	CPUUsage       int                `json:"cpu_usage"`       // 峰值CPU占用率
	MemoryUsage    int                `json:"memory_usage"`    // 峰值内存使用率
	OnlineDuration int                `json:"online_duration"` // 设备在线时间
	Status         int                `json:"status"`          // 设备状态
	GatewayDevices []GatewayDeviceLog `json:"gateway_devices"` // 网关设备列表
}

// GatewayDeviceLog 网关设备在日志中的格式
type GatewayDeviceLog struct {
	DeviceID       int           `json:"device_id"`       // 网关设备id
	CPUUsage       int           `json:"cpu_usage"`       // 峰值CPU占用率
	MemoryUsage    int           `json:"memory_usage"`    // 峰值内存使用率
	OnlineDuration int           `json:"online_duration"` // 设备在线时间
	Status         int           `json:"status"`          // 设备状态
	Users          []UserInfoLog `json:"users"`           // 用户列表
}

// SecurityDevice 安全接入管理设备（不映射到数据库表）
type SecurityDevice struct {
	DeviceID       int             `json:"device_id"`       // 安全接入管理设备id
	CPUUsage       int             `json:"cpu_usage"`       // 峰值CPU占用率
	MemoryUsage    int             `json:"memory_usage"`    // 峰值内存使用率
	OnlineDuration int             `json:"online_duration"` // 设备在线时间
	Status         int             `json:"status"`          // 设备状态
	GatewayDevices []GatewayDevice `json:"gateway_devices"` // 网关设备列表
}

// ToSecurityDeviceLog 将SecurityDevice转换为SecurityDeviceLog
func (sd *SecurityDevice) ToSecurityDeviceLog() SecurityDeviceLog {
	securityDeviceLog := SecurityDeviceLog{
		DeviceID:       sd.DeviceID,
		CPUUsage:       sd.CPUUsage,
		MemoryUsage:    sd.MemoryUsage,
		OnlineDuration: sd.OnlineDuration,
		Status:         sd.Status,
		GatewayDevices: make([]GatewayDeviceLog, 0, len(sd.GatewayDevices)),
	}

	// 转换每个网关设备
	for _, gateway := range sd.GatewayDevices {
		securityDeviceLog.GatewayDevices = append(
			securityDeviceLog.GatewayDevices,
			gateway.ToGatewayDeviceLog(),
		)
	}

	return securityDeviceLog
}

// GatewayDevice 网关设备（不映射到数据库表）
type GatewayDevice struct {
	DeviceID       int        `json:"device_id"`       // 网关设备id
	CPUUsage       int        `json:"cpu_usage"`       // 峰值CPU占用率
	MemoryUsage    int        `json:"memory_usage"`    // 峰值内存使用率
	OnlineDuration int        `json:"online_duration"` // 设备在线时间
	Status         int        `json:"status"`          // 设备状态
	Users          []UserInfo `json:"users"`           // 用户列表
}

// ToGatewayDeviceLog 将GatewayDevice转换为GatewayDeviceLog
func (gd *GatewayDevice) ToGatewayDeviceLog() GatewayDeviceLog {
	gatewayDeviceLog := GatewayDeviceLog{
		DeviceID:       gd.DeviceID,
		CPUUsage:       gd.CPUUsage,
		MemoryUsage:    gd.MemoryUsage,
		OnlineDuration: gd.OnlineDuration,
		Status:         gd.Status,
		Users:          make([]UserInfoLog, 0, len(gd.Users)),
	}

	// 转换每个用户
	for _, user := range gd.Users {
		gatewayDeviceLog.Users = append(
			gatewayDeviceLog.Users,
			user.ToUserInfoLog(),
		)
	}

	return gatewayDeviceLog
}
