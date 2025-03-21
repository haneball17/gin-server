# 证书绑定功能测试工具

## 简介

这是一个用于测试证书绑定功能的命令行工具，它可以：

1. 创建测试用户和设备
2. 生成测试用的证书和密钥文件
3. 测试用户和设备的证书/密钥绑定功能
4. 测试获取证书信息功能
5. 验证数据库记录是否正确创建

## 环境要求

- Go 1.15+
- MySQL 数据库（配置在 config.go 中）
- 服务器已运行且可访问

## 使用方法

### 编译

```bash
cd test
go build -o cert_test cert_test.go
```

### 运行

测试工具提供了多种命令行参数来控制测试行为：

```bash
# 查看帮助
./cert_test -help

# 创建测试设备和用户，并运行所有测试
./cert_test -create-device -create-user -test-all

# 使用指定的用户ID和设备ID测试
./cert_test -user-id 1001 -device-id DEV123456 -test-all

# 仅测试用户证书绑定
./cert_test -user-id 1001 -device-id DEV123456 -test-user-cert

# 指定服务器地址（默认为http://localhost:8080）
./cert_test -server http://192.168.1.100:8080 -create-device -create-user -test-all
```

### 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-server` | 服务器URL | http://localhost:8080 |
| `-create-user` | 创建测试用户 | false |
| `-create-device` | 创建测试设备 | false |
| `-user-name` | 测试用户名 | testUser |
| `-user-type` | 测试用户类型 | 1 |
| `-device-name` | 测试设备名 | testDevice |
| `-device-type` | 测试设备类型 | 1 |
| `-user-id` | 要测试的用户ID | 0 |
| `-device-id` | 要测试的设备ID | "" |
| `-test-user-cert` | 运行用户证书绑定测试 | false |
| `-test-user-key` | 运行用户密钥绑定测试 | false |
| `-test-device-cert` | 运行设备证书绑定测试 | false |
| `-test-device-key` | 运行设备密钥绑定测试 | false |
| `-test-get-cert` | 运行获取证书信息测试 | false |
| `-test-all` | 运行所有测试 | false |

## 测试流程

1. 初始化测试环境
   - 连接数据库
   - 创建临时目录用于存放测试文件

2. 创建测试数据（可选）
   - 创建测试设备
   - 创建测试用户

3. 执行测试
   - 生成测试用的证书和密钥文件
   - 发送HTTP请求到相应的API
   - 检查响应是否正确
   - 验证数据库中是否正确保存了记录

4. 清理测试环境
   - 删除临时文件
   - 关闭数据库连接

## 示例用法

### 完整测试流程

```bash
# 创建测试设备和用户，然后运行所有测试
./cert_test -create-device -create-user -test-all
```

这个命令会：
1. 创建一个测试设备
2. 创建一个测试用户，关联到刚创建的设备
3. 测试用户证书绑定
4. 测试用户密钥绑定
5. 测试设备证书绑定
6. 测试设备密钥绑定
7. 测试获取用户证书信息
8. 测试获取设备证书信息

### 测试已有数据

```bash
# 使用现有用户和设备测试证书绑定
./cert_test -user-id 1001 -device-id DEV123456 -test-user-cert -test-device-cert
```

这个命令会：
1. 使用ID为1001的用户测试证书绑定
2. 使用ID为DEV123456的设备测试证书绑定

## 注意事项

1. 测试工具会生成测试用的证书和密钥文件，内容是固定的样例数据
2. 创建的测试用户和设备会永久保存在数据库中
3. 运行测试前请确保服务器正在运行且可访问
4. 测试工具使用的数据库配置与服务器相同，来自config.go中的配置 