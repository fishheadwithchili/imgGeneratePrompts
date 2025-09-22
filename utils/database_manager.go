package utils

import (
	"fmt"
	"imgGeneratePrompts/config"
	"imgGeneratePrompts/models"
	"log"
)

// DatabaseManager 数据库管理工具
type DatabaseManager struct{}

// NewDatabaseManager 创建数据库管理工具实例
func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{}
}

// InitializeDatabase 初始化数据库 - 这是你要的"写数据库的方法"
func (dm *DatabaseManager) InitializeDatabase() error {
	log.Println("开始初始化数据库...")

	// 1. 创建数据库（如果不存在）
	if err := config.CreateDatabase(); err != nil {
		return fmt.Errorf("创建数据库失败: %v", err)
	}

	// 2. 初始化数据库连接并迁移表结构
	if err := config.InitDB(); err != nil {
		return fmt.Errorf("初始化数据库连接失败: %v", err)
	}

	log.Println("数据库初始化完成")
	return nil
}

// CreateSampleData 创建示例数据 (已优化)
func (dm *DatabaseManager) CreateSampleData() error {
	log.Println("开始创建示例数据...")

	db := config.GetDB()

	// 检查是否已有示例数据
	var promptCount int64
	db.Model(&models.Prompt{}).Count(&promptCount)
	if promptCount > 0 {
		log.Printf("检测到已有 %d 条提示词数据，跳过示例数据创建", promptCount)
		return nil
	}

	// 创建示例标签
	sampleTagNames := []string{
		"风景", "人物", "动物", "科幻", "复古",
		"现代", "暖色调", "冷色调", "高质量", "4K",
	}

	// 使用 map 来存储创建的标签，便于通过名字查找
	createdTagsMap := make(map[string]*models.Tag)
	for _, name := range sampleTagNames {
		tag := models.Tag{Name: name}
		if err := db.Create(&tag).Error; err != nil {
			log.Printf("创建标签 %s 失败: %v", name, err)
			continue
		}
		createdTagsMap[name] = &tag
	}

	// 创建示例提示词
	samplePrompts := []struct {
		Prompt   models.Prompt
		TagNames []string // ✨ 优化点：使用标签名而不是下标
	}{
		{
			Prompt: models.Prompt{
				PromptText:            "a beautiful sunset over mountains, golden hour, cinematic lighting, high quality",
				NegativePrompt:        "ugly, blurry, low quality, pixelated, noise",
				ModelName:             "stable-diffusion-v1-5",
				OutputImageURL:        "/uploads/sample_sunset.jpg",
				IsPublic:              true,
				StyleDescription:      "风景摄影风格，温暖的金色调",
				UsageScenario:         "适用于自然风光、旅游宣传、背景图片",
				AtmosphereDescription: "宁静、温暖、壮观的黄昏氛围",
				ExpressiveIntent:      "表现大自然的壮美和宁静",
				StructureAnalysis:     []byte(`{"主体":"山峰日落","光照":"黄金时刻","质量":"高质量","风格":"电影感"}`),
			},
			TagNames: []string{"风景", "暖色调", "高质量", "4K"},
		},
		{
			Prompt: models.Prompt{
				PromptText:            "portrait of a cat, professional photography, studio lighting, detailed fur texture",
				NegativePrompt:        "cartoon, anime, low resolution, distorted",
				ModelName:             "stable-diffusion-v1-5",
				OutputImageURL:        "/uploads/sample_cat.jpg",
				IsPublic:              true,
				StyleDescription:      "专业摄影风格，细致的毛发质感",
				UsageScenario:         "适用于宠物摄影、动物主题设计",
				AtmosphereDescription: "温馨、可爱、专业的摄影氛围",
				ExpressiveIntent:      "突出动物的可爱特征和毛发细节",
				StructureAnalysis:     []byte(`{"主体":"猫咪肖像","技法":"专业摄影","光照":"工作室灯光","细节":"毛发质感"}`),
			},
			TagNames: []string{"动物", "现代", "高质量"},
		},
		{
			Prompt: models.Prompt{
				PromptText:            "futuristic city skyline, neon lights, cyberpunk style, night scene, high-tech architecture",
				NegativePrompt:        "old, vintage, daylight, low quality",
				ModelName:             "stable-diffusion-xl",
				OutputImageURL:        "/uploads/sample_cyberpunk.jpg",
				IsPublic:              true,
				StyleDescription:      "赛博朋克风格，霓虹灯光效果",
				UsageScenario:         "适用于科幻题材、游戏背景、未来主题设计",
				AtmosphereDescription: "神秘、科技感十足的未来夜景",
				ExpressiveIntent:      "展现未来科技城市的繁华与神秘",
				StructureAnalysis:     []byte(`{"主体":"未来城市","风格":"赛博朋克","光效":"霓虹灯","时间":"夜景"}`),
			},
			TagNames: []string{"科幻", "现代", "冷色调", "高质量"},
		},
	}

	for i, sample := range samplePrompts {
		// ✨ 优化点：通过名字从 map 中查找并关联标签
		var tagsToAssociate []*models.Tag
		for _, tagName := range sample.TagNames {
			if tag, ok := createdTagsMap[tagName]; ok {
				tagsToAssociate = append(tagsToAssociate, tag)
			}
		}
		sample.Prompt.Tags = tagsToAssociate

		if err := db.Create(&sample.Prompt).Error; err != nil {
			log.Printf("创建示例提示词 %d 失败: %v", i+1, err)
			continue
		}
	}

	log.Println("示例数据创建完成")
	return nil
}

// ResetDatabase 重置数据库（危险操作）
func (dm *DatabaseManager) ResetDatabase() error {
	log.Println("警告：正在重置数据库！")
	return config.ResetDatabase()
}

// GetDatabaseStats 获取数据库统计信息
func (dm *DatabaseManager) GetDatabaseStats() (map[string]interface{}, error) {
	stats := config.GetTableStats()

	// 添加更多统计信息
	db := config.GetDB()

	// 获取公开/私有提示词统计
	var publicCount, privateCount int64
	db.Model(&models.Prompt{}).Where("is_public = ?", true).Count(&publicCount)
	db.Model(&models.Prompt{}).Where("is_public = ?", false).Count(&privateCount)

	// 获取最近7天的数据
	var recentCount int64
	db.Model(&models.Prompt{}).Where("created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)").Count(&recentCount)

	result := map[string]interface{}{
		"tables": stats,
		"prompts": map[string]interface{}{
			"total":         stats["prompts"],
			"public":        publicCount,
			"private":       privateCount,
			"recent_7_days": recentCount,
		},
		"tags": map[string]interface{}{
			"total": stats["tags"],
		},
	}

	return result, nil
}

// BackupData 备份数据（生成SQL）
func (dm *DatabaseManager) BackupData() (string, error) {
	// 这里可以实现数据备份逻辑
	// 返回SQL语句或保存到文件
	return "备份功能待实现", nil
}

// ValidateData 验证数据完整性
func (dm *DatabaseManager) ValidateData() error {
	db := config.GetDB()

	// 检查孤儿数据
	var orphanTags int64
	db.Table("prompt_tags").
		Joins("LEFT JOIN prompts ON prompts.id = prompt_tags.prompt_id").
		Where("prompts.id IS NULL").
		Count(&orphanTags)

	if orphanTags > 0 {
		return fmt.Errorf("发现 %d 个孤儿标签关联记录", orphanTags)
	}

	log.Println("数据完整性验证通过")
	return nil
}

// MigrateData 数据迁移工具
func (dm *DatabaseManager) MigrateData() error {
	// 这里可以实现数据迁移逻辑
	log.Println("数据迁移功能待实现")
	return nil
}

// WriteDatabase 写数据库的主要方法 - 这是你要的核心方法
func (dm *DatabaseManager) WriteDatabase() error {
	log.Println("开始写入数据库...")

	// 1. 初始化数据库
	if err := dm.InitializeDatabase(); err != nil {
		return fmt.Errorf("初始化数据库失败: %v", err)
	}

	// 2. 创建示例数据
	if err := dm.CreateSampleData(); err != nil {
		return fmt.Errorf("创建示例数据失败: %v", err)
	}

	// 3. 验证数据
	if err := dm.ValidateData(); err != nil {
		log.Printf("数据验证警告: %v", err)
		// 不返回错误，只记录警告
	}

	// 4. 显示统计信息
	stats, err := dm.GetDatabaseStats()
	if err != nil {
		log.Printf("获取统计信息失败: %v", err)
	} else {
		log.Printf("数据库统计信息: %+v", stats)
	}

	log.Println("数据库写入完成！")
	return nil
}
