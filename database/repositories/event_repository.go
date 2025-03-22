package repositories

import (
	"gin-server/database/models"
	"time"

	"gorm.io/gorm"
)

// EventRepository 事件仓库接口
type EventRepository interface {
	Repository
	// FindByID 根据ID查找事件
	FindByID(id int64) (*models.Event, error)
	// FindByTimeRange 查找指定时间范围内的事件
	FindByTimeRange(startTime, endTime time.Time) ([]models.Event, int64, error)
	// FindByTypeAndTimeRange 查找指定类型和时间范围内的事件
	FindByTypeAndTimeRange(eventType models.EventType, startTime, endTime time.Time) ([]models.Event, int64, error)
	// Create 创建事件
	Create(event *models.Event) error
	// Update 更新事件
	Update(event *models.Event) error
	// Delete 删除事件
	Delete(id int64) error
}

// eventRepository 事件仓库实现
type eventRepository struct {
	*BaseRepository
}

// NewEventRepository 创建事件仓库实例
func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// WithTx 使用事务进行操作
func (r *eventRepository) WithTx(tx *gorm.DB) Repository {
	return &eventRepository{
		BaseRepository: r.BaseRepository.WithTx(tx),
	}
}

// FindByID 根据ID查找事件
func (r *eventRepository) FindByID(id int64) (*models.Event, error) {
	var event models.Event
	if err := r.GetDB().Where("event_id = ?", id).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// FindByTimeRange 查找指定时间范围内的事件
func (r *eventRepository) FindByTimeRange(startTime, endTime time.Time) ([]models.Event, int64, error) {
	var events []models.Event
	var count int64

	// 查询总数
	if err := r.GetDB().Model(&models.Event{}).Where("event_time BETWEEN ? AND ?", startTime, endTime).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := r.GetDB().Where("event_time BETWEEN ? AND ?", startTime, endTime).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, count, nil
}

// FindByTypeAndTimeRange 查找指定类型和时间范围内的事件
func (r *eventRepository) FindByTypeAndTimeRange(eventType models.EventType, startTime, endTime time.Time) ([]models.Event, int64, error) {
	var events []models.Event
	var count int64

	// 查询总数
	if err := r.GetDB().Model(&models.Event{}).Where("event_type = ? AND event_time BETWEEN ? AND ?", eventType, startTime, endTime).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := r.GetDB().Where("event_type = ? AND event_time BETWEEN ? AND ?", eventType, startTime, endTime).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, count, nil
}

// Create 创建事件
func (r *eventRepository) Create(event *models.Event) error {
	return r.GetDB().Create(event).Error
}

// Update 更新事件
func (r *eventRepository) Update(event *models.Event) error {
	return r.GetDB().Save(event).Error
}

// Delete 删除事件
func (r *eventRepository) Delete(id int64) error {
	return r.GetDB().Where("event_id = ?", id).Delete(&models.Event{}).Error
}
