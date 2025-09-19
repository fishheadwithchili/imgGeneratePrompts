package controllers

import (
	"imgGeneratePrompts/models"
	"imgGeneratePrompts/services"
	"imgGeneratePrompts/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TagController 标签控制器
type TagController struct {
	tagService *services.TagService
}

// NewTagController 创建标签控制器实例
func NewTagController() *TagController {
	return &TagController{
		tagService: services.NewTagService(),
	}
}

// CreateTag 创建标签
func (tc *TagController) CreateTag(c *gin.Context) {
	var req models.CreateTagRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	tag, err := tc.tagService.CreateTag(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", tag.ToResponse())
}

// GetTag 获取单个标签
func (tc *TagController) GetTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的ID")
		return
	}

	tag, err := tc.tagService.GetTagByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, tag.ToResponse())
}

// GetAllTags 获取所有标签
func (tc *TagController) GetAllTags(c *gin.Context) {
	tags, err := tc.tagService.GetAllTags()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	// 转换为响应格式
	responses := make([]models.TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = tag.ToResponse()
	}

	utils.SuccessResponse(c, responses)
}

// SearchTags 搜索标签
func (tc *TagController) SearchTags(c *gin.Context) {
	keyword := c.Query("keyword")

	tags, err := tc.tagService.SearchTags(keyword)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	// 转换为响应格式
	responses := make([]models.TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = tag.ToResponse()
	}

	utils.SuccessResponse(c, responses)
}

// DeleteTag 删除标签
func (tc *TagController) DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的ID")
		return
	}

	if err := tc.tagService.DeleteTag(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

// GetTagStats 获取标签统计信息
func (tc *TagController) GetTagStats(c *gin.Context) {
	stats, err := tc.tagService.GetTagStats()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, stats)
}
