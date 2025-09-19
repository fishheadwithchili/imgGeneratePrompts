package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SaveUploadedFile 保存上传的文件
func SaveUploadedFile(file *multipart.FileHeader, uploadPath string) (string, error) {
	// 创建上传目录
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("创建上传目录失败: %v", err)
	}

	// 检查文件类型
	if !IsValidImageType(file.Filename) {
		return "", fmt.Errorf("不支持的文件类型")
	}

	// 生成唯一文件名
	filename := GenerateUniqueFilename(file.Filename)
	filePath := filepath.Join(uploadPath, filename)

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开上传文件失败: %v", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("保存文件失败: %v", err)
	}

	return filename, nil
}

// IsValidImageType 检查是否为有效的图片类型
func IsValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

// GenerateUniqueFilename 生成唯一的文件名
func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	nameWithoutExt := strings.TrimSuffix(originalFilename, ext)

	// 清理文件名，移除特殊字符
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, " ", "_")
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, "(", "")
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, ")", "")

	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%d%s", nameWithoutExt, timestamp, ext)
}

// GetFileURL 获取文件的访问URL
func GetFileURL(c *gin.Context, filename string) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/uploads/%s", scheme, c.Request.Host, filename)
}

// DeleteFile 删除文件
func DeleteFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // 文件不存在，认为删除成功
	}
	return os.Remove(filePath)
}

// GetFileSize 获取文件大小
func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// FileExists 检查文件是否存在
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
