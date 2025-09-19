package controllers

import (
	"imgGeneratePrompts/config"
	"imgGeneratePrompts/models"
	"imgGeneratePrompts/services"
	"imgGeneratePrompts/utils"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// PromptController 提示词控制器
type PromptController struct {
	promptService *services.PromptService
}

// NewPromptController 创建提示词控制器实例
func NewPromptController() *PromptController {
	return &PromptController{
		promptService: services.NewPromptService(),
	}
}

// CreatePrompt 创建提示词
func (pc *PromptController) CreatePrompt(c *gin.Context) {
	var req models.CreatePromptRequest

	// 绑定JSON数据
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// TODO: 这里应该处理图片生成逻辑
	// 现在先用一个占位符URL
	imageURL := "/uploads/placeholder.jpg"

	// 创建提示词
	prompt, err := pc.promptService.CreatePrompt(&req, imageURL)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", prompt.ToResponse())
}

// UploadAndCreatePrompt 上传图片并创建提示词
func (pc *PromptController) UploadAndCreatePrompt(c *gin.Context) {
	// 解析表单数据
	var req models.CreatePromptRequest

	// 绑定表单数据
	if err := c.ShouldBind(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// 处理标签字符串（如果是通过表单提交的逗号分隔字符串）
	if tagNamesStr := c.PostForm("tag_names"); tagNamesStr != "" {
		req.TagNames = strings.Split(tagNamesStr, ",")
		for i := range req.TagNames {
			req.TagNames[i] = strings.TrimSpace(req.TagNames[i])
		}
	}

	// 处理文件上传
	file, err := c.FormFile("image")
	if err != nil {
		utils.BadRequestResponse(c, "图片上传失败: "+err.Error())
		return
	}

	// 检查文件大小
	if file.Size > config.AppConfig.Server.MaxFileSize {
		utils.BadRequestResponse(c, "文件大小超出限制")
		return
	}

	// 保存文件
	filename, err := utils.SaveUploadedFile(file, config.AppConfig.Server.UploadPath)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	// 生成图片URL
	imageURL := utils.GetFileURL(c, filename)

	// 创建提示词
	prompt, err := pc.promptService.CreatePrompt(&req, imageURL)
	if err != nil {
		// 如果创建失败，删除已上传的文件
		utils.DeleteFile(filepath.Join(config.AppConfig.Server.UploadPath, filename))
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", prompt.ToResponse())
}

// GetPrompt 获取单个提示词
func (pc *PromptController) GetPrompt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的ID")
		return
	}

	prompt, err := pc.promptService.GetPromptByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, prompt.ToResponse())
}

// UpdatePrompt 更新提示词
func (pc *PromptController) UpdatePrompt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的ID")
		return
	}

	var req models.UpdatePromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	prompt, err := pc.promptService.UpdatePrompt(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", prompt.ToResponse())
}

// DeletePrompt 删除提示词
func (pc *PromptController) DeletePrompt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的ID")
		return
	}

	if err := pc.promptService.DeletePrompt(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

// GetPrompts 获取提示词列表
func (pc *PromptController) GetPrompts(c *gin.Context) {
	var query models.PromptQuery

	// 绑定查询参数
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// 处理标签查询（支持逗号分隔的标签名）
	if tagNamesStr := c.Query("tag_names"); tagNamesStr != "" {
		query.TagNames = strings.Split(tagNamesStr, ",")
		for i := range query.TagNames {
			query.TagNames[i] = strings.TrimSpace(query.TagNames[i])
		}
	}

	// 设置默认值
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	prompts, total, err := pc.promptService.GetPrompts(&query)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	// 转换为响应格式
	responses := make([]models.PromptResponse, len(prompts))
	for i, prompt := range prompts {
		responses[i] = prompt.ToResponse()
	}

	utils.PaginationResponse(c, responses, query.Page, query.PageSize, total)
}

// GetPublicPrompts 获取公开的提示词列表
func (pc *PromptController) GetPublicPrompts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	prompts, total, err := pc.promptService.GetPublicPrompts(page, pageSize)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	// 转换为响应格式
	responses := make([]models.PromptResponse, len(prompts))
	for i, prompt := range prompts {
		responses[i] = prompt.ToResponse()
	}

	utils.PaginationResponse(c, responses, page, pageSize, total)
}

// SearchPromptsByTags 根据标签搜索提示词
func (pc *PromptController) SearchPromptsByTags(c *gin.Context) {
	tagNamesStr := c.Query("tags")
	if tagNamesStr == "" {
		utils.BadRequestResponse(c, "请提供标签名称")
		return
	}

	tagNames := strings.Split(tagNamesStr, ",")
	for i := range tagNames {
		tagNames[i] = strings.TrimSpace(tagNames[i])
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	prompts, total, err := pc.promptService.SearchPromptsByTags(tagNames, page, pageSize)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	// 转换为响应格式
	responses := make([]models.PromptResponse, len(prompts))
	for i, prompt := range prompts {
		responses[i] = prompt.ToResponse()
	}

	utils.PaginationResponse(c, responses, page, pageSize, total)
}

// GetPromptStats 获取提示词统计信息
func (pc *PromptController) GetPromptStats(c *gin.Context) {
	stats, err := pc.promptService.GetPromptStats()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, stats)
}

// GetRecentPrompts 获取最近的提示词
func (pc *PromptController) GetRecentPrompts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit > 50 {
		limit = 50 // 限制最大数量
	}

	prompts, err := pc.promptService.GetRecentPrompts(limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	// 转换为响应格式
	responses := make([]models.PromptResponse, len(prompts))
	for i, prompt := range prompts {
		responses[i] = prompt.ToResponse()
	}

	utils.SuccessResponse(c, responses)
}

// CheckDuplicate 检查重复提示词
func (pc *PromptController) CheckDuplicate(c *gin.Context) {
	promptText := c.Query("prompt_text")
	if promptText == "" {
		utils.BadRequestResponse(c, "请提供提示词文本")
		return
	}

	prompts, err := pc.promptService.DuplicateCheck(promptText)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	// 转换为响应格式
	responses := make([]models.PromptResponse, len(prompts))
	for i, prompt := range prompts {
		responses[i] = prompt.ToResponse()
	}

	result := map[string]interface{}{
		"is_duplicate": len(prompts) > 0,
		"count":        len(prompts),
		"prompts":      responses,
	}

	utils.SuccessResponse(c, result)
}
