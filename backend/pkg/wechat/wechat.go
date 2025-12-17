package wechat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SessionResponse 微信登录凭证校验响应
type SessionResponse struct {
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符（需满足UnionID下发条件）
	ErrCode    int    `json:"errcode"`     // 错误码
	ErrMsg     string `json:"errmsg"`      // 错误信息
}

// Client 微信API客户端
type Client struct {
	AppID     string
	AppSecret string
	client    *http.Client
}

// NewClient 创建微信API客户端
func NewClient(appID, appSecret string) *Client {
	return &Client{
		AppID:     appID,
		AppSecret: appSecret,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Code2Session 登录凭证校验
// 通过 wx.login 接口获得临时登录凭证 code 后传到开发者服务器调用此接口完成登录流程
// 文档：https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
func (c *Client) Code2Session(code string) (*SessionResponse, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		c.AppID,
		c.AppSecret,
		code,
	)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("微信API请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result SessionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查微信返回的错误
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("微信API错误[%d]: %s", result.ErrCode, result.ErrMsg)
	}

	return &result, nil
}
