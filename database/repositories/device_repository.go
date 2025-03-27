package repositories

import (
	"gin-server/database/models"

	"gorm.io/gorm"
)

// DeviceRepository 设备仓库接口
type DeviceRepository interface {
	Repository
	// FindByID 根据ID查找设备
	FindByID(id uint) (*models.Device, error)
	// FindByDeviceID 根据设备ID查找设备
	FindByDeviceID(deviceID int) (*models.Device, error)
	// FindByDeviceName 根据设备名称查找设备
	FindByDeviceName(deviceName string) (*models.Device, error)
	// FindAll 查找所有设备
	FindAll() ([]models.Device, error)
	// Create 创建设备
	Create(device *models.Device) error
	// Update 更新设备
	Update(device *models.Device) error
	// Delete 删除设备
	Delete(id uint) error
}

// deviceRepository 设备仓库实现
type deviceRepository struct {
	*BaseRepository
}

// NewDeviceRepository 创建设备仓库实例
func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// WithTx 使用事务进行操作
func (r *deviceRepository) WithTx(tx *gorm.DB) Repository {
	return &deviceRepository{
		BaseRepository: r.BaseRepository.WithTx(tx),
	}
}

// FindByID 根据ID查找设备
func (r *deviceRepository) FindByID(id uint) (*models.Device, error) {
	var device models.Device
	if err := r.GetDB().First(&device, id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

// FindByDeviceID 根据设备ID查找设备
func (r *deviceRepository) FindByDeviceID(deviceID int) (*models.Device, error) {
	var device models.Device
	if err := r.GetDB().Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

// FindByDeviceName 根据设备名称查找设备
func (r *deviceRepository) FindByDeviceName(deviceName string) (*models.Device, error) {
	var device models.Device
	if err := r.GetDB().Where("device_name = ?", deviceName).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

// FindAll 查找所有设备
func (r *deviceRepository) FindAll() ([]models.Device, error) {
	var devices []models.Device
	if err := r.GetDB().Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

// Create 创建设备
func (r *deviceRepository) Create(device *models.Device) error {
	return r.GetDB().Create(device).Error
}

// Update 更新设备
func (r *deviceRepository) Update(device *models.Device) error {
	return r.GetDB().Save(device).Error
}

// Delete 删除设备
func (r *deviceRepository) Delete(id uint) error {
	return r.GetDB().Delete(&models.Device{}, id).Error
}
