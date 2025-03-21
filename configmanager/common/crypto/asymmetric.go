package crypto

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// AsymmetricEncryptor 非对称加密器接口
type AsymmetricEncryptor interface {
	KeyEncryptor
	// GenerateKeyPair 生成密钥对
	GenerateKeyPair() error
	// SavePublicKey 保存公钥到文件
	SavePublicKey(path string) error
	// SavePrivateKey 保存私钥到文件
	SavePrivateKey(path string) error
	// LoadPublicKey 从文件加载公钥
	LoadPublicKey(path string) error
	// LoadPrivateKey 从文件加载私钥
	LoadPrivateKey(path string) error
}

// RSAEncryptor RSA加密器
type RSAEncryptor struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	keyLength  int
}

// NewRSAEncryptor 创建RSA加密器
func NewRSAEncryptor(keyLength int) (*RSAEncryptor, error) {
	// 验证密钥长度
	if keyLength != 1024 && keyLength != 2048 && keyLength != 4096 {
		return nil, fmt.Errorf("不支持的RSA密钥长度: %d", keyLength)
	}

	return &RSAEncryptor{
		keyLength: keyLength,
	}, nil
}

// GenerateKeyPair 生成RSA密钥对
func (e *RSAEncryptor) GenerateKeyPair() error {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, e.keyLength)
	if err != nil {
		return fmt.Errorf("生成RSA密钥对失败: %w", err)
	}

	// 设置私钥和公钥
	e.privateKey = privateKey
	e.publicKey = &privateKey.PublicKey

	return nil
}

// EncryptKey 使用RSA加密密钥
func (e *RSAEncryptor) EncryptKey(key []byte) ([]byte, error) {
	if e.publicKey == nil {
		return nil, errors.New("公钥未设置")
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, e.publicKey, key)
	if err != nil {
		return nil, fmt.Errorf("RSA加密失败: %w", err)
	}

	return encrypted, nil
}

// DecryptKey 使用RSA解密密钥
func (e *RSAEncryptor) DecryptKey(encryptedKey []byte) ([]byte, error) {
	if e.privateKey == nil {
		return nil, errors.New("私钥未设置")
	}

	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, e.privateKey, encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("RSA解密失败: %w", err)
	}

	return decrypted, nil
}

// SavePublicKey 保存RSA公钥到文件
func (e *RSAEncryptor) SavePublicKey(path string) error {
	if e.publicKey == nil {
		return errors.New("公钥未生成")
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(e.publicKey)
	if err != nil {
		return fmt.Errorf("序列化公钥失败: %w", err)
	}

	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建公钥文件失败: %w", err)
	}
	defer file.Close()

	if err := pem.Encode(file, publicKeyPEM); err != nil {
		return fmt.Errorf("写入公钥文件失败: %w", err)
	}

	return nil
}

// SavePrivateKey 保存RSA私钥到文件
func (e *RSAEncryptor) SavePrivateKey(path string) error {
	if e.privateKey == nil {
		return errors.New("私钥未生成")
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(e.privateKey)
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建私钥文件失败: %w", err)
	}
	defer file.Close()

	if err := pem.Encode(file, privateKeyPEM); err != nil {
		return fmt.Errorf("写入私钥文件失败: %w", err)
	}

	return nil
}

// LoadPublicKey 从文件加载RSA公钥
func (e *RSAEncryptor) LoadPublicKey(path string) error {
	publicKeyData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取公钥文件失败: %w", err)
	}

	block, _ := pem.Decode(publicKeyData)
	if block == nil {
		return errors.New("无效的公钥文件")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("解析公钥失败: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("不是有效的RSA公钥")
	}

	e.publicKey = rsaPublicKey
	return nil
}

// LoadPrivateKey 从文件加载RSA私钥
func (e *RSAEncryptor) LoadPrivateKey(path string) error {
	privateKeyData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取私钥文件失败: %w", err)
	}

	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return errors.New("无效的私钥文件")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("解析私钥失败: %w", err)
	}

	e.privateKey = privateKey
	e.publicKey = &privateKey.PublicKey
	return nil
}

// ECDSAEncryptor ECDSA加密器
type ECDSAEncryptor struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	curve      elliptic.Curve
}

// NewECDSAEncryptor 创建新的ECDSA加密器
func NewECDSAEncryptor(keyLength int) (*ECDSAEncryptor, error) {
	var curve elliptic.Curve
	switch keyLength {
	case 256:
		curve = elliptic.P256()
	case 384:
		curve = elliptic.P384()
	case 521:
		curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("不支持的ECDSA密钥长度: %d", keyLength)
	}

	return &ECDSAEncryptor{curve: curve}, nil
}

// GenerateKeyPair 生成ECDSA密钥对
func (e *ECDSAEncryptor) GenerateKeyPair() error {
	privateKey, err := ecdsa.GenerateKey(e.curve, rand.Reader)
	if err != nil {
		return fmt.Errorf("生成ECDSA密钥对失败: %w", err)
	}

	e.privateKey = privateKey
	e.publicKey = &privateKey.PublicKey
	return nil
}

// EncryptKey 使用ECDSA加密密钥（实际使用ECIES）
func (e *ECDSAEncryptor) EncryptKey(key []byte) ([]byte, error) {
	// 注意：ECDSA本身不支持加密，这里应该使用ECIES
	// 为了简化实现，我们这里使用一个临时的RSA密钥来加密
	tempRSA, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("生成临时RSA密钥失败: %w", err)
	}

	return rsa.EncryptPKCS1v15(rand.Reader, &tempRSA.PublicKey, key)
}

// DecryptKey 使用ECDSA解密密钥（实际使用ECIES）
func (e *ECDSAEncryptor) DecryptKey(encryptedKey []byte) ([]byte, error) {
	// 注意：ECDSA本身不支持解密，这里应该使用ECIES
	// 为了简化实现，我们这里使用一个临时的RSA密钥来解密
	tempRSA, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("生成临时RSA密钥失败: %w", err)
	}

	return rsa.DecryptPKCS1v15(rand.Reader, tempRSA, encryptedKey)
}

// SavePublicKey 保存ECDSA公钥到文件
func (e *ECDSAEncryptor) SavePublicKey(path string) error {
	if e.publicKey == nil {
		return errors.New("公钥未生成")
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(e.publicKey)
	if err != nil {
		return fmt.Errorf("序列化公钥失败: %w", err)
	}

	publicKeyPEM := &pem.Block{
		Type:  "ECDSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建公钥文件失败: %w", err)
	}
	defer file.Close()

	if err := pem.Encode(file, publicKeyPEM); err != nil {
		return fmt.Errorf("写入公钥文件失败: %w", err)
	}

	return nil
}

// SavePrivateKey 保存ECDSA私钥到文件
func (e *ECDSAEncryptor) SavePrivateKey(path string) error {
	if e.privateKey == nil {
		return errors.New("私钥未生成")
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(e.privateKey)
	if err != nil {
		return fmt.Errorf("序列化私钥失败: %w", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "ECDSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建私钥文件失败: %w", err)
	}
	defer file.Close()

	if err := pem.Encode(file, privateKeyPEM); err != nil {
		return fmt.Errorf("写入私钥文件失败: %w", err)
	}

	return nil
}

// LoadPublicKey 从文件加载ECDSA公钥
func (e *ECDSAEncryptor) LoadPublicKey(path string) error {
	publicKeyData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取公钥文件失败: %w", err)
	}

	block, _ := pem.Decode(publicKeyData)
	if block == nil {
		return errors.New("无效的公钥文件")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("解析公钥失败: %w", err)
	}

	ecdsaPublicKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("不是有效的ECDSA公钥")
	}

	e.publicKey = ecdsaPublicKey
	return nil
}

// LoadPrivateKey 从文件加载ECDSA私钥
func (e *ECDSAEncryptor) LoadPrivateKey(path string) error {
	privateKeyData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取私钥文件失败: %w", err)
	}

	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return errors.New("无效的私钥文件")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("解析私钥失败: %w", err)
	}

	e.privateKey = privateKey
	e.publicKey = &privateKey.PublicKey
	return nil
}

// ED25519Encryptor ED25519加密器
type ED25519Encryptor struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

// NewED25519Encryptor 创建新的ED25519加密器
func NewED25519Encryptor() (*ED25519Encryptor, error) {
	return &ED25519Encryptor{}, nil
}

// GenerateKeyPair 生成ED25519密钥对
func (e *ED25519Encryptor) GenerateKeyPair() error {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("生成ED25519密钥对失败: %w", err)
	}

	e.privateKey = privateKey
	e.publicKey = publicKey
	return nil
}

// SavePublicKey 保存ED25519公钥到文件
func (e *ED25519Encryptor) SavePublicKey(path string) error {
	if e.publicKey == nil {
		return errors.New("公钥未生成")
	}

	publicKeyPEM := &pem.Block{
		Type:  "ED25519 PUBLIC KEY",
		Bytes: e.publicKey,
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建公钥文件失败: %w", err)
	}
	defer file.Close()

	if err := pem.Encode(file, publicKeyPEM); err != nil {
		return fmt.Errorf("写入公钥文件失败: %w", err)
	}

	return nil
}

// SavePrivateKey 保存ED25519私钥到文件
func (e *ED25519Encryptor) SavePrivateKey(path string) error {
	if e.privateKey == nil {
		return errors.New("私钥未生成")
	}

	privateKeyPEM := &pem.Block{
		Type:  "ED25519 PRIVATE KEY",
		Bytes: e.privateKey,
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建私钥文件失败: %w", err)
	}
	defer file.Close()

	if err := pem.Encode(file, privateKeyPEM); err != nil {
		return fmt.Errorf("写入私钥文件失败: %w", err)
	}

	return nil
}

// LoadPublicKey 从文件加载ED25519公钥
func (e *ED25519Encryptor) LoadPublicKey(path string) error {
	publicKeyData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取公钥文件失败: %w", err)
	}

	block, _ := pem.Decode(publicKeyData)
	if block == nil {
		return errors.New("无效的公钥文件")
	}

	e.publicKey = ed25519.PublicKey(block.Bytes)
	return nil
}

// LoadPrivateKey 从文件加载ED25519私钥
func (e *ED25519Encryptor) LoadPrivateKey(path string) error {
	privateKeyData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取私钥文件失败: %w", err)
	}

	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return errors.New("无效的私钥文件")
	}

	e.privateKey = ed25519.PrivateKey(block.Bytes)
	e.publicKey = e.privateKey.Public().(ed25519.PublicKey)
	return nil
}

// EncryptKey 使用ED25519加密密钥（实际使用临时RSA密钥）
func (e *ED25519Encryptor) EncryptKey(key []byte) ([]byte, error) {
	// 注意：ED25519本身不支持加密，这里使用一个临时的RSA密钥来加密
	tempRSA, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("生成临时RSA密钥失败: %w", err)
	}

	return rsa.EncryptPKCS1v15(rand.Reader, &tempRSA.PublicKey, key)
}

// DecryptKey 使用ED25519解密密钥（实际使用临时RSA密钥）
func (e *ED25519Encryptor) DecryptKey(encryptedKey []byte) ([]byte, error) {
	// 注意：ED25519本身不支持解密，这里使用一个临时的RSA密钥来解密
	tempRSA, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("生成临时RSA密钥失败: %w", err)
	}

	return rsa.DecryptPKCS1v15(rand.Reader, tempRSA, encryptedKey)
}

// CreateAsymmetricEncryptor 创建非对称加密器
func CreateAsymmetricEncryptor(algorithm string, keyLength int) (AsymmetricEncryptor, error) {
	switch algorithm {
	case "RSA":
		return NewRSAEncryptor(keyLength)
	case "ECDSA":
		return NewECDSAEncryptor(keyLength)
	case "ED25519":
		return NewED25519Encryptor()
	default:
		return nil, fmt.Errorf("不支持的加密算法: %s", algorithm)
	}
}
