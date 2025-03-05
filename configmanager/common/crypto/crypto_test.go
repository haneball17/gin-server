package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gin-server/config"
)

func TestMain(m *testing.M) {
	// 设置测试环境
	os.Setenv("DEBUG_LEVEL", "true")
	config.InitConfig()

	// 运行测试
	code := m.Run()

	// 清理测试环境
	os.Exit(code)
}

func TestNewAESEncryptor(t *testing.T) {
	tests := []struct {
		name      string
		keyLength int
		wantErr   bool
	}{
		{
			name:      "有效的128位密钥",
			keyLength: 128,
			wantErr:   false,
		},
		{
			name:      "有效的192位密钥",
			keyLength: 192,
			wantErr:   false,
		},
		{
			name:      "有效的256位密钥",
			keyLength: 256,
			wantErr:   false,
		},
		{
			name:      "无效的密钥长度",
			keyLength: 512,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptor, err := NewAESEncryptor(tt.keyLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAESEncryptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && encryptor == nil {
				t.Error("NewAESEncryptor() returned nil encryptor")
			}
		})
	}
}

func TestAESEncryptorWithKey(t *testing.T) {
	tests := []struct {
		name    string
		key     []byte
		wantErr bool
	}{
		{
			name:    "有效的128位密钥",
			key:     make([]byte, 16), // 128位 = 16字节
			wantErr: false,
		},
		{
			name:    "有效的192位密钥",
			key:     make([]byte, 24), // 192位 = 24字节
			wantErr: false,
		},
		{
			name:    "有效的256位密钥",
			key:     make([]byte, 32), // 256位 = 32字节
			wantErr: false,
		},
		{
			name:    "无效的密钥长度",
			key:     make([]byte, 64),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptor, err := NewAESEncryptorWithKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAESEncryptorWithKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && encryptor == nil {
				t.Error("NewAESEncryptorWithKey() returned nil encryptor")
			}
		})
	}
}

func TestAESEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "加密解密空数据",
			data:    []byte{},
			wantErr: false,
		},
		{
			name:    "加密解密普通文本",
			data:    []byte("Hello, World!"),
			wantErr: false,
		},
		{
			name:    "加密解密二进制数据",
			data:    []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
			wantErr: false,
		},
		{
			name:    "加密解密大量数据",
			data:    bytes.Repeat([]byte("Large data "), 1000),
			wantErr: false,
		},
	}

	encryptor, err := NewAESEncryptor(256)
	if err != nil {
		t.Fatalf("创建AES加密器失败: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 加密
			encrypted, err := encryptor.Encrypt(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 确保加密后的数据与原始数据不同
				if bytes.Equal(encrypted, tt.data) {
					t.Error("加密后的数据与原始数据相同")
				}

				// 解密
				decrypted, err := encryptor.Decrypt(encrypted)
				if err != nil {
					t.Errorf("Decrypt() error = %v", err)
					return
				}

				// 验证解密后的数据与原始数据相同
				if !bytes.Equal(decrypted, tt.data) {
					t.Error("解密后的数据与原始数据不匹配")
				}
			}
		})
	}
}

func TestAESDecryptInvalidData(t *testing.T) {
	encryptor, err := NewAESEncryptor(256)
	if err != nil {
		t.Fatalf("创建AES加密器失败: %v", err)
	}

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "解密空数据",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "解密无效数据",
			data:    []byte("invalid data"),
			wantErr: true,
		},
		{
			name:    "解密损坏的加密数据",
			data:    make([]byte, 100),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encryptor.Decrypt(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// createTestRSAPublicKeyPEM 创建测试用的RSA公钥PEM文件
func createTestRSAPublicKeyPEM(t *testing.T) (string, func()) {
	// 生成RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("生成RSA密钥对失败: %v", err)
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "crypto_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}

	// 创建PEM文件路径
	pemFile := filepath.Join(tempDir, "public.pem")

	// 将公钥编码为PKIX格式
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("编码公钥失败: %v", err)
	}

	// 创建PEM块
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	// 创建文件
	file, err := os.Create(pemFile)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("创建PEM文件失败: %v", err)
	}
	defer file.Close()

	// 写入PEM数据
	if err := pem.Encode(file, pemBlock); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("写入PEM数据失败: %v", err)
	}

	// 返回文件路径和清理函数
	return pemFile, func() {
		os.RemoveAll(tempDir)
	}
}

func TestRSAPublicKeyEncryptor(t *testing.T) {
	// 创建测试用的RSA公钥PEM文件
	pemFile, cleanup := createTestRSAPublicKeyPEM(t)
	defer cleanup()

	// 测试从PEM文件创建加密器
	encryptor, err := NewRSAPublicKeyEncryptorFromPEM(pemFile)
	if err != nil {
		t.Fatalf("创建RSA公钥加密器失败: %v", err)
	}
	if encryptor == nil {
		t.Fatal("加密器为空")
	}

	// 测试加密不同长度的密钥
	tests := []struct {
		name    string
		key     []byte
		wantErr bool
	}{
		{
			name:    "加密16字节密钥",
			key:     make([]byte, 16),
			wantErr: false,
		},
		{
			name:    "加密32字节密钥",
			key:     make([]byte, 32),
			wantErr: false,
		},
		{
			name:    "加密过长的密钥",
			key:     make([]byte, 512),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用随机数填充密钥
			if _, err := rand.Read(tt.key); err != nil {
				t.Fatalf("生成随机密钥失败: %v", err)
			}

			// 加密密钥
			encrypted, err := encryptor.EncryptKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 验证加密后的数据不为空且与原始数据不同
				if len(encrypted) == 0 {
					t.Error("加密后的数据为空")
				}
				if bytes.Equal(encrypted, tt.key) {
					t.Error("加密后的数据与原始数据相同")
				}
			}
		})
	}
}

func TestRSAPublicKeyEncryptorWithInvalidPEM(t *testing.T) {
	tests := []struct {
		name     string
		pemData  []byte
		wantErr  bool
		errorMsg string
	}{
		{
			name:     "不存在的PEM文件",
			pemData:  nil,
			wantErr:  true,
			errorMsg: "",
		},
		{
			name:     "无效的PEM数据",
			pemData:  []byte("invalid pem data"),
			wantErr:  true,
			errorMsg: "无效的PEM数据",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pemFile string
			if tt.pemData != nil {
				// 创建临时文件
				tempDir, err := os.MkdirTemp("", "crypto_test")
				if err != nil {
					t.Fatalf("创建临时目录失败: %v", err)
				}
				defer os.RemoveAll(tempDir)

				pemFile = filepath.Join(tempDir, "invalid.pem")
				if err := os.WriteFile(pemFile, tt.pemData, 0644); err != nil {
					t.Fatalf("写入测试数据失败: %v", err)
				}
			} else {
				pemFile = "nonexistent.pem"
			}

			_, err := NewRSAPublicKeyEncryptorFromPEM(pemFile)
			if err == nil {
				t.Error("期望错误但没有发生")
				return
			}
			if !tt.wantErr {
				t.Errorf("NewRSAPublicKeyEncryptorFromPEM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("错误消息不匹配，got = %v, want 包含 %v", err.Error(), tt.errorMsg)
			}
		})
	}
}
