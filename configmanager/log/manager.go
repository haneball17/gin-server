package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	encryptor service.LogEncryptor
	uploader  *service.UploadManager
	alerter   alert.Alerter
	stopChan  chan struct{}
	isRunning bool
}

// NewLogManager 创建日志管理器
func NewLogManager(cfg *config.Config, db *gorm.DB) (*LogManager, error) {
	// 创建告警器
	alerter := alert.GetDefaultAlerter()

	// 创建上传管理器
	uploader, err := service.NewUploadManager(cfg, alerter)
	if err != nil {
		return nil, fmt.Errorf("创建上传管理器失败: %v", err)
	}

	return &LogManager{
		config:    cfg,
		db:        db,
		generator: service.NewGenerator(db, alerter),
		encryptor: service.NewLogEncryptor(cfg, alerter),
		uploader:  uploader,
		alerter:   alerter,
		stopChan:  make(chan struct{}),
		isRunning: false,
	}, nil
}

// Start 实现Manager接口
func (m *LogManager) Start() error {
	if m.isRunning {
		return fmt.Errorf("日志管理器已经在运行")
	}

	m.isRunning = true

	// 启动时先生成一次日志
	if err := m.GenerateLog(); err != nil {
		m.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogGenerate,
			Message: "初始生成日志失败",
			Error:   err,
			Module:  "LogManager",
		})
	}

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

	// 关闭上传管理器
	if err := m.uploader.Close(); err != nil {
		m.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogUpload,
			Message: "关闭上传管理器失败",
			Error:   err,
			Module:  "LogManager",
		})
	}

	return nil
}

// GenerateLog 实现Manager接口
func (m *LogManager) GenerateLog() error {
	startTime := time.Now().Add(-5 * time.Minute)
	duration := int64(5 * 60) // 5分钟，单位：秒

	// 生成文件路径
	timestamp := time.Now().Format("20060102150405")
	filePath := strings.ReplaceAll(filepath.Join("logs", timestamp, "log.json"), "\\", "/")

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

	// 处理日志文件（可能包含加密）
	processedPath, keyPath, err := m.encryptor.ProcessLog(filePath)
	if err != nil {
		m.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogEncrypt,
			Message: "加密日志文件失败",
			Error:   err,
			Module:  "LogManager",
		})
		return fmt.Errorf("加密日志文件失败: %v", err)
	}

	// 规范化路径
	processedPath = strings.ReplaceAll(processedPath, "\\", "/")
	if keyPath != "" {
		keyPath = strings.ReplaceAll(keyPath, "\\", "/")
	}

	// 上传处理后的文件
	if err := m.uploader.Upload(processedPath, keyPath); err != nil {
		m.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogUpload,
			Message: "上传日志文件失败",
			Error:   err,
			Module:  "LogManager",
		})
		return fmt.Errorf("上传日志文件失败: %v", err)
	}

	return nil
}

// run 运行日志管理器
func (m *LogManager) run() {
	interval := time.Duration(m.config.ConfigManager.LogManager.GenerateInterval) * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	if m.config.DebugLevel == "true" {
		log.Printf("日志管理器开始运行，生成间隔: %v\n", interval)
	}

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
			} else if m.config.DebugLevel == "true" {
				log.Println("成功生成并上传日志")
			}
		case <-m.stopChan:
			return
		}
	}
}

// ListRemoteLogFiles 列出远程日志文件
func (m *LogManager) ListRemoteLogFiles() ([]service.File, error) {
	uploadDir := strings.ReplaceAll(m.config.ConfigManager.LogManager.UploadDir, "\\", "/")
	if !strings.HasSuffix(uploadDir, "/") {
		uploadDir += "/"
	}
	return m.uploader.ListFiles(uploadDir)
}

// DownloadLogFile 下载日志文件
func (m *LogManager) DownloadLogFile(file service.File) ([]byte, error) {
	uploadDir := strings.ReplaceAll(m.config.ConfigManager.LogManager.UploadDir, "\\", "/")
	if !strings.HasSuffix(uploadDir, "/") {
		uploadDir += "/"
	}
	file.Path = strings.ReplaceAll(file.Path, "\\", "/")
	return m.uploader.DownloadFile(uploadDir + file.Path)
}

func (m *LogManager) uploadLog(logPath string) error {
	uploadDir := strings.ReplaceAll(m.config.ConfigManager.LogManager.UploadDir, "\\", "/")
	if !strings.HasSuffix(uploadDir, "/") {
		uploadDir += "/"
	}

	// 读取日志文件
	data, err := os.ReadFile(logPath)
	if err != nil {
		return fmt.Errorf("读取日志文件失败: %v", err)
	}

	// 获取文件名
	fileName := filepath.Base(logPath)
	remotePath := uploadDir + fileName

	// 上传文件
	err = m.uploader.UploadFile(remotePath, data)
	if err != nil {
		return fmt.Errorf("上传日志文件失败: %v", err)
	}

	return nil
}
