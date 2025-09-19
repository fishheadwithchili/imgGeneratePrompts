package services

import (
	"errors"
	"fmt"
	"imgGeneratePrompts/config"
	"imgGeneratePrompts/models"
	"strings"

	"gorm.io/gorm"
)

// PromptService 提示词服务
type PromptService struct {
	db         *gorm.DB
	tagService *TagService
}

// NewPromptService 创建提示词服务实例
func NewPromptService() *PromptService {
	return &PromptService{
		db:         config.GetDB(),
		tagService: NewTagService(),
	}
}

// CreatePrompt 创建提示词
func (s *PromptService) CreatePrompt(req *models.CreatePromptRequest, imageURL string) (*models.Prompt, error) {
	// 处理标签
	tags, err := s.tagService.GetOrCreateTags(req.TagNames)
	if err != nil {
		return nil, fmt.Errorf("处理标签失败: %v", err)
	}

	prompt := &models.Prompt{
		PromptText:            req.PromptText,
		NegativePrompt:        req.NegativePrompt,
		ModelName:             req.ModelName,
		ImageURL:              imageURL,
		IsPublic:              req.IsPublic,
		StyleDescription:      req.StyleDescription,
		UsageScenario:         req.UsageScenario,
		AtmosphereDescription: req.AtmosphereDescription,
		ExpressiveIntent:      req.ExpressiveIntent,
		StructureAnalysis:     req.StructureAnalysis,
		Tags:                  tags,
	}

	result := s.db.Create(prompt)
	if result.Error != nil {
		return nil, fmt.Errorf("创建提示词失败: %v", result.Error)
	}

	// 预加载标签信息
	if err := s.db.Preload("Tags").First(prompt, prompt.ID).Error; err != nil {
		return nil, fmt.Errorf("获取创建的提示词失败: %v", err)
	}

	return prompt, nil
}

// GetPromptByID 根据ID获取提示词
func (s *PromptService) GetPromptByID(id uint) (*models.Prompt, error) {
	var prompt models.Prompt
	result := s.db.Preload("Tags").First(&prompt, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("提示词不存在")
		}
		return nil, fmt.Errorf("获取提示词失败: %v", result.Error)
	}
	return &prompt, nil
}

// UpdatePrompt 更新提示词
func (s *PromptService) UpdatePrompt(id uint, req *models.UpdatePromptRequest) (*models.Prompt, error) {
	prompt, err := s.GetPromptByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	updates := make(map[string]interface{})

	if req.PromptText != nil {
		updates["prompt_text"] = *req.PromptText
	}
	if req.NegativePrompt != nil {
		updates["negative_prompt"] = *req.NegativePrompt
	}
	if req.ModelName != nil {
		updates["model_name"] = *req.ModelName
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	if req.StyleDescription != nil {
		updates["style_description"] = *req.StyleDescription
	}
	if req.UsageScenario != nil {
		updates["usage_scenario"] = *req.UsageScenario
	}
	if req.AtmosphereDescription != nil {
		updates["atmosphere_description"] = *req.AtmosphereDescription
	}
	if req.ExpressiveIntent != nil {
		updates["expressive_intent"] = *req.ExpressiveIntent
	}
	if req.StructureAnalysis != nil {
		updates["structure_analysis"] = *req.StructureAnalysis
	}

	if len(updates) > 0 {
		result := s.db.Model(prompt).Updates(updates)
		if result.Error != nil {
			return nil, fmt.Errorf("更新提示词失败: %v", result.Error)
		}
	}

	// 处理标签更新
	if req.TagNames != nil {
		tags, err := s.tagService.GetOrCreateTags(req.TagNames)
		if err != nil {
			return nil, fmt.Errorf("处理标签失败: %v", err)
		}

		// 替换关联的标签
		if err := s.db.Model(prompt).Association("Tags").Replace(tags); err != nil {
			return nil, fmt.Errorf("更新标签关联失败: %v", err)
		}
	}

	// 重新获取更新后的数据
	return s.GetPromptByID(id)
}

// DeletePrompt 删除提示词（软删除）
func (s *PromptService) DeletePrompt(id uint) error {
	result := s.db.Delete(&models.Prompt{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除提示词失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("提示词不存在")
	}
	return nil
}

// GetPrompts 获取提示词列表
func (s *PromptService) GetPrompts(query *models.PromptQuery) ([]models.Prompt, int64, error) {
	var prompts []models.Prompt
	var total int64

	// 构建查询
	db := s.db.Model(&models.Prompt{}).Preload("Tags")

	// 添加过滤条件
	if query.ModelName != "" {
		db = db.Where("model_name = ?", query.ModelName)
	}
	if query.IsPublic != nil {
		db = db.Where("is_public = ?", *query.IsPublic)
	}
	if query.Keyword != "" {
		keyword := "%" + strings.TrimSpace(query.Keyword) + "%"
		db = db.Where("prompt_text LIKE ? OR negative_prompt LIKE ? OR style_description LIKE ? OR usage_scenario LIKE ?",
			keyword, keyword, keyword, keyword)
	}

	// 标签过滤
	if len(query.TagNames) > 0 {
		db = db.Joins("JOIN prompt_tags ON prompts.id = prompt_tags.prompt_id").
			Joins("JOIN tags ON prompt_tags.tag_id = tags.id").
			Where("tags.name IN ?", query.TagNames).
			Group("prompts.id")
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取总数失败: %v", err)
	}

	// 排序
	orderBy := "created_at desc"
	if query.SortBy != "" {
		sortOrder := "desc"
		if query.SortOrder == "asc" {
			sortOrder = "asc"
		}
		switch query.SortBy {
		case "created_at":
			orderBy = fmt.Sprintf("%s %s", query.SortBy, sortOrder)
		}
	}
	db = db.Order(orderBy)

	// 分页
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// 执行查询
	if err := db.Find(&prompts).Error; err != nil {
		return nil, 0, fmt.Errorf("获取提示词列表失败: %v", err)
	}

	return prompts, total, nil
}

// GetPublicPrompts 获取公开的提示词列表
func (s *PromptService) GetPublicPrompts(page, pageSize int) ([]models.Prompt, int64, error) {
	isPublic := true
	query := &models.PromptQuery{
		Page:     page,
		PageSize: pageSize,
		IsPublic: &isPublic,
	}
	return s.GetPrompts(query)
}

// SearchPromptsByTags 根据标签搜索提示词
func (s *PromptService) SearchPromptsByTags(tagNames []string, page, pageSize int) ([]models.Prompt, int64, error) {
	query := &models.PromptQuery{
		Page:     page,
		PageSize: pageSize,
		TagNames: tagNames,
	}
	return s.GetPrompts(query)
}

// GetPromptStats 获取提示词统计信息
func (s *PromptService) GetPromptStats() (map[string]interface{}, error) {
	var totalPrompts, publicPrompts int64

	// 总提示词数量
	if err := s.db.Model(&models.Prompt{}).Count(&totalPrompts).Error; err != nil {
		return nil, fmt.Errorf("获取提示词总数失败: %v", err)
	}

	// 公开提示词数量
	if err := s.db.Model(&models.Prompt{}).Where("is_public = ?", true).Count(&publicPrompts).Error; err != nil {
		return nil, fmt.Errorf("获取公开提示词数量失败: %v", err)
	}

	// 按模型分组统计
	type ModelStats struct {
		ModelName string `json:"model_name"`
		Count     int64  `json:"count"`
	}

	var modelStats []ModelStats
	err := s.db.Model(&models.Prompt{}).
		Select("model_name, COUNT(*) as count").
		Where("model_name != ''").
		Group("model_name").
		Order("count DESC").
		Scan(&modelStats).Error

	if err != nil {
		return nil, fmt.Errorf("获取模型统计失败: %v", err)
	}

	return map[string]interface{}{
		"total_prompts":   totalPrompts,
		"public_prompts":  publicPrompts,
		"private_prompts": totalPrompts - publicPrompts,
		"model_stats":     modelStats,
	}, nil
}

// GetRecentPrompts 获取最近的提示词
func (s *PromptService) GetRecentPrompts(limit int) ([]models.Prompt, error) {
	var prompts []models.Prompt
	result := s.db.Preload("Tags").
		Where("is_public = ?", true).
		Order("created_at DESC").
		Limit(limit).
		Find(&prompts)

	if result.Error != nil {
		return nil, fmt.Errorf("获取最近提示词失败: %v", result.Error)
	}

	return prompts, nil
}

// DuplicateCheck 检查重复的提示词
func (s *PromptService) DuplicateCheck(promptText string) ([]models.Prompt, error) {
	var prompts []models.Prompt
	result := s.db.Preload("Tags").
		Where("prompt_text = ?", promptText).
		Find(&prompts)

	if result.Error != nil {
		return nil, fmt.Errorf("检查重复提示词失败: %v", result.Error)
	}

	return prompts, nil
}
