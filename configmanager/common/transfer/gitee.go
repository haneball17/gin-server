package transfer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gin-server/config"
)

const (
	giteeAPIBase = "https://gitee.com/api/v5"
)

// GiteeTransporter Gitee文件传输实现
type GiteeTransporter struct {
	accessToken string
	owner       string
	repo        string
	branch      string
	client      *http.Client
}

// GiteeFile Gitee文件信息
type GiteeFile struct {
	Path    string `json:"path"`
	Type    string `json:"type"`
	Size    int64  `json:"size"`
	Name    string `json:"name"`
	SHA     string `json:"sha"`
	Content string `json:"content"`
}

// NewGiteeTransporter 创建Gitee传输器
func NewGiteeTransporter(cfg *config.GiteeConfig) *GiteeTransporter {
	return &GiteeTransporter{
		accessToken: cfg.AccessToken,
		owner:       cfg.Owner,
		repo:        cfg.Repo,
		branch:      cfg.Branch,
		client:      &http.Client{Timeout: 30 * time.Second},
	}
}

// Upload 上传文件到Gitee
func (t *GiteeTransporter) Upload(localPath, remotePath string) error {
	cfg := config.GetConfig()

	// 读取本地文件
	content, err := os.ReadFile(localPath)
	if err != nil {
		return NewTransferError("upload", remotePath, fmt.Errorf("读取本地文件失败: %w", err))
	}

	if cfg.DebugLevel == "true" {
		log.Printf("准备上传文件 %s 到 %s\n", localPath, remotePath)
	}

	// 获取现有文件信息
	existingFile, err := t.getFile(remotePath)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("获取文件信息失败: %v\n", err)
		}
		if !strings.Contains(err.Error(), "404") {
			return NewTransferError("upload", remotePath, fmt.Errorf("获取文件信息失败: %w", err))
		}
	}

	if cfg.DebugLevel == "true" {
		if existingFile != nil {
			log.Printf("文件已存在，SHA: %s\n", existingFile.SHA)
		} else {
			log.Println("文件不存在，将创建新文件")
		}
	}

	// 准备基本请求数据
	requestData := map[string]interface{}{
		"access_token": t.accessToken,
		"content":      base64.StdEncoding.EncodeToString(content),
		"branch":       t.branch,
	}

	// 根据文件是否存在设置不同的消息
	if existingFile != nil {
		requestData["message"] = fmt.Sprintf("更新文件: %s", remotePath)
		requestData["sha"] = existingFile.SHA
	} else {
		requestData["message"] = fmt.Sprintf("创建文件: %s", remotePath)
	}

	// 编码请求数据
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return NewTransferError("upload", remotePath, fmt.Errorf("编码请求数据失败: %w", err))
	}

	if cfg.DebugLevel == "true" {
		log.Printf("请求数据: %s\n", string(jsonData))
	}

	// 构造API URL
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", giteeAPIBase, t.owner, t.repo, remotePath)

	if cfg.DebugLevel == "true" {
		log.Printf("API URL: %s\n", url)
	}

	// 根据文件是否存在选择不同的HTTP方法
	method := http.MethodPut
	if existingFile == nil {
		method = http.MethodPost
	}

	// 创建请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return NewTransferError("upload", remotePath, fmt.Errorf("创建请求失败: %w", err))
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := t.client.Do(req)
	if err != nil {
		return NewTransferError("upload", remotePath, fmt.Errorf("发送请求失败: %w", err))
	}
	defer resp.Body.Close()

	// 读取响应内容（用于调试）
	respBody, _ := io.ReadAll(resp.Body)

	if cfg.DebugLevel == "true" {
		log.Printf("响应状态码: %d\n", resp.StatusCode)
		log.Printf("响应内容: %s\n", string(respBody))
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return NewTransferError("upload", remotePath, fmt.Errorf("上传文件失败: HTTP %d - %s", resp.StatusCode, string(respBody)))
	}

	return nil
}

// Download 从Gitee下载文件
func (t *GiteeTransporter) Download(remotePath, localPath string) error {
	// 获取文件信息
	file, err := t.getFile(remotePath)
	if err != nil {
		return NewTransferError("download", remotePath, err)
	}

	if file == nil {
		return NewTransferError("download", remotePath, fmt.Errorf("文件不存在"))
	}

	// 解码文件内容
	content, err := base64.StdEncoding.DecodeString(file.Content)
	if err != nil {
		return NewTransferError("download", remotePath, fmt.Errorf("解码文件内容失败: %w", err))
	}

	// 创建本地目录
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return NewTransferError("download", localPath, fmt.Errorf("创建本地目录失败: %w", err))
	}

	// 写入本地文件
	if err := os.WriteFile(localPath, content, 0644); err != nil {
		return NewTransferError("download", localPath, fmt.Errorf("写入本地文件失败: %w", err))
	}

	return nil
}

// List 列出目录内容
func (t *GiteeTransporter) List(remotePath string) ([]FileInfo, error) {
	// 构建API URL
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s?access_token=%s&ref=%s",
		giteeAPIBase, t.owner, t.repo, remotePath, t.accessToken, t.branch)

	// 发送请求
	resp, err := t.client.Get(url)
	if err != nil {
		return nil, NewTransferError("list", remotePath, fmt.Errorf("发送请求失败: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []FileInfo{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, NewTransferError("list", remotePath, fmt.Errorf("获取目录列表失败: HTTP %d - %s", resp.StatusCode, string(body)))
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewTransferError("list", remotePath, fmt.Errorf("读取响应失败: %w", err))
	}

	// 尝试解析为文件数组
	var giteeFiles []GiteeFile
	if err := json.Unmarshal(body, &giteeFiles); err != nil {
		// 如果解析数组失败，尝试解析单个文件
		var singleFile GiteeFile
		if err := json.Unmarshal(body, &singleFile); err != nil {
			return nil, NewTransferError("list", remotePath, fmt.Errorf("解析响应失败: %w", err))
		}
		giteeFiles = []GiteeFile{singleFile}
	}

	// 转换为FileInfo
	files := make([]FileInfo, 0, len(giteeFiles))
	for _, file := range giteeFiles {
		files = append(files, FileInfo{
			Name:  file.Name,
			Size:  file.Size,
			IsDir: file.Type == "dir",
			Path:  file.Path,
		})
	}

	return files, nil
}

// Delete 删除文件
func (t *GiteeTransporter) Delete(remotePath string) error {
	// 获取文件信息
	file, err := t.getFile(remotePath)
	if err != nil {
		return NewTransferError("delete", remotePath, err)
	}

	if file == nil {
		return NewTransferError("delete", remotePath, fmt.Errorf("文件不存在"))
	}

	// 准备请求数据
	data := map[string]interface{}{
		"access_token": t.accessToken,
		"sha":          file.SHA,
		"branch":       t.branch,
		"message":      fmt.Sprintf("Delete %s", remotePath),
	}

	// 编码请求数据
	jsonData, err := json.Marshal(data)
	if err != nil {
		return NewTransferError("delete", remotePath, fmt.Errorf("编码请求数据失败: %w", err))
	}

	// 构建API URL
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", giteeAPIBase, t.owner, t.repo, remotePath)

	// 创建请求
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return NewTransferError("delete", remotePath, fmt.Errorf("创建请求失败: %w", err))
	}

	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := t.client.Do(req)
	if err != nil {
		return NewTransferError("delete", remotePath, fmt.Errorf("发送请求失败: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return NewTransferError("delete", remotePath, fmt.Errorf("删除文件失败: HTTP %d - %s", resp.StatusCode, string(body)))
	}

	return nil
}

// LastModified 获取文件最后修改时间
func (t *GiteeTransporter) LastModified(remotePath string) (time.Time, error) {
	// 获取文件提交历史
	url := fmt.Sprintf("%s/repos/%s/%s/commits?access_token=%s&path=%s&page=1&per_page=1",
		giteeAPIBase, t.owner, t.repo, t.accessToken, remotePath)

	resp, err := t.client.Get(url)
	if err != nil {
		return time.Time{}, NewTransferError("lastModified", remotePath, fmt.Errorf("发送请求失败: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return time.Time{}, NewTransferError("lastModified", remotePath, fmt.Errorf("获取提交历史失败: HTTP %d - %s", resp.StatusCode, string(body)))
	}

	var commits []struct {
		Commit struct {
			Author struct {
				Date time.Time `json:"date"`
			} `json:"author"`
		} `json:"commit"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return time.Time{}, NewTransferError("lastModified", remotePath, fmt.Errorf("解析响应失败: %w", err))
	}

	if len(commits) == 0 {
		return time.Time{}, NewTransferError("lastModified", remotePath, fmt.Errorf("文件不存在"))
	}

	return commits[0].Commit.Author.Date, nil
}

// Close 关闭传输器
func (t *GiteeTransporter) Close() error {
	return nil
}

// getFile 获取文件信息
func (t *GiteeTransporter) getFile(path string) (*GiteeFile, error) {
	cfg := config.GetConfig()
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s?access_token=%s&ref=%s",
		giteeAPIBase, t.owner, t.repo, path, t.accessToken, t.branch)

	if cfg.DebugLevel == "true" {
		log.Printf("获取文件信息，URL: %s\n", url)
	}

	resp, err := t.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}
	defer resp.Body.Close()

	// 如果文件不存在，返回nil和nil
	if resp.StatusCode == http.StatusNotFound {
		if cfg.DebugLevel == "true" {
			log.Printf("文件不存在: %s\n", path)
		}
		return nil, nil
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Printf("响应状态码: %d\n", resp.StatusCode)
		log.Printf("响应内容: %s\n", string(body))
	}

	// 如果状态码不是200，返回错误
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d - %s", resp.StatusCode, string(body))
	}

	// 首先尝试解析为单个文件
	var file GiteeFile
	if err := json.Unmarshal(body, &file); err != nil {
		// 如果解析单个文件失败，尝试解析为数组
		var files []GiteeFile
		if err := json.Unmarshal(body, &files); err != nil {
			return nil, fmt.Errorf("解析响应内容失败: %w", err)
		}

		if cfg.DebugLevel == "true" {
			log.Printf("获取到 %d 个文件\n", len(files))
		}

		// 在数组中查找匹配的文件
		for _, f := range files {
			if cfg.DebugLevel == "true" {
				log.Printf("检查文件: %s\n", f.Path)
			}
			if f.Path == path {
				return &f, nil
			}
		}
		// 如果没有找到匹配的文件，返回nil
		if cfg.DebugLevel == "true" {
			log.Printf("未找到匹配的文件: %s\n", path)
		}
		return nil, nil
	}

	return &file, nil
}
