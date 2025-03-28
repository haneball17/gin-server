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
- 日志管理
  - 日志生成、加密和上传
  - 日志查询接口

## API 接口说明

### 用户管理接口

#### 1. 用户注册

- **接口**: `POST /regist/users`
- **功能**: 注册新用户
- **请求格式**: JSON
- **请求参数**:
  ```json
  {
    "user_name": "string",        // 用户名，必填，4-20字符
    "pass_wd": "string",          // 密码，必填，最少8字符
    "user_id": int,               // 用户唯一标识，必填
    "user_type": int,             // 用户类型，必填
    "gateway_device_id": int      // 用户所属网关设备ID，必填，注意：用户注册之前需先进行设备注册，获取到真实设备id之后才可以进行用户注册
  }
  ```
- **响应格式**: JSON
- **响应示例 (成功)**:
  ```json
  {
    "code": 200,
    "message": "用户注册成功",
    "data": {
      "ID": 1,
      "CreatedAt": "2025-03-22T23:26:54.346+08:00",
      "UpdatedAt": "2025-03-22T23:26:54.346+08:00",
      "DeletedAt": null,
      "username": "testuser01",
      "user_id": 10001,
      "user_type": 1,
      "gateway_device_id": 2001,
      "status": null,
      "online_duration": 0,
      "cert_id": "",
      "key_id": "",
      "email": ""
    }
  }
  ```
- **响应示例 (失败)**:
  ```json
  {
    "error": "用户ID已存在"
  }
  ```

#### 2. 获取用户列表

- **接口**: `GET /search/users`
- **功能**: 获取所有用户信息
- **请求格式**: 无参数
- **响应格式**: JSON
- **响应示例**:
  ```json
  {
    "users": [
      {
        "ID": 1,
        "CreatedAt": "2025-03-22T23:26:54.346+08:00",
        "UpdatedAt": "2025-03-22T23:26:54.346+08:00",
        "DeletedAt": null,
        "username": "user_1001_1",
        "user_id": 10000,
        "user_type": 2,
        "gateway_device_id": 1001,
        "status": null,
        "online_duration": 768,
        "cert_id": "",
        "key_id": "",
        "email": "GjWRSP1R@example.com",
        "permission_mask": "11100010",
        "last_login_timestamp": null,
        "offline_timestamp": null,
        "login_ip": "212.193.3.138",
        "illegal_login_times": null
      }
    ]
  }
  ```
- **响应字段说明**:

| 字段名               | 类型     | 描述                                       |
| -------------------- | -------- | ------------------------------------------ |
| ID                   | INT      | 数据库自增主键                             |
| CreatedAt            | DATETIME | 记录创建时间                               |
| UpdatedAt            | DATETIME | 记录最后更新时间                           |
| DeletedAt            | DATETIME | 记录删除时间（null表示未删除）             |
| username             | STRING   | 用户名                                     |
| user_id              | INT      | 用户唯一标识                               |
| user_type            | INT      | 用户类型                                   |
| gateway_device_id    | INT      | 用户所属网关设备ID                         |
| status               | INT      | 用户状态（1:在线，2:离线，3:冻结，4:注销） |
| online_duration      | INT      | 用户在线时长（秒）                         |
| cert_id              | STRING   | 证书ID                                     |
| key_id               | STRING   | 密钥ID                                     |
| email                | STRING   | 邮箱地址                                   |
| permission_mask      | STRING   | 权限位掩码                                 |
| last_login_timestamp | DATETIME | 最后登录时间                               |
| offline_timestamp    | DATETIME | 最后离线时间                               |
| login_ip             | STRING   | 用户登录IP                                 |
| illegal_login_times  | INT      | 非法登录尝试次数                           |

#### 3. 指定用户查找

- **接口**: `GET /search/user`
- **功能**: 根据ID查询指定用户信息
- **请求格式**: URL查询参数
- **请求参数**:
  - id: 用户ID，必填，整数类型
- **请求示例**: `http://localhost:8080/search/user?id=10001`
- **响应格式**: JSON
- **响应示例 (成功)**:
  ```json
  {
    "code": 200,
    "message": "用户查询成功",
    "data": {
      "ID": 1,
      "Username": "用户名1",
      "UserID": 10001,
      "UserType": 1,
      "GatewayDeviceID": 1001,
      "Status": 1,
      "OnlineDuration": 3600,
      "CertID": "cert_id_1",
      "KeyID": "key_id_1",
      "Email": "user1@example.com",
      "CreatedAt": "2024-03-05T15:00:00Z"
    }
  }
  ```
- **响应示例 (失败)**:
  ```json
  {
    "error": "用户不存在"
  }
  ```

#### 4. 更新用户信息

- **接口**: `PUT /update/users/:id`
- **功能**: 更新指定用户的信息
- **路径参数**: id - 用户ID
- **请求格式**: JSON
- **请求参数**: 与注册接口相同，字段可选
- **请求示例**:
  ```json
  {
    "userName": "新用户名",
    "passWD": "新密码",
    "userType": 2,
    "email": "newemail@example.com"
  }
  ```
- **响应格式**: JSON
- **响应示例 (成功)**:
  ```json
  {
    "message": "用户信息更新成功"
  }
  ```
- **响应示例 (失败)**:
  ```json
  {
    "error": "用户不存在"
  }
  ```

### 设备管理接口

#### 1. 设备注册

- **接口**: `POST /regist/devices`
- **功能**: 注册新设备
- **请求格式**: JSON
- **请求参数**:
  ```json
  {
    "device_name": "string",          // 设备名称，必填，4-50字符
    "device_type": int,               // 设备类型，必填，1-4分别代表不同类型设备
    "pass_wd": "string",              // 设备登录口令，必填，最少8字符
    "device_id": int,                 // 设备唯一标识，必填
    "superior_device_id": int         // 上级设备ID，必填（安全接入管理设备为0）
  }
  ```
- **响应格式**: JSON
- **响应示例 (成功)**:
  ```json
  {
    "code": 200,
    "message": "设备注册成功",
    "data": {
      "ID": 1,
      "CreatedAt": "2025-03-22T23:26:54.327+08:00",
      "UpdatedAt": "2025-03-22T23:26:54.327+08:00",
      "DeletedAt": null,
      "device_name": "TestDevice01",
      "device_type": 1,
      "password": "password123",
      "device_id": 2001,
      "superior_device_id": 0,
      "device_status": 2,
      "peak_cpu_usage": 0,
      "peak_memory_usage": 0,
      "online_duration": 0,
      "cert_id": "",
      "key_id": "",
      "register_ip": "127.0.0.1",
      "email": "",
      "hardware_fingerprint": "",
      "anonymous_user": ""
    }
  }
  ```
- **响应示例 (失败)**:
  ```json
  {
    "error": "设备ID已存在"
  }
  ```

#### 2. 获取设备列表

- **接口**: `GET /search/devices`
- **功能**: 获取所有设备信息
- **请求格式**: 无参数
- **响应格式**: JSON
- **响应示例**:
  ```json
  {
    "code": 200,
    "message": "获取设备列表成功",
    "data": [
      {
        "ID": 1,
        "CreatedAt": "2025-03-22T23:26:54.327+08:00",
        "UpdatedAt": "2025-03-22T23:26:54.327+08:00",
        "DeletedAt": null,
        "device_name": "安全接入管理设备",
        "device_type": 4,
        "password": "admin123456",
        "device_id": 1000,
        "superior_device_id": 0,
        "device_status": 1,
        "peak_cpu_usage": 0,
        "peak_memory_usage": 0,
        "online_duration": 0,
        "cert_id": "",
        "key_id": "",
        "register_ip": "220.42.76.214",
        "email": "VCSNayol@company.net",
        "hardware_fingerprint": "",
        "anonymous_user": ""
      }
    ]
  }
  ```
- **响应字段说明**:

| 字段名               | 类型     | 描述                                                                        |
| -------------------- | -------- | --------------------------------------------------------------------------- |
| ID                   | INT      | 数据库自增主键                                                              |
| CreatedAt            | DATETIME | 记录创建时间                                                                |
| UpdatedAt            | DATETIME | 记录最后更新时间                                                            |
| DeletedAt            | DATETIME | 记录删除时间（null表示未删除）                                              |
| device_name          | STRING   | 设备名称                                                                    |
| device_type          | INT      | 设备类型（1:网关设备A型，2:网关设备B型，3:网关设备C型，4:安全接入管理设备） |
| password             | STRING   | 设备登录口令                                                                |
| device_id            | INT      | 设备唯一标识                                                                |
| superior_device_id   | INT      | 上级设备ID（安全接入管理设备为0）                                           |
| device_status        | INT      | 设备状态（1:在线，2:离线，3:冻结，4:注销）                                  |
| peak_cpu_usage       | INT      | 峰值CPU使用率（百分比）                                                     |
| peak_memory_usage    | INT      | 峰值内存使用率（百分比）                                                    |
| online_duration      | INT      | 设备在线时长（秒）                                                          |
| cert_id              | STRING   | 证书ID                                                                      |
| key_id               | STRING   | 密钥ID                                                                      |
| register_ip          | STRING   | 设备注册IP                                                                  |
| email                | STRING   | 联系邮箱                                                                    |
| hardware_fingerprint | STRING   | 设备硬件指纹                                                                |
| anonymous_user       | STRING   | 匿名用户                                                                    |

#### 3. 指定设备查找

- **接口**: `GET /search/device`
- **功能**: 根据ID查询指定设备信息
- **请求格式**: URL查询参数
- **请求参数**:
  - id: 设备ID，必填，整数类型
- **请求示例**: `http://localhost:8080/search/device?id=1001`
- **响应格式**: JSON
- **响应示例 (成功)**:
  ```json
  {
    "code": 200,
    "message": "设备查询成功",
    "data": {
      "ID": 1,
      "DeviceName": "设备1",
      "DeviceType": 1,
      "DeviceID": 1001,
      "SuperiorDeviceID": 0,
      "DeviceStatus": 1,
      "CertID": "cert_id_1",
      "KeyID": "key_id_1",
      "RegisterIP": "192.168.1.100",
      "Email": "device1@example.com",
      "CreatedAt": "2024-03-05T15:00:00Z"
    }
  }
  ```
- **响应示例 (失败)**:
  ```json
  {
    "error": "设备不存在"
  }
  ```

#### 4. 更新设备信息

- **接口**: `PUT /update/devices/:id`
- **功能**: 更新指定设备的信息
- **路径参数**: id - 设备ID
- **请求格式**: JSON
- **请求参数**: 与注册接口相同，字段可选
- **请求示例**:
  ```json
  {
    "deviceName": "新设备名",
    "passWD": "新密码",
    "deviceStatus": 2,
    "email": "newemail@example.com"
  }
  ```
- **响应格式**: JSON
- **响应示例 (成功)**:
  ```json
  {
    "message": "设备信息更新成功"
  }
  ```
- **响应示例 (失败)**:
  ```json
  {
    "error": "设备不存在"
  }
  ```

### 证书管理接口

#### 1. 绑定用户证书

- **接口**: `POST /bind/users/:id/cert`
- **功能**: 上传并绑定用户的证书文件
- **路径参数**: id - 用户ID
- **请求格式**: multipart/form-data
- **请求参数**:

  - cert: 证书文件（.pem格式）
- **请求示例 (curl)**:

  ```bash
  curl -X POST http://localhost:8080/bind/users/1001/cert \
    -F "cert=@/path/to/user_certificate.pem"
  ```
- **响应格式**: JSON
- **响应示例 (成功)**:

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
- **响应示例 (失败)**:

  ```json
  {
    "error": "用户不存在"
  }
  ```

#### 2. 绑定用户密钥

- **接口**: `POST /bind/users/:id/key`
- **功能**: 上传并绑定用户的密钥文件
- **路径参数**: id - 用户ID
- **请求格式**: multipart/form-data
- **请求参数**:
  - key: 密钥文件（.pem格式）
- **请求示例 (curl)**:
  ```bash
  curl -X POST http://localhost:8080/bind/users/1001/key \
    -F "key=@/path/to/user_private_key.pem"
  ```
- **响应格式**: JSON
- **响应示例 (成功)**:
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
- **响应示例 (失败)**:
  ```json
  {
    "error": "无法保存密钥文件"
  }
  ```

#### 3. 绑定设备证书

- **接口**: `POST /bind/devices/:id/cert`
- **功能**: 上传并绑定设备的证书文件
- **路径参数**: id - 设备ID
- **请求格式**: multipart/form-data
- **请求参数**:
  - cert: 证书文件（.pem格式）
- **请求示例 (curl)**:
  ```bash
  curl -X POST http://localhost:8080/bind/devices/1001/cert \
    -F "cert=@/path/to/device_certificate.pem"
  ```
- **响应格式**: JSON
- **响应示例 (成功)**:
  ```json
  {
    "code": 200,
    "message": "证书绑定成功",
    "data": {
      "deviceID": "1001",
      "certPath": "/d/code/golang/git/gin-server/regist/certs/certs/device_1001.pem"
    }
  }
  ```
- **响应示例 (失败)**:
  ```json
  {
    "error": "设备不存在"
  }
  ```

#### 4. 绑定设备密钥

- **接口**: `POST /bind/devices/:id/key`
- **功能**: 上传并绑定设备的密钥文件
- **路径参数**: id - 设备ID
- **请求格式**: multipart/form-data
- **请求参数**:
  - key: 密钥文件（.pem格式）
- **请求示例 (curl)**:
  ```bash
  curl -X POST http://localhost:8080/bind/devices/1001/key \
    -F "key=@/path/to/device_private_key.pem"
  ```
- **响应格式**: JSON
- **响应示例 (成功)**:
  ```json
  {
    "code": 200,
    "message": "密钥绑定成功",
    "data": {
      "deviceID": "1001",
      "keyPath": "/d/code/golang/git/gin-server/regist/certs/keys/device_1001.pem"
    }
  }
  ```
- **响应示例 (失败)**:
  ```json
  {
    "error": "文件大小超过限制"
  }
  ```

#### 5. 获取证书信息

- **接口**: `GET /cert/info`
- **功能**: 获取指定实体（用户或设备）的证书信息
- **请求格式**: URL查询参数
- **请求参数**:
  - type: 实体类型（user或device）
  - id: 实体ID
- **请求示例**: `http://localhost:8080/cert/info?type=user&id=1001`
- **响应格式**: JSON
- **响应示例 (成功)**:
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
- **响应示例 (未找到记录)**:
  ```json
  {
    "code": 200,
    "message": "未找到证书记录",
    "data": null
  }
  ```
- **响应示例 (参数错误)**:
  ```json
  {
    "error": "缺少必要的参数"
  }
  ```

### 认证管理接口

#### 1. 获取认证记录

- **接口**: `GET /auth/records`
- **功能**: 查询认证记录
- **请求格式**: URL查询参数
- **请求参数**:
  - username: 用户名（可选）
  - reply: 认证结果（可选）
  - start_date: 开始日期（可选）
  - end_date: 结束日期（可选）
  - page: 页码，默认1
  - page_size: 每页记录数，默认10
  - class: 认证类别（可选）
- **请求示例**: `http://localhost:8080/auth/records?username=user1&reply=Access-Accept&page=1&page_size=20`
- **响应格式**: JSON
- **响应示例**:
  ```json
  {
    "code": 200,
    "message": "Success",
    "data": {
      "total": 35,
      "total_pages": 2,
      "page": 1,
      "page_size": 20,
      "records": [
        {
          "id": 1,
          "username": "user1",
          "reply": "Access-Accept",
          "authdate": "2024-03-07T10:00:00Z",
          "class": "VPN"
        },
        {
          "id": 2,
          "username": "user1",
          "reply": "Access-Accept",
          "authdate": "2024-03-07T11:30:00Z",
          "class": "VPN"
        }
        // 更多记录...
      ]
    }
  }
  ```

### 日志管理接口

#### 1. 获取最新日志

- **接口**: `GET /logs/latest`
- **功能**: 获取系统中最新生成的日志文件内容
- **请求格式**: 无参数
- **响应格式**: JSON
- **响应示例 (成功)**:
  ```json
  {
    "timeRange": {
      "startTime": "2024-03-07T10:00:00Z",
      "duration": 300
    },
    "securityEvents": {
      "events": [
        {
          "eventId": 1001,
          "deviceId": "1001",
          "eventTime": "2024-03-07T10:02:00Z",
          "eventType": 1,
          "eventCode": "SEC_001",
          "eventDesc": "异常登录尝试",
          "createdAt": "2024-03-07T10:02:01Z"
        }
      ]
    },
    "performanceEvents": {
      "securityDevices": [
        {
          "deviceId": "1001",
          "cpuUsage": 45,
          "memoryUsage": 60,
          "onlineDuration": 3600,
          "status": 1,
          "gatewayDevices": [
            {
              "deviceId": "1002",
              "cpuUsage": 30,
              "memoryUsage": 40,
              "onlineDuration": 3600,
              "status": 1,
              "users": [
                {
                  "userId": 10001,
                  "status": 1,
                  "onlineDuration": 1800,
                  "behaviors": [
                    {
                      "time": "2024-03-07T10:01:00Z",
                      "type": 1,
                      "dataType": 1,
                      "dataSize": 1024
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
          "eventId": 2001,
          "deviceId": "1002",
          "eventTime": "2024-03-07T10:03:00Z",
          "eventType": 2,
          "eventCode": "FAULT_001",
          "eventDesc": "设备离线",
          "createdAt": "2024-03-07T10:03:01Z"
        }
      ]
    }
  }
  ```
- **响应示例 (失败)**:
  ```json
  {
    "error": "查找最新日志文件失败: 日志目录为空，没有日志文件"
  }
  ```

#### 2. 手动生成日志

- **接口**: `POST /logs/generate`
- **功能**: 手动触发生成新的日志文件
- **请求格式**: 无参数
- **响应格式**: JSON
- **响应示例 (成功)**:
  ```json
  {
    "message": "日志生成成功"
  }
  ```
- **响应示例 (失败)**:
  ```json
  {
    "error": "生成日志文件失败: 获取设备信息错误"
  }
  ```

#### 3. 获取远程日志文件列表

- **接口**: `GET /logs/files`
- **功能**: 获取已上传到远程存储的日志文件列表
- **请求格式**: 无参数
- **响应格式**: JSON
- **响应示例**:
  ```json
  [
    {
      "name": "log_20240307_100000.json",
      "path": "log_20240307_100000.json",
      "size": 10240,
      "modified_time": "2024-03-07T10:30:00Z"
    }
  ]
  ```

#### 4. 根据时间范围查询日志文件

- **接口**: `GET /logs/files/search`
- **功能**: 查询指定时间范围内的日志文件记录
- **请求格式**: URL查询参数
- **请求参数**:
  - start_time: 开始时间（RFC3339格式），可选，默认为24小时前
  - end_time: 结束时间（RFC3339格式），可选，默认为当前时间
- **请求示例**: `http://localhost:8080/logs/files/search?start_time=2024-03-06T00:00:00Z&end_time=2024-03-07T23:59:59Z`
- **响应格式**: JSON
- **响应示例**:
  ```json
  {
    "total": 24,
    "files": [
      {
        "id": 1,
        "file_name": "log_20240307_100000.json",
        "file_path": "logs/20240307/log.json",
        "file_size": 10240,
        "start_time": "2024-03-07T10:00:00Z",
        "end_time": "2024-03-07T10:10:00Z", 
        "is_encrypted": true,
        "is_uploaded": true,
        "remote_path": "log/log_20240307_100000.json",
        "uploaded_time": "2024-03-07T10:10:30Z"
      }
    ]
  }
  ```

#### 5. 创建事件记录

- **接口**: `POST /logs/events`
- **功能**: 创建新的事件记录
- **请求格式**: JSON
- **请求参数**:
  ```json
  {
    "event_code": "SEC_001",  // 事件代码，必填
    "event_desc": "异常登录尝试", // 事件描述，必填
    "device_id": 1001,  // 设备ID，必填
    "event_type": 1  // 事件类型，必填，1:安全事件, 2:故障事件
  }
  ```
- **响应格式**: JSON
- **响应示例 (成功)**:
  ```json
  {
    "ID": 1,
    "CreatedAt": "2024-03-07T15:04:05Z",
    "UpdatedAt": "2024-03-07T15:04:05Z",
    "event_id": 1678234445000000000,
    "device_id": 1001,
    "event_time": "2024-03-07T15:04:05Z",
    "event_type": 1,
    "event_code": "SEC_001",
    "event_desc": "异常登录尝试"
  }
  ```

#### 6. 查询事件记录

- **接口**: `GET /logs/events/search`
- **功能**: 查询指定时间范围内的事件记录
- **请求格式**: URL查询参数
- **请求参数**:
  - start_time: 开始时间（RFC3339格式），可选，默认为1小时前
  - end_time: 结束时间（RFC3339格式），可选，默认为当前时间
- **请求示例**: `http://localhost:8080/logs/events/search?start_time=2024-03-07T00:00:00Z&end_time=2024-03-07T23:59:59Z`
- **响应格式**: JSON
- **响应示例**:
  ```json
  {
    "total": 5,
    "events": [
      {
        "ID": 1,
        "CreatedAt": "2024-03-07T15:04:05Z",
        "UpdatedAt": "2024-03-07T15:04:05Z",
        "event_id": 1678234445000000000,
        "device_id": 1001,
        "event_time": "2024-03-07T15:04:05Z",
        "event_type": 1,
        "event_code": "SEC_001",
        "event_desc": "异常登录尝试"
      }
    ]
  }
  ```

#### 7. 记录用户行为

- **接口**: `POST /logs/behaviors`
- **功能**: 记录用户行为数据
- **请求格式**: JSON
- **请求参数**:
  ```json
  {
    "user_id": 10001,  // 用户ID，必填
    "behavior_type": 1,  // 行为类型，必填，1:发送，2:接收
    "data_type": 1,  // 数据类型，必填，1:文件，2:消息
    "data_size": 1024  // 数据大小（字节），必填
  }
  ```
- **响应格式**: JSON
- **响应示例 (成功)**:
  ```json
  {
    "ID": 1,
    "CreatedAt": "2024-03-07T15:04:05Z",
    "UpdatedAt": "2024-03-07T15:04:05Z",
    "user_id": 10001,
    "behavior_time": "2024-03-07T15:04:05Z",
    "behavior_type": 1,
    "data_type": 1,
    "data_size": 1024
  }
  ```

#### 8. 查询用户行为

- **接口**: `GET /logs/behaviors/:user_id`
- **功能**: 查询指定用户的行为记录
- **路径参数**: user_id - 用户ID
- **请求示例**: `http://localhost:8080/logs/behaviors/10001`
- **响应格式**: JSON
- **响应示例**:
  ```json
  {
    "total": 10,
    "behaviors": [
      {
        "ID": 1,
        "CreatedAt": "2024-03-07T15:04:05Z",
        "UpdatedAt": "2024-03-07T15:04:05Z",
        "user_id": 10001,
        "behavior_time": "2024-03-07T15:04:05Z",
        "behavior_type": 1,
        "data_type": 1,
        "data_size": 1024
      }
    ],
    "user_id": 10001
  }
  ```

## 日志管理模块详细说明

### 日志文件结构

系统生成的日志文件采用JSON格式，包含以下主要部分：

#### 1. 时间范围（time_range）

- **start_time**: 统计起始时间，ISO8601格式（例如："2025-03-27T20:14:40.6911388+08:00"）
- **duration**: 统计时长（秒），即日志覆盖的时间段

#### 2. 安全事件（security_events）

- **events**: 安全事件列表，包含安全相关的告警和事件
  - **event_id**: 事件ID，唯一标识一个事件
  - **device_id**: 设备ID，事件发生的设备
  - **event_time**: 事件发生时间
  - **event_type**: 事件类型，1表示安全事件
  - **event_code**: 事件代码，例如"SEC_001"
  - **event_desc**: 事件描述

#### 3. 性能事件（performance_events）

- **security_devices**: 安全接入管理设备列表
  - **device_id**: 安全接入管理设备ID
  - **cpu_usage**: CPU使用率（百分比，0-100）
  - **memory_usage**: 内存使用率（百分比，0-100）
  - **online_duration**: 设备在线时长（秒）
  - **status**: 设备状态（1:在线，2:离线，3:冻结，4:注销）
  - **gateway_devices**: 网关设备列表
    - **device_id**: 网关设备ID
    - **cpu_usage**: CPU使用率（百分比，0-100）
    - **memory_usage**: 内存使用率（百分比，0-100）
    - **online_duration**: 设备在线时长（秒）
    - **status**: 设备状态（1:在线，2:离线，3:冻结，4:注销）
    - **users**: 用户列表
      - **user_id**: 用户ID
      - **status**: 用户状态（1:在线，2:离线，3:冻结，4:注销）
      - **online_duration**: 用户在线时长（秒）
      - **behaviors**: 用户行为列表
        - **time**: 行为发生时间
        - **type**: 行为类型（1:发送，2:接收）
        - **data_type**: 数据类型（1:文件，2:消息）
        - **data_size**: 数据大小（字节）

#### 4. 故障事件（fault_events）

- **events**: 故障事件列表，包含系统故障和错误
  - **event_id**: 事件ID，唯一标识一个事件
  - **device_id**: 设备ID，事件发生的设备
  - **event_time**: 事件发生时间
  - **event_type**: 事件类型，2表示故障事件
  - **event_code**: 事件代码，例如"FAULT_001"
  - **event_desc**: 事件描述

### 日志管理模块工作流程

1. **定时生成**：
   - 系统根据配置的生成间隔（默认10分钟）自动生成日志文件
   - 生成时会收集该时间段内的所有事件、设备性能和用户行为数据
   - 日志文件以 `YYYYMMDDHHMMSS.json` 格式命名（如：`20240515100000.json`）

2. **加密处理**：

   - 如果启用加密（配置项 `enable_encryption`设为true），日志文件会使用AES密钥加密
   - 加密后的文件存储在原目录下的 `encrypted`子目录中
   - 加密使用的密钥通过配置的公钥加密后存储

3. **远程上传**：

   - 日志文件生成并加密后会自动上传到远程存储
   - 支持两种存储方式：Gitee仓库和FTP服务器
   - 上传成功后，本地数据库会记录日志文件的元数据

4. **查询访问**：

   - 通过REST API接口提供日志查询服务
   - 支持按时间范围查询日志文件和事件记录
   - 可获取最新的日志文件内容

### 日志文件存储位置

- **本地存储**：

  - 未加密：`logs/YYYYMMDD/YYYYMMDDHHMMSS.json`
  - 加密后：`logs/YYYYMMDD/encrypted/YYYYMMDDHHMMSS.json`
- **远程存储**：

  - 统一存储在仓库分支的 `/log`目录下
  - 文件格式：`/log/YYYYMMDDHHMMSS.tar.gz` (包含加密的日志文件和密钥)

### 注意事项

1. **性能考虑**：

   - 日志生成间隔不宜设置过短，建议不少于5分钟
   - 日志文件可能随时间累积变大，请确保有足够的存储空间
   - 对于大型部署，应考虑定期清理或归档旧日志文件

2. **安全考虑**：

   - 强烈建议启用日志加密功能以保护敏感信息
   - 密钥应妥善保管，避免泄露
   - 应定期更换加密密钥

3. **存储选择**：

   - Gitee存储适合小型部署和测试环境
   - FTP存储更适合生产环境，提供更好的性能和更大的存储容量
   - 无论选择哪种存储方式，都应确保有适当的备份策略

4. **配置建议**：

   - 本地日志目录应有足够的磁盘空间
   - 应正确配置公钥和私钥路径
   - 远程存储的访问凭证应定期更新

## 错误码说明

| 状态码 | 描述       | 说明                               |
| ------ | ---------- | ---------------------------------- |
| 200    | 成功       | 请求处理成功                       |
| 201    | 创建成功   | 资源创建成功                       |
| 400    | 请求错误   | 请求参数错误或格式不正确           |
| 401    | 未授权     | 缺少或无效的认证信息               |
| 403    | 禁止访问   | 没有权限访问所请求的资源           |
| 404    | 资源不存在 | 请求的资源不存在                   |
| 409    | 资源冲突   | 请求的资源与服务器上已有资源冲突   |
| 500    | 服务器错误 | 服务器内部错误，无法完成请求的处理 |

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
├── database/       # 数据库管理
│   ├── connections/  # 数据库连接管理
│   ├── migrations/   # 数据库迁移
│   ├── models/       # 数据模型
│   ├── repositories/ # 数据仓库
│   └── testdata/     # 测试数据
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

## 证书管理说明

### 证书存储结构

证书和密钥文件按照以下结构存储：

```
regist/certs/
├── certs/             # 存放证书文件
│   ├── user_1001.pem  # 用户证书文件(命名格式: user_<userID>.pem)
│   └── device_1001.pem  # 设备证书文件(命名格式: device_<deviceID>.pem)
└── keys/              # 存放密钥文件
    ├── user_1001.pem  # 用户密钥文件
    └── device_1001.pem  # 设备密钥文件
```

### 安全考虑

- 证书目录：权限设置为0755
- 密钥目录：权限设置为0700（更严格的权限控制）
- 证书文件：权限默认为0644
- 密钥文件：权限设置为0600（只有所有者可读写）

### 证书和密钥文件要求

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

## 数据库结构说明

项目使用了两个主要的数据库：

### 1. 主数据库 (gin_server)

主数据库包含以下表：

#### 1.1 设备表 (devices)

| 字段名              | 类型         | 描述                                                                      |
| ------------------- | ------------ | ------------------------------------------------------------------------- |
| id                  | INT          | 自增主键                                                                  |
| deviceName          | VARCHAR(50)  | 设备名称                                                                  |
| deviceType          | INT          | 设备类型，1:网关设备A型，2:网关设备B型，3:网关设备C型，4:安全接入管理设备 |
| password            | VARCHAR(255) | 设备登录口令                                                              |
| deviceID            | INT          | 设备唯一标识                                                              |
| superiorDeviceID    | INT          | 上级设备ID                                                                |
| deviceStatus        | INT          | 设备状态，1:在线，2:离线，3:冻结，4:注销                                  |
| certID              | VARCHAR(64)  | 证书ID                                                                    |
| keyID               | VARCHAR(64)  | 密钥ID                                                                    |
| registerIP          | VARCHAR(24)  | 注册IP                                                                    |
| email               | VARCHAR(32)  | 联系邮箱                                                                  |
| hardwareFingerprint | VARCHAR(128) | 设备硬件指纹                                                              |
| anonymousUser       | VARCHAR(50)  | 匿名用户                                                                  |

#### 1.2 用户表 (users)

| 字段名          | 类型         | 描述                                     |
| --------------- | ------------ | ---------------------------------------- |
| id              | INT          | 自增主键                                 |
| username        | VARCHAR(20)  | 用户名                                   |
| password        | VARCHAR(255) | 密码                                     |
| userID          | INT          | 用户唯一标识                             |
| userType        | INT          | 用户类型                                 |
| gatewayDeviceID | INT          | 用户所属网关设备ID                       |
| status          | INT          | 用户状态，1:在线，2:离线，3:冻结，4:注销 |
| onlineDuration  | INT          | 在线时长                                 |
| certID          | VARCHAR(64)  | 证书ID                                   |
| keyID           | VARCHAR(64)  | 密钥ID                                   |
| email           | VARCHAR(32)  | 邮箱                                     |

#### 1.3 证书表 (certs)

| 字段名      | 类型         | 描述                  |
| ----------- | ------------ | --------------------- |
| id          | INT          | 自增主键              |
| entity_type | VARCHAR(10)  | 实体类型(user/device) |
| entity_id   | VARCHAR(20)  | 实体ID                |
| cert_path   | VARCHAR(255) | 证书文件路径          |
| key_path    | VARCHAR(255) | 密钥文件路径          |
| upload_time | DATETIME     | 上传时间              |

### 2. Radius认证数据库 (radius)

#### 2.1 认证记录表 (radpostauth)

| 字段名   | 类型         | 描述     |
| -------- | ------------ | -------- |
| id       | INT          | 自增主键 |
| username | VARCHAR(64)  | 用户名   |
| pass     | VARCHAR(64)  | 密码     |
| reply    | VARCHAR(32)  | 认证响应 |
| authdate | TIMESTAMP(6) | 认证时间 |
| class    | VARCHAR(64)  | 认证类型 |

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

6. 证书和密钥安全

   - 证书和密钥应当妥善保管
   - 密钥文件权限设置为限制性权限(0600)
   - 定期更新证书和密钥
   - 对证书的有效性进行验证

## 许可证

[MIT License](LICENSE)

## 接口更新日志

### 2024-03-25 接口命名规范调整

为了提高代码一致性和遵循最佳实践，系统接口已经进行了以下调整：

1. **字段命名规范化**：

   - 所有接口请求和响应字段统一采用蛇形命名法（snake_case）
   - 如：`userName` 更改为 `user_name`，`deviceID` 更改为 `device_id`

2. **设备注册接口简化**：

   - 设备注册接口简化为5个必要字段
   - 其他非必要字段（如 `device_status`、`register_ip`、`email` 等）将由系统自动填充

3. **响应格式统一**：

   - 所有接口响应统一使用 `code`、`message`、`data` 三字段结构
   - 成功响应包含完整的对象数据

这些更改旨在简化API使用并提高系统的可维护性。客户端应用需要相应更新以适应新的命名格式。

### 近期功能更新 (2024-05-15)

为提高日志管理的安全性和规范性，日志管理功能已完成以下更新：

1. **日志文件命名格式修改**：
   - 原格式：`log_YYYYMMDD_HHmmss.json`
   - 新格式：`YYYYMMDDHHMMSS.json`（如：`20240515100000.json`）
   - 更简洁的命名方式便于排序和识别

2. **AES密钥生成优化**：
   - 每次上传日志时都会重新生成一次AES密钥用于加密
   - 通过`crypto.NewAESEncryptor`方法生成随机密钥，基于Go的`crypto/rand`包
   - 密钥长度根据配置可以是128位、192位或256位
   - 提高了密钥安全性，避免长期使用相同密钥可能带来的风险

3. **启用打包功能**：
   - 修改实现：使用`UploadManager.Upload`方法替代原有的`UploadFile`方法
   - 通过`CompressStep`组件自动将日志和密钥文件打包为单一`tar.gz`文件
   - 打包格式为`YYYYMMDDHHMMSS.tar.gz`（如：`20240515100000.tar.gz`）
   - 归档文件内含加密的日志文件和RSA公钥加密的AES密钥
   - 简化了管理并确保日志和其对应密钥始终一起存储

4. **优化上传路径**：
   - 统一设置上传目录为`/log`，确保所有日志文件存储在规范位置
   - 修改实现：在`LogManager.uploadLog`方法中强制指定上传路径
   - 不再依赖配置文件中的`UploadDir`设置，避免误配置
   - 便于自动化工具进行日志收集和分析
