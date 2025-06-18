package response

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
	Page    int         `json:"page,omitempty"`
	Size    int         `json:"size,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// SuccessWithPagination 带分页的成功响应
func SuccessWithPagination(c *gin.Context, data interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden 403错误
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

// NotFound 404错误
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// Conflict 409错误
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message)
}

// InternalServerError 500错误
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

// ValidationError 验证错误
func ValidationError(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "验证失败: "+message)
}

// DatabaseError 数据库错误
func DatabaseError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "数据库操作失败: "+message)
}

// FileUploadError 文件上传错误
func FileUploadError(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "文件上传失败: "+message)
}

// GetPaginationParams 获取分页参数
func GetPaginationParams(c *gin.Context) (page, size int) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, _ = strconv.Atoi(pageStr)
	size, _ = strconv.Atoi(sizeStr)

	// 参数验证
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	return page, size
}

// GetOffset 获取偏移量
func GetOffset(page, size int) int {
	return (page - 1) * size
} 