package compress

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	testDataDir  = "testdata"
	testSrcDir   = filepath.Join(testDataDir, "src")
	testDestDir  = filepath.Join(testDataDir, "dest")
	testFiles    = []string{"file1.txt", "file2.txt", "dir1/file3.txt"}
	testContents = []string{"content1", "content2", "content3"}
)

func TestMain(m *testing.M) {
	// 设置
	setup()
	// 运行测试
	code := m.Run()
	// 清理
	teardown()
	// 退出
	os.Exit(code)
}

func setup() {
	// 创建测试目录
	os.MkdirAll(testSrcDir, 0755)
	os.MkdirAll(testDestDir, 0755)

	// 创建测试文件
	for i, file := range testFiles {
		path := filepath.Join(testSrcDir, file)
		os.MkdirAll(filepath.Dir(path), 0755)
		os.WriteFile(path, []byte(testContents[i]), 0644)
	}
}

func teardown() {
	// 清理测试目录
	os.RemoveAll(testDataDir)
}

func TestTarGzCompressor_Compress(t *testing.T) {
	compressor := NewTarGzCompressor()
	archivePath := filepath.Join(testDestDir, "test.tar.gz")

	// 测试压缩目录
	err := compressor.Compress(testSrcDir, archivePath,
		WithCompressionLevel(6),
		WithBufferSize(1024),
		WithProgressCallback(func(current, total int64) {
			// 进度回调测试
		}),
	)
	if err != nil {
		t.Fatalf("压缩失败: %v", err)
	}

	// 验证压缩文件是否存在
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		t.Error("压缩文件未创建")
	}
}

func TestTarGzCompressor_Decompress(t *testing.T) {
	compressor := NewTarGzCompressor()
	archivePath := filepath.Join(testDestDir, "test.tar.gz")
	extractPath := filepath.Join(testDestDir, "extract")

	// 先压缩
	err := compressor.Compress(testSrcDir, archivePath)
	if err != nil {
		t.Fatalf("压缩失败: %v", err)
	}

	// 测试解压
	err = compressor.Decompress(archivePath, extractPath,
		WithBufferSize(1024),
		WithProgressCallback(func(current, total int64) {
			// 进度回调测试
		}),
	)
	if err != nil {
		t.Fatalf("解压失败: %v", err)
	}

	// 验证解压文件
	for i, file := range testFiles {
		path := filepath.Join(extractPath, file)
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("读取解压文件失败 %s: %v", file, err)
			continue
		}
		if string(content) != testContents[i] {
			t.Errorf("文件内容不匹配 %s: 期望 %s, 实际 %s", file, testContents[i], string(content))
		}
	}
}

func TestTarGzCompressor_CompressWithPatterns(t *testing.T) {
	compressor := NewTarGzCompressor()
	archivePath := filepath.Join(testDestDir, "test_patterns.tar.gz")
	extractPath := filepath.Join(testDestDir, "extract_patterns")

	// 测试包含/排除模式
	err := compressor.Compress(testSrcDir, archivePath,
		WithIncludePatterns([]string{"*.txt"}),
		WithExcludePatterns([]string{"dir1/*"}),
	)
	if err != nil {
		t.Fatalf("压缩失败: %v", err)
	}

	// 解压并验证
	err = compressor.Decompress(archivePath, extractPath)
	if err != nil {
		t.Fatalf("解压失败: %v", err)
	}

	// 验证文件
	for _, file := range []string{"file1.txt", "file2.txt"} {
		if _, err := os.Stat(filepath.Join(extractPath, file)); os.IsNotExist(err) {
			t.Errorf("文件应该存在: %s", file)
		}
	}

	// dir1/file3.txt 应该被排除
	if _, err := os.Stat(filepath.Join(extractPath, "dir1/file3.txt")); !os.IsNotExist(err) {
		t.Error("文件不应该存在: dir1/file3.txt")
	}
}

func TestCompressError(t *testing.T) {
	compressor := NewTarGzCompressor()

	// 测试压缩不存在的文件
	err := compressor.Compress("nonexistent", "output.tar.gz")
	if err == nil {
		t.Error("应该返回错误")
	}

	// 测试解压不存在的文件
	err = compressor.Decompress("nonexistent.tar.gz", "output")
	if err == nil {
		t.Error("应该返回错误")
	}

	// 验证错误类型
	var compressErr *CompressError
	if !IsCompressError(err, &compressErr) {
		t.Error("错误类型应该是 CompressError")
	}
}

// IsCompressError 检查错误是否为压缩错误
func IsCompressError(err error, target **CompressError) bool {
	if err == nil {
		return false
	}
	compressErr, ok := err.(*CompressError)
	if !ok {
		return false
	}
	if target != nil {
		*target = compressErr
	}
	return true
}
