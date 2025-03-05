# Gin Server

基于Gin框架的服务器应用，提供用户认证、设备管理和配置管理功能。

## 功能特性

- 用户认证管理
  - 用户注册与管理
  - 认证记录查询
  - 权限控制
- 设备管理
  - 设备注册
  - 设备信息更新
  - 设备状态监控
- 配置管理
  - 远程配置同步
  - 文件变更监控
  - 加密传输
  - 多存储方式支持（Gitee/FTP）

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