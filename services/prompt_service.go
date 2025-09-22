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

// CreatePrompt 创建提示词（兼容旧版本）
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
		OutputImageURL:        req.OutputImageURL,
		IsPublic:              req.IsPublic,
		StyleDescription:      req.StyleDescription,
		UsageScenario:         req.UsageScenario,
		AtmosphereDescription: req.AtmosphereDescription,
		ExpressiveIntent:      req.ExpressiveIntent,
		// --- FIX: Convert string from request to []byte for json.RawMessage ---
		StructureAnalysis: []byte(req.StructureAnalysis),
		Tags:              tags,
	}

	// 设置输入图片URLs
	prompt.SetInputImageURLs(req.InputImageURLs)

	// 如果没有指定输出图片，使用传入的imageURL作为输出图片
	if prompt.OutputImageURL == "" {
		prompt.OutputImageURL = imageURL
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

// CreatePromptWithImages 创建提示词（新版本，支持多图片）
func (s *PromptService) CreatePromptWithImages(req *models.CreatePromptRequest) (*models.Prompt, error) {
	// 处理标签
	tags, err := s.tagService.GetOrCreateTags(req.TagNames)
	if err != nil {
		return nil, fmt.Errorf("处理标签失败: %v", err)
	}

	prompt := &models.Prompt{
		PromptText:            req.PromptText,
		NegativePrompt:        req.NegativePrompt,
		ModelName:             req.ModelName,
		OutputImageURL:        req.OutputImageURL,
		IsPublic:              req.IsPublic,
		StyleDescription:      req.StyleDescription,
		UsageScenario:         req.UsageScenario,
		AtmosphereDescription: req.AtmosphereDescription,
		ExpressiveIntent:      req.ExpressiveIntent,
		// --- FIX: Convert string from request to []byte for json.RawMessage ---
		StructureAnalysis: []byte(req.StructureAnalysis),
		Tags:              tags,
	}

	// 设置输入图片URLs（多个图片以逗号分隔存储）
	prompt.SetInputImageURLs(req.InputImageURLs)

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

// AnalyzePromptData AI分析图片和提示词，返回建议内容
func (s *PromptService) AnalyzePromptData(promptText, modelName, outputImageBase64 string, referenceImagesBase64 []string) (*models.AnalyzePromptResponse, error) {
	// 这里应该调用真实的AI API（如Google Gemini、OpenAI等）
	// 现在返回模拟数据用于测试

	// 构建AI请求（示例）
	// aiPrompt := fmt.Sprintf(`
	// 请分析以下内容并返回JSON格式的建议：
	// 正向提示词：%s
	// 模型名称：%s
	// 输出图片：[Base64数据]
	// 参考图片数量：%d
	//
	// 请返回以下JSON格式的内容：
	// {
	//   "negative_prompt": "负面提示词",
	//   "style_description": "风格描述",
	//   "usage_scenario": "适用场景",
	//   "atmosphere_description": "氛围描述",
	//   "expressive_intent": "表现意图",
	//   "structure_analysis": "结构分析JSON",
	//   "tag_names": ["标签1", "标签2"]
	// }
	// `, promptText, modelName, len(referenceImagesBase64))

	// TODO: 实际调用AI API
	// response := callAIAPI(aiPrompt, outputImageBase64, referenceImagesBase64)

	// 返回模拟数据
	mockResponse := &models.AnalyzePromptResponse{
		NegativePrompt:        "ugly, blurry, low quality, watermark, text, distorted, deformed, bad anatomy",
		StyleDescription:      "数字艺术插画风格，色彩鲜艳，高对比度，富有想象力，细节丰富。",
		UsageScenario:         "适用于社交媒体帖子、博客文章配图、个人艺术项目、数字艺术展示。",
		AtmosphereDescription: "梦幻、超现实、充满活力的氛围，带有神秘和魔幻的色彩。",
		ExpressiveIntent:      "旨在通过超现实主义的视觉效果，激发观众的想象力和好奇心，传达一种奇幻的美感。",
		StructureAnalysis:     `{"主体":"根据提示词生成的核心元素","背景":"环境和场景描述","光照":"光源和光影效果","构图":"画面布局和元素排列","色彩":"主要色调和色彩搭配"}`,
		TagNames: []string{
			"数字艺术",
			"插画",
			"超现实主义",
			"鲜艳色彩",
			"奇幻",
			"创意设计",
		},
	}

	// 根据实际的提示词内容调整模拟数据
	if strings.Contains(strings.ToLower(promptText), "portrait") || strings.Contains(strings.ToLower(promptText), "人物") {
		mockResponse.TagNames = append(mockResponse.TagNames, "人物", "肖像")
		mockResponse.StyleDescription = "写实或半写实的人物肖像风格，注重细节表现和情感传达。"
	}

	if strings.Contains(strings.ToLower(promptText), "landscape") || strings.Contains(strings.ToLower(promptText), "风景") {
		mockResponse.TagNames = append(mockResponse.TagNames, "风景", "自然")
		mockResponse.StyleDescription = "自然风景画风格，强调大自然的美丽和壮观。"
	}

	if modelName != "" {
		mockResponse.TagNames = append(mockResponse.TagNames, modelName)
	}

	return mockResponse, nil
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
		// --- FIX: Convert string pointer from request to []byte for json.RawMessage ---
		updates["structure_analysis"] = []byte(*req.StructureAnalysis)
	}
	if req.OutputImageURL != nil {
		updates["output_image_url"] = *req.OutputImageURL
	}
	if len(req.InputImageURLs) > 0 {
		// 创建一个临时Prompt对象来使用SetInputImageURLs方法
		tempPrompt := &models.Prompt{}
		tempPrompt.SetInputImageURLs(req.InputImageURLs)
		updates["input_image_url"] = tempPrompt.InputImageURL
	}

	if len(updates) > 0 {
		result := s.db.Model(prompt).Updates(updates)
		if result.Error != nil {
			return nil, fmt.Errorf("更新提示词失败: %v", result.Error)
		}
	}

	// 处理标签更新
	if req.TagNames != nil && len(req.TagNames) > 0 {
		tags, err := s.tagService.GetOrCreateTags(req.TagNames)
		if err != nil {
			return nil, fmt.Errorf("处理标签失败: %v", err)
		}

		// 替换关联的标签
		if err := s.db.Model(prompt).Association("Tags").Replace(tags); err != nil {
			return nil, fmt.Errorf("更新标签关联失败: %v", err)
		}
	} else if req.TagNames != nil && len(req.TagNames) == 0 {
		// 如果传入空数组，清除所有标签
		if err := s.db.Model(prompt).Association("Tags").Clear(); err != nil {
			return nil, fmt.Errorf("清除标签关联失败: %v", err)
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
