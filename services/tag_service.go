package services

import (
	"errors"
	"fmt"
	"imgGeneratePrompts/config"
	"imgGeneratePrompts/models"
	"strings"

	"gorm.io/gorm"
)

// TagService 标签服务
type TagService struct {
	db *gorm.DB
}

// NewTagService 创建标签服务实例
func NewTagService() *TagService {
	return &TagService{
		db: config.GetDB(),
	}
}

// CreateTag 创建标签
func (s *TagService) CreateTag(req *models.CreateTagRequest) (*models.Tag, error) {
	// 检查标签是否已存在
	var existingTag models.Tag
	if err := s.db.Where("name = ?", req.Name).First(&existingTag).Error; err == nil {
		return &existingTag, nil // 如果已存在，直接返回
	}

	tag := &models.Tag{
		Name: req.Name,
	}

	result := s.db.Create(tag)
	if result.Error != nil {
		return nil, fmt.Errorf("创建标签失败: %v", result.Error)
	}

	return tag, nil
}

// GetTagByID 根据ID获取标签
func (s *TagService) GetTagByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	result := s.db.First(&tag, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("标签不存在")
		}
		return nil, fmt.Errorf("获取标签失败: %v", result.Error)
	}
	return &tag, nil
}

// GetTagByName 根据名称获取标签
func (s *TagService) GetTagByName(name string) (*models.Tag, error) {
	var tag models.Tag
	result := s.db.Where("name = ?", name).First(&tag)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("标签不存在")
		}
		return nil, fmt.Errorf("获取标签失败: %v", result.Error)
	}
	return &tag, nil
}

// GetAllTags 获取所有标签
func (s *TagService) GetAllTags() ([]models.Tag, error) {
	var tags []models.Tag
	result := s.db.Order("name ASC").Find(&tags)
	if result.Error != nil {
		return nil, fmt.Errorf("获取标签列表失败: %v", result.Error)
	}
	return tags, nil
}

// GetOrCreateTags 获取或创建标签（批量）
func (s *TagService) GetOrCreateTags(tagNames []string) ([]*models.Tag, error) {
	if len(tagNames) == 0 {
		return []*models.Tag{}, nil
	}

	var tags []*models.Tag

	for _, name := range tagNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		// 尝试获取现有标签
		tag, err := s.GetTagByName(name)
		if err != nil {
			// 如果不存在，创建新标签
			tag, err = s.CreateTag(&models.CreateTagRequest{Name: name})
			if err != nil {
				return nil, fmt.Errorf("创建标签 %s 失败: %v", name, err)
			}
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// DeleteTag 删除标签
func (s *TagService) DeleteTag(id uint) error {
	// 检查是否有提示词使用此标签
	var count int64
	s.db.Table("prompt_tags").Where("tag_id = ?", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("无法删除标签，还有 %d 个提示词在使用此标签", count)
	}

	result := s.db.Delete(&models.Tag{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除标签失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("标签不存在")
	}
	return nil
}

// SearchTags 搜索标签
func (s *TagService) SearchTags(keyword string) ([]models.Tag, error) {
	var tags []models.Tag
	query := s.db.Model(&models.Tag{})

	if keyword != "" {
		keyword = "%" + strings.TrimSpace(keyword) + "%"
		query = query.Where("name LIKE ?", keyword)
	}

	result := query.Order("name ASC").Find(&tags)
	if result.Error != nil {
		return nil, fmt.Errorf("搜索标签失败: %v", result.Error)
	}

	return tags, nil
}

// GetTagStats 获取标签统计信息
func (s *TagService) GetTagStats() (map[string]interface{}, error) {
	var totalTags int64
	if err := s.db.Model(&models.Tag{}).Count(&totalTags).Error; err != nil {
		return nil, fmt.Errorf("获取标签总数失败: %v", err)
	}

	// 获取最受欢迎的标签（被使用次数最多）
	type TagUsage struct {
		TagID    uint   `json:"tag_id"`
		TagName  string `json:"tag_name"`
		UseCount int64  `json:"use_count"`
	}

	var popularTags []TagUsage
	err := s.db.Table("prompt_tags").
		Select("tag_id, tags.name as tag_name, COUNT(*) as use_count").
		Joins("LEFT JOIN tags ON tags.id = prompt_tags.tag_id").
		Group("tag_id").
		Order("use_count DESC").
		Limit(10).
		Scan(&popularTags).Error

	if err != nil {
		return nil, fmt.Errorf("获取热门标签失败: %v", err)
	}

	return map[string]interface{}{
		"total_tags":   totalTags,
		"popular_tags": popularTags,
	}, nil
}
