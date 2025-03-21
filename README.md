# Gin Server

基于Gin框架的服务器应用，提供用户认证、设备管理、配置管理和证书管理功能。

## 功能特性

- 用户管理
  - 用户注册
  - 用户信息查询
  - 用户信息更新
- 设备管理
  - 设备注册
  - 设备信息查询
  - 设备信息更新
- 证书管理（新增功能）
  - 用户证书/密钥绑定
  - 设备证书/密钥绑定
  - 证书信息查询
  - 自动创建证书目录
  - 安全存储和权限控制
- 认证管理
  - 认证记录查询
  - 权限控制
- 配置管理
  - 远程配置同步
  - 文件变更监控
  - 加密传输
  - 多存储方式支持（Gitee/FTP）
- 告警管理
  - 多级别告警（INFO/WARNING/ERROR/FATAL）
  - 多类型告警支持
  - 可扩展的告警处理机制
  - 告警日志记录

## API 接口说明

### 用户管理接口

#### 1. 用户注册

- **接口**: POST `/regist/users`
- **功能**: 注册新用户
- **请求参数**:
  ```json
  {
    "userName": "string",        // 用户名，必填，4-20字符
    "passWD": "string",         // 密码，必填，最少8字符
    "userID": int,              // 用户唯一标识，必填
    "userType": int,            // 用户类型，必填
    "gatewayDeviceID": "string" // 用户所属网关设备ID，必填,注意：用户注册之前需先进行设备注册，获取到真实设备id之后才可以进行用户注册。
  }
  ```
- **响应示例**:
  ```json
  {
    "code": 201,
    "message": "User created",
    "data": {
      "userName": "string",
      "userID": int,
      "created_at": "2024-03-05T15:00:00Z"
    }
  }
  ```

#### 2. 获取用户列表

- **接口**: GET `/search/users`
- **功能**: 获取所有用户信息
- **响应示例**:
  ```json
  {
    "users": [
      {
        "userName": "string",
        "userID": int,
        "userType": int,
        "gatewayDeviceID": "string",
        // ... 其他用户信息
      }
    ]
  }
  ```

#### 3. 更新用户信息

- **接口**: PUT `/update/users/:id`
- **功能**: 更新指定用户的信息
- **路径参数**: id - 用户ID
- **请求参数**: 与注册接口相同，字段可选
- **响应示例**:
  ```json
  {
    "code": 200,
    "message": "User updated successfully",
    "data": {
      // 更新的字段信息
    }
  }
  ```

### 设备管理接口

#### 1. 设备注册

- **接口**: POST `/regist/devices`
- **功能**: 注册新设备
- **请求参数**:
  ```json
  {
    "deviceName": "string",      // 设备名称，必填，4-50字符
    "deviceType": int,           // 设备类型，必填，1-4分别代表不同类型设备
    "passWD": "string",         // 设备登录口令，必填，最少8字符
    "deviceID": "string",       // 设备唯一标识，必填
    "superiorDeviceID": "string", // 上级设备ID，必填（安全接入管理设备可为空）
    "certID": "string",         // 证书ID，可选
    "keyID": "string"           // 密钥ID，可选
  }
  ```
- **响应示例**:
  ```json
  {
    "code": 201,
    "message": "Device registered",
    "data": {
      "deviceName": "string",
      "deviceID": "string",
      "registered_at": "2024-03-05T15:00:00Z"
    }
  }
  ```

#### 2. 获取设备列表

- **接口**: GET `/search/devices`
- **功能**: 获取所有设备信息
- **响应示例**:
  ```json
  {
    "devices": [
      {
        "deviceName": "string",
        "deviceType": int,
        "deviceID": "string",
        "superiorDeviceID": "string",
        // ... 其他设备信息
      }
    ]
  }
  ```

#### 3. 更新设备信息

- **接口**: PUT `/update/devices/:id`
- **功能**: 更新指定设备的信息
- **路径参数**: id - 设备ID
- **请求参数**: 与注册接口相同，字段可选
- **响应示例**:
  ```json
  {
    "code": 200,
    "message": "Device updated successfully",
    "data": {
      // 更新的字段信息
    }
  }
  ```

### 证书管理接口（新增）

#### 1. 绑定用户证书

- **接口**: POST `/bind/users/:id/cert`
- **功能**: 上传并绑定用户的证书文件
- **路径参数**: id - 用户ID
- **请求格式**: multipart/form-data
- **请求参数**:
  - cert: 证书文件（.pem格式）
- **响应示例**:
  ```json
  {
    "code": 200,
    "message": "证书绑定成功",
    "data": {
      "userID": "1001",
      "certPath": "/path/to/cert/file.pem"
    }
  }
  ```

#### 2. 绑定用户密钥

- **接口**: POST `/bind/users/:id/key`
- **功能**: 上传并绑定用户的密钥文件
- **路径参数**: id - 用户ID
- **请求格式**: multipart/form-data
- **请求参数**:
  - key: 密钥文件（.pem格式）
- **响应示例**:
  ```json
  {
    "code": 200,
    "message": "密钥绑定成功",
    "data": {
      "userID": "1001",
      "keyPath": "/path/to/key/file.pem"
    }
  }
  ```

#### 3. 绑定设备证书

- **接口**: POST `/bind/devices/:id/cert`
- **功能**: 上传并绑定设备的证书文件
- **路径参数**: id - 设备ID
- **请求格式**: multipart/form-data
- **请求参数**:
  - cert: 证书文件（.pem格式）
- **响应示例**:
  ```json
  {
    "code": 200,
    "message": "证书绑定成功",
    "data": {
      "deviceID": "DEV123456",
      "certPath": "/path/to/cert/file.pem"
    }
  }
  ```

#### 4. 绑定设备密钥

- **接口**: POST `/bind/devices/:id/key`
- **功能**: 上传并绑定设备的密钥文件
- **路径参数**: id - 设备ID
- **请求格式**: multipart/form-data
- **请求参数**:
  - key: 密钥文件（.pem格式）
- **响应示例**:
  ```json
  {
    "code": 200,
    "message": "密钥绑定成功",
    "data": {
      "deviceID": "DEV123456",
      "keyPath": "/path/to/key/file.pem"
    }
  }
  ```

#### 5. 获取证书信息

- **接口**: GET `/cert/info`
- **功能**: 获取指定实体（用户或设备）的证书信息
- **查询参数**:
  - type: 实体类型（user或device）
  - id: 实体ID
- **响应示例**:
  ```json
  {
    "code": 200,
    "message": "获取证书信息成功",
    "data": {
      "id": 1,
      "entity_type": "user",
      "entity_id": "1001",
      "cert_path": "/path/to/cert/file.pem",
      "key_path": "/path/to/key/file.pem",
      "upload_time": "2024-03-20T10:15:30Z"
    }
  }
  ```

### 认证管理接口

#### 1. 获取认证记录

- **接口**: GET `/auth/records`
- **功能**: 查询认证记录
- **查询参数**:
  - username: 用户名（可选）
  - reply: 认证结果（可选）
  - start_date: 开始日期（可选）
  - end_date: 结束日期（可选）
  - page: 页码，默认1
  - page_size: 每页记录数，默认10
  - class: 认证类别（可选）
- **响应示例**:
  ```json
  {
    "code": 200,
    "message": "Success",
    "data": {
      "total": int,
      "total_pages": int,
      "page": int,
      "page_size": int,
      "records": [
        {
          "id": int,
          "username": "string",
          "reply": "string",
          "authdate": "string",
          "class": "string"
        }
      ]
    }
  }
  ```

## 错误码说明

- 200: 成功
- 201: 创建成功
- 400: 请求参数错误
- 401: 未授权
- 403: 禁止访问
- 404: 资源不存在
- 409: 资源冲突
- 500: 服务器内部错误

## 系统要求

- Go 1.23 或更高版本
- MySQL/MariaDB 数据库
- Windows/Linux 操作系统

## 快速开始

1. 克隆项目

```bash
git clone https://github.com/yourusername/gin-server.git
cd gin-server
```

2. 安装依赖

```bash
go mod download
```

3. 配置环境变量

```bash
# 服务器配置
export SERVER_PORT=8080
export DEBUG_LEVEL=true

# 数据库配置
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=your_user
export DB_PASSWORD=your_password
export DB_NAME=your_database

# Radius数据库配置
export RADIUS_DB_HOST=localhost
export RADIUS_DB_PORT=3306
export RADIUS_DB_USER=radius_user
export RADIUS_DB_PASSWORD=radius_password
export RADIUS_DB_NAME=radius

# 存储配置
export STORAGE_TYPE=gitee  # 或 ftp

# Gitee配置（如果使用Gitee存储）
export GITEE_ACCESS_TOKEN=your_token
export GITEE_OWNER=your_username
export GITEE_REPO=your_repo
export GITEE_BRANCH=master

# FTP配置（如果使用FTP存储）
export FTP_HOST=ftp.example.com
export FTP_PORT=21
export FTP_USERNAME=your_username
export FTP_PASSWORD=your_password
```

4. 运行服务

```bash
go run main.go
```

## 项目结构

```
gin-server/
├── auth/           # 认证模块
│   ├── handler/    # 请求处理器
│   ├── model/      # 数据模型
│   └── router/     # 路由配置
├── config/         # 配置管理
├── configmanager/  # 配置管理模块
│   ├── common/     # 公共组件
│   └── ...
├── regist/        # 注册模块
│   ├── handler/   # 请求处理器
│   ├── model/     # 数据模型
│   └── router/    # 路由配置
├── test/          # 测试工具
│   ├── cert_test.go  # 证书绑定测试工具
│   ├── curl_test.sh  # curl命令测试脚本
│   └── README.md     # 测试工具说明
└── main.go        # 主程序入口
```

## 配置说明

### 路径配置注意事项

1. 路径分隔符

   - 配置文件中统一使用正斜杠(/)作为路径分隔符
   - 程序会自动处理Windows和Linux系统的路径差异
2. 远程仓库目录命名

   - 时间格式：YYYYMMDDHHmmss
   - 示例：20240305150000

### 加密配置

1. 密钥长度选择

   - AES: 推荐使用256位
   - RSA: 推荐使用2048位或4096位
   - ECDSA: 推荐使用256位或384位
2. 密钥文件存放

   - 建议将密钥文件放在单独的安全目录
   - 确保密钥文件具有适当的访问权限

## 证书管理说明（新增）

### 证书存储结构

证书和密钥文件按照以下结构存储：

```
regist/certs/
├── certs/             # 存放证书文件
│   ├── user_1001.pem  # 用户证书文件(命名格式: user_<userID>.pem)
│   └── device_DEV123456.pem  # 设备证书文件(命名格式: device_<deviceID>.pem)
└── keys/              # 存放密钥文件
    ├── user_1001.pem  # 用户密钥文件
    └── device_DEV123456.pem  # 设备密钥文件
```

### 安全考虑

- 证书目录：权限设置为0755
- 密钥目录：权限设置为0700（更严格的权限控制）
- 证书文件：权限默认为0644
- 密钥文件：权限设置为0600（只有所有者可读写）

### 数据库结构

证书信息存储在`certs`表中，包含以下字段：

| 字段名      | 类型      | 描述                | 约束                |
|-------------|-----------|-------------------|---------------------|
| id          | INT       | 自增主键           | AUTO_INCREMENT PRIMARY KEY |
| entity_type | VARCHAR(10) | 实体类型(user/device) | NOT NULL           |
| entity_id   | VARCHAR(20) | 实体ID           | NOT NULL           |
| cert_path   | VARCHAR(255) | 证书文件路径      | NULL               |
| key_path    | VARCHAR(255) | 密钥文件路径      | NULL               |
| upload_time | DATETIME  | 上传时间           | NOT NULL           |

索引：
- 联合索引：`(entity_type, entity_id)`

### 测试工具

为测试证书功能，我们提供了两种测试工具：

1. Go测试工具 (`test/cert_test.go`)
   - 完整功能测试，支持创建测试用户和设备
   - 自动测试所有证书绑定API
   - 验证数据库记录

   在Windows系统中使用方法：
   ```powershell
   cd D:\code\golang\git\gin-server\test
   go build -o cert_test.exe cert_test.go
   .\cert_test.exe -create-device -create-user -test-all
   ```

2. Shell脚本工具 (`test/curl_test.sh`)
   - 使用curl命令测试API
   - 生成测试证书和密钥
   - 适用于Linux/Mac或Windows的Git Bash环境

   使用方法：
   ```bash
   chmod +x test/curl_test.sh
   ./test/curl_test.sh
   ```

详细的测试工具说明请参考 `test/README.md` 文件。

## 证书管理接口请求示例

### 1. 绑定用户证书

**接口**: `POST /bind/users/:id/cert`

#### cURL 请求示例

```bash
curl -X POST http://localhost:8080/bind/users/1001/cert \
  -F "cert=@/path/to/user_certificate.pem"
```

#### 请求说明

- 使用multipart/form-data格式上传文件
- `:id` 替换为实际的用户ID (例如：1001)
- 表单字段名必须为 `cert`
- 文件必须是PEM格式的证书文件
- 文件大小不得超过8MB

#### 响应示例 (成功)

```json
{
  "code": 200,
  "message": "证书绑定成功",
  "data": {
    "userID": "1001",
    "certPath": "/d/code/golang/git/gin-server/regist/certs/certs/user_1001.pem"
  }
}
```

#### 响应示例 (失败)

```json
{
  "error": "用户不存在"
}
```

### 2. 绑定用户密钥

**接口**: `POST /bind/users/:id/key`

#### cURL 请求示例

```bash
curl -X POST http://localhost:8080/bind/users/1001/key \
  -F "key=@/path/to/user_private_key.pem"
```

#### 请求说明

- 使用multipart/form-data格式上传文件
- `:id` 替换为实际的用户ID (例如：1001)
- 表单字段名必须为 `key`
- 文件必须是PEM格式的私钥文件
- 文件大小不得超过8MB

#### 响应示例 (成功)

```json
{
  "code": 200,
  "message": "密钥绑定成功",
  "data": {
    "userID": "1001",
    "keyPath": "/d/code/golang/git/gin-server/regist/certs/keys/user_1001.pem"
  }
}
```

#### 响应示例 (失败)

```json
{
  "error": "无法保存密钥文件"
}
```

### 3. 绑定设备证书

**接口**: `POST /bind/devices/:id/cert`

#### cURL 请求示例

```bash
curl -X POST http://localhost:8080/bind/devices/DEV123456/cert \
  -F "cert=@/path/to/device_certificate.pem"
```

#### 请求说明

- 使用multipart/form-data格式上传文件
- `:id` 替换为实际的设备ID (例如：DEV123456)
- 表单字段名必须为 `cert`
- 文件必须是PEM格式的证书文件
- 文件大小不得超过8MB

#### 响应示例 (成功)

```json
{
  "code": 200,
  "message": "证书绑定成功",
  "data": {
    "deviceID": "DEV123456",
    "certPath": "/d/code/golang/git/gin-server/regist/certs/certs/device_DEV123456.pem"
  }
}
```

#### 响应示例 (失败)

```json
{
  "error": "设备不存在"
}
```

### 4. 绑定设备密钥

**接口**: `POST /bind/devices/:id/key`

#### cURL 请求示例

```bash
curl -X POST http://localhost:8080/bind/devices/DEV123456/key \
  -F "key=@/path/to/device_private_key.pem"
```

#### 请求说明

- 使用multipart/form-data格式上传文件
- `:id` 替换为实际的设备ID (例如：DEV123456)
- 表单字段名必须为 `key`
- 文件必须是PEM格式的私钥文件
- 文件大小不得超过8MB

#### 响应示例 (成功)

```json
{
  "code": 200,
  "message": "密钥绑定成功",
  "data": {
    "deviceID": "DEV123456",
    "keyPath": "/d/code/golang/git/gin-server/regist/certs/keys/device_DEV123456.pem"
  }
}
```

#### 响应示例 (失败)

```json
{
  "error": "文件大小超过限制"
}
```

### 5. 获取证书信息

**接口**: `GET /cert/info`

#### cURL 请求示例 (获取用户证书信息)

```bash
curl -X GET "http://localhost:8080/cert/info?type=user&id=1001"
```

#### cURL 请求示例 (获取设备证书信息)

```bash
curl -X GET "http://localhost:8080/cert/info?type=device&id=DEV123456"
```

#### 请求说明

- 使用查询参数传递实体类型和ID
- `type` 参数必须为 `user` 或 `device`
- `id` 参数为实体的唯一标识

#### 响应示例 (成功)

```json
{
  "code": 200,
  "message": "获取证书信息成功",
  "data": {
    "id": 1,
    "entity_type": "user",
    "entity_id": "1001",
    "cert_path": "/d/code/golang/git/gin-server/regist/certs/certs/user_1001.pem",
    "key_path": "/d/code/golang/git/gin-server/regist/certs/keys/user_1001.pem",
    "upload_time": "2023-08-15T14:30:15Z"
  }
}
```

#### 响应示例 (未找到记录)

```json
{
  "code": 200,
  "message": "未找到证书记录",
  "data": null
}
```

#### 响应示例 (参数错误)

```json
{
  "error": "缺少必要的参数"
}
```

### 请求示例中的参数说明

#### 1. 实体ID格式

- 用户ID：整数类型，例如 1001, 1002 等
- 设备ID：字符串类型，例如 DEV123456, DEV789012 等

#### 2. 证书和密钥文件要求

- 格式：必须是PEM编码的X.509证书和私钥
- 大小：最大8MB
- 证书示例：
  ```
  -----BEGIN CERTIFICATE-----
  MIIDazCCAlOgAwIBAgIUECPZA...（中间内容省略）...BtNJ9AQKBgQD
  -----END CERTIFICATE-----
  ```
- 密钥示例：
  ```
  -----BEGIN PRIVATE KEY-----
  MIIEvQIBADANBgkqhkiG9w0BA...（中间内容省略）...QA5BICbW1gtrIByvXbpPuPE=
  -----END PRIVATE KEY-----
  ```

#### 3. 常见错误处理

| 错误情况 | HTTP状态码 | 错误信息 |
|---------|-----------|---------|
| 缺少文件 | 400 | "未找到证书文件" 或 "未找到密钥文件" |
| 实体不存在 | 404 | "用户不存在" 或 "设备不存在" |
| 文件过大 | 413 | "文件大小超过限制" |
| 文件格式错误 | 400 | "无效的证书格式" 或 "无效的密钥格式" |
| 服务器错误 | 500 | "保存文件失败" 或 "数据库操作失败" |

### 批量测试示例

可以使用提供的Shell脚本进行批量测试：

```bash
# 在Windows系统使用Git Bash
cd /d/code/golang/git/gin-server
./test/curl_test.sh

# 或在Linux/Mac系统
chmod +x test/curl_test.sh
./test/curl_test.sh
```

或使用Go测试工具：

```bash
# 在Windows系统使用PowerShell
cd D:\code\golang\git\gin-server\test
go build -o cert_test.exe cert_test.go
.\cert_test.exe -create-device -create-user -test-all

# 或在Linux/Mac系统
cd /path/to/gin-server/test
go build -o cert_test cert_test.go
./cert_test -create-device -create-user -test-all
```

## 注意事项

1. 系统兼容性

   - 代码已做跨平台兼容处理
   - Windows和Linux系统下的路径会自动转换
   - 文件操作使用统一的接口
2. 数据库配置

   - 确保数据库字符集为UTF-8
   - 建议使用独立的数据库用户
   - 定期备份数据库
3. 安全性

   - 及时更新依赖包
   - 定期更换密钥
   - 不要在代码中硬编码敏感信息
   - 使用环境变量或配置文件管理敏感配置
4. 性能优化

   - 合理设置轮询间隔
   - 监控文件数量不宜过多
   - 注意日志文件大小控制
5. 调试模式

   - 生产环境建议关闭调试模式
   - 调试日志可能包含敏感信息
6. 证书和密钥安全（新增）

   - 证书和密钥应当妥善保管
   - 密钥文件权限设置为限制性权限(0600)
   - 定期更新证书和密钥
   - 对证书的有效性进行验证

## 常见问题

1. 路径问题

   - Q: Windows下路径出现异常？
   - A: 检查是否使用了正斜杠(/)，程序会自动处理转换
2. 权限问题

   - Q: 无法创建或访问文件？
   - A: 检查程序运行用户的权限和文件系统权限
3. 数据库连接

   - Q: 数据库连接失败？
   - A: 检查数据库配置和网络连接
4. 证书相关问题（新增）

   - Q: 上传证书失败？
   - A: 确保证书格式为PEM，大小不超过8MB
   - Q: 证书目录不存在？
   - A: 系统会在启动时自动创建证书目录结构
   - Q: Windows下运行测试工具失败？
   - A: 请使用`.\cert_test.exe`而不是`./cert_test`命令

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

[MIT License](LICENSE)

## 告警模块说明

### 告警级别

系统支持四个告警级别：

- INFO：信息级别，用于记录普通操作信息
- WARNING：警告级别，用于提示潜在问题
- ERROR：错误级别，用于记录严重问题
- FATAL：致命级别，用于记录可能导致系统崩溃的问题

### 告警类型

系统支持以下告警类型：

- LOG_GENERATE：日志生成相关告警
- LOG_ENCRYPT：日志加密相关告警
- LOG_UPLOAD：日志上传相关告警
- STRATEGY_SYNC：策略同步相关告警
- STRATEGY_APPLY：策略应用相关告警
- CERT_BINDING：证书绑定相关告警（新增）

### 告警信息结构

每条告警包含以下信息：

```json
{
  "level": "string",      // 告警级别
  "type": "string",       // 告警类型
  "message": "string",    // 告警消息
  "error": "string",      // 错误信息（可选）
  "retryCount": int,      // 重试次数
  "timestamp": "string",  // 时间戳
  "module": "string"      // 告警模块
}
```

### 使用示例

```go
// 创建告警器
alerter := alert.NewLogAlerter()

// 创建告警信息
alertInfo := &alert.Alert{
    Level:      alert.AlertLevelInfo,
    Type:       alert.AlertTypeLogGenerate,
    Message:    "日志生成成功",
    Error:      nil,
    RetryCount: 0,
    Timestamp:  time.Now(),
    Module:     "LogModule",
}

// 发送告警
err := alerter.Alert(alertInfo)
```

## 日志生成模块说明

### 日志生成周期

系统默认每5分钟生成一次日志文件，记录最近5分钟内的系统运行状态，包括：

- 设备状态（CPU使用率、内存使用率、在线时长等）
- 用户行为（文件传输、消息发送等）
- 安全事件
- 故障事件

### 日志文件结构

日志文件（log.json）包含以下主要字段：

```json
{
  "timeRange": {
    "startTime": "2024-03-07T10:00:00Z",  // 统计起始时间
    "duration": 300                        // 统计时长（秒）
  },
  "securityEvents": {
    "events": [
      {
        "eventId": 1001,                   // 事件ID
        "deviceId": "SEC00000001",         // 设备ID
        "eventTime": "2024-03-07T10:02:00Z", // 事件发生时间
        "eventType": 1,                    // 事件类型（1:安全事件）
        "eventCode": "SEC_001",            // 事件代码
        "eventDesc": "异常登录尝试",         // 事件描述
        "createdAt": "2024-03-07T10:02:01Z" // 记录创建时间
      }
    ]
  },
  "performanceEvents": {
    "securityDevices": [
      {
        "deviceId": "SEC00000001",         // 安全接入管理设备ID
        "cpuUsage": 45,                    // CPU使用率峰值(%)
        "memoryUsage": 60,                 // 内存使用率峰值(%)
        "onlineDuration": 3600,            // 在线时长(秒)
        "status": 1,                       // 设备状态(1:在线,2:离线)
        "gatewayDevices": [
          {
            "deviceId": "GWA00000001",     // 网关设备ID
            "cpuUsage": 30,                // CPU使用率峰值(%)
            "memoryUsage": 40,             // 内存使用率峰值(%)
            "onlineDuration": 3600,        // 在线时长(秒)
            "status": 1,                   // 设备状态
            "users": [
              {
                "userId": 10001,           // 用户ID
                "status": 1,               // 用户状态(1:在线,2:离线)
                "onlineDuration": 1800,    // 在线时长(秒)
                "behaviors": [
                  {
                    "time": "2024-03-07T10:01:00Z", // 行为发生时间
                    "type": 1,             // 行为类型(1:发送,2:接收)
                    "dataType": 1,         // 数据类型(1:文件,2:消息)
                    "dataSize": 1024       // 数据大小(字节)
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  },
  "faultEvents": {
    "events": [
      {
        "eventId": 2001,                   // 事件ID
        "deviceId": "GWA00000001",         // 设备ID
        "eventTime": "2024-03-07T10:03:00Z", // 事件发生时间
        "eventType": 2,                    // 事件类型（2:故障事件）
        "eventCode": "FAULT_001",          // 事件代码
        "eventDesc": "设备离线",            // 事件描述
        "createdAt": "2024-03-07T10:03:01Z" // 记录创建时间
      }
    ]
  }
}
```

## 数据库结构说明

项目使用了两个主要的数据库：

### 1. 主数据库 (gin_server)

主数据库包含以下表：

#### 1.1 设备表 (devices)

| 字段名                    | 类型         | 描述                                                                      | 约束                         |
| ------------------------- | ------------ | ------------------------------------------------------------------------- | ---------------------------- |
| id                        | INT          | 自增主键                                                                  | AUTO_INCREMENT PRIMARY KEY   |
| deviceName                | VARCHAR(50)  | 设备名称                                                                  | NOT NULL                     |
| deviceType                | INT          | 设备类型，1:网关设备A型，2:网关设备B型，3:网关设备C型，4:安全接入管理设备 | NOT NULL                     |
| passWD                    | VARCHAR(255) | 设备登录口令                                                              | NOT NULL                     |
| deviceID                  | CHAR(12)     | 设备唯一标识                                                              | NOT NULL                     |
| superiorDeviceID          | CHAR(12)     | 上级设备ID                                                                | NOT NULL                     |
| deviceStatus              | INT          | 设备状态，1:在线，2:离线，3:冻结，4:注销                                  | DEFAULT 2                    |
| peakCPUUsage              | INT          | 峰值CPU使用率                                                             | DEFAULT 0                    |
| peakMemoryUsage           | INT          | 峰值内存使用率                                                            | DEFAULT 0                    |
| onlineDuration            | INT          | 在线时长                                                                  | DEFAULT 0                    |
| certID                    | VARCHAR(64)  | 证书ID                                                                    | NULL                         |
| keyID                     | VARCHAR(64)  | 密钥ID                                                                    | NULL                         |
| registerIP                | VARCHAR(24)  | 上级设备IP                                                                | NULL                         |
| email                     | VARCHAR(32)  | 联系邮箱                                                                  | NULL                         |
| deviceHardwareFingerprint | CHAR(128)    | 设备硬件指纹                                                              | NULL                         |
| anonymousUser             | VARCHAR(50)  | 匿名用户                                                                  | NULL                         |
| created_at                | TIMESTAMP(3) | 创建时间                                                                  | DEFAULT CURRENT_TIMESTAMP(3) |

索引：

- `idx_deviceid` (deviceID)
- `idx_devicename` (deviceName)
- `idx_superiordeviceid` (superiorDeviceID)

#### 1.2 用户表 (users)

| 字段名             | 类型         | 描述                                     | 约束                         |
| ------------------ | ------------ | ---------------------------------------- | ---------------------------- |
| id                 | INT          | 自增主键                                 | AUTO_INCREMENT PRIMARY KEY   |
| userName           | VARCHAR(20)  | 用户名                                   | NOT NULL                     |
| passWD             | VARCHAR(255) | 密码                                     | NOT NULL                     |
| userID             | INT          | 用户唯一标识                             | NOT NULL                     |
| userType           | INT          | 用户类型                                 | NOT NULL                     |
| gatewayDeviceID    | VARCHAR(12)  | 用户所属网关设备ID                       | NOT NULL，FOREIGN KEY        |
| status             | INT          | 用户状态，1:在线，2:离线，3:冻结，4:注销 | NULL                         |
| onlineDuration     | INT          | 在线时长                                 | NULL, DEFAULT 0              |
| certID             | VARCHAR(64)  | 证书ID                                   | NULL                         |
| keyID              | VARCHAR(64)  | 密钥ID                                   | NULL                         |
| email              | VARCHAR(32)  | 邮箱                                     | NULL                         |
| permissionMask     | CHAR(8)      | 权限位掩码                               | NULL                         |
| lastLoginTimeStamp | DATETIME(3)  | 登录时间戳                               | NULL                         |
| offLineTimeStamp   | DATETIME(3)  | 离线时间戳                               | NULL                         |
| loginIP            | CHAR(24)     | 用户登录IP                               | NULL                         |
| illegalLoginTimes  | INT          | 用户本次的非法登录次数                   | NULL                         |
| created_at         | TIMESTAMP(3) | 创建时间                                 | DEFAULT CURRENT_TIMESTAMP(3) |

索引：

- `idx_userid` (userID)
- `idx_username` (userName)
- `idx_email` (email)
- `idx_gatewaydeviceid` (gatewayDeviceID)

外键关系：

- `gatewayDeviceID` 关联 `devices(deviceID)` 表

#### 1.3 用户行为表 (user_behaviors)

| 字段名       | 类型         | 描述                     | 约束                         |
| ------------ | ------------ | ------------------------ | ---------------------------- |
| behaviorID   | INT          | 行为ID                   | AUTO_INCREMENT PRIMARY KEY   |
| userID       | INT          | 用户ID                   | NOT NULL，FOREIGN KEY        |
| behaviorTime | DATETIME(3)  | 行为开始时间             | NOT NULL                     |
| behaviorType | INT          | 行为类型，1:发送，2:接收 | NOT NULL                     |
| dataType     | INT          | 数据类型，1:文件，2:消息 | NOT NULL                     |
| dataSize     | BIGINT       | 数据大小                 | NOT NULL                     |
| created_at   | TIMESTAMP(3) | 创建时间                 | DEFAULT CURRENT_TIMESTAMP(3) |

索引：

- `idx_userid` (userID)
- `idx_behaviortime` (behaviorTime)
- `idx_behaviortype` (behaviorType)

外键关系：

- `userID` 关联 `users(userID)` 表

#### 1.4 事件表 (events)

| 字段名    | 类型         | 描述                             | 约束                      |
| --------- | ------------ | -------------------------------- | ------------------------- |
| eventId   | BIGINT       | 事件ID                           | PRIMARY KEY               |
| deviceId  | VARCHAR(12)  | 设备ID                           | NOT NULL                  |
| eventTime | DATETIME     | 事件发生时间                     | NOT NULL                  |
| eventType | INT          | 事件类型，1:安全事件，2:故障事件 | NOT NULL                  |
| eventCode | VARCHAR(20)  | 事件代码                         | NOT NULL                  |
| eventDesc | VARCHAR(255) | 事件描述                         | NOT NULL                  |
| createdAt | TIMESTAMP    | 创建时间                         | DEFAULT CURRENT_TIMESTAMP |

#### 1.5 证书表 (certs) （新增）

| 字段名      | 类型        | 描述                | 约束                      |
| ----------- | ----------- | ------------------- | ------------------------- |
| id          | INT         | 自增主键            | AUTO_INCREMENT PRIMARY KEY |
| entity_type | VARCHAR(10) | 实体类型(user/device) | NOT NULL                |
| entity_id   | VARCHAR(20) | 实体ID              | NOT NULL                |
| cert_path   | VARCHAR(255) | 证书文件路径        | NULL                    |
| key_path    | VARCHAR(255) | 密钥文件路径        | NULL                    |
| upload_time | DATETIME    | 上传时间            | NOT NULL                |

索引：

- `idx_entity` (entity_type, entity_id)

### 2. Radius认证数据库 (radius)

#### 2.1 认证记录表 (radpostauth)

| 字段名   | 类型         | 描述     | 约束                         |
| -------- | ------------ | -------- | ---------------------------- |
| id       | INT          | 自增主键 | AUTO_INCREMENT PRIMARY KEY   |
| username | VARCHAR(64)  | 用户名   | NOT NULL                     |
| pass     | VARCHAR(64)  | 密码     | NOT NULL                     |
| reply    | VARCHAR(32)  | 认证响应 | NOT NULL                     |
| authdate | TIMESTAMP(6) | 认证时间 | DEFAULT CURRENT_TIMESTAMP(6) |
| class    | VARCHAR(64)  | 认证类型 | NULL                         |

索引：

- `idx_username` (username)
- `idx_class` (class)
