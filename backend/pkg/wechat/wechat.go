package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
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

	mu           sync.Mutex
	accessToken  string
	accessExpire time.Time
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

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

// GetAccessToken 获取小程序全局接口调用凭据（带简单内存缓存）
// 文档：https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-access-token/getAccessToken.html
func (c *Client) GetAccessToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 预留 60 秒缓冲
	if c.accessToken != "" && time.Now().Before(c.accessExpire.Add(-60*time.Second)) {
		return c.accessToken, nil
	}

	url := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		c.AppID,
		c.AppSecret,
	)

	resp, err := c.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("微信API请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result AccessTokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信API错误[%d]: %s", result.ErrCode, result.ErrMsg)
	}
	if result.AccessToken == "" {
		return "", fmt.Errorf("微信access_token为空")
	}

	c.accessToken = result.AccessToken
	c.accessExpire = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	return c.accessToken, nil
}

type SubscribeMessageRequest struct {
	ToUser     string                 `json:"touser"`
	TemplateID string                 `json:"template_id"`
	Page       string                 `json:"page,omitempty"`
	Data       map[string]interface{} `json:"data"`
}

type SubscribeMessageResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// SendSubscribeMessage 发送订阅消息
// 文档：https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-message-subscribe/sendMessage.html
func (c *Client) SendSubscribeMessage(req *SubscribeMessageRequest) error {
	if req == nil {
		return fmt.Errorf("请求不能为空")
	}
	if req.ToUser == "" || req.TemplateID == "" {
		return fmt.Errorf("touser/template_id不能为空")
	}

	accessToken, err := c.GetAccessToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s", accessToken)

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("构造请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("微信API请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	var result SubscribeMessageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}
	if result.ErrCode != 0 {
		return fmt.Errorf("微信API错误[%d]: %s", result.ErrCode, result.ErrMsg)
	}

	return nil
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
