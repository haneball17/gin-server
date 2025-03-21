package crypto

import (
	"fmt"
	"gin-server/config"
	"gin-server/configmanager/common/fileutil"
	"os"
	"path/filepath"
)

// KeyManager 密钥管理器
type KeyManager struct {
	config *config.Config
}

// NewKeyManager 创建密钥管理器
func NewKeyManager(cfg *config.Config) *KeyManager {
	return &KeyManager{
		config: cfg,
	}
}

// EnsureKeyPair 确保密钥对存在，如果不存在则生成
func (km *KeyManager) EnsureKeyPair() error {
	// 检查密钥目录是否存在
	keyDir := "keys"
	if err := fileutil.EnsureDir(keyDir); err != nil {
		return fmt.Errorf("创建密钥目录失败: %w", err)
	}

	// 检查公钥是否存在
	publicKeyPath := filepath.Join(keyDir, "public.pem")
	if !fileutil.IsFileExists(publicKeyPath) {
		// 创建非对称加密器
		encryptor, err := CreateAsymmetricEncryptor(
			km.config.ConfigManager.LogManager.Encryption.PublicKeyAlgorithm,
			km.config.ConfigManager.LogManager.Encryption.PublicKeyLength,
		)
		if err != nil {
			return fmt.Errorf("创建非对称加密器失败: %w", err)
		}

		// 生成密钥对
		if err := encryptor.GenerateKeyPair(); err != nil {
			return fmt.Errorf("生成密钥对失败: %w", err)
		}

		// 保存公钥
		if err := encryptor.SavePublicKey(publicKeyPath); err != nil {
			return fmt.Errorf("保存公钥失败: %w", err)
		}

		// 保存私钥
		privateKeyPath := filepath.Join(keyDir, "private.pem")
		if err := encryptor.SavePrivateKey(privateKeyPath); err != nil {
			// 如果保存私钥失败，删除已生成的公钥
			os.Remove(publicKeyPath)
			return fmt.Errorf("保存私钥失败: %w", err)
		}

		// 更新配置中的密钥路径
		km.config.ConfigManager.LogManager.Encryption.PublicKeyPath = publicKeyPath
		km.config.ConfigManager.LogManager.Encryption.PrivateKeyPath = privateKeyPath
	}

	return nil
}
