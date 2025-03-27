package log

import (
	"time"

	"gin-server/database/models"
	"gin-server/database/repositories"
)

// LogService 日志服务接口
type LogService interface {
	// CreateLogFile 创建日志文件记录
	CreateLogFile(fileName string, fileSize int64, filePath string) (*models.LogFile, error)

	// MarkLogFileAsUploaded 标记日志文件为已上传
	MarkLogFileAsUploaded(id uint, remotePath string) error

	// GetLatestLogFile 获取最新的日志文件
	GetLatestLogFile() (*models.LogFile, error)

	// GetLogFilesByTimeRange 根据时间范围获取日志文件
	GetLogFilesByTimeRange(startTime, endTime time.Time) ([]models.LogFile, int64, error)

	// CreateEvent 创建事件记录
	CreateEvent(eventCode string, eventDesc string, deviceID int, eventType models.EventType) (*models.Event, error)

	// GetEventsByTimeRange 根据时间范围获取事件
	GetEventsByTimeRange(startTime, endTime time.Time) ([]models.Event, int64, error)

	// LogUserBehavior 记录用户行为
	LogUserBehavior(userID int, behaviorType int, dataType int, dataSize int64) (*models.UserBehavior, error)

	// GetUserBehaviorsByUserID 获取用户行为记录
	GetUserBehaviorsByUserID(userID int) ([]models.UserBehavior, int64, error)
}

// logService 日志服务实现
type logService struct {
	repoFactory repositories.RepositoryFactory
}

// NewLogService 创建日志服务实例
func NewLogService(repoFactory repositories.RepositoryFactory) LogService {
	return &logService{
		repoFactory: repoFactory,
	}
}

// CreateLogFile 创建日志文件记录
func (s *logService) CreateLogFile(fileName string, fileSize int64, filePath string) (*models.LogFile, error) {
	repo := s.repoFactory.GetLogFileRepository()

	// 创建日志文件记录
	logFile := &models.LogFile{
		FileName:  fileName,
		FilePath:  filePath,
		FileSize:  fileSize,
		StartTime: time.Now().Add(-24 * time.Hour), // 默认为24小时前的数据
		EndTime:   time.Now(),
	}

	err := repo.Create(logFile)
	if err != nil {
		return nil, err
	}

	return logFile, nil
}

// MarkLogFileAsUploaded 标记日志文件为已上传
func (s *logService) MarkLogFileAsUploaded(id uint, remotePath string) error {
	return s.repoFactory.GetLogFileRepository().MarkAsUploaded(id, remotePath)
}

// GetLatestLogFile 获取最新的日志文件
func (s *logService) GetLatestLogFile() (*models.LogFile, error) {
	return s.repoFactory.GetLogFileRepository().FindLatest()
}

// GetLogFilesByTimeRange 根据时间范围获取日志文件
func (s *logService) GetLogFilesByTimeRange(startTime, endTime time.Time) ([]models.LogFile, int64, error) {
	return s.repoFactory.GetLogFileRepository().FindByTimeRange(startTime, endTime)
}

// CreateEvent 创建事件记录
func (s *logService) CreateEvent(eventCode string, eventDesc string, deviceID int, eventType models.EventType) (*models.Event, error) {
	repo := s.repoFactory.GetEventRepository()

	// 创建事件记录
	event := &models.Event{
		EventID:   time.Now().UnixNano(), // 使用纳秒级时间戳作为事件ID
		DeviceID:  deviceID,
		EventTime: time.Now(),
		EventType: eventType,
		EventCode: eventCode,
		EventDesc: eventDesc,
	}

	err := repo.Create(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// GetEventsByTimeRange 根据时间范围获取事件
func (s *logService) GetEventsByTimeRange(startTime, endTime time.Time) ([]models.Event, int64, error) {
	return s.repoFactory.GetEventRepository().FindByTimeRange(startTime, endTime)
}

// LogUserBehavior 记录用户行为
func (s *logService) LogUserBehavior(userID int, behaviorType int, dataType int, dataSize int64) (*models.UserBehavior, error) {
	repo := s.repoFactory.GetUserBehaviorRepository()

	// 创建用户行为记录
	behavior := &models.UserBehavior{
		UserID:       userID,
		BehaviorTime: time.Now(),
		BehaviorType: behaviorType,
		DataType:     dataType,
		DataSize:     dataSize,
	}

	err := repo.Create(behavior)
	if err != nil {
		return nil, err
	}

	return behavior, nil
}

// GetUserBehaviorsByUserID 获取用户行为记录
func (s *logService) GetUserBehaviorsByUserID(userID int) ([]models.UserBehavior, int64, error) {
	return s.repoFactory.GetUserBehaviorRepository().FindByUserID(userID)
}
