package config

import (
	"fmt"
	"imgGeneratePrompts/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接（保持向后兼容，包含迁移）
func InitDB() error {
	return InitDBWithMigration()
}

// InitDBWithMigration 初始化数据库连接并执行迁移
func InitDBWithMigration() error {
	if err := connectToDatabase(); err != nil {
		return err
	}

	// 自动迁移数据表
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("数据表迁移失败: %v", err)
	}

	return nil
}

// InitDBWithoutMigration 初始化数据库连接但不执行迁移
func InitDBWithoutMigration() error {
	return connectToDatabase()
}

// connectToDatabase 连接到数据库的基础函数
func connectToDatabase() error {
	dsn := AppConfig.GetDSN()

	// 配置GORM日志
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 获取底层的sql.DB，配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置了连接可复用的最大时间

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	log.Println("数据库连接成功")
	return nil
}

// autoMigrate 自动迁移数据表
func autoMigrate() error {
	// 按顺序迁移所有模型
	err := DB.AutoMigrate(
		&models.Tag{},    // 先迁移标签表
		&models.Prompt{}, // 再迁移提示词表（包含外键关系）
	)

	if err != nil {
		return fmt.Errorf("自动迁移失败: %v", err)
	}

	log.Println("数据表迁移完成")

	// 插入初始数据
	if err := insertInitialData(); err != nil {
		log.Printf("插入初始数据失败: %v", err)
		// 不返回错误，因为可能是数据已存在
	}

	return nil
}

// insertInitialData 插入初始数据
func insertInitialData() error {
	// 检查是否已有数据
	var count int64
	DB.Model(&models.Tag{}).Count(&count)
	if count > 0 {
		log.Println("检测到已有标签数据，跳过初始化")
		return nil
	}

	// 创建一些初始标签
	initialTags := []models.Tag{
		{Name: "风景"},
		{Name: "人物"},
		{Name: "动物"},
		{Name: "建筑"},
		{Name: "抽象"},
		{Name: "科幻"},
		{Name: "复古"},
		{Name: "现代"},
		{Name: "暖色调"},
		{Name: "冷色调"},
	}

	for _, tag := range initialTags {
		if err := DB.Create(&tag).Error; err != nil {
			return fmt.Errorf("创建初始标签失败: %v", err)
		}
	}

	log.Println("初始标签数据插入成功")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("获取数据库实例失败: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("关闭数据库连接失败: %v", err)
		} else {
			log.Println("数据库连接已关闭")
		}
	}
}

// CreateDatabase 创建数据库（如果不存在）
func CreateDatabase() error {
	// 构建不包含数据库名的DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=%s&parseTime=true&loc=Local",
		AppConfig.Database.User,
		AppConfig.Database.Password,
		AppConfig.Database.Host,
		AppConfig.Database.Port,
		AppConfig.Database.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接MySQL服务器失败: %v", err)
	}

	// 创建数据库
	createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET %s COLLATE %s_unicode_ci",
		AppConfig.Database.DBName,
		AppConfig.Database.Charset,
		AppConfig.Database.Charset,
	)

	if err := db.Exec(createSQL).Error; err != nil {
		return fmt.Errorf("创建数据库失败: %v", err)
	}

	log.Printf("数据库 %s 创建成功或已存在", AppConfig.Database.DBName)

	// 关闭连接
	sqlDB, _ := db.DB()
	sqlDB.Close()

	return nil
}

// ResetDatabase 重置数据库（危险操作，仅用于开发环境）
func ResetDatabase() error {
	log.Println("警告：正在重置数据库，所有数据将被删除！")

	// 删除所有表
	if err := DB.Migrator().DropTable(&models.Prompt{}, &models.Tag{}); err != nil {
		return fmt.Errorf("删除表失败: %v", err)
	}

	// 重新创建表
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("重新创建表失败: %v", err)
	}

	log.Println("数据库重置完成")
	return nil
}

// GetTableStats 获取表统计信息
func GetTableStats() map[string]int64 {
	stats := make(map[string]int64)

	var promptCount, tagCount int64
	DB.Model(&models.Prompt{}).Count(&promptCount)
	DB.Model(&models.Tag{}).Count(&tagCount)

	stats["prompts"] = promptCount
	stats["tags"] = tagCount

	return stats
}
