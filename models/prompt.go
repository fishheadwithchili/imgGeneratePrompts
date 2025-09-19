package models

import (
	"time"

	"gorm.io/gorm"
)

// Tag 标签模型 - 对应 tags 表
type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(100);unique;not null;comment:标签名称"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	Prompts   []*Prompt `json:"-" gorm:"many2many:prompt_tags;"` // 在JSON中隐藏反向关联，避免循环引用
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

// Prompt 提示词模型 - 对应 prompts 表
type Prompt struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index;comment:软删除时间"` // 在JSON中隐藏
	PromptText     string         `json:"prompt_text" gorm:"type:text;not null;comment:正面提示词"`
	NegativePrompt string         `json:"negative_prompt" gorm:"type:text;comment:负面提示词"`
	ModelName      string         `json:"model_name" gorm:"type:varchar(100);comment:使用的AI模型名称"`
	ImageURL       string         `json:"image_url" gorm:"type:varchar(500);not null;comment:生成图片的存储路径或URL"`
	IsPublic       bool           `json:"is_public" gorm:"default:false;comment:是否公开"`

	// 新增描述字段
	StyleDescription      string `json:"style_description" gorm:"type:varchar(500);comment:风格描述"`
	UsageScenario         string `json:"usage_scenario" gorm:"type:varchar(500);comment:适用场景描述"`
	AtmosphereDescription string `json:"atmosphere_description" gorm:"type:varchar(500);comment:氛围描述"`
	ExpressiveIntent      string `json:"expressive_intent" gorm:"type:varchar(500);comment:表现意图描述"`
	StructureAnalysis     string `json:"structure_analysis" gorm:"type:json;comment:提示词结构分析"`

	// 多对多关系字段
	Tags []*Tag `json:"tags" gorm:"many2many:prompt_tags;"`
}

// TableName 指定表名
func (Prompt) TableName() string {
	return "prompts"
}

// BeforeCreate 创建前的钩子
func (p *Prompt) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// BeforeUpdate 更新前的钩子
func (p *Prompt) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

// PromptResponse 响应结构体（用于API返回）
type PromptResponse struct {
	ID                    uint      `json:"id"`
	CreatedAt             time.Time `json:"created_at"`
	PromptText            string    `json:"prompt_text"`
	NegativePrompt        string    `json:"negative_prompt"`
	ModelName             string    `json:"model_name"`
	ImageURL              string    `json:"image_url"`
	IsPublic              bool      `json:"is_public"`
	StyleDescription      string    `json:"style_description"`
	UsageScenario         string    `json:"usage_scenario"`
	AtmosphereDescription string    `json:"atmosphere_description"`
	ExpressiveIntent      string    `json:"expressive_intent"`
	StructureAnalysis     string    `json:"structure_analysis"`
	Tags                  []*Tag    `json:"tags"`
}

// ToResponse 转换为响应结构体
func (p *Prompt) ToResponse() PromptResponse {
	return PromptResponse{
		ID:                    p.ID,
		CreatedAt:             p.CreatedAt,
		PromptText:            p.PromptText,
		NegativePrompt:        p.NegativePrompt,
		ModelName:             p.ModelName,
		ImageURL:              p.ImageURL,
		IsPublic:              p.IsPublic,
		StyleDescription:      p.StyleDescription,
		UsageScenario:         p.UsageScenario,
		AtmosphereDescription: p.AtmosphereDescription,
		ExpressiveIntent:      p.ExpressiveIntent,
		StructureAnalysis:     p.StructureAnalysis,
		Tags:                  p.Tags,
	}
}

// CreatePromptRequest 创建提示词的请求结构体
type CreatePromptRequest struct {
	PromptText            string   `json:"prompt_text" binding:"required"`
	NegativePrompt        string   `json:"negative_prompt"`
	ModelName             string   `json:"model_name"`
	IsPublic              bool     `json:"is_public"`
	StyleDescription      string   `json:"style_description"`
	UsageScenario         string   `json:"usage_scenario"`
	AtmosphereDescription string   `json:"atmosphere_description"`
	ExpressiveIntent      string   `json:"expressive_intent"`
	StructureAnalysis     string   `json:"structure_analysis"`
	TagNames              []string `json:"tag_names"` // 标签名称列表
}

// UpdatePromptRequest 更新提示词的请求结构体
type UpdatePromptRequest struct {
	PromptText            *string  `json:"prompt_text"`
	NegativePrompt        *string  `json:"negative_prompt"`
	ModelName             *string  `json:"model_name"`
	IsPublic              *bool    `json:"is_public"`
	StyleDescription      *string  `json:"style_description"`
	UsageScenario         *string  `json:"usage_scenario"`
	AtmosphereDescription *string  `json:"atmosphere_description"`
	ExpressiveIntent      *string  `json:"expressive_intent"`
	StructureAnalysis     *string  `json:"structure_analysis"`
	TagNames              []string `json:"tag_names"` // 标签名称列表
}

// PromptQuery 查询参数结构体
type PromptQuery struct {
	Page      int      `form:"page" binding:"min=1"`
	PageSize  int      `form:"page_size" binding:"min=1,max=100"`
	ModelName string   `form:"model_name"`
	IsPublic  *bool    `form:"is_public"`
	Keyword   string   `form:"keyword"`
	TagNames  []string `form:"tag_names"`  // 标签名称过滤
	SortBy    string   `form:"sort_by"`    // created_at
	SortOrder string   `form:"sort_order"` // asc, desc
}

// CreateTagRequest 创建标签的请求结构体
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,max=100"`
}

// TagResponse 标签响应结构体
type TagResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse 转换为响应结构体
func (t *Tag) ToResponse() TagResponse {
	return TagResponse{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
	}
}
