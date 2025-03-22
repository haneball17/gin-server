package repositories

import (
	"gin-server/database/models"
	"time"

	"gorm.io/gorm"
)

// RadiusAuthRepository Radius认证仓库接口
type RadiusAuthRepository interface {
	Repository
	// FindByID 根据ID查找认证记录
	FindByID(id int) (*models.RadPostAuth, error)
	// FindByConditions 根据条件查询认证记录
	FindByConditions(query models.RadPostAuthQuery) ([]models.RadPostAuth, int64, error)
	// Create 创建认证记录
	Create(auth *models.RadPostAuth) error
	// CreateTable 确保表存在
	CreateTable() error
}

// radiusAuthRepository Radius认证仓库实现
type radiusAuthRepository struct {
	*BaseRepository
}

// NewRadiusAuthRepository 创建Radius认证仓库实例
func NewRadiusAuthRepository(db *gorm.DB) RadiusAuthRepository {
	return &radiusAuthRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// WithTx 使用事务进行操作
func (r *radiusAuthRepository) WithTx(tx *gorm.DB) Repository {
	return &radiusAuthRepository{
		BaseRepository: r.BaseRepository.WithTx(tx),
	}
}

// FindByID 根据ID查找认证记录
func (r *radiusAuthRepository) FindByID(id int) (*models.RadPostAuth, error) {
	var auth models.RadPostAuth
	if err := r.GetDB().First(&auth, id).Error; err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindByConditions 根据条件查询认证记录
func (r *radiusAuthRepository) FindByConditions(query models.RadPostAuthQuery) ([]models.RadPostAuth, int64, error) {
	db := r.GetDB().Model(&models.RadPostAuth{})

	// 添加查询条件
	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.Reply != "" {
		db = db.Where("reply = ?", query.Reply)
	}
	if query.StartDate != "" {
		startTime, err := time.Parse("2006-01-02", query.StartDate)
		if err == nil {
			db = db.Where("authdate >= ?", startTime)
		}
	}
	if query.EndDate != "" {
		endTime, err := time.Parse("2006-01-02", query.EndDate)
		if err == nil {
			// 设置为当天的结束时间
			endTime = endTime.Add(24*time.Hour - time.Second)
			db = db.Where("authdate <= ?", endTime)
		}
	}
	if query.Class != "" {
		db = db.Where("class = ?", query.Class)
	}

	// 获取总记录数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	offset := (query.Page - 1) * query.PageSize

	var auths []models.RadPostAuth
	if err := db.Order("authdate DESC").Limit(query.PageSize).Offset(offset).Find(&auths).Error; err != nil {
		return nil, 0, err
	}

	return auths, total, nil
}

// Create 创建认证记录
func (r *radiusAuthRepository) Create(auth *models.RadPostAuth) error {
	return r.GetDB().Create(auth).Error
}

// CreateTable 确保表存在
func (r *radiusAuthRepository) CreateTable() error {
	// 检查表是否存在，不存在则创建
	return r.GetDB().AutoMigrate(&models.RadPostAuth{})
}
