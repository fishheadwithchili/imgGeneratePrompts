package controllers

import (
	"encoding/base64"
	"encoding/json"
	"imgGeneratePrompts/config"
	"imgGeneratePrompts/models"
	"imgGeneratePrompts/services"
	"imgGeneratePrompts/utils"
	"io/ioutil"
	"log"
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

	// 检查Content-Type并正确绑定数据
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/json") {
		// 绑定JSON数据
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ValidationErrorResponse(c, err)
			return
		}
	} else {
		// 尝试绑定表单数据
		if err := c.ShouldBind(&req); err != nil {
			utils.ValidationErrorResponse(c, err)
			return
		}
		// 处理表单中的tag_names字符串
		if tagNamesStr := c.PostForm("tag_names"); tagNamesStr != "" {
			req.TagNames = strings.Split(tagNamesStr, ",")
			for i := range req.TagNames {
				req.TagNames[i] = strings.TrimSpace(req.TagNames[i])
			}
		}
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

// UploadAndCreatePrompt 上传图片并创建提示词（支持多图片）
func (pc *PromptController) UploadAndCreatePrompt(c *gin.Context) {
	log.Println("--- DEBUG: [UploadAndCreatePrompt] Received a request to create prompt with image. ---")
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

	// --- DEBUG: 检查收到的 multipart form ---
	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("--- DEBUG: Error parsing multipart form: %v ---", err)
	}

	if form != nil && form.File != nil {
		log.Printf("--- DEBUG: Received multipart form with %d file fields.", len(form.File))
		for key, files := range form.File {
			log.Printf("--- DEBUG:   - Received field key: '%s' with %d file(s).", key, len(files))
			for i, fileHeader := range files {
				log.Printf("--- DEBUG:     - File #%d: Name='%s', Size=%d bytes", i+1, fileHeader.Filename, fileHeader.Size)
			}
		}
	} else {
		log.Println("--- DEBUG: No files received in multipart form. ---")
	}
	// --- END DEBUG ---

	// 处理图片上传（支持多个输入图片）
	if err == nil && form != nil {
		inputImageURLs := []string{}

		// 处理输入参考图片（支持多个）
		if files, ok := form.File["input_images"]; ok {
			log.Printf("--- DEBUG: Found '%d' files under the key 'input_images'. Processing them now.", len(files))
			for _, file := range files {
				if file.Size > config.AppConfig.Server.MaxFileSize {
					utils.BadRequestResponse(c, "文件大小超出限制: "+file.Filename)
					return
				}
				filename, err := utils.SaveUploadedFile(file, config.AppConfig.Server.UploadPath)
				if err != nil {
					utils.InternalServerErrorResponse(c, err.Error())
					return
				}
				imageURL := utils.GetFileURL(c, filename)
				inputImageURLs = append(inputImageURLs, imageURL)
			}
		}

		// 兼容旧版本：处理reference_images字段
		if files, ok := form.File["reference_images"]; ok {
			log.Printf("--- DEBUG: Found '%d' files under the legacy key 'reference_images'.", len(files))
			for _, file := range files {
				filename, err := utils.SaveUploadedFile(file, config.AppConfig.Server.UploadPath)
				if err != nil {
					utils.InternalServerErrorResponse(c, err.Error())
					return
				}
				imageURL := utils.GetFileURL(c, filename)
				inputImageURLs = append(inputImageURLs, imageURL)
			}
		}
		req.InputImageURLs = inputImageURLs

		// 处理输出图片（单个）
		if files, ok := form.File["output_image"]; ok && len(files) > 0 {
			file := files[0] // 只取第一个
			log.Printf("--- DEBUG: Found file under the key 'output_image': %s", file.Filename)
			if file.Size > config.AppConfig.Server.MaxFileSize {
				utils.BadRequestResponse(c, "输出图片文件大小超出限制")
				return
			}
			filename, err := utils.SaveUploadedFile(file, config.AppConfig.Server.UploadPath)
			if err != nil {
				utils.InternalServerErrorResponse(c, err.Error())
				return
			}
			req.OutputImageURL = utils.GetFileURL(c, filename)
		}

		// 兼容旧版本的单个image字段
		if files, ok := form.File["image"]; ok && len(files) > 0 {
			log.Printf("--- DEBUG: Found file under the legacy key 'image'.")
			file := files[0]
			if req.OutputImageURL == "" {
				filename, err := utils.SaveUploadedFile(file, config.AppConfig.Server.UploadPath)
				if err != nil {
					utils.InternalServerErrorResponse(c, err.Error())
					return
				}
				req.OutputImageURL = utils.GetFileURL(c, filename)
			}
		}
	}

	log.Printf("--- DEBUG: Request object before saving to DB: %+v", req)
	// 创建提示词
	prompt, err := pc.promptService.CreatePromptWithImages(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", prompt.ToResponse())
}

// AnalyzePrompt 智能生成接口 - 分析图片并返回建议内容
func (pc *PromptController) AnalyzePrompt(c *gin.Context) {
	var req models.AnalyzePromptRequest

	// 绑定表单数据
	if err := c.ShouldBind(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// 准备图片Base64编码
	var outputImageBase64 string
	var inputImageBase64 []string

	// 获取表单
	form, err := c.MultipartForm()
	if err != nil {
		utils.BadRequestResponse(c, "请提供图片文件")
		return
	}

	// 处理输出图片（必需）
	if files, ok := form.File["output_image"]; ok && len(files) > 0 {
		file := files[0]

		// 打开文件
		src, err := file.Open()
		if err != nil {
			utils.InternalServerErrorResponse(c, "无法读取输出图片")
			return
		}
		defer src.Close()

		// 读取文件内容
		data, err := ioutil.ReadAll(src)
		if err != nil {
			utils.InternalServerErrorResponse(c, "无法读取输出图片内容")
			return
		}

		// 转换为Base64
		outputImageBase64 = base64.StdEncoding.EncodeToString(data)
	} else {
		utils.BadRequestResponse(c, "请提供输出图片")
		return
	}

	// 处理输入参考图片（可选）
	if files, ok := form.File["input_images"]; ok {
		for _, file := range files {
			// 打开文件
			src, err := file.Open()
			if err != nil {
				continue
			}
			defer src.Close()

			// 读取文件内容
			data, err := ioutil.ReadAll(src)
			if err != nil {
				continue
			}

			// 转换为Base64
			base64Str := base64.StdEncoding.EncodeToString(data)
			inputImageBase64 = append(inputImageBase64, base64Str)
		}
	}

	// 兼容旧版本：处理reference_images字段
	if files, ok := form.File["reference_images"]; ok {
		for _, file := range files {
			// 打开文件
			src, err := file.Open()
			if err != nil {
				continue
			}
			defer src.Close()

			// 读取文件内容
			data, err := ioutil.ReadAll(src)
			if err != nil {
				continue
			}

			// 转换为Base64
			base64Str := base64.StdEncoding.EncodeToString(data)
			inputImageBase64 = append(inputImageBase64, base64Str)
		}
	}

	// 调用服务层进行AI分析
	response, err := pc.promptService.AnalyzePromptData(req.PromptText, req.ModelName, outputImageBase64, inputImageBase64)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "分析成功", response)
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
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ValidationErrorResponse(c, err)
			return
		}
	} else {
		if err := c.ShouldBind(&req); err != nil {
			utils.ValidationErrorResponse(c, err)
			return
		}
		// 处理表单中的tag_names字符串
		if tagNamesStr := c.PostForm("tag_names"); tagNamesStr != "" {
			tags := strings.Split(tagNamesStr, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			req.TagNames = tags
		}
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

	// --- DEBUG: Enhanced Logging ---
	log.Printf("准备转换 %d 条提示词为响应格式...", len(prompts))
	responses := make([]models.PromptResponse, len(prompts))
	for i, prompt := range prompts {
		responseObj := prompt.ToResponse()

		// 尝试序列化每个对象，并打印结果
		jsonBytes, jsonErr := json.Marshal(responseObj)
		if jsonErr != nil {
			log.Printf("!!!!!!!!!!!! 错误：序列化提示词 ID %d 时失败: %v", prompt.ID, jsonErr)
			// 打印出有问题的对象结构，帮助调试
			log.Printf("有问题的对象数据: %+v", responseObj)
			utils.InternalServerErrorResponse(c, "服务器在处理数据时遇到问题")
			return
		}

		// 打印序列化后的 JSON 字符串
		log.Printf("--- 序列化 ID %d 结果: %s", prompt.ID, string(jsonBytes))

		responses[i] = responseObj
	}
	log.Println("--- DEBUG: 所有提示词转换完成 ---")

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
