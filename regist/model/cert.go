package model

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gin-server/config"
)

// 证书文件存储的基础路径
const (
	BaseCertPath = "regist/certs"
	CertsDir     = "certs"
	KeysDir      = "keys"
)

// 确保证书目录存在
func EnsureCertDirsExist() error {
	cfg := config.GetConfig()

	// 创建基础目录
	basePath := BaseCertPath
	if err := os.MkdirAll(basePath, 0755); err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("创建基础证书目录失败: %v\n", err)
		}
		return fmt.Errorf("创建基础证书目录失败: %w", err)
	}

	// 创建证书目录
	certPath := filepath.Join(basePath, CertsDir)
	if err := os.MkdirAll(certPath, 0755); err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("创建证书目录失败: %v\n", err)
		}
		return fmt.Errorf("创建证书目录失败: %w", err)
	}

	// 创建密钥目录
	keyPath := filepath.Join(basePath, KeysDir)
	if err := os.MkdirAll(keyPath, 0700); err != nil { // 密钥目录权限更严格
		if cfg.DebugLevel == "true" {
			log.Printf("创建密钥目录失败: %v\n", err)
		}
		return fmt.Errorf("创建密钥目录失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Println("证书目录结构已创建")
	}

	return nil
}

// 获取文件的绝对路径
func GetAbsolutePath(relativePath string) string {
	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		log.Printf("获取工作目录失败: %v\n", err)
		return relativePath // 如果失败，返回相对路径
	}

	// 构建绝对路径
	absPath := filepath.Join(workDir, relativePath)

	// 根据操作系统规范化路径
	if runtime.GOOS == "windows" {
		// Windows 下使用反斜杠
		absPath = strings.ReplaceAll(absPath, "/", "\\")
	}

	return absPath
}

// SaveCertFile 保存证书文件
func SaveCertFile(entityType, entityID string, fileReader io.Reader, isKey bool) (string, error) {
	cfg := config.GetConfig()

	// 确保目录存在
	if err := EnsureCertDirsExist(); err != nil {
		return "", err
	}

	// 确定文件类型和目录
	dirType := CertsDir
	if isKey {
		dirType = KeysDir
	}

	// 构建文件名和路径
	fileName := fmt.Sprintf("%s_%s.pem", entityType, entityID)
	filePath := filepath.Join(BaseCertPath, dirType, fileName)

	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("创建文件失败: %v\n", err)
		}
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer out.Close()

	// 写入文件内容
	_, err = io.Copy(out, fileReader)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("写入文件失败: %v\n", err)
		}
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	// 设置合适的权限
	if isKey {
		// 密钥文件权限更严格
		if err := os.Chmod(filePath, 0600); err != nil {
			if cfg.DebugLevel == "true" {
				log.Printf("设置密钥文件权限失败: %v\n", err)
			}
			// 不中断流程，只记录错误
			log.Printf("警告: 设置密钥文件权限失败: %v\n", err)
		}
	}

	// 获取并返回绝对路径
	absPath := GetAbsolutePath(filePath)
	if cfg.DebugLevel == "true" {
		log.Printf("文件保存成功: %s\n", absPath)
	}

	return absPath, nil
}

// AddCertRecord 添加证书记录到数据库
func AddCertRecord(entityType, entityID, certPath, keyPath string) error {
	cfg := config.GetConfig()
	db := GetDB()

	// 检查是否已存在记录
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM certs WHERE entity_type = ? AND entity_id = ?", entityType, entityID).Scan(&count)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("检查证书记录失败: %v\n", err)
		}
		return fmt.Errorf("检查证书记录失败: %w", err)
	}

	// 当前时间
	now := time.Now()

	// 如果记录已存在，更新记录
	if count > 0 {
		query := "UPDATE certs SET upload_time = ?"
		params := []interface{}{now}

		// 只更新非空字段
		if certPath != "" {
			query += ", cert_path = ?"
			params = append(params, certPath)
		}
		if keyPath != "" {
			query += ", key_path = ?"
			params = append(params, keyPath)
		}

		// 添加 WHERE 条件
		query += " WHERE entity_type = ? AND entity_id = ?"
		params = append(params, entityType, entityID)

		// 执行更新
		_, err := db.Exec(query, params...)
		if err != nil {
			if cfg.DebugLevel == "true" {
				log.Printf("更新证书记录失败: %v\n", err)
			}
			return fmt.Errorf("更新证书记录失败: %w", err)
		}
	} else {
		// 创建新记录
		_, err := db.Exec(
			"INSERT INTO certs (entity_type, entity_id, cert_path, key_path, upload_time) VALUES (?, ?, ?, ?, ?)",
			entityType, entityID, certPath, keyPath, now,
		)
		if err != nil {
			if cfg.DebugLevel == "true" {
				log.Printf("添加证书记录失败: %v\n", err)
			}
			return fmt.Errorf("添加证书记录失败: %w", err)
		}
	}

	if cfg.DebugLevel == "true" {
		log.Printf("证书记录已保存: %s/%s\n", entityType, entityID)
	}
	return nil
}

// UpdateUserCertInfo 更新用户表中的证书信息
func UpdateUserCertInfo(userID string, certID, keyID string) error {
	cfg := config.GetConfig()
	db := GetDB()

	// 构建更新语句
	query := "UPDATE users SET"
	params := []interface{}{}

	if certID != "" {
		query += " certID = ?"
		params = append(params, certID)
	}

	if keyID != "" {
		if len(params) > 0 {
			query += ","
		}
		query += " keyID = ?"
		params = append(params, keyID)
	}

	// 如果没有需要更新的字段，直接返回
	if len(params) == 0 {
		return nil
	}

	// 添加 WHERE 条件
	query += " WHERE userID = ?"
	params = append(params, userID)

	// 执行更新
	_, err := db.Exec(query, params...)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("更新用户证书信息失败: %v\n", err)
		}
		return fmt.Errorf("更新用户证书信息失败: %w", err)
	}

	return nil
}

// UpdateDeviceCertInfo 更新设备表中的证书信息
func UpdateDeviceCertInfo(deviceID, certID, keyID string) error {
	cfg := config.GetConfig()
	db := GetDB()

	// 构建更新语句
	query := "UPDATE devices SET"
	params := []interface{}{}

	if certID != "" {
		query += " certID = ?"
		params = append(params, certID)
	}

	if keyID != "" {
		if len(params) > 0 {
			query += ","
		}
		query += " keyID = ?"
		params = append(params, keyID)
	}

	// 如果没有需要更新的字段，直接返回
	if len(params) == 0 {
		return nil
	}

	// 添加 WHERE 条件
	query += " WHERE deviceID = ?"
	params = append(params, deviceID)

	// 执行更新
	_, err := db.Exec(query, params...)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("更新设备证书信息失败: %v\n", err)
		}
		return fmt.Errorf("更新设备证书信息失败: %w", err)
	}

	return nil
}

// GetCertInfo 获取指定实体的证书信息
func GetCertInfo(entityType, entityID string) (*Cert, error) {
	cfg := config.GetConfig()
	db := GetDB()

	var cert Cert
	var uploadTime string

	err := db.QueryRow(
		"SELECT id, entity_type, entity_id, cert_path, key_path, upload_time FROM certs WHERE entity_type = ? AND entity_id = ?",
		entityType, entityID,
	).Scan(&cert.ID, &cert.EntityType, &cert.EntityID, &cert.CertPath, &cert.KeyPath, &uploadTime)

	if err != nil {
		if err == sql.ErrNoRows {
			// 没有找到记录
			if cfg.DebugLevel == "true" {
				log.Printf("未找到证书记录: %s/%s\n", entityType, entityID)
			}
			return nil, nil // 返回nil表示没有找到记录
		}
		// 其他数据库错误
		if cfg.DebugLevel == "true" {
			log.Printf("获取证书信息失败: %v\n", err)
		}
		return nil, fmt.Errorf("获取证书信息失败: %w", err)
	}

	// 解析时间
	cert.UploadTime, err = time.Parse("2006-01-02 15:04:05", uploadTime)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("解析上传时间失败: %v\n", err)
		}
		// 使用当前时间作为后备
		cert.UploadTime = time.Now()
	}

	return &cert, nil
}
