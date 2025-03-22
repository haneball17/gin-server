package repositories

import (
	"gin-server/database/models"
	"time"

	"gorm.io/gorm"
)

// CertRepository 证书仓库接口
type CertRepository interface {
	Repository
	// FindByID 根据ID查找证书
	FindByID(id uint) (*models.Cert, error)
	// FindByEntity 根据实体类型和ID查找证书
	FindByEntity(entityType, entityID string) (*models.Cert, error)
	// FindAll 查找所有证书
	FindAll() ([]models.Cert, error)
	// Create 创建证书
	Create(cert *models.Cert) error
	// Update 更新证书
	Update(cert *models.Cert) error
	// Delete 删除证书
	Delete(id uint) error
	// UpdateCertPath 更新证书路径
	UpdateCertPath(entityType, entityID, certPath string) error
	// UpdateKeyPath 更新密钥路径
	UpdateKeyPath(entityType, entityID, keyPath string) error
}

// certRepository 证书仓库实现
type certRepository struct {
	*BaseRepository
}

// NewCertRepository 创建证书仓库实例
func NewCertRepository(db *gorm.DB) CertRepository {
	return &certRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// WithTx 使用事务进行操作
func (r *certRepository) WithTx(tx *gorm.DB) Repository {
	return &certRepository{
		BaseRepository: r.BaseRepository.WithTx(tx),
	}
}

// FindByID 根据ID查找证书
func (r *certRepository) FindByID(id uint) (*models.Cert, error) {
	var cert models.Cert
	if err := r.GetDB().First(&cert, id).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

// FindByEntity 根据实体类型和ID查找证书
func (r *certRepository) FindByEntity(entityType, entityID string) (*models.Cert, error) {
	var cert models.Cert
	if err := r.GetDB().Where("entity_type = ? AND entity_id = ?", entityType, entityID).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

// FindAll 查找所有证书
func (r *certRepository) FindAll() ([]models.Cert, error) {
	var certs []models.Cert
	if err := r.GetDB().Find(&certs).Error; err != nil {
		return nil, err
	}
	return certs, nil
}

// Create 创建证书
func (r *certRepository) Create(cert *models.Cert) error {
	// 确保上传时间已设置
	if cert.UploadTime.IsZero() {
		cert.UploadTime = time.Now()
	}
	return r.GetDB().Create(cert).Error
}

// Update 更新证书
func (r *certRepository) Update(cert *models.Cert) error {
	return r.GetDB().Save(cert).Error
}

// Delete 删除证书
func (r *certRepository) Delete(id uint) error {
	return r.GetDB().Delete(&models.Cert{}, id).Error
}

// UpdateCertPath 更新证书路径
func (r *certRepository) UpdateCertPath(entityType, entityID, certPath string) error {
	// 查找证书记录
	var cert models.Cert
	result := r.GetDB().Where("entity_type = ? AND entity_id = ?", entityType, entityID).First(&cert)

	// 如果记录不存在，创建新记录
	if result.Error == gorm.ErrRecordNotFound {
		cert = models.Cert{
			EntityType: entityType,
			EntityID:   entityID,
			CertPath:   certPath,
			UploadTime: time.Now(),
		}
		return r.Create(&cert)
	} else if result.Error != nil {
		return result.Error
	}

	// 更新证书路径
	cert.CertPath = certPath
	cert.UploadTime = time.Now()
	return r.Update(&cert)
}

// UpdateKeyPath 更新密钥路径
func (r *certRepository) UpdateKeyPath(entityType, entityID, keyPath string) error {
	// 查找证书记录
	var cert models.Cert
	result := r.GetDB().Where("entity_type = ? AND entity_id = ?", entityType, entityID).First(&cert)

	// 如果记录不存在，创建新记录
	if result.Error == gorm.ErrRecordNotFound {
		cert = models.Cert{
			EntityType: entityType,
			EntityID:   entityID,
			KeyPath:    keyPath,
			UploadTime: time.Now(),
		}
		return r.Create(&cert)
	} else if result.Error != nil {
		return result.Error
	}

	// 更新密钥路径
	cert.KeyPath = keyPath
	cert.UploadTime = time.Now()
	return r.Update(&cert)
}
