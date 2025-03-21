package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gin-server/config"

	_ "github.com/go-sql-driver/mysql"
)

// 证书和密钥的示例内容
const (
	certContent = `
-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUECPZAZ4aIZnKzuhc9fE4/BUNdnIwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMzA0MTAwODM5NTlaFw0yNDA0
MDkwODM5NTlaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQC8YobzCNZs0jF8G1X3GhJm3HkLZ2s9F5X8O5zQFG+e
88JxCJYgLPYdd9yaiaCSXvYRWMtDcQgDmA3S7YUCh5p1U/dUibMzjnCOgZ7QQS1K
vY7qwTFEHgL+Qo6mW8UAcwraTTY7Avpj5j+Mp3pEEGQzHS7H94PnzPBqpUmIcYgj
2zVbHvLkaIxuxaLGT+a6QMLtZbYQKNvtlqRRZeGO+Zn7/jpzL6/GwSgXfIb/XaUE
L6d67xxUibNI3KtQYZU9V7iCgwdQ5t9xttO0ZqrOlPv4QUs/t3wZmD+TiDkEj1kC
hHbbcYOqRRmFQo/8UUz7tseTYAEkSyFonXmwGHgzAgMBAAGjUzBRMB0GA1UdDgQW
BBTFVyKGaKCYz4NfJeK5YuT0PCxRUDAfBgNVHSMEGDAWgBTFVyKGaKCYz4NfJeK5
YuT0PCxRUDAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCkwRtL
lHH8D9uY9OKkGqPPhGJAZX9kkVX1w2pp/aIJJi8j0IsbjhrdGFUVk8UbQgj+Gub0
/NMGqSqkWQT3aNZ41S9wbZSsYQHnvlY6A7+WArcPQ0PXeFakUjwGYR8hPXPzRTlG
YnBHwxuLCRvO9XBCbN3GbUJjL2vlHwJyqw0QmvAi3xTUm/8pnVnVWGE+W/wFbAnc
5/Vp0mGq8GRxHIz0UYEZVYLNBPzoQoGYbJLYwLN5+jQpAIpZdi03SGPLfznGJ6jY
FuMIPI0NLvkW+7Ur+fkG55bNH0hx8jVcDx+Vqq/GFqvvxN10uGVB3v2MrI3uDEbv
PB4oFHHUdE4TUrfr
-----END CERTIFICATE-----
`

	keyContent = `
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC8YobzCNZs0jF8
G1X3GhJm3HkLZ2s9F5X8O5zQFG+e88JxCJYgLPYdd9yaiaCSXvYRWMtDcQgDmA3S
7YUCh5p1U/dUibMzjnCOgZ7QQS1KvY7qwTFEHgL+Qo6mW8UAcwraTTY7Avpj5j+M
p3pEEGQzHS7H94PnzPBqpUmIcYgj2zVbHvLkaIxuxaLGT+a6QMLtZbYQKNvtlqRR
ZeGO+Zn7/jpzL6/GwSgXfIb/XaUEL6d67xxUibNI3KtQYZU9V7iCgwdQ5t9xttO0
ZqrOlPv4QUs/t3wZmD+TiDkEj1kChHbbcYOqRRmFQo/8UUz7tseTYAEkSyFonXmw
GHgzAgMBAAECggEALEf3AnmEQ8vvtTvTWyf7rDpB7xx5BTCxq6qszdyHCJXtVgdX
WFC2Y3kKWj2ATm0SiWHRSOKM7QvjNr8Nb5YdYhcD4mMTX2hkSKFuMviGh5xWnPKS
AqHa3kGPO89gQvQQkrJJh6GUhhaGXAarR8jzIkEY+TlVpMRnEyzL/2CY6ZHGS3e3
c2GQqRlczQuFJQOxJJO3Tn86ZVtLS1FE33QdAkXx1JLz5fKFRULvW67ZZ5AHhjrF
bEBZxRXjS0q4hcRNi2AHgKKx/bK6xyCTKFoHjnHp2JMqvBz/FCDFKrPGAHMQvfES
N5GRwEVCjk/bmOxYXcKcRJIYJcO6tSGxSO2BtNJ9AQKBgQD/D0h9QpXyXtyBPg+Z
4q+W4CK/DoXSrOJU1jj7XE7aIxXHR5s2/dx6GQjvwvgZXgHJpGVJBvB/UxGhKb2m
C17xx3r+fWuOZDBXdkv0C9BoeuEUK2qJbuU1PDXV2iuTAi5GEwKDCrHiBfUNP+Nn
lVLdyn/fA8IQwkXIhdEurLTCXwKBgQC9DOCHzHVSzCnDLGgH5jEVbCdLAZEjIrbZ
c96QceSYfwQZqFxnDyDXiwUnUYOV3H1+pgWmQJEe2/9Jf9dMiJGsXxPxNpXCaqGj
Uo1vJgAQyCKQpPfWNcKRnmGLNf5Uy44NQDgYV1gJl+WRjOH+SPRtmcYvNrEqAUbL
VFcKvSTqnQKBgQDQsaeR6a0W0T76lxpUJU1dX/O1iRtgQK6t0Ib+42UPxJnr8OFU
WVb8wgKo7KNVEGskS3JOEiIQPBnqRpW2jJJfHEya60q3o16zC7Rc0BxkqQfTt+2e
kdrndfWVj1Xp45CiaFikG0YhMvwM6JtKM9Uy09xGiJzBAK2+Q3R0vwJYbQKBgE9a
G3MclXYK7sD5mQzPztQu63zJ+OFG9t4XQaxC9oj41HHXbVqevbxLMYgxSsWGa8+o
ZY7dQntF6xsToSu3TsW8tG3aLXELOEU6aKkY2t2JMQAqKJoXrwXfWogC3ytfUnfZ
1qPuoPrq8iqNtmKGxanUtESrpXoAkJwZDjpGLnvtAoGAC3HpgfDQn2ma1eBUiMWf
rh5iQJXgxn99lDh5QxuHSZjQqQXZwYDYYHQfXf+8vE7YC3Bj0rJ4acekpsBz36gn
G1FWmwAAEVu4qvBcC6DDMSctpIbJEyXZOSBvBGxWZ/gL3ooSpCfMaIPKKl3Invrx
QA5BICbW1gtrIByvXbpPuPE=
-----END PRIVATE KEY-----
`
)

// 证书测试工具结构
type CertTest struct {
	serverURL string
	db        *sql.DB
	client    *http.Client
	tempDir   string
}

// 初始化测试工具
func NewCertTest(serverURL string) (*CertTest, error) {
	// 创建临时目录保存测试文件
	tempDir, err := ioutil.TempDir("", "cert_test")
	if err != nil {
		return nil, fmt.Errorf("创建临时目录失败: %w", err)
	}

	// 初始化配置
	cfg := config.GetConfig()

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&allowNativePasswords=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("测试数据库连接失败: %w", err)
	}

	return &CertTest{
		serverURL: serverURL,
		db:        db,
		client:    &http.Client{Timeout: 10 * time.Second},
		tempDir:   tempDir,
	}, nil
}

// 清理临时目录
func (ct *CertTest) Cleanup() {
	if ct.db != nil {
		ct.db.Close()
	}
	if ct.tempDir != "" {
		os.RemoveAll(ct.tempDir)
	}
}

// 创建测试证书文件
func (ct *CertTest) createTestCertFile() (string, error) {
	filePath := filepath.Join(ct.tempDir, "test_cert.pem")
	err := ioutil.WriteFile(filePath, []byte(certContent), 0644)
	if err != nil {
		return "", fmt.Errorf("创建测试证书文件失败: %w", err)
	}
	return filePath, nil
}

// 创建测试密钥文件
func (ct *CertTest) createTestKeyFile() (string, error) {
	filePath := filepath.Join(ct.tempDir, "test_key.pem")
	err := ioutil.WriteFile(filePath, []byte(keyContent), 0644)
	if err != nil {
		return "", fmt.Errorf("创建测试密钥文件失败: %w", err)
	}
	return filePath, nil
}

// 创建并发送文件上传请求
func (ct *CertTest) sendFileUploadRequest(url, filePath, fieldName string) (*http.Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("创建表单字段失败: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("复制文件内容失败: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("关闭表单写入器失败: %w", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return ct.client.Do(req)
}

// 检查证书表记录
func (ct *CertTest) checkCertRecord(entityType, entityID string) (bool, error) {
	var count int
	err := ct.db.QueryRow("SELECT COUNT(*) FROM certs WHERE entity_type = ? AND entity_id = ?", entityType, entityID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("查询证书记录失败: %w", err)
	}
	return count > 0, nil
}

// 测试用户证书绑定
func (ct *CertTest) TestBindUserCert(userID int) error {
	fmt.Printf("测试用户证书绑定 (userID=%d)...\n", userID)

	// 检查用户是否存在
	var count int
	err := ct.db.QueryRow("SELECT COUNT(*) FROM users WHERE userID = ?", userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("检查用户是否存在失败: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("用户不存在，请先创建用户")
	}

	// 创建测试证书文件
	certFile, err := ct.createTestCertFile()
	if err != nil {
		return err
	}

	// 发送请求
	url := fmt.Sprintf("%s/bind/users/%d/cert", ct.serverURL, userID)
	resp, err := ct.sendFileUploadRequest(url, certFile, "cert")
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应内容失败: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(body))
	}

	// 检查数据库记录
	exists, err := ct.checkCertRecord("user", strconv.Itoa(userID))
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("证书记录未添加到数据库")
	}

	fmt.Println("用户证书绑定测试成功！")
	fmt.Printf("响应: %s\n", string(body))
	return nil
}

// 测试用户密钥绑定
func (ct *CertTest) TestBindUserKey(userID int) error {
	fmt.Printf("测试用户密钥绑定 (userID=%d)...\n", userID)

	// 检查用户是否存在
	var count int
	err := ct.db.QueryRow("SELECT COUNT(*) FROM users WHERE userID = ?", userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("检查用户是否存在失败: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("用户不存在，请先创建用户")
	}

	// 创建测试密钥文件
	keyFile, err := ct.createTestKeyFile()
	if err != nil {
		return err
	}

	// 发送请求
	url := fmt.Sprintf("%s/bind/users/%d/key", ct.serverURL, userID)
	resp, err := ct.sendFileUploadRequest(url, keyFile, "key")
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应内容失败: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(body))
	}

	// 检查数据库记录
	exists, err := ct.checkCertRecord("user", strconv.Itoa(userID))
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("密钥记录未添加到数据库")
	}

	fmt.Println("用户密钥绑定测试成功！")
	fmt.Printf("响应: %s\n", string(body))
	return nil
}

// 测试设备证书绑定
func (ct *CertTest) TestBindDeviceCert(deviceID string) error {
	fmt.Printf("测试设备证书绑定 (deviceID=%s)...\n", deviceID)

	// 检查设备是否存在
	var count int
	err := ct.db.QueryRow("SELECT COUNT(*) FROM devices WHERE deviceID = ?", deviceID).Scan(&count)
	if err != nil {
		return fmt.Errorf("检查设备是否存在失败: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("设备不存在，请先创建设备")
	}

	// 创建测试证书文件
	certFile, err := ct.createTestCertFile()
	if err != nil {
		return err
	}

	// 发送请求
	url := fmt.Sprintf("%s/bind/devices/%s/cert", ct.serverURL, deviceID)
	resp, err := ct.sendFileUploadRequest(url, certFile, "cert")
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应内容失败: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(body))
	}

	// 检查数据库记录
	exists, err := ct.checkCertRecord("device", deviceID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("证书记录未添加到数据库")
	}

	fmt.Println("设备证书绑定测试成功！")
	fmt.Printf("响应: %s\n", string(body))
	return nil
}

// 测试设备密钥绑定
func (ct *CertTest) TestBindDeviceKey(deviceID string) error {
	fmt.Printf("测试设备密钥绑定 (deviceID=%s)...\n", deviceID)

	// 检查设备是否存在
	var count int
	err := ct.db.QueryRow("SELECT COUNT(*) FROM devices WHERE deviceID = ?", deviceID).Scan(&count)
	if err != nil {
		return fmt.Errorf("检查设备是否存在失败: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("设备不存在，请先创建设备")
	}

	// 创建测试密钥文件
	keyFile, err := ct.createTestKeyFile()
	if err != nil {
		return err
	}

	// 发送请求
	url := fmt.Sprintf("%s/bind/devices/%s/key", ct.serverURL, deviceID)
	resp, err := ct.sendFileUploadRequest(url, keyFile, "key")
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应内容失败: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(body))
	}

	// 检查数据库记录
	exists, err := ct.checkCertRecord("device", deviceID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("密钥记录未添加到数据库")
	}

	fmt.Println("设备密钥绑定测试成功！")
	fmt.Printf("响应: %s\n", string(body))
	return nil
}

// 测试获取证书信息
func (ct *CertTest) TestGetCertInfo(entityType, entityID string) error {
	fmt.Printf("测试获取证书信息 (type=%s, id=%s)...\n", entityType, entityID)

	// 发送请求
	url := fmt.Sprintf("%s/cert/info?type=%s&id=%s", ct.serverURL, entityType, entityID)
	resp, err := ct.client.Get(url)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应内容失败: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(body))
	}

	fmt.Println("获取证书信息测试成功！")
	fmt.Printf("响应: %s\n", string(body))
	return nil
}

// 创建测试用户
func (ct *CertTest) CreateTestUser(userName string, userType int, gatewayDeviceID string) (int, error) {
	fmt.Printf("创建测试用户 (userName=%s)...\n", userName)

	// 检查设备是否存在
	var count int
	err := ct.db.QueryRow("SELECT COUNT(*) FROM devices WHERE deviceID = ?", gatewayDeviceID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("检查设备是否存在失败: %w", err)
	}
	if count == 0 {
		return 0, fmt.Errorf("网关设备不存在，请先创建设备: %s", gatewayDeviceID)
	}

	// 生成随机用户ID
	userID := int(time.Now().Unix()%10000) + 1000

	// 插入测试用户
	_, err = ct.db.Exec(
		"INSERT INTO users (userName, passWD, userID, userType, gatewayDeviceID, created_at) VALUES (?, ?, ?, ?, ?, NOW())",
		userName, "password123", userID, userType, gatewayDeviceID,
	)
	if err != nil {
		return 0, fmt.Errorf("创建测试用户失败: %w", err)
	}

	fmt.Printf("成功创建测试用户，userID=%d\n", userID)
	return userID, nil
}

// 创建测试设备
func (ct *CertTest) CreateTestDevice(deviceName string, deviceType int) (string, error) {
	fmt.Printf("创建测试设备 (deviceName=%s)...\n", deviceName)

	// 生成随机设备ID (12位字符)
	timestamp := time.Now().Unix() % 1000000
	deviceID := fmt.Sprintf("DEV%06d", timestamp)

	// 插入测试设备，上级设备ID与设备ID相同
	_, err := ct.db.Exec(
		"INSERT INTO devices (deviceName, deviceType, passWD, deviceID, superiorDeviceID, created_at) VALUES (?, ?, ?, ?, ?, NOW())",
		deviceName, deviceType, "password123", deviceID, deviceID,
	)
	if err != nil {
		return "", fmt.Errorf("创建测试设备失败: %w", err)
	}

	fmt.Printf("成功创建测试设备，deviceID=%s\n", deviceID)
	return deviceID, nil
}

// RunAllTests 运行所有测试
func (ct *CertTest) RunAllTests(userID int, deviceID string) error {
	tests := []struct {
		name string
		test func() error
	}{
		{"用户证书绑定", func() error { return ct.TestBindUserCert(userID) }},
		{"用户密钥绑定", func() error { return ct.TestBindUserKey(userID) }},
		{"设备证书绑定", func() error { return ct.TestBindDeviceCert(deviceID) }},
		{"设备密钥绑定", func() error { return ct.TestBindDeviceKey(deviceID) }},
		{"获取用户证书信息", func() error { return ct.TestGetCertInfo("user", strconv.Itoa(userID)) }},
		{"获取设备证书信息", func() error { return ct.TestGetCertInfo("device", deviceID) }},
	}

	failed := false
	for _, test := range tests {
		fmt.Printf("\n=== 运行测试: %s ===\n", test.name)
		if err := test.test(); err != nil {
			fmt.Printf("测试失败: %v\n", err)
			failed = true
		}
	}

	if failed {
		return fmt.Errorf("部分测试失败")
	}
	return nil
}

func main() {
	// 命令行参数
	serverURL := flag.String("server", "http://localhost:8080", "服务器URL")
	createUser := flag.Bool("create-user", false, "创建测试用户")
	createDevice := flag.Bool("create-device", false, "创建测试设备")
	userName := flag.String("user-name", "testUser", "测试用户名")
	userType := flag.Int("user-type", 1, "测试用户类型")
	deviceName := flag.String("device-name", "testDevice", "测试设备名")
	deviceType := flag.Int("device-type", 1, "测试设备类型")
	userID := flag.Int("user-id", 0, "要测试的用户ID")
	deviceID := flag.String("device-id", "", "要测试的设备ID")
	runUserCertTest := flag.Bool("test-user-cert", false, "运行用户证书绑定测试")
	runUserKeyTest := flag.Bool("test-user-key", false, "运行用户密钥绑定测试")
	runDeviceCertTest := flag.Bool("test-device-cert", false, "运行设备证书绑定测试")
	runDeviceKeyTest := flag.Bool("test-device-key", false, "运行设备密钥绑定测试")
	runGetCertInfoTest := flag.Bool("test-get-cert", false, "运行获取证书信息测试")
	runAllTests := flag.Bool("test-all", false, "运行所有测试")

	flag.Parse()

	// 初始化测试工具
	certTest, err := NewCertTest(*serverURL)
	if err != nil {
		log.Fatalf("初始化测试工具失败: %v", err)
	}
	defer certTest.Cleanup()

	// 创建测试设备
	createdDeviceID := *deviceID
	if *createDevice {
		createdDeviceID, err = certTest.CreateTestDevice(*deviceName, *deviceType)
		if err != nil {
			log.Fatalf("创建测试设备失败: %v", err)
		}
	}

	// 如果没有指定设备ID，使用创建的设备ID
	if createdDeviceID == "" {
		log.Fatalf("未指定设备ID，请使用 -device-id 参数指定，或者使用 -create-device 创建新设备")
	}

	// 创建测试用户
	createdUserID := *userID
	if *createUser {
		createdUserID, err = certTest.CreateTestUser(*userName, *userType, createdDeviceID)
		if err != nil {
			log.Fatalf("创建测试用户失败: %v", err)
		}
	}

	// 如果没有指定用户ID，使用创建的用户ID
	if createdUserID == 0 {
		log.Fatalf("未指定用户ID，请使用 -user-id 参数指定，或者使用 -create-user 创建新用户")
	}

	// 运行测试
	if *runAllTests {
		fmt.Println("\n=== 运行所有测试 ===")
		if err := certTest.RunAllTests(createdUserID, createdDeviceID); err != nil {
			log.Fatalf("测试失败: %v", err)
		}
		fmt.Println("\n=== 所有测试通过！===")
		return
	}

	// 运行单个测试
	if *runUserCertTest {
		if err := certTest.TestBindUserCert(createdUserID); err != nil {
			log.Fatalf("用户证书绑定测试失败: %v", err)
		}
	}

	if *runUserKeyTest {
		if err := certTest.TestBindUserKey(createdUserID); err != nil {
			log.Fatalf("用户密钥绑定测试失败: %v", err)
		}
	}

	if *runDeviceCertTest {
		if err := certTest.TestBindDeviceCert(createdDeviceID); err != nil {
			log.Fatalf("设备证书绑定测试失败: %v", err)
		}
	}

	if *runDeviceKeyTest {
		if err := certTest.TestBindDeviceKey(createdDeviceID); err != nil {
			log.Fatalf("设备密钥绑定测试失败: %v", err)
		}
	}

	if *runGetCertInfoTest {
		fmt.Println("输入要查询的实体类型 (user/device):")
		var entityType string
		fmt.Scanln(&entityType)

		fmt.Println("输入要查询的实体ID:")
		var entityID string
		fmt.Scanln(&entityID)

		if entityType == "" || entityID == "" {
			entityType = "user"
			entityID = strconv.Itoa(createdUserID)
			fmt.Printf("使用默认值: type=%s, id=%s\n", entityType, entityID)
		}

		if err := certTest.TestGetCertInfo(entityType, entityID); err != nil {
			log.Fatalf("获取证书信息测试失败: %v", err)
		}
	}

	// 如果没有指定任何测试，显示帮助
	if !(*runUserCertTest || *runUserKeyTest || *runDeviceCertTest || *runDeviceKeyTest || *runGetCertInfoTest || *runAllTests || *createUser || *createDevice) {
		fmt.Println("请指定要运行的测试，例如:")
		fmt.Println("  -create-device -create-user -test-all: 创建测试设备和用户，然后运行所有测试")
		fmt.Println("  -device-id DEV123456 -user-id 1001 -test-user-cert: 使用指定的设备和用户测试用户证书绑定")
		fmt.Println("  -device-id DEV123456 -user-id 1001 -test-device-key: 使用指定的设备和用户测试设备密钥绑定")
		fmt.Println("运行 -help 查看所有选项")
	}
}
