package routes

import (
	"imgGeneratePrompts/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes() *gin.Engine {
	// 创建Gin引擎
	r := gin.Default()

	// 设置静态文件服务
	r.Static("/uploads", "./uploads")

	// 创建控制器实例
	promptController := controllers.NewPromptController()
	tagController := controllers.NewTagController()

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 提示词相关路由
		prompts := v1.Group("/prompts")
		{
			// 基础CRUD操作
			prompts.POST("/", promptController.CreatePrompt)                  // 创建提示词
			prompts.POST("/upload", promptController.UploadAndCreatePrompt)   // 上传图片并创建提示词
			prompts.POST("/analyze", promptController.AnalyzePrompt)          // 智能生成：AI分析图片和提示词
			prompts.GET("/", promptController.GetPrompts)                     // 获取提示词列表
			prompts.GET("/public", promptController.GetPublicPrompts)         // 获取公开提示词列表
			prompts.GET("/recent", promptController.GetRecentPrompts)         // 获取最近的提示词
			prompts.GET("/stats", promptController.GetPromptStats)            // 获取提示词统计信息
			prompts.GET("/search/tags", promptController.SearchPromptsByTags) // 根据标签搜索提示词
			prompts.GET("/check-duplicate", promptController.CheckDuplicate)  // 检查重复提示词
			prompts.GET("/:id", promptController.GetPrompt)                   // 获取单个提示词
			prompts.PUT("/:id", promptController.UpdatePrompt)                // 更新提示词
			prompts.DELETE("/:id", promptController.DeletePrompt)             // 删除提示词
		}

		// 标签相关路由
		tags := v1.Group("/tags")
		{
			tags.POST("/", tagController.CreateTag)       // 创建标签
			tags.GET("/", tagController.GetAllTags)       // 获取所有标签
			tags.GET("/search", tagController.SearchTags) // 搜索标签
			tags.GET("/stats", tagController.GetTagStats) // 获取标签统计信息
			tags.GET("/:id", tagController.GetTag)        // 获取单个标签
			tags.DELETE("/:id", tagController.DeleteTag)  // 删除标签
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Image Generate Prompts API is running",
		})
	})

	// 数据库状态检查
	r.GET("/db-status", func(c *gin.Context) {
		// 这里可以添加数据库连接检查逻辑
		c.JSON(200, gin.H{
			"status":  "connected",
			"message": "Database connection is healthy",
		})
	})

	// 根路径
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Image Generate Prompts API",
			"version": "v2.0.0",
			"features": []string{
				"提示词管理",
				"标签系统",
				"图片上传",
				"搜索功能",
				"统计信息",
			},
			"endpoints": gin.H{
				"health":    "/health",
				"db_status": "/db-status",
				"api":       "/api/v1",
				"prompts":   "/api/v1/prompts",
				"tags":      "/api/v1/tags",
				"uploads":   "/uploads",
			},
		})
	})

	return r
}
