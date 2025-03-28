package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gin-server/configmanager/common/alert"
	"gin-server/configmanager/common/fileutil"
	"gin-server/database/models"
	"gin-server/database/repositories"

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
	db                 *gorm.DB
	alerter            alert.Alerter
	eventRepository    repositories.EventRepository
	deviceRepository   repositories.DeviceRepository
	userRepository     repositories.UserRepository
	behaviorRepository repositories.UserBehaviorRepository
}

// NewGenerator 创建日志生成器实例
func NewGenerator(db *gorm.DB, alerter alert.Alerter) *Generator {
	repoFactory := repositories.NewRepositoryFactory(db)
	return &Generator{
		db:                 db,
		alerter:            alerter,
		eventRepository:    repoFactory.GetEventRepository(),
		deviceRepository:   repoFactory.GetDeviceRepository(),
		userRepository:     repoFactory.GetUserRepository(),
		behaviorRepository: repoFactory.GetUserBehaviorRepository(),
	}
}

// Generate 生成日志
// startTime: 日志的起始时间，通常为上次生成日志的时间
// duration: 日志覆盖的时间范围（单位：秒）
func (g *Generator) Generate(startTime time.Time, duration int64) (*models.LogContent, error) {
	// 计算结束时间
	endTime := startTime.Add(time.Duration(duration) * time.Second)

	// 创建日志内容结构
	logContent := &models.LogContent{}

	// 设置时间范围
	logContent.TimeRange.StartTime = startTime
	logContent.TimeRange.Duration = duration

	// 记录生成过程的日志
	log.Printf("开始收集日志数据，时间范围: %v - %v (持续时间: %d秒)\n",
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
		duration)

	// 在时间差为0的情况下，生成空日志，但保留时间戳信息
	if duration == 0 {
		log.Printf("时间范围为0，生成空日志\n")
		// 仅设置空集合，不执行查询
		logContent.SecurityEvents.Events = []models.Event{}
		logContent.FaultEvents.Events = []models.Event{}
		logContent.PerformanceEvents.SecurityDevices = []models.SecurityDevice{}
		return logContent, nil
	}

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
	logContent.SecurityEvents.Events = securityEvents

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
	logContent.FaultEvents.Events = faultEvents

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
	logContent.PerformanceEvents.SecurityDevices = securityDevices

	// 记录日志数据量统计
	log.Printf("收集完成 - 安全事件: %d, 故障事件: %d, 安全设备: %d\n",
		len(securityEvents),
		len(faultEvents),
		len(securityDevices))

	return logContent, nil
}

// GenerateToFile 生成日志并写入文件
// startTime: 日志的起始时间，通常为上次生成日志的时间
// duration: 日志覆盖的时间范围（单位：秒）
// filePath: 保存日志文件的路径
func (g *Generator) GenerateToFile(startTime time.Time, duration int64, filePath string) error {
	// 生成日志
	logContent, err := g.Generate(startTime, duration)
	if err != nil {
		return err
	}

	// 转换为标准日志格式
	logContentLog := logContent.ToLogContentLog()

	// 转换为JSON
	data, err := json.MarshalIndent(logContentLog, "", "  ")
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
func (g *Generator) getSecurityEvents(startTime, endTime time.Time) ([]models.Event, error) {
	events, count, err := g.eventRepository.FindByTypeAndTimeRange(models.EventTypeSecurity, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("查询安全事件失败: %w", err)
	}

	// 日志记录查询结果
	if count == 0 {
		log.Printf("在时间范围 %v - %v 内未找到安全事件\n",
			startTime.Format(time.RFC3339),
			endTime.Format(time.RFC3339))
	} else {
		log.Printf("在时间范围 %v - %v 内找到 %d 条安全事件\n",
			startTime.Format(time.RFC3339),
			endTime.Format(time.RFC3339),
			count)
	}

	return events, nil
}

// getFaultEvents 获取故障事件
func (g *Generator) getFaultEvents(startTime, endTime time.Time) ([]models.Event, error) {
	events, count, err := g.eventRepository.FindByTypeAndTimeRange(models.EventTypeFault, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("查询故障事件失败: %w", err)
	}

	// 日志记录查询结果
	if count == 0 {
		log.Printf("在时间范围 %v - %v 内未找到故障事件\n",
			startTime.Format(time.RFC3339),
			endTime.Format(time.RFC3339))
	} else {
		log.Printf("在时间范围 %v - %v 内找到 %d 条故障事件\n",
			startTime.Format(time.RFC3339),
			endTime.Format(time.RFC3339),
			count)
	}

	return events, nil
}

// getSecurityDevices 获取安全接入管理设备
func (g *Generator) getSecurityDevices(startTime, endTime time.Time) ([]models.SecurityDevice, error) {
	var securityDevices []models.SecurityDevice

	// 查询所有安全接入管理设备
	var devices []models.Device
	err := g.db.Where("device_type = ?", DeviceTypeSecurityMgmt).Find(&devices).Error
	if err != nil {
		return nil, fmt.Errorf("查询安全接入管理设备失败: %w", err)
	}

	// 处理每个安全接入管理设备
	for _, device := range devices {
		// 获取网关设备
		gatewayDevices, err := g.getGatewayDevices(device.DeviceID, startTime, endTime)
		if err != nil {
			return nil, err
		}

		// 创建安全设备对象
		securityDevice := models.SecurityDevice{
			DeviceID:       device.DeviceID,
			CPUUsage:       device.PeakCPUUsage,
			MemoryUsage:    device.PeakMemoryUsage,
			OnlineDuration: device.OnlineDuration,
			Status:         device.DeviceStatus,
			GatewayDevices: gatewayDevices,
		}

		securityDevices = append(securityDevices, securityDevice)
	}

	return securityDevices, nil
}

// getGatewayDevices 获取网关设备
func (g *Generator) getGatewayDevices(securityDeviceID int, startTime, endTime time.Time) ([]models.GatewayDevice, error) {
	var gatewayDevices []models.GatewayDevice

	// 查询所有隶属于指定安全接入管理设备的网关设备
	var devices []models.Device
	err := g.db.Where("superior_device_id = ? AND (device_type = ? OR device_type = ? OR device_type = ?)",
		securityDeviceID, DeviceTypeGatewayA, DeviceTypeGatewayB, DeviceTypeGatewayC).Find(&devices).Error
	if err != nil {
		return nil, fmt.Errorf("查询网关设备失败: %w", err)
	}

	// 处理每个网关设备
	for _, device := range devices {
		// 获取用户
		users, err := g.getUsers(device.DeviceID, startTime, endTime)
		if err != nil {
			return nil, err
		}

		// 创建网关设备对象
		gatewayDevice := models.GatewayDevice{
			DeviceID:       device.DeviceID,
			CPUUsage:       device.PeakCPUUsage,
			MemoryUsage:    device.PeakMemoryUsage,
			OnlineDuration: device.OnlineDuration,
			Status:         device.DeviceStatus,
			Users:          users,
		}

		gatewayDevices = append(gatewayDevices, gatewayDevice)
	}

	return gatewayDevices, nil
}

// getUsers 获取用户
func (g *Generator) getUsers(gatewayDeviceID int, startTime, endTime time.Time) ([]models.UserInfo, error) {
	var userInfos []models.UserInfo

	// 查询所有隶属于指定网关设备的用户
	var users []models.User
	err := g.db.Where("gateway_device_id = ?", gatewayDeviceID).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 处理每个用户
	for _, user := range users {
		// 获取用户行为
		behaviors, err := g.getUserBehaviors(user.UserID, startTime, endTime)
		if err != nil {
			return nil, err
		}

		// 设置状态，确保指针值安全
		status := 2 // 默认离线
		if user.Status != nil {
			status = *user.Status
		}

		// 创建用户信息对象
		userInfo := models.UserInfo{
			UserID:         user.UserID,
			Status:         status,
			OnlineDuration: user.OnlineDuration,
			Behaviors:      behaviors,
		}

		userInfos = append(userInfos, userInfo)
	}

	return userInfos, nil
}

// getUserBehaviors 获取用户行为
func (g *Generator) getUserBehaviors(userID int, startTime, endTime time.Time) ([]models.UserBehavior, error) {
	behaviors, count, err := g.behaviorRepository.FindByUserIDAndTimeRange(userID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("查询用户行为失败: %w", err)
	}

	// 仅在调试模式下记录详细信息，避免日志过多
	if count > 0 {
		log.Printf("用户(ID: %d) 在时间范围 %v - %v 内找到 %d 条行为记录\n",
			userID,
			startTime.Format(time.RFC3339),
			endTime.Format(time.RFC3339),
			count)
	}

	return behaviors, nil
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
