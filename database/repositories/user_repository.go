package repositories

import (
	"gin-server/database/models"

	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	Repository
	// FindByID 根据ID查找用户
	FindByID(id uint) (*models.User, error)
	// FindByUserID 根据用户唯一标识查找用户
	FindByUserID(userID int) (*models.User, error)
	// FindByUsername 根据用户名查找用户
	FindByUsername(username string) (*models.User, error)
	// FindByEmail 根据邮箱查找用户
	FindByEmail(email string) (*models.User, error)
	// FindAll 查找所有用户
	FindAll() ([]models.User, error)
	// Create 创建用户
	Create(user *models.User) error
	// Update 更新用户
	Update(user *models.User) error
	// Delete 删除用户
	Delete(id uint) error
	// UpdateLastLogin 更新最后登录信息
	UpdateLastLogin(id uint, ip string) error
}

// userRepository 用户仓库实现
type userRepository struct {
	*BaseRepository
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// WithTx 使用事务进行操作
func (r *userRepository) WithTx(tx *gorm.DB) Repository {
	return &userRepository{
		BaseRepository: r.BaseRepository.WithTx(tx),
	}
}

// FindByID 根据ID查找用户
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.GetDB().First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUserID 根据用户唯一标识查找用户
func (r *userRepository) FindByUserID(userID int) (*models.User, error) {
	var user models.User
	if err := r.GetDB().Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.GetDB().Where("user_name = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.GetDB().Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll 查找所有用户
func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	if err := r.GetDB().Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Create 创建用户
func (r *userRepository) Create(user *models.User) error {
	return r.GetDB().Create(user).Error
}

// Update 更新用户
func (r *userRepository) Update(user *models.User) error {
	return r.GetDB().Save(user).Error
}

// Delete 删除用户
func (r *userRepository) Delete(id uint) error {
	return r.GetDB().Delete(&models.User{}, id).Error
}

// UpdateLastLogin 更新最后登录信息
func (r *userRepository) UpdateLastLogin(id uint, ip string) error {
	return r.GetDB().Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_time_stamp": gorm.Expr("NOW()"),
		"login_ip":              ip,
	}).Error
}
