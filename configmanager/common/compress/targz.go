package compress

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

// TarGzCompressor tar.gz格式压缩器
type TarGzCompressor struct{}

// NewTarGzCompressor 创建tar.gz格式压缩器
func NewTarGzCompressor() *TarGzCompressor {
	return &TarGzCompressor{}
}

// Compress 压缩文件或目录
func (c *TarGzCompressor) Compress(src string, dest string, opts ...Option) error {
	logDebug("开始压缩: %s -> %s", src, dest)
	options := processOptions(opts...)

	// 创建目标文件
	destFile, err := os.Create(dest)
	if err != nil {
		return NewCompressError("create", dest, err)
	}
	defer destFile.Close()

	// 创建gzip写入器
	gw, err := gzip.NewWriterLevel(destFile, options.CompressionLevel)
	if err != nil {
		return NewCompressError("gzip", "", err)
	}
	defer gw.Close()

	// 创建tar写入器
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// 获取源文件信息
	fi, err := os.Stat(src)
	if err != nil {
		return NewCompressError("stat", src, err)
	}

	// 如果是目录，遍历并添加文件
	if fi.IsDir() {
		err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 检查是否匹配包含/排除模式
			if !shouldInclude(path, src, options.IncludePatterns, options.ExcludePatterns) {
				return nil
			}

			// 获取相对路径
			relPath, err := filepath.Rel(src, path)
			if err != nil {
				return err
			}

			// 添加到tar文件
			return addToTar(tw, path, relPath, info, options.BufferSize, options.ProgressCallback)
		})
		if err != nil {
			return NewCompressError("walk", src, err)
		}
	} else {
		// 单个文件直接添加
		err = addToTar(tw, src, filepath.Base(src), fi, options.BufferSize, options.ProgressCallback)
		if err != nil {
			return NewCompressError("add", src, err)
		}
	}

	logDebug("压缩完成: %s -> %s", src, dest)
	return nil
}

// Decompress 解压文件
func (c *TarGzCompressor) Decompress(src string, dest string, opts ...Option) error {
	logDebug("开始解压: %s -> %s", src, dest)
	options := processOptions(opts...)

	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return NewCompressError("open", src, err)
	}
	defer srcFile.Close()

	// 创建gzip读取器
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return NewCompressError("gzip", src, err)
	}
	defer gr.Close()

	// 创建tar读取器
	tr := tar.NewReader(gr)

	// 确保目标目录存在
	if err := os.MkdirAll(dest, 0755); err != nil {
		return NewCompressError("mkdir", dest, err)
	}

	// 遍历tar文件
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return NewCompressError("read", src, err)
		}

		// 检查是否匹配包含/排除模式
		if !shouldInclude(header.Name, "", options.IncludePatterns, options.ExcludePatterns) {
			continue
		}

		// 构建目标路径
		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// 创建目录
			if err := os.MkdirAll(target, 0755); err != nil {
				return NewCompressError("mkdir", target, err)
			}
		case tar.TypeReg:
			// 创建文件
			if err := extractFile(tr, target, header.Size, options.BufferSize, options.ProgressCallback); err != nil {
				return NewCompressError("extract", target, err)
			}
		}
	}

	logDebug("解压完成: %s -> %s", src, dest)
	return nil
}

// addToTar 添加文件到tar
func addToTar(tw *tar.Writer, path, relPath string, info os.FileInfo, bufSize int, callback ProgressCallback) error {
	// 创建header
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = relPath

	// 写入header
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	// 如果是普通文件，写入内容
	if info.Mode().IsRegular() {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := copyWithProgress(tw, file, info.Size(), callback, bufSize); err != nil {
			return err
		}
	}

	return nil
}

// extractFile 解压文件
func extractFile(tr io.Reader, path string, size int64, bufSize int, callback ProgressCallback) error {
	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// 创建目标文件
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return copyWithProgress(file, tr, size, callback, bufSize)
}

// shouldInclude 检查文件是否应该包含
func shouldInclude(path, base string, includePatterns, excludePatterns []string) bool {
	// 如果没有包含和排除模式，则包含所有文件
	if len(includePatterns) == 0 && len(excludePatterns) == 0 {
		return true
	}

	// 获取相对路径
	relPath := path
	if base != "" {
		var err error
		relPath, err = filepath.Rel(base, path)
		if err != nil {
			return false
		}
	}

	// 检查排除模式
	for _, pattern := range excludePatterns {
		if matched, _ := filepath.Match(pattern, relPath); matched {
			return false
		}
	}

	// 如果没有包含模式，则包含所有未被排除的文件
	if len(includePatterns) == 0 {
		return true
	}

	// 检查包含模式
	for _, pattern := range includePatterns {
		if matched, _ := filepath.Match(pattern, relPath); matched {
			return true
		}
	}

	return false
}
