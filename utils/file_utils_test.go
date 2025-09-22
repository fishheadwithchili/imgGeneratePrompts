package utils_test

import (
	"crypto/tls"
	"imgGeneratePrompts/utils"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGenerateUniqueFilename 测试生成唯一文件名
func TestGenerateUniqueFilename(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		expectedPart string
	}{
		{
			name:         "简单文件名",
			input:        "test.jpg",
			expectedPart: "test_",
		},
		{
			name:         "带空格的文件名",
			input:        "my beautiful photo.png",
			expectedPart: "my_beautiful_photo_",
		},
		{
			name:         "带括号的文件名",
			input:        "image(1).gif",
			expectedPart: "image1_",
		},
		{
			name:         "无后缀的文件名",
			input:        "myfile",
			expectedPart: "myfile_",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := utils.GenerateUniqueFilename(tc.input)
			assert.Contains(t, filename, tc.expectedPart, "文件名应包含预期的部分")

			originalExt := filepath.Ext(tc.input)
			if originalExt != "" {
				assert.True(t, strings.HasSuffix(filename, originalExt), "文件名应保留原始后缀")
			}
		})
	}

	// 测试唯一性
	filename1 := utils.GenerateUniqueFilename("unique.txt")
	time.Sleep(1 * time.Second) // 增加1秒延时确保时间戳不同
	filename2 := utils.GenerateUniqueFilename("unique.txt")
	assert.NotEqual(t, filename1, filename2, "两次生成的文件名应该不同")
}

// TestIsValidImageType 测试图片类型验证
func TestIsValidImageType(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		expected bool
	}{
		{name: "有效 JPG", filename: "image.jpg", expected: true},
		{name: "有效 JPEG", filename: "image.jpeg", expected: true},
		{name: "有效 PNG", filename: "image.png", expected: true},
		{name: "有效 GIF", filename: "image.gif", expected: true},
		{name: "有效 BMP", filename: "image.bmp", expected: true},
		{name: "有效 WEBP", filename: "image.webp", expected: true},
		{name: "有效 大写JPG", filename: "image.JPG", expected: true},
		{name: "无效 TXT", filename: "document.txt", expected: false},
		{name: "无效 PDF", filename: "report.pdf", expected: false},
		{name: "无后缀", filename: "image", expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.IsValidImageType(tc.filename)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestGetFileURL 测试获取文件URL
func TestGetFileURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("HTTP URL", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Host = "localhost:8080"

		url := utils.GetFileURL(c, "test.jpg")
		expected := "http://localhost:8080/uploads/test.jpg"
		assert.Equal(t, expected, url)
	})

	t.Run("HTTPS URL", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/test", nil)
		req.Host = "example.com"
		// 模拟一个简单的TLS状态
		req.TLS = &tls.ConnectionState{}
		c.Request = req

		url := utils.GetFileURL(c, "secure.png")
		expected := "https://example.com/uploads/secure.png"
		assert.Equal(t, expected, url)
	})
}
