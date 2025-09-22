package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"imgGeneratePrompts/config"
	"imgGeneratePrompts/models"
	"imgGeneratePrompts/routes"
	"imgGeneratePrompts/utils"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// APITestSuite 是我们的API测试套件
type APITestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
	cfg    *config.Config

	// 用于在测试之间共享的ID
	testTagID    uint
	testPromptID uint
}

// SetupSuite 在所有测试运行之前执行
func (s *APITestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// 加载配置
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}
	s.cfg = config.AppConfig

	// --- 安全措施：创建一个专用的测试数据库 ---
	originalDBName := s.cfg.Database.DBName
	s.cfg.Database.DBName = originalDBName + "_apitest" // e.g., img_generate_prompts_apitest

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=%s&parseTime=true&loc=Local",
		s.cfg.Database.User, s.cfg.Database.Password, s.cfg.Database.Host, s.cfg.Database.Port, s.cfg.Database.Charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到MySQL服务器: %v", err)
	}

	// 删除可能存在的旧测试数据库
	db.Exec("DROP DATABASE IF EXISTS " + s.cfg.Database.DBName)
	// 创建新的测试数据库
	err = db.Exec("CREATE DATABASE " + s.cfg.Database.DBName).Error
	if err != nil {
		log.Fatalf("无法创建测试数据库: %v", err)
	}

	// 连接到新的测试数据库并执行迁移
	err = config.InitDBWithMigration()
	if err != nil {
		log.Fatalf("无法连接到测试数据库或执行迁移: %v", err)
	}
	s.db = config.GetDB()

	// 确保上传目录存在
	os.MkdirAll(s.cfg.Server.UploadPath, os.ModePerm)

	// 设置路由
	s.router = routes.SetupRoutes()
}

// TearDownSuite 在所有测试运行之后执行
func (s *APITestSuite) TearDownSuite() {
	// --- 安全措施：删除测试数据库和上传文件，清理环境 ---
	dbName := s.cfg.Database.DBName
	if dbName != "" && dbName != "img_generate_prompts" { // 双重检查以防误删
		s.db.Exec("DROP DATABASE " + dbName)
	}
	config.CloseDB()
	os.RemoveAll(s.cfg.Server.UploadPath)
}

// SetupTest 在每个测试方法运行之前执行
func (s *APITestSuite) SetupTest() {
	// 清理所有表，确保每个测试都在干净的环境中运行
	s.db.Exec("DELETE FROM prompt_tags")
	s.db.Exec("DELETE FROM prompts")
	s.db.Exec("DELETE FROM tags")
	// 重新插入初始标签
	config.GetDB().AutoMigrate(&models.Tag{}, &models.Prompt{})
}

// TestAPISuite 运行测试套件
func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip("在短模式下跳过API测试")
	}
	suite.Run(t, new(APITestSuite))
}

// performRequest 辅助函数，用于执行HTTP请求
func (s *APITestSuite) performRequest(method, path string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	return w
}

// TestSystemRoutes 测试系统级路由
func (s *APITestSuite) TestSystemRoutes() {
	w := s.performRequest("GET", "/health", nil, nil)
	assert.Equal(s.T(), http.StatusOK, w.Code)

	w = s.performRequest("GET", "/db-status", nil, nil)
	assert.Equal(s.T(), http.StatusOK, w.Code)

	w = s.performRequest("GET", "/", nil, nil)
	assert.Equal(s.T(), http.StatusOK, w.Code)
}

// TestTagsAPI 测试标签接口的完整生命周期
func (s *APITestSuite) TestTagsAPI() {
	// 1. Create a new tag
	createBody := bytes.NewBufferString(`{"name": "测试标签"}`)
	w := s.performRequest("POST", "/api/v1/tags/", createBody, map[string]string{"Content-Type": "application/json"})
	assert.Equal(s.T(), http.StatusOK, w.Code)
	var createResponse utils.ResponseData
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	tagData := createResponse.Data.(map[string]interface{})
	s.testTagID = uint(tagData["id"].(float64))
	assert.Equal(s.T(), "测试标签", tagData["name"])

	// 2. Get all tags
	w = s.performRequest("GET", "/api/v1/tags/", nil, nil)
	assert.Equal(s.T(), http.StatusOK, w.Code)

	// 3. Search tags
	w = s.performRequest("GET", "/api/v1/tags/search?keyword=测试", nil, nil)
	assert.Equal(s.T(), http.StatusOK, w.Code)

	// 4. Get tag stats
	w = s.performRequest("GET", "/api/v1/tags/stats", nil, nil)
	assert.Equal(s.T(), http.StatusOK, w.Code)

	// 5. Delete the tag
	deleteURL := fmt.Sprintf("/api/v1/tags/%d", s.testTagID)
	w = s.performRequest("DELETE", deleteURL, nil, nil)
	assert.Equal(s.T(), http.StatusOK, w.Code)
}

// TestPromptsAPI 测试提示词接口的完整生命周期
func (s *APITestSuite) TestPromptsAPI() {
	// 1. Create a prompt
	promptBody := `{"prompt_text": "一个宁静的湖景", "tag_names": ["风景", "宁静"]}`
	w := s.performRequest("POST", "/api/v1/prompts/", bytes.NewBufferString(promptBody), map[string]string{"Content-Type": "application/json"})
	assert.Equal(s.T(), http.StatusOK, w.Code)
	var createResponse utils.ResponseData
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	promptData := createResponse.Data.(map[string]interface{})
	s.testPromptID = uint(promptData["id"].(float64))
	assert.Equal(s.T(), "一个宁静的湖景", promptData["prompt_text"])

	// 2. Get all prompts
	w = s.performRequest("GET", "/api/v1/prompts/?page=1&page_size=10", nil, nil)
	assert.Equal(s.T(), http.StatusOK, w.Code)

	// 3. Get the specific prompt
	getURL := fmt.Sprintf("/api/v1/prompts/%d", s.testPromptID)
	w = s.performRequest("GET", getURL, nil, nil)
	assert.Equal(s.T(), http.StatusOK, w.Code)

	// 4. Update the prompt
	updateBody := `{"prompt_text": "一个更新后的湖景", "is_public": true}`
	w = s.performRequest("PUT", getURL, bytes.NewBufferString(updateBody), map[string]string{"Content-Type": "application/json"})
	assert.Equal(s.T(), http.StatusOK, w.Code)

	// 5. Check duplicate
	w = s.performRequest("GET", "/api/v1/prompts/check-duplicate?prompt_text=一个更新后的湖景", nil, nil)
	assert.Equal(s.T(), http.StatusOK, w.Code)

	// 6. Delete the prompt
	w = s.performRequest("DELETE", getURL, nil, nil)
	assert.Equal(s.T(), http.StatusOK, w.Code)
}

// TestUploadAndAnalyzeAPI 测试上传和分析接口
func (s *APITestSuite) TestUploadAndAnalyzeAPI() {
	// 创建一个假的图片文件内容
	imageContent := []byte("这是一个假的图片")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("output_image", "test.jpg")
	part.Write(imageContent)
	writer.WriteField("prompt_text", "测试上传")
	writer.WriteField("tag_names", "上传,测试")
	writer.Close()

	// 1. Test Upload
	w := s.performRequest("POST", "/api/v1/prompts/upload", body, map[string]string{"Content-Type": writer.FormDataContentType()})
	assert.Equal(s.T(), http.StatusOK, w.Code, "Upload应该成功")
	var uploadResponse utils.ResponseData
	json.Unmarshal(w.Body.Bytes(), &uploadResponse)
	promptID := uint(uploadResponse.Data.(map[string]interface{})["id"].(float64))

	// 2. Test Analyze
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	part, _ = writer.CreateFormFile("output_image", "test2.jpg")
	part.Write(imageContent)
	writer.WriteField("prompt_text", "测试分析")
	writer.Close()
	w = s.performRequest("POST", "/api/v1/prompts/analyze", body, map[string]string{"Content-Type": writer.FormDataContentType()})
	assert.Equal(s.T(), http.StatusOK, w.Code, "Analyze应该成功")

	// Cleanup
	s.db.Delete(&models.Prompt{}, promptID)
	os.RemoveAll(filepath.Join(s.cfg.Server.UploadPath))
}

// TestPromptSearchAndFilterAPI 测试提示词的搜索和过滤
func (s *APITestSuite) TestPromptSearchAndFilterAPI() {
	// 创建一些测试数据
	s.db.Create(&models.Prompt{PromptText: "公开的提示词", IsPublic: true})
	s.db.Create(&models.Prompt{PromptText: "最近的提示词"})

	// 1. Test Get Public
	w := s.performRequest("GET", "/api/v1/prompts/public", nil, nil)
	assert.Equal(s.T(), http.StatusOK, w.Code)

	// 2. Test Get Recent
	w = s.performRequest("GET", "/api/v1/prompts/recent?limit=5", nil, nil)
	assert.Equal(s.T(), http.StatusOK, w.Code)

	// 3. Test Get Stats
	w = s.performRequest("GET", "/api/v1/prompts/stats", nil, nil)
	assert.Equal(s.T(), http.StatusOK, w.Code)
}
