package services_test

import (
	"fmt"
	"imgGeneratePrompts/config"
	"imgGeneratePrompts/models"
	"imgGeneratePrompts/services"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// PromptServiceTestSuite 是 PromptService 的测试套件
type PromptServiceTestSuite struct {
	suite.Suite
	db      *gorm.DB
	service *services.PromptService
	tagSvc  *services.TagService
}

// SetupSuite 在测试套件开始前运行
func (s *PromptServiceTestSuite) SetupSuite() {
	// 查找配置文件
	paths := []string{"../apikey/database.env", "./apikey/database.env"}
	var configPath string
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			configPath = path
			break
		}
	}
	if configPath == "" {
		s.T().Fatal("找不到数据库配置文件")
	}

	// 加载配置
	// 临时的配置加载方式，仅用于测试
	os.Setenv("DB_CONFIG_PATH", configPath)
	err := config.LoadConfig()
	if err != nil {
		s.T().Fatalf("加载配置失败: %v", err)
	}

	// 修改数据库名称用于测试
	testDBName := "img_generate_prompts_svctest"
	config.AppConfig.Database.DBName = testDBName

	// 创建测试数据库
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.Database.User, config.AppConfig.Database.Password, config.AppConfig.Database.Host, config.AppConfig.Database.Port)

	db, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{})
	s.Require().NoError(err, "连接MySQL失败")
	db.Exec("DROP DATABASE IF EXISTS " + testDBName)
	err = db.Exec("CREATE DATABASE " + testDBName).Error
	s.Require().NoError(err, "创建测试数据库失败")

	// 连接到测试数据库
	err = config.InitDBWithMigration()
	s.Require().NoError(err, "连接测试数据库失败")
	s.db = config.GetDB()
	s.service = services.NewPromptService()
	s.tagSvc = services.NewTagService()
}

// TearDownSuite 在测试套件结束后运行
func (s *PromptServiceTestSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	sqlDB.Close()

	// 删除测试数据库
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.Database.User, config.AppConfig.Database.Password, config.AppConfig.Database.Host, config.AppConfig.Database.Port)
	db, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{})
	s.Require().NoError(err)
	err = db.Exec("DROP DATABASE " + config.AppConfig.Database.DBName).Error
	s.Require().NoError(err, "删除测试数据库失败")
}

// SetupTest 在每个测试方法运行前清理数据库
func (s *PromptServiceTestSuite) SetupTest() {
	s.db.Exec("DELETE FROM prompt_tags")
	s.db.Exec("DELETE FROM prompts")
	s.db.Exec("DELETE FROM tags")
}

// TestCreateAndGetPrompt 测试创建和获取提示词
func (s *PromptServiceTestSuite) TestCreateAndGetPrompt() {
	req := &models.CreatePromptRequest{
		PromptText: "测试提示词",
		TagNames:   []string{"标签1", "标签2"},
	}

	prompt, err := s.service.CreatePromptWithImages(req)
	s.NoError(err)
	s.NotNil(prompt)
	s.Equal("测试提示词", prompt.PromptText)
	s.Len(prompt.Tags, 2)

	fetchedPrompt, err := s.service.GetPromptByID(prompt.ID)
	s.NoError(err)
	s.NotNil(fetchedPrompt)
	s.Equal(prompt.ID, fetchedPrompt.ID)
	s.Equal("测试提示词", fetchedPrompt.PromptText)
	s.Len(fetchedPrompt.Tags, 2, "获取的提示词应该包含标签信息")
}

// TestUpdatePrompt 测试更新提示词
func (s *PromptServiceTestSuite) TestUpdatePrompt() {
	// 1. 先创建一个
	req := &models.CreatePromptRequest{PromptText: "原始提示词", TagNames: []string{"原始标签"}}
	prompt, _ := s.service.CreatePromptWithImages(req)

	// 2. 更新它
	newText := "更新后的提示词"
	updateReq := &models.UpdatePromptRequest{
		PromptText: &newText,
		TagNames:   []string{"新标签1", "新标签2"},
	}
	updatedPrompt, err := s.service.UpdatePrompt(prompt.ID, updateReq)

	s.NoError(err)
	s.Equal(newText, updatedPrompt.PromptText)
	s.Len(updatedPrompt.Tags, 2)
	s.Equal("新标签1", updatedPrompt.Tags[0].Name)
}

// TestDeletePrompt 测试删除提示词 (软删除)
func (s *PromptServiceTestSuite) TestDeletePrompt() {
	req := &models.CreatePromptRequest{PromptText: "待删除的提示词"}
	prompt, _ := s.service.CreatePromptWithImages(req)

	err := s.service.DeletePrompt(prompt.ID)
	s.NoError(err)

	// 应该找不到了
	_, err = s.service.GetPromptByID(prompt.ID)
	s.Error(err)

	// 但在数据库中应该还存在 (软删除)
	var count int64
	s.db.Model(&models.Prompt{}).Unscoped().Where("id = ?", prompt.ID).Count(&count)
	s.Equal(int64(1), count)
}

// TestGetPrompts 测试复杂的列表查询
func (s *PromptServiceTestSuite) TestGetPrompts() {
	// 准备数据
	tagPublic, _ := s.tagSvc.CreateTag(&models.CreateTagRequest{Name: "公开"})
	s.service.CreatePromptWithImages(&models.CreatePromptRequest{PromptText: "公开提示词1", ModelName: "SD1.5", IsPublic: true, TagNames: []string{"公开"}})
	s.service.CreatePromptWithImages(&models.CreatePromptRequest{PromptText: "私有提示词1", ModelName: "SDXL", IsPublic: false})

	// 1. 测试按 IsPublic 过滤
	isPublic := true
	query := &models.PromptQuery{IsPublic: &isPublic, Page: 1, PageSize: 10}
	prompts, total, err := s.service.GetPrompts(query)
	s.NoError(err)
	s.Equal(int64(1), total)
	s.Equal("公开提示词1", prompts[0].PromptText)

	// 2. 测试按 ModelName 过滤
	query = &models.PromptQuery{ModelName: "SDXL", Page: 1, PageSize: 10}
	prompts, total, err = s.service.GetPrompts(query)
	s.NoError(err)
	s.Equal(int64(1), total)
	s.Equal("私有提示词1", prompts[0].PromptText)

	// 3. 测试按 TagNames 过滤
	query = &models.PromptQuery{TagNames: []string{tagPublic.Name}, Page: 1, PageSize: 10}
	prompts, total, err = s.service.GetPrompts(query)
	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(prompts[0].Tags, 1)
}

// TestAnalyzePromptData 测试AI分析模拟函数
func (s *PromptServiceTestSuite) TestAnalyzePromptData() {
	res, err := s.service.AnalyzePromptData("test", "test_model", "base64data", nil)
	s.NoError(err)
	s.NotEmpty(res.NegativePrompt)
	s.NotEmpty(res.StyleDescription)
	s.NotEmpty(res.TagNames)
	s.Contains(res.TagNames, "test_model")
}

// TestPromptService runs the test suite
func TestPromptService(t *testing.T) {
	suite.Run(t, new(PromptServiceTestSuite))
}
