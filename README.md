# Gin Server

基于Gin框架的服务器应用，提供用户认证、设备管理和配置管理功能。

## 功能特性

- 用户管理
  - 用户注册
  - 用户信息查询
  - 用户信息更新
- 设备管理
  - 设备注册
  - 设备信息查询
  - 设备信息更新
- 认证管理
  - 认证记录查询
  - 权限控制
- 配置管理
  - 远程配置同步
  - 文件变更监控
  - 加密传输
  - 多存储方式支持（Gitee/FTP）

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
    "gatewayDeviceID": "string", // 用户所属网关设备ID，必填,注意：用户注册之前需先进行设备注册，获取到真实设备id之后才可以进行用户注册。
    "certID": "string",         // 证书ID，可选
    "keyID": "string"           // 密钥ID，可选
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

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

[MIT License](LICENSE)
