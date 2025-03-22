# 数据库代码迁移指南

本文档提供关于如何从旧的数据库操作方式迁移到新的、更优化的数据库操作模式的指导。

## 1. 数据库连接管理器变更

### 旧方式(不推荐)

```go
// 初始化数据库连接管理器
dbConfig := &database.DBConfig{
    Host:     cfg.DBHost,
    Port:     cfg.DBPort,
    User:     cfg.DBUser,
    Password: cfg.DBPassword,
    DBName:   cfg.DBName,
    Debug:    cfg.Debug,
}

// 创建数据库连接管理器
dbManager, err := database.NewManager(dbConfig)
if err != nil {
    stdlog.Fatalf("创建数据库连接管理器失败: %v", err)
}
defer dbManager.Close()

// 获取数据库连接
db := dbManager.GetDB()

// 初始化Radius数据库连接
if err := dbManager.InitRadiusDB(cfg); err != nil {
    stdlog.Fatalf("初始化Radius数据库失败: %v", err)
}
radiusDB := dbManager.GetRadiusDB()
```

### 新方式(推荐)

```go
// 初始化数据库连接管理器
database.Initialize(cfg)
defer database.CloseAll() // 程序结束时关闭所有数据库连接

// 获取主数据库连接
db, err := database.GetDB()
if err != nil {
    stdlog.Fatalf("获取数据库连接失败: %v", err)
}

// 获取Radius数据库连接(如果需要)
radiusDB, err := database.GetRadiusDB()
if err != nil {
    stdlog.Fatalf("获取Radius数据库连接失败: %v", err)
}
```

## 2. 数据库模型

为了保持代码的一致性和可维护性，所有数据库模型都应该定义在 `database/models` 目录下。

### 旧模型(不推荐)

`regist/model/db.go` 中定义了一些数据模型：

```go
// User 结构体定义用户信息
type User struct {
    UserName        string `json:"userName"`        // 用户名
    PassWD          string `json:"passWD"`          // 密码
    UserID          int    `json:"userID"`          // 用户唯一标识
    // ...
}

// Device 结构体定义设备信息
type Device struct {
    // ...
}
```

### 新模型(推荐)

`database/models` 目录下的模型：

```go
// User 用户信息
type User struct {
    gorm.Model
    Username    string    `json:"username" gorm:"column:username;uniqueIndex;not null"`
    Password    string    `json:"-" gorm:"column:password;not null"`
    // ...
}

// Device 设备信息
type Device struct {
    gorm.Model
    // ...
}
```

## 3. 数据库操作

### 旧方式(不推荐)

直接在模型文件中定义数据库操作函数：

```go
// GetUserByID 根据ID获取用户
func GetUserByID(userID int) (User, error) {
    // ...
}

// UpdateUser 更新用户信息
func UpdateUser(user User) error {
    // ...
}
```

### 新方式(推荐)

使用仓库模式进行数据库操作：

```go
// 获取用户仓库
userRepo := repoFactory.GetUserRepository()

// 查询用户
user, err := userRepo.FindByID(userID)

// 更新用户
user.Username = "新用户名"
err = userRepo.Update(user)
```

## 4. 数据库迁移

### 旧方式(不推荐)

```go
// 初始化数据库连接
model.InitDB()

// 手动迁移
if err := logModel.MigrateDatabase(db); err != nil {
    // ...
}
```

### 新方式(推荐)

```go
// 运行数据库迁移
if err := migrations.AutoMigrate(db); err != nil {
    // ...
}
```

## 5. 数据库事务

### 旧方式(不推荐)

```go
tx := db.Begin()
// 执行操作
if err != nil {
    tx.Rollback()
    return err
}
return tx.Commit()
```

### 新方式(推荐)

```go
// 使用仓库工厂获取事务
err := db.Transaction(func(tx *gorm.DB) error {
    txRepoFactory := repoFactory.WithTx(tx)
    userRepo := txRepoFactory.GetUserRepository()
  
    // 执行操作
    if err := userRepo.Create(user); err != nil {
        return err // 自动回滚
    }
  
    return nil // 自动提交
})
```

## 迁移计划

1. 首先替换所有的数据库连接管理器相关代码
2. 然后用新的模型替换旧的模型引用
3. 使用仓库模式替换直接的数据库操作
4. 最后更新数据库迁移代码

## 注意事项

- 确保在迁移过程中不会丢失数据
- 在生产环境中执行迁移前，先在测试环境中充分测试
- 迁移应该分阶段进行，每个阶段都要经过充分测试
