package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveUploadedFile 保存上传的文件
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	// 创建目标文件
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// 复制文件内容
	_, err = io.Copy(out, src)
	return err
}

// GenerateUniqueFilename 生成唯一的文件名
func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	name := originalFilename[:len(originalFilename)-len(ext)]
	
	// 使用时间戳和原文件名的 MD5 生成唯一文件名
	timestamp := time.Now().Format("20060102150405")
	hash := md5.Sum([]byte(name + timestamp))
	hashStr := hex.EncodeToString(hash[:])[:8]
	
	return fmt.Sprintf("%s_%s_%s%s", name, timestamp, hashStr, ext)
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// RemoveFile 删除文件（忽略错误）
func RemoveFile(path string) {
	_ = os.Remove(path)
}

// GetFileSize 获取文件大小
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// EnsureDirExists 确保目录存在
func EnsureDirExists(dir string) error {
	return os.MkdirAll(dir, 0755)
}
