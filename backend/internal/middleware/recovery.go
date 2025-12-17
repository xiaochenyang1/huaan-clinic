package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/logger"
	"huaan-medical/pkg/response"
)

// Recovery panic恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误日志
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("ip", c.ClientIP()),
					zap.String("stack", string(debug.Stack())),
				)

				// 返回错误响应
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Code:    errorcode.ErrInternalServer,
					Message: errorcode.GetMessage(errorcode.ErrInternalServer),
				})
			}
		}()

		c.Next()
	}
}
