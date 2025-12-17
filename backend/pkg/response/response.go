package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"huaan-medical/pkg/errorcode"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 分页数据结构
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errorcode.Success,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errorcode.Success,
		Message: message,
		Data:    data,
	})
}

// SuccessWithPage 分页数据成功响应
func SuccessWithPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    errorcode.Success,
		Message: "success",
		Data: PageData{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: errorcode.GetMessage(code),
	})
}

// FailWithMessage 带自定义消息的失败响应
func FailWithMessage(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// FailWithError 错误响应
func FailWithError(c *gin.Context, err error) {
	if appErr, ok := err.(*errorcode.AppError); ok {
		c.JSON(http.StatusOK, Response{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    errorcode.ErrInternalServer,
		Message: errorcode.GetMessage(errorcode.ErrInternalServer),
	})
}

// BadRequest 参数错误响应 (HTTP 400)
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    errorcode.ErrInvalidParams,
		Message: message,
	})
}

// Unauthorized 未授权响应 (HTTP 401)
func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    errorcode.ErrUnauthorized,
		Message: errorcode.GetMessage(errorcode.ErrUnauthorized),
	})
}

// Forbidden 禁止访问响应 (HTTP 403)
func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Code:    errorcode.ErrForbidden,
		Message: errorcode.GetMessage(errorcode.ErrForbidden),
	})
}

// NotFound 资源不存在响应 (HTTP 404)
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, Response{
		Code:    errorcode.ErrNotFound,
		Message: errorcode.GetMessage(errorcode.ErrNotFound),
	})
}

// InternalError 服务器内部错误响应 (HTTP 500)
func InternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    errorcode.ErrInternalServer,
		Message: errorcode.GetMessage(errorcode.ErrInternalServer),
	})
}
