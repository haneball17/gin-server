package repositories

import (
	"fmt"
	"gin-server/database"
	"gin-server/database/models"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

// LogFileRepository 日志文件仓库接口
type LogFileRepository interface {
	Repository
	// FindByID 根据ID查找日志文件
	FindByID(id uint) (*models.LogFile, error)
	// FindByFileName 根据文件名查找日志文件
	FindByFileName(fileName string) (*models.LogFile, error)
	// FindLatest 查找最新的日志文件
	FindLatest() (*models.LogFile, error)
	// FindByTimeRange 查找指定时间范围内的日志文件
	FindByTimeRange(startTime, endTime time.Time) ([]models.LogFile, int64, error)
	// Create 创建日志文件记录
	Create(logFile *models.LogFile) error
	// Update 更新日志文件记录
	Update(logFile *models.LogFile) error
	// Delete 删除日志文件记录
	Delete(id uint) error
	// MarkAsUploaded 标记日志文件为已上传
	MarkAsUploaded(id uint, remotePath string) error
}

// logFileRepository 日志文件仓库实现
type logFileRepository struct {
	*BaseRepository
}

// NewLogFileRepository 创建日志文件仓库实例
func NewLogFileRepository(db *gorm.DB) LogFileRepository {
	return &logFileRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// WithTx 使用事务进行操作
func (r *logFileRepository) WithTx(tx *gorm.DB) Repository {
	return &logFileRepository{
		BaseRepository: r.BaseRepository.WithTx(tx),
	}
}

// FindByID 根据ID查找日志文件
func (r *logFileRepository) FindByID(id uint) (*models.LogFile, error) {
	var logFile models.LogFile
	if err := r.GetDB().First(&logFile, id).Error; err != nil {
		return nil, err
	}
	return &logFile, nil
}

// FindByFileName 根据文件名查找日志文件
func (r *logFileRepository) FindByFileName(fileName string) (*models.LogFile, error) {
	var logFile models.LogFile
	if err := r.GetDB().Where("file_name = ?", fileName).First(&logFile).Error; err != nil {
		return nil, err
	}
	return &logFile, nil
}

// FindLatest 查找最新的日志文件
func (r *logFileRepository) FindLatest() (*models.LogFile, error) {
	var logFile models.LogFile
	err := r.GetDB().Order("created_at DESC").First(&logFile).Error
	if err != nil {
		// 如果是表不存在错误
		if strings.Contains(err.Error(), "Table 'gin_server.log_files' doesn't exist") {
			log.Printf("发现log_files表不存在，尝试初始化...")

			// 使用新的初始化机制，调用数据库模块的表初始化函数
			if initErr := database.InitLogFilesTable(); initErr != nil {
				log.Printf("初始化log_files表失败: %v", initErr)
				return nil, fmt.Errorf("表不存在且初始化失败: %w (初始化错误: %v)", err, initErr)
			}

			log.Printf("log_files表初始化成功，重试查询...")
			// 再次尝试查询
			retryErr := r.GetDB().Order("created_at DESC").First(&logFile).Error
			if retryErr != nil {
				// 如果是记录不存在的错误，这在业务逻辑上是可接受的
				if retryErr == gorm.ErrRecordNotFound {
					return nil, retryErr
				}
				return nil, fmt.Errorf("表初始化后查询失败: %w", retryErr)
			}
			return &logFile, nil
		}

		// 如果是记录不存在的错误，这是正常的业务逻辑情况
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}

		// 其他类型的错误
		return nil, fmt.Errorf("查询最新日志文件失败: %w", err)
	}
	return &logFile, nil
}

// FindByTimeRange 查找指定时间范围内的日志文件
func (r *logFileRepository) FindByTimeRange(startTime, endTime time.Time) ([]models.LogFile, int64, error) {
	var logFiles []models.LogFile
	var count int64

	// 查询总数
	query := r.GetDB().Model(&models.LogFile{}).Where("(start_time BETWEEN ? AND ?) OR (end_time BETWEEN ? AND ?)",
		startTime, endTime, startTime, endTime)
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := query.Find(&logFiles).Error; err != nil {
		return nil, 0, err
	}

	return logFiles, count, nil
}

// Create 创建日志文件记录
func (r *logFileRepository) Create(logFile *models.LogFile) error {
	return r.GetDB().Create(logFile).Error
}

// Update 更新日志文件记录
func (r *logFileRepository) Update(logFile *models.LogFile) error {
	return r.GetDB().Save(logFile).Error
}

// Delete 删除日志文件记录
func (r *logFileRepository) Delete(id uint) error {
	return r.GetDB().Delete(&models.LogFile{}, id).Error
}

// MarkAsUploaded 标记日志文件为已上传
func (r *logFileRepository) MarkAsUploaded(id uint, remotePath string) error {
	return r.GetDB().Model(&models.LogFile{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_uploaded":   true,
		"remote_path":   remotePath,
		"uploaded_time": time.Now(),
	}).Error
}
