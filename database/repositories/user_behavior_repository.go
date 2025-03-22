package repositories

import (
	"gin-server/database/models"
	"time"

	"gorm.io/gorm"
)

// UserBehaviorRepository 用户行为仓库接口
type UserBehaviorRepository interface {
	Repository
	// FindByID 根据ID查找用户行为
	FindByID(id uint) (*models.UserBehavior, error)
	// FindByUserID 查找指定用户的行为
	FindByUserID(userID int) ([]models.UserBehavior, int64, error)
	// FindByTimeRange 查找指定时间范围内的用户行为
	FindByTimeRange(startTime, endTime time.Time) ([]models.UserBehavior, int64, error)
	// FindByUserIDAndTimeRange 查找指定用户和时间范围内的行为
	FindByUserIDAndTimeRange(userID int, startTime, endTime time.Time) ([]models.UserBehavior, int64, error)
	// Create 创建用户行为记录
	Create(behavior *models.UserBehavior) error
	// Update 更新用户行为记录
	Update(behavior *models.UserBehavior) error
	// Delete 删除用户行为记录
	Delete(id uint) error
}

// userBehaviorRepository 用户行为仓库实现
type userBehaviorRepository struct {
	*BaseRepository
}

// NewUserBehaviorRepository 创建用户行为仓库实例
func NewUserBehaviorRepository(db *gorm.DB) UserBehaviorRepository {
	return &userBehaviorRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// WithTx 使用事务进行操作
func (r *userBehaviorRepository) WithTx(tx *gorm.DB) Repository {
	return &userBehaviorRepository{
		BaseRepository: r.BaseRepository.WithTx(tx),
	}
}

// FindByID 根据ID查找用户行为
func (r *userBehaviorRepository) FindByID(id uint) (*models.UserBehavior, error) {
	var behavior models.UserBehavior
	if err := r.GetDB().First(&behavior, id).Error; err != nil {
		return nil, err
	}
	return &behavior, nil
}

// FindByUserID 查找指定用户的行为
func (r *userBehaviorRepository) FindByUserID(userID int) ([]models.UserBehavior, int64, error) {
	var behaviors []models.UserBehavior
	var count int64

	// 查询总数
	if err := r.GetDB().Model(&models.UserBehavior{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := r.GetDB().Where("user_id = ?", userID).Find(&behaviors).Error; err != nil {
		return nil, 0, err
	}

	return behaviors, count, nil
}

// FindByTimeRange 查找指定时间范围内的用户行为
func (r *userBehaviorRepository) FindByTimeRange(startTime, endTime time.Time) ([]models.UserBehavior, int64, error) {
	var behaviors []models.UserBehavior
	var count int64

	// 查询总数
	if err := r.GetDB().Model(&models.UserBehavior{}).Where("behavior_time BETWEEN ? AND ?", startTime, endTime).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := r.GetDB().Where("behavior_time BETWEEN ? AND ?", startTime, endTime).Find(&behaviors).Error; err != nil {
		return nil, 0, err
	}

	return behaviors, count, nil
}

// FindByUserIDAndTimeRange 查找指定用户和时间范围内的行为
func (r *userBehaviorRepository) FindByUserIDAndTimeRange(userID int, startTime, endTime time.Time) ([]models.UserBehavior, int64, error) {
	var behaviors []models.UserBehavior
	var count int64

	// 查询总数
	if err := r.GetDB().Model(&models.UserBehavior{}).
		Where("user_id = ? AND behavior_time BETWEEN ? AND ?", userID, startTime, endTime).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := r.GetDB().Where("user_id = ? AND behavior_time BETWEEN ? AND ?", userID, startTime, endTime).
		Find(&behaviors).Error; err != nil {
		return nil, 0, err
	}

	return behaviors, count, nil
}

// Create 创建用户行为记录
func (r *userBehaviorRepository) Create(behavior *models.UserBehavior) error {
	return r.GetDB().Create(behavior).Error
}

// Update 更新用户行为记录
func (r *userBehaviorRepository) Update(behavior *models.UserBehavior) error {
	return r.GetDB().Save(behavior).Error
}

// Delete 删除用户行为记录
func (r *userBehaviorRepository) Delete(id uint) error {
	return r.GetDB().Delete(&models.UserBehavior{}, id).Error
}
