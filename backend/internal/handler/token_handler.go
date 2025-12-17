package handler

import (
	"github.com/gin-gonic/gin"

	"huaan-medical/internal/middleware"
	"huaan-medical/internal/service"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

// TokenHandler Token处理器
type TokenHandler struct {
	service *service.TokenService
}

// NewTokenHandler 创建Token处理器实例
func NewTokenHandler() *TokenHandler {
	return &TokenHandler{
		service: service.NewTokenService(),
	}
}

// IdempotentTokenResponse 幂等Token响应
type IdempotentTokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"` // 过期时间（秒）
}

// GetIdempotentToken 获取幂等Token
// @Summary 获取幂等Token
// @Description 获取用于防止重复提交的幂等Token
// @Tags Token
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=IdempotentTokenResponse}
// @Router /api/token/idempotent [get]
func (h *TokenHandler) GetIdempotentToken(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	token, expiresIn, err := h.service.GenerateIdempotentToken(userID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	resp := &IdempotentTokenResponse{
		Token:     token,
		ExpiresIn: expiresIn,
	}

	response.Success(c, resp)
}
