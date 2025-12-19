package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"huaan-medical/internal/model"
	"huaan-medical/pkg/database"
	"huaan-medical/pkg/logger"
)

// responseWriter 自定义响应写入器，用于捕获响应内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装响应写入器
		blw := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)

		// 构建日志字段
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		}

		// 添加用户ID（如果存在）
		if userID := GetUserID(c); userID > 0 {
			fields = append(fields, zap.Int64("user_id", userID))
		}
		if adminID := GetAdminID(c); adminID > 0 {
			fields = append(fields, zap.Int64("admin_id", adminID))
		}

		// 添加错误信息
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		// 根据状态码选择日志级别
		status := c.Writer.Status()
		switch {
		case status >= 500:
			logger.Error("Server Error", fields...)
		case status >= 400:
			logger.Warn("Client Error", fields...)
		default:
			logger.Info("Request", fields...)
		}

		// 写入操作日志（仅记录管理后台接口）
		if strings.HasPrefix(path, "/api/admin") {
			writeOperationLogToDB(c, path, query, requestBody, blw.body.Bytes(), status, latency)
		}
	}
}

func writeOperationLogToDB(c *gin.Context, path, query string, requestBody, responseBody []byte, status int, latency time.Duration) {
	db := database.GetDB()
	if db == nil {
		return
	}

	adminID := GetAdminID(c)
	adminName := GetAdminUsername(c)

	module, action := inferModuleAction(c.Request.Method, path)

	bodyStr := string(requestBody)
	respStr := string(responseBody)
	errStr := ""
	if len(c.Errors) > 0 {
		errStr = c.Errors.String()
	}

	if isSensitivePath(path) {
		bodyStr = ""
		respStr = ""
	}

	bodyStr = truncateString(bodyStr, 2000)
	respStr = truncateString(respStr, 2000)
	errStr = truncateString(errStr, 1000)

	log := &model.OperationLog{
		AdminID:   adminID,
		AdminName: adminName,
		Module:    module,
		Action:    action,
		Method:    c.Request.Method,
		Path:      path,
		Query:     truncateString(query, 1000),
		Body:      bodyStr,
		Response:  respStr,
		IP:        c.ClientIP(),
		UserAgent: truncateString(c.Request.UserAgent(), 512),
		Status:    status,
		Latency:   latency.Milliseconds(),
		ErrorMsg:  errStr,
	}

	_ = db.Create(log).Error
}

func isSensitivePath(path string) bool {
	switch path {
	case "/api/admin/login",
		"/api/admin/password",
		"/api/user/login",
		"/api/user/login/password",
		"/api/user/login/phone",
		"/api/auth/refresh",
		"/api/sms/send":
		return true
	default:
		return false
	}
}

func inferModuleAction(method, path string) (module string, action string) {
	trimmed := strings.Trim(path, "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) >= 3 && parts[0] == "api" && parts[1] == "admin" {
		module = parts[2]
	}

	// 特殊动作
	if strings.Contains(path, "/export") {
		return module, "export"
	}
	if path == "/api/admin/login" {
		return "admin", "login"
	}

	switch method {
	case "GET":
		action = "view"
	case "POST":
		action = "create"
	case "PUT", "PATCH":
		action = "update"
	case "DELETE":
		action = "delete"
	default:
		action = "other"
	}
	return module, action
}

func truncateString(s string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}
