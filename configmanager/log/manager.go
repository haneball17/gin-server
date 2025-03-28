package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gin-server/config"
	"gin-server/configmanager/common/alert"
	"gin-server/configmanager/log/service"
	"gin-server/database/models"
	"gin-server/database/repositories"

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
	config           *config.Config
	db               *gorm.DB
	generator        *service.Generator
	encryptor        service.LogEncryptor
	uploader         *service.UploadManager
	alerter          alert.Alerter
	logService       LogService
	stopChan         chan struct{}
	isRunning        bool
	lastGenerateTime time.Time // 上次生成日志的时间，用于下次生成的起始时间
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

	// 创建仓库工厂
	repoFactory := repositories.NewRepositoryFactory(db)

	// 创建日志服务
	logService := NewLogService(repoFactory)

	// 创建日志管理器实例
	manager := &LogManager{
		config:     cfg,
		db:         db,
		generator:  service.NewGenerator(db, alerter),
		encryptor:  service.NewLogEncryptor(cfg, alerter),
		uploader:   uploader,
		alerter:    alerter,
		logService: logService,
		stopChan:   make(chan struct{}),
		isRunning:  false,
	}

	// 初始化 lastGenerateTime
	// 尝试从数据库获取最新的日志文件记录
	latestLogFile, err := logService.GetLatestLogFile()
	if err == nil && latestLogFile != nil && !latestLogFile.EndTime.IsZero() {
		// 使用最新日志文件的结束时间作为上次生成时间
		manager.lastGenerateTime = latestLogFile.EndTime
		if cfg.DebugLevel == "true" {
			log.Printf("从数据库加载上次日志生成时间: %v，日志文件ID: %d，文件名: %s\n",
				manager.lastGenerateTime.Format(time.RFC3339),
				latestLogFile.ID,
				latestLogFile.FileName)
		}
	} else {
		// 如果没有找到日志记录或时间为零值
		intervalMinutes := cfg.ConfigManager.LogManager.GenerateInterval

		// 确保时间间隔至少为1分钟
		if intervalMinutes < 1 {
			intervalMinutes = 1
			if cfg.DebugLevel == "true" {
				log.Printf("配置的时间间隔小于1分钟，已调整为1分钟\n")
			}
		}

		manager.lastGenerateTime = time.Now().Add(-time.Duration(intervalMinutes) * time.Minute)

		if cfg.DebugLevel == "true" {
			if err != nil {
				log.Printf("获取最新日志记录失败: %v\n", err)
			} else if latestLogFile == nil {
				log.Printf("未找到历史日志记录\n")
			} else if latestLogFile.EndTime.IsZero() {
				log.Printf("最新日志记录的结束时间为零值\n")
			}

			log.Printf("使用默认的上次生成时间: %v（当前时间减去%d分钟）\n",
				manager.lastGenerateTime.Format(time.RFC3339),
				intervalMinutes)
		}
	}

	return manager, nil
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

// GenerateLog 生成日志
func (m *LogManager) GenerateLog() error {
	// 获取当前时间作为生成的截止时间
	now := time.Now()

	// 使用上次生成时间作为起始时间，确保时间连续性
	startTime := m.lastGenerateTime

	// 记录调试信息
	if m.config.DebugLevel == "true" {
		log.Printf("开始生成日志，起始时间: %v，当前时间: %v\n",
			startTime.Format(time.RFC3339),
			now.Format(time.RFC3339))
	}

	// 如果上次生成时间为零值（首次生成），使用配置的生成间隔计算起始时间
	if startTime.IsZero() {
		// 获取配置的生成间隔（分钟）
		intervalMinutes := m.config.ConfigManager.LogManager.GenerateInterval

		// 确保时间间隔至少为1分钟
		if intervalMinutes < 1 {
			intervalMinutes = 1
		}

		// 计算开始时间（当前时间减去间隔时间）
		startTime = now.Add(-time.Duration(intervalMinutes) * time.Minute)

		if m.config.DebugLevel == "true" {
			log.Printf("首次生成日志，使用配置的时间间隔，起始时间: %v（当前时间减去%d分钟）\n",
				startTime.Format(time.RFC3339),
				intervalMinutes)
		}
	} else if startTime.After(now) {
		// 处理上次生成时间异常的情况（未来时间）
		log.Printf("警告: 上次生成时间(%v)晚于当前时间(%v)，这是不正常的\n",
			startTime.Format(time.RFC3339),
			now.Format(time.RFC3339))

		// 使用当前时间作为起始时间，避免异常
		startTime = now
	}

	// 计算持续时间（秒）
	durationSeconds := int64(now.Sub(startTime).Seconds())

	// 即使时间差为0，也允许生成空日志以保持连续性
	if durationSeconds == 0 && m.config.DebugLevel == "true" {
		log.Printf("时间差为0秒，将生成空日志以保持连续性\n")
	}

	// 构建日志文件名（格式：YYYYMMDDHHMMSS.json）
	fileName := fmt.Sprintf("%s.json", startTime.Format("20060102150405"))

	// 确保日志目录存在
	logDir := m.config.ConfigManager.LogManager.LogDir
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 构建日志文件完整路径
	logPath := filepath.Join(logDir, fileName)

	// 生成日志内容并写入文件
	if err := m.generator.GenerateToFile(startTime, durationSeconds, logPath); err != nil {
		m.alerter.Alert(&alert.Alert{
			Level: alert.AlertLevelError,
			Type:  alert.AlertTypeLogGenerate,
			Message: fmt.Sprintf("生成日志文件失败，时间范围: %v - %v",
				startTime.Format(time.RFC3339),
				now.Format(time.RFC3339)),
			Error:  err,
			Module: "LogManager",
		})
		return fmt.Errorf("生成日志文件失败: %v", err)
	}

	m.alerter.Alert(&alert.Alert{
		Level: alert.AlertLevelInfo,
		Type:  alert.AlertTypeLogGenerate,
		Message: fmt.Sprintf("成功生成日志文件: %s，时间范围: %v - %v",
			fileName,
			startTime.Format(time.RFC3339),
			now.Format(time.RFC3339)),
		Module: "LogManager",
	})

	// 处理日志文件（加密）
	processedLogPath, keyPath, err := m.encryptor.ProcessLog(logPath)
	if err != nil {
		return fmt.Errorf("处理日志文件失败: %v", err)
	}

	// 更新生成的日志文件路径和密钥路径到配置
	m.config.ConfigManager.LogManager.ProcessedLogPath = processedLogPath
	m.config.ConfigManager.LogManager.ProcessedKeyPath = keyPath

	// 上传日志文件
	if err := m.uploadLog(processedLogPath); err != nil {
		return fmt.Errorf("上传日志文件失败: %v", err)
	}

	// 更新上次生成时间为当前时间，确保下次生成时的连续性
	m.lastGenerateTime = now

	if m.config.DebugLevel == "true" {
		log.Printf("日志生成完成，时间范围: %v - %v，更新上次生成时间为: %v\n",
			startTime.Format(time.RFC3339),
			now.Format(time.RFC3339),
			m.lastGenerateTime.Format(time.RFC3339))
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

// uploadLog 上传指定的日志文件
func (m *LogManager) uploadLog(logPath string) error {
	// 确保上传目录为/log
	uploadDir := "/log/"

	// 设置新的上传目录到配置
	m.config.ConfigManager.LogManager.UploadDir = uploadDir

	// 获取密钥文件路径
	keyPath := m.config.ConfigManager.LogManager.ProcessedKeyPath

	// 使用Upload方法上传（会自动压缩打包）
	err := m.uploader.Upload(logPath, keyPath)
	if err != nil {
		return fmt.Errorf("上传日志文件失败: %v", err)
	}

	// 获取文件信息并创建或更新日志记录
	fileInfo, err := os.Stat(logPath)
	if err != nil {
		return fmt.Errorf("获取日志文件信息失败: %v", err)
	}

	// 获取文件名
	fileName := filepath.Base(logPath)

	// 从文件名中提取时间信息 (格式：YYYYMMDDHHMMSS.json)
	startTimeStr := strings.TrimSuffix(fileName, ".json")
	logStartTime, parseErr := time.ParseInLocation("20060102150405", startTimeStr, time.Local)

	// 默认使用当前时间作为结束时间
	logEndTime := time.Now()

	// 如果解析成功，使用从文件名解析出的时间作为起始时间
	if parseErr != nil && m.config.DebugLevel == "true" {
		log.Printf("无法从文件名解析时间: %v, 使用当前配置的时间范围", parseErr)
		// 使用当前记录的时间范围，不做其他处理
	}

	if m.config.DebugLevel == "true" {
		log.Printf("创建日志文件记录，文件名: %s, 时间范围: %v - %v\n",
			fileName,
			logStartTime.Format(time.RFC3339),
			logEndTime.Format(time.RFC3339))
	}

	// 创建日志文件记录，使用准确的时间范围
	logFile, err := m.logService.CreateLogFileWithTimeRange(fileName, fileInfo.Size(), logPath, logStartTime, logEndTime)
	if err != nil {
		return fmt.Errorf("创建日志文件记录失败: %v", err)
	}

	// 构建远程路径并标记为已上传
	remotePath := uploadDir + strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".tar.gz"
	err = m.logService.MarkLogFileAsUploaded(logFile.ID, remotePath)
	if err != nil {
		return fmt.Errorf("更新日志文件上传状态失败: %v", err)
	}

	if m.config.DebugLevel == "true" {
		log.Printf("日志文件记录已创建并标记为已上传，ID: %d, 远程路径: %s\n",
			logFile.ID, remotePath)
	}

	return nil
}

// GetLatestLogContent 获取最新的日志文件内容
func (m *LogManager) GetLatestLogContent() ([]byte, error) {
	// 获取最新的日志文件路径
	latestLogPath, err := m.findLatestLogFile()
	if err != nil {
		return nil, fmt.Errorf("查找最新日志文件失败: %w", err)
	}

	// 读取日志文件内容
	content, err := os.ReadFile(latestLogPath)
	if err != nil {
		return nil, fmt.Errorf("读取日志文件失败: %w", err)
	}

	return content, nil
}

// findLatestLogFile 查找最新的日志文件
func (m *LogManager) findLatestLogFile() (string, error) {
	logDir := "logs"

	// 读取logs目录
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return "", fmt.Errorf("读取日志目录失败: %w", err)
	}

	if len(entries) == 0 {
		return "", fmt.Errorf("日志目录为空，没有日志文件")
	}

	// 按照时间戳排序目录条目（目录名为时间戳格式）
	var logDirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			logDirs = append(logDirs, entry.Name())
		}
	}

	if len(logDirs) == 0 {
		return "", fmt.Errorf("日志目录中没有日志子目录")
	}

	// 按时间戳排序（降序）
	sort.Slice(logDirs, func(i, j int) bool {
		return logDirs[i] > logDirs[j] // 降序排列
	})

	// 获取最新的目录
	latestDir := filepath.Join(logDir, logDirs[0])

	// 检查log.json文件是否存在
	logFilePath := filepath.Join(latestDir, "log.json")
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		// 如果没有找到log.json文件，检查是否有加密文件夹
		encryptedDir := filepath.Join(latestDir, "encrypted")
		encryptedLogPath := filepath.Join(encryptedDir, "log.json")

		if _, err := os.Stat(encryptedLogPath); os.IsNotExist(err) {
			return "", fmt.Errorf("最新的日志目录中找不到日志文件")
		}

		// 返回加密的日志文件路径
		return encryptedLogPath, nil
	}

	return logFilePath, nil
}

// GetEventsByTimeRange 根据时间范围获取事件记录
func (m *LogManager) GetEventsByTimeRange(startTime, endTime time.Time) ([]models.Event, int64, error) {
	return m.logService.GetEventsByTimeRange(startTime, endTime)
}

// CreateEvent 创建事件记录
func (m *LogManager) CreateEvent(eventCode string, eventDesc string, deviceID int, eventType models.EventType) (*models.Event, error) {
	return m.logService.CreateEvent(eventCode, eventDesc, deviceID, eventType)
}

// LogUserBehavior 记录用户行为
func (m *LogManager) LogUserBehavior(userID int, behaviorType int, dataType int, dataSize int64) (*models.UserBehavior, error) {
	return m.logService.LogUserBehavior(userID, behaviorType, dataType, dataSize)
}

// GetUserBehaviorsByUserID 获取用户行为记录
func (m *LogManager) GetUserBehaviorsByUserID(userID int) ([]models.UserBehavior, int64, error) {
	return m.logService.GetUserBehaviorsByUserID(userID)
}

// GetLogFilesByTimeRange 根据时间范围获取日志文件
func (m *LogManager) GetLogFilesByTimeRange(startTime, endTime time.Time) ([]models.LogFile, int64, error) {
	return m.logService.GetLogFilesByTimeRange(startTime, endTime)
}
