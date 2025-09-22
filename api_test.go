package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"imgGeneratePrompts/config"
	"imgGeneratePrompts/routes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine
var testTagID uint
var testPromptID uint

func TestMain(m *testing.M) {
	// 设置测试环境
	setup()

	// 运行测试
	code := m.Run()

	// 清理测试数据
	teardown()

	os.Exit(code)
}

func setup() {
	gin.SetMode(gin.TestMode)

	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	if err := config.InitDBWithoutMigration(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 确保上传目录存在
	os.MkdirAll(config.AppConfig.Server.UploadPath, os.ModePerm)

	// 设置路由
	router = routes.SetupRoutes()
}

func teardown() {
	// 清理测试数据
	db := config.GetDB()

	// 删除测试创建的数据
	if testPromptID > 0 {
		db.Exec("DELETE FROM prompt_tags WHERE prompt_id = ?", testPromptID)
		db.Exec("DELETE FROM prompts WHERE id = ?", testPromptID)
	}

	if testTagID > 0 {
		db.Exec("DELETE FROM tags WHERE name LIKE 'test_%'")
	}

	// 关闭数据库
	config.CloseDB()
}

// 测试健康检查接口
func TestHealthCheck(t *testing.T) {
	w := performRequest("GET", "/health", nil)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "ok" {
		t.Errorf("健康检查失败")
	}
}

// 测试创建标签（JSON格式）
func TestCreateTag_JSON(t *testing.T) {
	body := map[string]string{
		"name": "test_tag_json",
	}

	w := performJSONRequest("POST", "/api/v1/tags/", body)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"] != float64(200) {
		t.Errorf("创建标签失败: %s", w.Body.String())
	}

	// 保存标签ID供后续测试使用
	if data, ok := response["data"].(map[string]interface{}); ok {
		if id, ok := data["id"].(float64); ok {
			testTagID = uint(id)
		}
	}
}

// 测试创建标签（表单格式）
func TestCreateTag_Form(t *testing.T) {
	formData := map[string]string{
		"name": "test_tag_form",
	}

	w := performFormRequest("POST", "/api/v1/tags/", formData)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", w.Code, w.Body.String())
	}
}

// 测试获取所有标签
func TestGetAllTags(t *testing.T) {
	w := performRequest("GET", "/api/v1/tags/", nil)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"] != float64(200) {
		t.Errorf("获取标签列表失败")
	}
}

// 测试创建提示词（JSON格式）
func TestCreatePrompt_JSON(t *testing.T) {
	body := map[string]interface{}{
		"prompt_text":     "A beautiful landscape",
		"negative_prompt": "ugly, blurry",
		"model_name":      "test_model",
		"is_public":       true,
		"tag_names":       []string{"test_tag1", "test_tag2"},
	}

	w := performJSONRequest("POST", "/api/v1/prompts/", body)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"] != float64(200) {
		t.Errorf("创建提示词失败: %s", w.Body.String())
	}

	// 保存提示词ID供后续测试使用
	if data, ok := response["data"].(map[string]interface{}); ok {
		if id, ok := data["id"].(float64); ok {
			testPromptID = uint(id)
		}
	}
}

// 测试创建提示词（表单格式）
func TestCreatePrompt_Form(t *testing.T) {
	formData := map[string]string{
		"prompt_text":     "A beautiful sunset",
		"negative_prompt": "dark, gloomy",
		"model_name":      "test_model_form",
		"is_public":       "true",
		"tag_names":       "sunset,nature,beautiful",
	}

	w := performFormRequest("POST", "/api/v1/prompts/", formData)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", w.Code, w.Body.String())
	}
}

// 测试上传图片并创建提示词
func TestUploadAndCreatePrompt(t *testing.T) {
	// 创建测试图片文件内容
	imageContent := []byte("fake image content for testing")

	fields := map[string]string{
		"prompt_text":     "Test prompt with image",
		"negative_prompt": "blurry",
		"model_name":      "test_model",
		"is_public":       "true",
		"tag_names":       "test,image,upload",
	}

	files := map[string][]byte{
		"output_image": imageContent,
	}

	w := performMultipartRequest("POST", "/api/v1/prompts/upload", fields, files)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", w.Code, w.Body.String())
	}
}

// 测试获取提示词列表
func TestGetPrompts(t *testing.T) {
	w := performRequest("GET", "/api/v1/prompts?page=1&page_size=10", nil)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"] != float64(200) {
		t.Errorf("获取提示词列表失败")
	}
}

// 测试更新提示词（JSON格式）
func TestUpdatePrompt_JSON(t *testing.T) {
	if testPromptID == 0 {
		t.Skip("没有可用的测试提示词ID")
	}

	body := map[string]interface{}{
		"prompt_text": "Updated prompt text",
		"tag_names":   []string{"updated_tag1", "updated_tag2"},
	}

	w := performJSONRequest("PUT", fmt.Sprintf("/api/v1/prompts/%d", testPromptID), body)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", w.Code, w.Body.String())
	}
}

// 测试更新提示词（表单格式）
func TestUpdatePrompt_Form(t *testing.T) {
	if testPromptID == 0 {
		t.Skip("没有可用的测试提示词ID")
	}

	formData := map[string]string{
		"prompt_text": "Updated via form",
		"tag_names":   "form_tag1,form_tag2",
	}

	w := performFormRequest("PUT", fmt.Sprintf("/api/v1/prompts/%d", testPromptID), formData)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", w.Code, w.Body.String())
	}
}

// 测试获取单个提示词
func TestGetPrompt(t *testing.T) {
	if testPromptID == 0 {
		t.Skip("没有可用的测试提示词ID")
	}

	w := performRequest("GET", fmt.Sprintf("/api/v1/prompts/%d", testPromptID), nil)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", w.Code, w.Body.String())
	}
}

// 测试检查重复提示词
func TestCheckDuplicate(t *testing.T) {
	w := performRequest("GET", "/api/v1/prompts/check-duplicate?prompt_text=test", nil)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", w.Code, w.Body.String())
	}
}

// 测试搜索标签
func TestSearchTags(t *testing.T) {
	w := performRequest("GET", "/api/v1/tags/search?keyword=test", nil)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", w.Code, w.Body.String())
	}
}

// 测试获取提示词统计信息
func TestGetPromptStats(t *testing.T) {
	w := performRequest("GET", "/api/v1/prompts/stats", nil)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", w.Code, w.Body.String())
	}
}

// 测试获取标签统计信息
func TestGetTagStats(t *testing.T) {
	w := performRequest("GET", "/api/v1/tags/stats", nil)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", w.Code, w.Body.String())
	}
}

// 测试分析提示词
func TestAnalyzePrompt(t *testing.T) {
	imageContent := []byte("fake image content")

	fields := map[string]string{
		"prompt_text": "Analyze this prompt",
		"model_name":  "test_model",
	}

	files := map[string][]byte{
		"output_image": imageContent,
	}

	w := performMultipartRequest("POST", "/api/v1/prompts/analyze", fields, files)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", w.Code, w.Body.String())
	}
}

// 辅助函数：执行普通请求
func performRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// 辅助函数：执行JSON请求
func performJSONRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	jsonBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// 辅助函数：执行表单请求
func performFormRequest(method, path string, formData map[string]string) *httptest.ResponseRecorder {
	form := make(map[string][]string)
	for key, value := range formData {
		form[key] = []string{value}
	}

	req, _ := http.NewRequest(method, path, bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = form

	// 构建表单数据
	formStr := ""
	for key, value := range formData {
		if formStr != "" {
			formStr += "&"
		}
		formStr += key + "=" + value
	}

	req, _ = http.NewRequest(method, path, bytes.NewBufferString(formStr))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// 辅助函数：执行multipart请求（文件上传）
func performMultipartRequest(method, path string, fields map[string]string, files map[string][]byte) *httptest.ResponseRecorder {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加普通字段
	for key, value := range fields {
		writer.WriteField(key, value)
	}

	// 添加文件
	for fieldname, content := range files {
		part, err := writer.CreateFormFile(fieldname, "test.png")
		if err != nil {
			panic(err)
		}
		part.Write(content)
	}

	writer.Close()

	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
