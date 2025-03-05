package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"gin-server/config"
)

// Encryptor 加密器接口
type Encryptor interface {
	// Encrypt 加密数据
	Encrypt(data []byte) ([]byte, error)
	// Decrypt 解密数据
	Decrypt(data []byte) ([]byte, error)
}

// KeyEncryptor 密钥加密器接口
type KeyEncryptor interface {
	// EncryptKey 加密密钥
	EncryptKey(key []byte) ([]byte, error)
	// DecryptKey 解密密钥
	DecryptKey(encryptedKey []byte) ([]byte, error)
}

// RSAPublicKeyEncryptor RSA公钥加密器
type RSAPublicKeyEncryptor struct {
	publicKey *rsa.PublicKey
}

// NewRSAPublicKeyEncryptorFromPEM 从PEM文件创建RSA公钥加密器
func NewRSAPublicKeyEncryptorFromPEM(pemFile string) (*RSAPublicKeyEncryptor, error) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("从PEM文件加载RSA公钥: %s\n", pemFile)
	}

	// 读取PEM文件
	pemData, err := os.ReadFile(pemFile)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("读取PEM文件失败: %v\n", err)
		}
		return nil, fmt.Errorf("读取PEM文件失败: %w", err)
	}

	// 解码PEM块
	block, _ := pem.Decode(pemData)
	if block == nil {
		if cfg.DebugLevel == "true" {
			log.Println("无效的PEM数据")
		}
		return nil, errors.New("无效的PEM数据")
	}

	// 解析公钥
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("解析公钥失败: %v\n", err)
		}
		return nil, fmt.Errorf("解析公钥失败: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		if cfg.DebugLevel == "true" {
			log.Println("不是RSA公钥")
		}
		return nil, errors.New("不是RSA公钥")
	}

	if cfg.DebugLevel == "true" {
		log.Printf("RSA公钥加载成功，密钥长度: %d\n", rsaPub.Size()*8)
	}
	return &RSAPublicKeyEncryptor{publicKey: rsaPub}, nil
}

// EncryptKey 使用RSA公钥加密密钥
func (e *RSAPublicKeyEncryptor) EncryptKey(key []byte) ([]byte, error) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("开始RSA加密密钥，密钥长度: %d\n", len(key))
	}

	// 使用PKCS1v15填充方案加密
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, e.publicKey, key)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("RSA加密失败: %v\n", err)
		}
		return nil, fmt.Errorf("RSA加密失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Printf("RSA加密完成，加密后长度: %d\n", len(encrypted))
	}
	return encrypted, nil
}

// AESEncryptor AES加密器
type AESEncryptor struct {
	key []byte
}

// NewAESEncryptor 创建新的AES加密器
func NewAESEncryptor(keyLength int) (*AESEncryptor, error) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("创建新的AES加密器，密钥长度: %d\n", keyLength)
	}

	if keyLength != 128 && keyLength != 192 && keyLength != 256 {
		return nil, fmt.Errorf("不支持的AES密钥长度: %d", keyLength)
	}

	// 生成随机密钥
	key := make([]byte, keyLength/8)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("生成AES密钥失败: %v\n", err)
		}
		return nil, fmt.Errorf("生成AES密钥失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Println("AES加密器创建成功")
	}
	return &AESEncryptor{key: key}, nil
}

// NewAESEncryptorWithKey 使用指定密钥创建AES加密器
func NewAESEncryptorWithKey(key []byte) (*AESEncryptor, error) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("使用指定密钥创建AES加密器")
	}

	keyLength := len(key) * 8
	if keyLength != 128 && keyLength != 192 && keyLength != 256 {
		return nil, fmt.Errorf("不支持的AES密钥长度: %d", keyLength)
	}

	if cfg.DebugLevel == "true" {
		log.Println("AES加密器创建成功")
	}
	return &AESEncryptor{key: key}, nil
}

// GetKey 获取AES密钥
func (e *AESEncryptor) GetKey() []byte {
	return e.key
}

// Encrypt 使用AES-GCM模式加密数据
func (e *AESEncryptor) Encrypt(data []byte) ([]byte, error) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("开始AES加密，数据长度: %d\n", len(data))
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("创建AES密码块失败: %v\n", err)
		}
		return nil, fmt.Errorf("创建AES密码块失败: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("创建GCM失败: %v\n", err)
		}
		return nil, fmt.Errorf("创建GCM失败: %w", err)
	}

	// 生成随机数
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("生成随机数失败: %v\n", err)
		}
		return nil, fmt.Errorf("生成随机数失败: %w", err)
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	if cfg.DebugLevel == "true" {
		log.Printf("AES加密完成，加密后数据长度: %d\n", len(ciphertext))
	}
	return ciphertext, nil
}

// Decrypt 使用AES-GCM模式解密数据
func (e *AESEncryptor) Decrypt(data []byte) ([]byte, error) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("开始AES解密，数据长度: %d\n", len(data))
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("创建AES密码块失败: %v\n", err)
		}
		return nil, fmt.Errorf("创建AES密码块失败: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("创建GCM失败: %v\n", err)
		}
		return nil, fmt.Errorf("创建GCM失败: %w", err)
	}

	if len(data) < gcm.NonceSize() {
		if cfg.DebugLevel == "true" {
			log.Println("加密数据长度不足")
		}
		return nil, errors.New("加密数据长度不足")
	}

	nonce := data[:gcm.NonceSize()]
	ciphertext := data[gcm.NonceSize():]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("解密数据失败: %v\n", err)
		}
		return nil, fmt.Errorf("解密数据失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Printf("AES解密完成，解密后数据长度: %d\n", len(plaintext))
	}
	return plaintext, nil
}
