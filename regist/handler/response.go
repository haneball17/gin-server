package handler

import (
	"gin-server/database/models"
)

// UserResponse 用户查询响应结构体
type UserResponse struct {
	ID                 uint    `json:"id"`
	Username           string  `json:"user_name"`
	UserID             int     `json:"user_id"`
	UserType           int     `json:"user_type"`
	GatewayDeviceID    int     `json:"gateway_device_id"`
	Status             *int    `json:"status"`
	OnlineDuration     int     `json:"online_duration"`
	CertID             string  `json:"cert_id"`
	KeyID              string  `json:"key_id"`
	Email              string  `json:"email"`
	PermissionMask     string  `json:"permission_mask,omitempty"`
	LastLoginTimeStamp *string `json:"last_login_timestamp,omitempty"`
	OffLineTimeStamp   *string `json:"offline_timestamp,omitempty"`
	LoginIP            string  `json:"login_ip,omitempty"`
	IllegalLoginTimes  *int    `json:"illegal_login_times,omitempty"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at,omitempty"`
}

// convertUserModelToResponse 将用户模型转换为响应结构体
func convertUserModelToResponse(user *models.User) UserResponse {
	response := UserResponse{
		ID:              user.ID,
		Username:        user.Username,
		UserID:          user.UserID,
		UserType:        user.UserType,
		GatewayDeviceID: user.GatewayDeviceID,
		Status:          user.Status,
		OnlineDuration:  user.OnlineDuration,
		CertID:          user.CertID,
		KeyID:           user.KeyID,
		Email:           user.Email,
		CreatedAt:       user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// 可选字段，如果存在则添加
	if !user.UpdatedAt.IsZero() {
		updatedAt := user.UpdatedAt.Format("2006-01-02T15:04:05Z")
		response.UpdatedAt = updatedAt
	}

	if user.PermissionMask != "" {
		response.PermissionMask = user.PermissionMask
	}

	if user.LastLoginTimeStamp != nil {
		lastLogin := user.LastLoginTimeStamp.Format("2006-01-02T15:04:05Z")
		response.LastLoginTimeStamp = &lastLogin
	}

	if user.OffLineTimeStamp != nil {
		offlineTime := user.OffLineTimeStamp.Format("2006-01-02T15:04:05Z")
		response.OffLineTimeStamp = &offlineTime
	}

	if user.LoginIP != "" {
		response.LoginIP = user.LoginIP
	}

	if user.IllegalLoginTimes != nil {
		response.IllegalLoginTimes = user.IllegalLoginTimes
	}

	return response
}

// DeviceResponse 设备查询响应结构体
type DeviceResponse struct {
	ID                  uint   `json:"id"`
	DeviceName          string `json:"device_name"`
	DeviceType          int    `json:"device_type"`
	DeviceID            int    `json:"device_id"`
	SuperiorDeviceID    int    `json:"superior_device_id"`
	DeviceStatus        int    `json:"device_status"`
	PeakCPUUsage        int    `json:"peak_cpu_usage,omitempty"`
	PeakMemoryUsage     int    `json:"peak_memory_usage,omitempty"`
	OnlineDuration      int    `json:"online_duration"`
	CertID              string `json:"cert_id"`
	KeyID               string `json:"key_id"`
	RegisterIP          string `json:"register_ip"`
	Email               string `json:"email"`
	HardwareFingerprint string `json:"hardware_fingerprint,omitempty"`
	AnonymousUser       string `json:"anonymous_user,omitempty"`
	LongAddress         string `json:"long_address,omitempty"`
	ShortAddress        string `json:"short_address,omitempty"`
	SESKey              string `json:"ses_key,omitempty"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at,omitempty"`
}

// convertDeviceModelToResponse 将设备模型转换为响应结构体
func convertDeviceModelToResponse(device *models.Device) DeviceResponse {
	response := DeviceResponse{
		ID:                  device.ID,
		DeviceName:          device.DeviceName,
		DeviceType:          device.DeviceType,
		DeviceID:            device.DeviceID,
		SuperiorDeviceID:    device.SuperiorDeviceID,
		DeviceStatus:        device.DeviceStatus,
		PeakCPUUsage:        device.PeakCPUUsage,
		PeakMemoryUsage:     device.PeakMemoryUsage,
		OnlineDuration:      device.OnlineDuration,
		CertID:              device.CertID,
		KeyID:               device.KeyID,
		RegisterIP:          device.RegisterIP,
		Email:               device.Email,
		HardwareFingerprint: device.HardwareFingerprint,
		AnonymousUser:       device.AnonymousUser,
		LongAddress:         device.LongAddress,
		ShortAddress:        device.ShortAddress,
		SESKey:              device.SESKey,
		CreatedAt:           device.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// 可选字段，如果存在则添加
	if !device.UpdatedAt.IsZero() {
		updatedAt := device.UpdatedAt.Format("2006-01-02T15:04:05Z")
		response.UpdatedAt = updatedAt
	}

	return response
}
