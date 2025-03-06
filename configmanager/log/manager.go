package log

import (
	"fmt"
	"path/filepath"
	"time"

	"gin-server/config"
	"gin-server/configmanager/common/alert"
	"gin-server/configmanager/log/service"

	"gorm.io/gorm"
)

// Manager 日志管理器接口
type Manager interface {
	// Start 启动日志管理器
	Start() error

	// Stop 停止日志管理器
	Stop() error

	// GenerateLog 生成日志
	GenerateLog() error
}

// LogManager 日志管理器实现
type LogManager struct {
	config    *config.Config
	db        *gorm.DB
	generator *service.Generator
	alerter   alert.Alerter
	stopChan  chan struct{}
	isRunning bool
}

// NewLogManager 创建日志管理器
func NewLogManager(cfg *config.Config, db *gorm.DB) *LogManager {
	return &LogManager{
		config:    cfg,
		db:        db,
		generator: service.NewGenerator(db, nil),
		alerter:   alert.GetDefaultAlerter(),
		stopChan:  make(chan struct{}),
		isRunning: false,
	}
}

// Start 实现Manager接口
func (m *LogManager) Start() error {
	if m.isRunning {
		return fmt.Errorf("日志管理器已经在运行")
	}

	m.isRunning = true
	go m.run()
	return nil
}

// Stop 实现Manager接口
func (m *LogManager) Stop() error {
	if !m.isRunning {
		return fmt.Errorf("日志管理器未在运行")
	}

	m.isRunning = false
	close(m.stopChan)
	return nil
}

// GenerateLog 实现Manager接口
func (m *LogManager) GenerateLog() error {
	startTime := time.Now().Add(-5 * time.Minute)
	duration := int64(5 * 60) // 5分钟，单位：秒

	// 生成文件路径
	timestamp := time.Now().Format("20060102150405")
	filePath := filepath.Join("logs", timestamp, "log.json")

	// 生成日志文件
	err := m.generator.GenerateToFile(startTime, duration, filePath)
	if err != nil {
		m.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogGenerate,
			Message: "生成日志文件失败",
			Error:   err,
			Module:  "LogManager",
		})
		return fmt.Errorf("生成日志文件失败: %v", err)
	}

	return nil
}

// run 运行日志管理器
func (m *LogManager) run() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := m.GenerateLog(); err != nil {
				m.alerter.Alert(&alert.Alert{
					Level:   alert.AlertLevelError,
					Type:    alert.AlertTypeLogGenerate,
					Message: "定时生成日志失败",
					Error:   err,
					Module:  "LogManager",
				})
			}
		case <-m.stopChan:
			return
		}
	}
}
