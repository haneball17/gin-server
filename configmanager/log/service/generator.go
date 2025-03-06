package service

import (
	"encoding/json"
	"fmt"
	"time"

	"gin-server/configmanager/common/alert"
	"gin-server/configmanager/common/fileutil"
	"gin-server/configmanager/log/model"

	"gorm.io/gorm"
)

const (
	// 设备类型常量
	DeviceTypeSecurityMgmt = 4 // 安全接入管理设备
	DeviceTypeGatewayA     = 1 // 网关设备A
	DeviceTypeGatewayB     = 2 // 网关设备B
	DeviceTypeGatewayC     = 3 // 网关设备C
)

// Generator 日志生成器
type Generator struct {
	db      *gorm.DB
	alerter alert.Alerter
}

// NewGenerator 创建日志生成器实例
func NewGenerator(db *gorm.DB, alerter alert.Alerter) *Generator {
	return &Generator{
		db:      db,
		alerter: alerter,
	}
}

// Generate 生成日志
func (g *Generator) Generate(startTime time.Time, duration int64) (*model.LogFile, error) {
	endTime := startTime.Add(time.Duration(duration) * time.Second)
	logFile := &model.LogFile{}

	// 设置时间范围
	logFile.TimeRange.StartTime = startTime
	logFile.TimeRange.Duration = duration

	// 获取安全事件
	securityEvents, err := g.getSecurityEvents(startTime, endTime)
	if err != nil {
		g.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogGenerate,
			Message: "生成日志文件失败",
			Error:   fmt.Errorf("获取安全事件失败: %w", err),
			Module:  "LogGenerator",
		})
		return nil, err
	}
	logFile.SecurityEvents.Events = securityEvents

	// 获取故障事件
	faultEvents, err := g.getFaultEvents(startTime, endTime)
	if err != nil {
		g.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogGenerate,
			Message: "生成日志文件失败",
			Error:   fmt.Errorf("获取故障事件失败: %w", err),
			Module:  "LogGenerator",
		})
		return nil, err
	}
	logFile.FaultEvents.Events = faultEvents

	// 获取性能事件
	securityDevices, err := g.getSecurityDevices(startTime, endTime)
	if err != nil {
		g.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogGenerate,
			Message: "生成日志文件失败",
			Error:   fmt.Errorf("获取性能事件失败: %w", err),
			Module:  "LogGenerator",
		})
		return nil, err
	}
	logFile.PerformanceEvents.SecurityDevices = securityDevices

	return logFile, nil
}

// GenerateToFile 生成日志并写入文件
func (g *Generator) GenerateToFile(startTime time.Time, duration int64, filePath string) error {
	// 生成日志
	logFile, err := g.Generate(startTime, duration)
	if err != nil {
		return err
	}

	// 转换为JSON
	data, err := json.MarshalIndent(logFile, "", "  ")
	if err != nil {
		g.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogGenerate,
			Message: "生成日志文件失败",
			Error:   fmt.Errorf("序列化日志数据失败: %w", err),
			Module:  "LogGenerator",
		})
		return err
	}

	// 写入文件
	if err := fileutil.WriteFile(filePath, data, 0644); err != nil {
		g.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogGenerate,
			Message: "生成日志文件失败",
			Error:   fmt.Errorf("写入日志文件失败: %w", err),
			Module:  "LogGenerator",
		})
		return err
	}

	return nil
}

// getSecurityEvents 获取安全事件
func (g *Generator) getSecurityEvents(startTime, endTime time.Time) ([]model.Event, error) {
	var events []model.Event
	err := g.retryOperation(func() error {
		return g.db.Where("eventType = ? AND eventTime BETWEEN ? AND ?",
			model.EventTypeSecurity, startTime, endTime).Find(&events).Error
	})
	return events, err
}

// getFaultEvents 获取故障事件
func (g *Generator) getFaultEvents(startTime, endTime time.Time) ([]model.Event, error) {
	var events []model.Event
	err := g.retryOperation(func() error {
		return g.db.Where("eventType = ? AND eventTime BETWEEN ? AND ?",
			model.EventTypeFault, startTime, endTime).Find(&events).Error
	})
	return events, err
}

// getSecurityDevices 获取安全接入管理设备
func (g *Generator) getSecurityDevices(startTime, endTime time.Time) ([]model.SecurityDevice, error) {
	var devices []struct {
		DeviceID        string `gorm:"column:deviceID"`
		DeviceStatus    int    `gorm:"column:deviceStatus"`
		PeakCPUUsage    int    `gorm:"column:peakCPUUsage"`
		PeakMemoryUsage int    `gorm:"column:peakMemoryUsage"`
		OnlineDuration  int    `gorm:"column:onlineDuration"`
	}
	err := g.retryOperation(func() error {
		return g.db.Table("devices").
			Select("deviceID, deviceStatus, peakCPUUsage, peakMemoryUsage, onlineDuration").
			Where("deviceType = ?", DeviceTypeSecurityMgmt).
			Find(&devices).Error
	})
	if err != nil {
		return nil, err
	}

	// 转换为SecurityDevice结构
	securityDevices := make([]model.SecurityDevice, len(devices))
	for i, device := range devices {
		securityDevices[i] = model.SecurityDevice{
			DeviceID:       device.DeviceID,
			CPUUsage:       device.PeakCPUUsage,
			MemoryUsage:    device.PeakMemoryUsage,
			OnlineDuration: device.OnlineDuration,
			Status:         device.DeviceStatus,
		}

		// 获取关联的网关设备
		gatewayDevices, err := g.getGatewayDevices(device.DeviceID, startTime, endTime)
		if err != nil {
			return nil, err
		}
		securityDevices[i].GatewayDevices = gatewayDevices
	}

	return securityDevices, nil
}

// getGatewayDevices 获取网关设备
func (g *Generator) getGatewayDevices(securityDeviceID string, startTime, endTime time.Time) ([]model.GatewayDevice, error) {
	var devices []struct {
		DeviceID        string `gorm:"column:deviceID"`
		DeviceStatus    int    `gorm:"column:deviceStatus"`
		PeakCPUUsage    int    `gorm:"column:peakCPUUsage"`
		PeakMemoryUsage int    `gorm:"column:peakMemoryUsage"`
		OnlineDuration  int    `gorm:"column:onlineDuration"`
	}
	err := g.retryOperation(func() error {
		return g.db.Table("devices").
			Select("deviceID, deviceStatus, peakCPUUsage, peakMemoryUsage, onlineDuration").
			Where("superiorDeviceID = ?", securityDeviceID).
			Find(&devices).Error
	})
	if err != nil {
		return nil, err
	}

	// 转换为GatewayDevice结构
	gatewayDevices := make([]model.GatewayDevice, len(devices))
	for i, device := range devices {
		gatewayDevices[i] = model.GatewayDevice{
			DeviceID:       device.DeviceID,
			CPUUsage:       device.PeakCPUUsage,
			MemoryUsage:    device.PeakMemoryUsage,
			OnlineDuration: device.OnlineDuration,
			Status:         device.DeviceStatus,
		}

		// 获取关联的用户
		users, err := g.getUsers(device.DeviceID, startTime, endTime)
		if err != nil {
			return nil, err
		}
		gatewayDevices[i].Users = users
	}

	return gatewayDevices, nil
}

// getUsers 获取用户
func (g *Generator) getUsers(gatewayDeviceID string, startTime, endTime time.Time) ([]model.User, error) {
	var users []struct {
		UserID         int `gorm:"column:userID"`
		Status         int `gorm:"column:status"`
		OnlineDuration int `gorm:"column:onlineDuration"`
	}
	err := g.retryOperation(func() error {
		return g.db.Table("users").
			Select("userID, status, onlineDuration").
			Where("gatewayDeviceID = ?", gatewayDeviceID).
			Find(&users).Error
	})
	if err != nil {
		return nil, err
	}

	// 转换为User结构
	result := make([]model.User, len(users))
	for i, user := range users {
		result[i] = model.User{
			UserID:         user.UserID,
			Status:         user.Status,
			OnlineDuration: user.OnlineDuration,
		}

		// 获取用户行为
		behaviors, err := g.getUserBehaviors(user.UserID, startTime, endTime)
		if err != nil {
			return nil, err
		}
		result[i].Behaviors = behaviors
	}

	return result, nil
}

// getUserBehaviors 获取用户行为
func (g *Generator) getUserBehaviors(userID int, startTime, endTime time.Time) ([]model.Behavior, error) {
	var behaviors []struct {
		BehaviorTime time.Time `gorm:"column:behaviorTime"`
		BehaviorType int       `gorm:"column:behaviorType"`
		DataType     int       `gorm:"column:dataType"`
		DataSize     int64     `gorm:"column:dataSize"`
	}
	err := g.retryOperation(func() error {
		return g.db.Table("user_behaviors").
			Select("behaviorTime, behaviorType, dataType, dataSize").
			Where("userID = ? AND behaviorTime BETWEEN ? AND ?",
				userID, startTime, endTime).
			Find(&behaviors).Error
	})
	if err != nil {
		return nil, err
	}

	// 转换为Behavior结构
	result := make([]model.Behavior, len(behaviors))
	for i, behavior := range behaviors {
		result[i] = model.Behavior{
			Time:     behavior.BehaviorTime,
			Type:     behavior.BehaviorType,
			DataType: behavior.DataType,
			DataSize: behavior.DataSize,
		}
	}

	return result, nil
}

// retryOperation 重试操作
func (g *Generator) retryOperation(operation func() error) error {
	var err error
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return err
}
