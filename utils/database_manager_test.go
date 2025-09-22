package utils_test

import (
	"fmt"
	"imgGeneratePrompts/config"
	"imgGeneratePrompts/models"
	"imgGeneratePrompts/utils"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DatabaseManagerTestSuite 是我们的测试套件结构
type DatabaseManagerTestSuite struct {
	suite.Suite
	db      *gorm.DB
	manager *utils.DatabaseManager
	cfg     *config.Config
}

// SetupSuite 在所有测试运行之前执行
func (s *DatabaseManagerTestSuite) SetupSuite() {
	// 加载配置
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}
	s.cfg = config.AppConfig
	s.manager = utils.NewDatabaseManager()

	// --- 安全措施：创建一个专用的测试数据库 ---
	originalDBName := s.cfg.Database.DBName
	s.cfg.Database.DBName = originalDBName + "_test" // e.g., img_generate_prompts_test

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

	// 连接到新的测试数据库
	err = config.InitDB() // config.InitDB 会使用修改后的 s.cfg
	if err != nil {
		log.Fatalf("无法连接到测试数据库: %v", err)
	}
	s.db = config.GetDB()
}

// TearDownSuite 在所有测试运行之后执行
func (s *DatabaseManagerTestSuite) TearDownSuite() {
	// --- 安全措施：删除测试数据库，清理环境 ---
	dbName := s.cfg.Database.DBName
	if dbName != "" && dbName != "img_generate_prompts" { // 双重检查以防误删
		s.db.Exec("DROP DATABASE " + dbName)
	}
	config.CloseDB()
}

// SetupTest 在每个测试方法运行之前执行
func (s *DatabaseManagerTestSuite) SetupTest() {
	// 清理所有表，确保每个测试都在干净的环境中运行
	s.db.Migrator().DropTable(&models.Prompt{}, &models.Tag{})
	s.db.AutoMigrate(&models.Prompt{}, &models.Tag{})
}

// TestDatabaseManagerTestSuite 运行测试套件
func TestDatabaseManagerTestSuite(t *testing.T) {
	// 跳过短模式下的测试
	if testing.Short() {
		t.Skip("在短模式下跳过数据库测试")
	}
	suite.Run(t, new(DatabaseManagerTestSuite))
}

// TestInitializeDatabase 测试数据库初始化
func (s *DatabaseManagerTestSuite) TestInitializeDatabase() {
	// Arrange
	// 清理环境
	s.db.Migrator().DropTable(&models.Prompt{}, &models.Tag{})

	// Act
	err := s.manager.InitializeDatabase()

	// Assert
	assert.NoError(s.T(), err, "InitializeDatabase 应该不会返回错误")
	// 检查表是否已创建
	assert.True(s.T(), s.db.Migrator().HasTable(&models.Prompt{}), "应该已创建 'prompts' 表")
	assert.True(s.T(), s.db.Migrator().HasTable(&models.Tag{}), "应该已创建 'tags' 表")
}

// TestCreateSampleData 测试创建示例数据
func (s *DatabaseManagerTestSuite) TestCreateSampleData() {
	// Arrange
	// 确保表存在
	s.manager.InitializeDatabase()

	// Act
	err := s.manager.CreateSampleData()

	// Assert
	assert.NoError(s.T(), err, "CreateSampleData 应该不会返回错误")

	var promptCount, tagCount int64
	s.db.Model(&models.Prompt{}).Count(&promptCount)
	s.db.Model(&models.Tag{}).Count(&tagCount)

	assert.Greater(s.T(), promptCount, int64(0), "应该已创建示例提示词")
	assert.Greater(s.T(), tagCount, int64(0), "应该已创建示例标签")
}

// TestResetDatabase 测试重置数据库
func (s *DatabaseManagerTestSuite) TestResetDatabase() {
	// Arrange
	// 先填充数据
	s.manager.InitializeDatabase()
	s.manager.CreateSampleData()

	var promptCountBefore, tagCountBefore int64
	s.db.Model(&models.Prompt{}).Count(&promptCountBefore)
	s.db.Model(&models.Tag{}).Count(&tagCountBefore)
	assert.Greater(s.T(), promptCountBefore, int64(0)) // 确认有数据

	// Act
	err := s.manager.ResetDatabase()

	// Assert
	assert.NoError(s.T(), err, "ResetDatabase 应该不会返回错误")

	var promptCountAfter, tagCountAfter int64
	s.db.Model(&models.Prompt{}).Count(&promptCountAfter)
	s.db.Model(&models.Tag{}).Count(&tagCountAfter)

	assert.Equal(s.T(), int64(0), promptCountAfter, "重置后 'prompts' 表应该为空")
	// 注意：ResetDatabase 之后会重新插入初始标签
	assert.Greater(s.T(), tagCountAfter, int64(0), "重置后 'tags' 表应该包含初始标签")
}

// TestGetDatabaseStats 测试获取统计信息
func (s *DatabaseManagerTestSuite) TestGetDatabaseStats() {
	// Arrange
	// 填充数据
	s.manager.InitializeDatabase()
	s.manager.CreateSampleData()

	// Act
	stats, err := s.manager.GetDatabaseStats()

	// Assert
	assert.NoError(s.T(), err, "GetDatabaseStats 应该不会返回错误")
	assert.NotNil(s.T(), stats)

	promptsStats, ok := stats["prompts"].(map[string]interface{})
	assert.True(s.T(), ok)
	assert.Equal(s.T(), int64(3), promptsStats["total"], "示例数据应该有3个提示词")
}

// TestWriteDatabase 测试完整的写入流程
func (s *DatabaseManagerTestSuite) TestWriteDatabase() {
	// Arrange
	// 清理环境
	s.db.Migrator().DropTable(&models.Prompt{}, &models.Tag{})

	// Act
	err := s.manager.WriteDatabase()

	// Assert
	assert.NoError(s.T(), err, "WriteDatabase 应该不会返回错误")
	var promptCount, tagCount int64
	s.db.Model(&models.Prompt{}).Count(&promptCount)
	s.db.Model(&models.Tag{}).Count(&tagCount)
	assert.Greater(s.T(), promptCount, int64(0), "WriteDatabase 后应该有提示词数据")
	assert.Greater(s.T(), tagCount, int64(0), "WriteDatabase 后应该有标签数据")
}

// TestMain 覆盖 Go 原生的 TestMain，确保我们的测试套件被执行
func TestMain(m *testing.M) {
	// 这个函数是可选的，但可以用来在所有测试前后执行全局设置
	// 在这里我们使用 suite 的 Setup/TearDown，所以这个函数可以保持简单
	os.Exit(m.Run())
}
