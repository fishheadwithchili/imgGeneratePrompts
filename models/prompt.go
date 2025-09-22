package models

import (
	"encoding/json"
	"strings"
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
	ID                    uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt             time.Time      `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt             time.Time      `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index;comment:软删除时间"` // 在JSON中隐藏
	PromptText            string         `json:"prompt_text" gorm:"type:text;not null;comment:正面提示词"`
	NegativePrompt        string         `json:"negative_prompt" gorm:"type:text;comment:负面提示词"`
	ModelName             string         `json:"model_name" gorm:"type:varchar(100);comment:使用的AI模型名称"`
	InputImageURL         string         `json:"input_image_url" gorm:"type:varchar(500);comment:输入的参照图片的存储路径或URL；可能多个图片"`
	OutputImageURL        string         `json:"output_image_url" gorm:"type:varchar(500);comment:输出的参照图片的存储路径或URL"`
	IsPublic              bool           `json:"is_public" gorm:"default:false;comment:是否公开"`
	StyleDescription      string         `json:"style_description" gorm:"type:varchar(500);comment:风格描述"`
	UsageScenario         string         `json:"usage_scenario" gorm:"type:varchar(500);comment:适用场景描述"`
	AtmosphereDescription string         `json:"atmosphere_description" gorm:"type:varchar(500);comment:氛围描述"`
	ExpressiveIntent      string         `json:"expressive_intent" gorm:"type:varchar(500);comment:表现意图描述"`
	// --- FIX: Changed type from string to json.RawMessage ---
	// This tells the JSON marshaller to treat this field as pre-formatted JSON
	// and embed it directly, avoiding double-encoding issues.
	StructureAnalysis json.RawMessage `json:"structure_analysis" gorm:"type:json;comment:提示词结构分析"`

	// 多对多关系字段
	Tags []*Tag `json:"tags" gorm:"many2many:prompt_tags;"`
}

// TableName 指定表名
func (Prompt) TableName() string {
	return "prompts"
}

// GetInputImageURLs 获取输入图片URL列表
func (p *Prompt) GetInputImageURLs() []string {
	if p.InputImageURL == "" {
		return []string{}
	}
	urls := strings.Split(p.InputImageURL, ",")
	result := make([]string, 0, len(urls))
	for _, url := range urls {
		if trimmed := strings.TrimSpace(url); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// SetInputImageURLs 设置输入图片URL列表
func (p *Prompt) SetInputImageURLs(urls []string) {
	if len(urls) == 0 {
		p.InputImageURL = ""
		return
	}
	// 过滤空字符串并组合
	filteredURLs := make([]string, 0, len(urls))
	for _, url := range urls {
		if trimmed := strings.TrimSpace(url); trimmed != "" {
			filteredURLs = append(filteredURLs, trimmed)
		}
	}
	p.InputImageURL = strings.Join(filteredURLs, ",")
}

// BeforeSave 在保存（创建或更新）前的钩子
func (p *Prompt) BeforeSave(tx *gorm.DB) (err error) {
	// --- FIX: Updated hook to handle json.RawMessage ([]byte) ---
	// Ensures StructureAnalysis is always a valid JSON object before saving.
	if len(p.StructureAnalysis) == 0 || string(p.StructureAnalysis) == `""` {
		p.StructureAnalysis = json.RawMessage("{}")
	} else {
		var js json.RawMessage
		if err := json.Unmarshal(p.StructureAnalysis, &js); err != nil {
			// If it's not valid JSON, set a default value to prevent database errors.
			p.StructureAnalysis = json.RawMessage("{}")
		}
	}
	return nil
}

// PromptResponse 响应结构体（用于API返回）
type PromptResponse struct {
	ID                    uint      `json:"id"`
	CreatedAt             time.Time `json:"created_at"`
	PromptText            string    `json:"prompt_text"`
	NegativePrompt        string    `json:"negative_prompt"`
	ModelName             string    `json:"model_name"`
	InputImageURLs        []string  `json:"input_image_urls"` // 解析后的输入图片URL数组
	OutputImageURL        string    `json:"output_image_url"`
	IsPublic              bool      `json:"is_public"`
	StyleDescription      string    `json:"style_description"`
	UsageScenario         string    `json:"usage_scenario"`
	AtmosphereDescription string    `json:"atmosphere_description"`
	ExpressiveIntent      string    `json:"expressive_intent"`
	// --- FIX: Changed type to json.RawMessage to match the model ---
	// This ensures the raw JSON is passed through to the final response correctly.
	StructureAnalysis json.RawMessage `json:"structure_analysis"`
	Tags              []*Tag          `json:"tags"`
}

// ToResponse 转换为响应结构体
func (p *Prompt) ToResponse() PromptResponse {
	return PromptResponse{
		ID:                    p.ID,
		CreatedAt:             p.CreatedAt,
		PromptText:            p.PromptText,
		NegativePrompt:        p.NegativePrompt,
		ModelName:             p.ModelName,
		InputImageURLs:        p.GetInputImageURLs(),
		OutputImageURL:        p.OutputImageURL,
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
	PromptText            string   `form:"prompt_text" json:"prompt_text" binding:"required"`
	NegativePrompt        string   `form:"negative_prompt" json:"negative_prompt"`
	ModelName             string   `form:"model_name" json:"model_name"`
	IsPublic              bool     `form:"is_public" json:"is_public"`
	StyleDescription      string   `form:"style_description" json:"style_description"`
	UsageScenario         string   `form:"usage_scenario" json:"usage_scenario"`
	AtmosphereDescription string   `form:"atmosphere_description" json:"atmosphere_description"`
	ExpressiveIntent      string   `form:"expressive_intent" json:"expressive_intent"`
	StructureAnalysis     string   `form:"structure_analysis" json:"structure_analysis"`
	InputImageURLs        []string `form:"input_image_urls" json:"input_image_urls"`
	OutputImageURL        string   `form:"output_image_url" json:"output_image_url"`
	TagNames              []string `form:"tag_names" json:"tag_names"`
}

// UpdatePromptRequest 更新提示词的请求结构体
type UpdatePromptRequest struct {
	PromptText            *string  `form:"prompt_text" json:"prompt_text"`
	NegativePrompt        *string  `form:"negative_prompt" json:"negative_prompt"`
	ModelName             *string  `form:"model_name" json:"model_name"`
	IsPublic              *bool    `form:"is_public" json:"is_public"`
	StyleDescription      *string  `form:"style_description" json:"style_description"`
	UsageScenario         *string  `form:"usage_scenario" json:"usage_scenario"`
	AtmosphereDescription *string  `form:"atmosphere_description" json:"atmosphere_description"`
	ExpressiveIntent      *string  `form:"expressive_intent" json:"expressive_intent"`
	StructureAnalysis     *string  `form:"structure_analysis" json:"structure_analysis"`
	InputImageURLs        []string `form:"input_image_urls" json:"input_image_urls"`
	OutputImageURL        *string  `form:"output_image_url" json:"output_image_url"`
	TagNames              []string `form:"tag_names" json:"tag_names"`
}

// AnalyzePromptRequest 分析请求结构体
type AnalyzePromptRequest struct {
	PromptText string `form:"prompt_text" binding:"required"`
	ModelName  string `form:"model_name"`
}

// AnalyzePromptResponse 智能生成响应结构体
type AnalyzePromptResponse struct {
	NegativePrompt        string   `json:"negative_prompt"`
	StyleDescription      string   `json:"style_description"`
	UsageScenario         string   `json:"usage_scenario"`
	AtmosphereDescription string   `json:"atmosphere_description"`
	ExpressiveIntent      string   `json:"expressive_intent"`
	StructureAnalysis     string   `json:"structure_analysis"`
	TagNames              []string `json:"tag_names"`
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
	Name string `form:"name" json:"name" binding:"required,max=100"`
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
