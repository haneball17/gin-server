package service

import (
	"fmt"
	"gin-server/config"
	"gin-server/configmanager/common/alert"
	"gin-server/configmanager/common/crypto"
	"gin-server/configmanager/common/fileutil"
	"path/filepath"
)

// LogEncryptor 日志加密器接口
type LogEncryptor interface {
	// ProcessLog 处理日志文件
	// 如果启用了加密，则加密日志文件并返回加密后的文件路径和密钥文件路径
	// 如果未启用加密，则直接返回原始文件路径
	ProcessLog(logPath string) (resultPath string, keyPath string, err error)
}

// DefaultLogEncryptor 默认日志加密器实现
type DefaultLogEncryptor struct {
	config  *config.Config
	alerter alert.Alerter
}

// NewLogEncryptor 创建日志加密器
func NewLogEncryptor(cfg *config.Config, alerter alert.Alerter) *DefaultLogEncryptor {
	return &DefaultLogEncryptor{
		config:  cfg,
		alerter: alerter,
	}
}

// ProcessLog 处理日志文件
func (e *DefaultLogEncryptor) ProcessLog(logPath string) (string, string, error) {
	// 如果未启用加密，直接返回原始文件路径
	if !e.config.ConfigManager.LogManager.EnableEncryption {
		return logPath, "", nil
	}

	// 创建加密文件的目标目录
	encryptedDir := filepath.Join(filepath.Dir(logPath), "encrypted")
	if err := fileutil.EnsureDir(encryptedDir); err != nil {
		e.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogEncrypt,
			Message: "创建加密文件目录失败",
			Error:   err,
			Module:  "LogEncryptor",
		})
		return "", "", fmt.Errorf("创建加密文件目录失败: %w", err)
	}

	// 加密文件路径
	encryptedPath := filepath.Join(encryptedDir, filepath.Base(logPath))
	keyPath := filepath.Join(encryptedDir, "key.txt")

	// 读取原始日志文件
	logData, err := fileutil.ReadFile(logPath)
	if err != nil {
		e.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogEncrypt,
			Message: "读取日志文件失败",
			Error:   err,
			Module:  "LogEncryptor",
		})
		return "", "", fmt.Errorf("读取日志文件失败: %w", err)
	}

	// 创建AES加密器
	aes, err := crypto.NewAESEncryptor(e.config.ConfigManager.LogManager.Encryption.AESKeyLength)
	if err != nil {
		e.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogEncrypt,
			Message: "创建AES加密器失败",
			Error:   err,
			Module:  "LogEncryptor",
		})
		return "", "", fmt.Errorf("创建AES加密器失败: %w", err)
	}

	// 加密日志数据
	encryptedData, err := aes.Encrypt(logData)
	if err != nil {
		e.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogEncrypt,
			Message: "AES加密失败",
			Error:   err,
			Module:  "LogEncryptor",
		})
		return "", "", fmt.Errorf("AES加密失败: %w", err)
	}

	// 保存加密后的日志文件
	if err := fileutil.WriteFile(encryptedPath, encryptedData, 0644); err != nil {
		e.alerter.Alert(&alert.Alert{
			Level:   alert.AlertLevelError,
			Type:    alert.AlertTypeLogEncrypt,
			Message: "保存加密日志文件失败",
			Error:   err,
			Module:  "LogEncryptor",
		})
		return "", "", fmt.Errorf("保存加密日志文件失败: %w", err)
	}

	// 如果配置了公钥加密
	if e.config.ConfigManager.LogManager.Encryption.PublicKeyPath != "" {
		// 创建非对称加密器
		asym, err := crypto.CreateAsymmetricEncryptor(
			e.config.ConfigManager.LogManager.Encryption.PublicKeyAlgorithm,
			e.config.ConfigManager.LogManager.Encryption.PublicKeyLength,
		)
		if err != nil {
			e.alerter.Alert(&alert.Alert{
				Level:   alert.AlertLevelError,
				Type:    alert.AlertTypeLogEncrypt,
				Message: "创建非对称加密器失败",
				Error:   err,
				Module:  "LogEncryptor",
			})
			return "", "", fmt.Errorf("创建非对称加密器失败: %w", err)
		}

		// 加载公钥
		if err := asym.LoadPublicKey(e.config.ConfigManager.LogManager.Encryption.PublicKeyPath); err != nil {
			e.alerter.Alert(&alert.Alert{
				Level:   alert.AlertLevelError,
				Type:    alert.AlertTypeLogEncrypt,
				Message: "加载公钥失败",
				Error:   err,
				Module:  "LogEncryptor",
			})
			return "", "", fmt.Errorf("加载公钥失败: %w", err)
		}

		// 加密AES密钥
		encryptedKey, err := asym.EncryptKey(aes.GetKey())
		if err != nil {
			e.alerter.Alert(&alert.Alert{
				Level:   alert.AlertLevelError,
				Type:    alert.AlertTypeLogEncrypt,
				Message: "加密AES密钥失败",
				Error:   err,
				Module:  "LogEncryptor",
			})
			return "", "", fmt.Errorf("加密AES密钥失败: %w", err)
		}

		// 保存加密后的密钥文件
		if err := fileutil.WriteFile(keyPath, encryptedKey, 0644); err != nil {
			e.alerter.Alert(&alert.Alert{
				Level:   alert.AlertLevelError,
				Type:    alert.AlertTypeLogEncrypt,
				Message: "保存密钥文件失败",
				Error:   err,
				Module:  "LogEncryptor",
			})
			return "", "", fmt.Errorf("保存密钥文件失败: %w", err)
		}

		return encryptedPath, keyPath, nil
	}

	return encryptedPath, "", nil
}
