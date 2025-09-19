package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseData 统一响应数据结构
type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginationData 分页数据结构
type PaginationData struct {
	Items      interface{} `json:"items"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, ResponseData{
		Code:    code,
		Message: message,
	})
}

// BadRequestResponse 400错误响应
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

// UnauthorizedResponse 401错误响应
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message)
}

// ForbiddenResponse 403错误响应
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message)
}

// NotFoundResponse 404错误响应
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
}

// InternalServerErrorResponse 500错误响应
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message)
}

// PaginationResponse 分页响应
func PaginationResponse(c *gin.Context, items interface{}, page, pageSize int, total int64) {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	data := PaginationData{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	SuccessResponse(c, data)
}

// ValidationErrorResponse 参数验证错误响应
func ValidationErrorResponse(c *gin.Context, err error) {
	BadRequestResponse(c, "参数验证失败: "+err.Error())
}

// BindJSONError 绑定JSON错误处理
func BindJSONError(c *gin.Context, err error) {
	BadRequestResponse(c, "JSON格式错误: "+err.Error())
}
