package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/internal/repository"
	"huaan-medical/pkg/config"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/jwt"
	"huaan-medical/pkg/utils"
	"huaan-medical/pkg/wechat"
)

// UserService 用户服务
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// WeChatLoginRequest 微信登录请求
type WeChatLoginRequest struct {
	Code string `json:"code" binding:"required"` // 微信登录凭证
}

// WeChatLoginResponse 微信登录响应
type WeChatLoginResponse struct {
	Token     string       `json:"token"`
	ExpiresIn int64        `json:"expires_in"`
	User      *model.UserVO `json:"user"`
	IsNew     bool         `json:"is_new"` // 是否新用户
}

// UpdateUserInfoRequest 更新用户信息请求
type UpdateUserInfoRequest struct {
	Nickname string `json:"nickname" binding:"max=64"`
	Avatar   string `json:"avatar" binding:"max=512"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
	Gender   int    `json:"gender" binding:"oneof=0 1 2"`
}

// WeChatLogin 微信登录
func (s *UserService) WeChatLogin(req *WeChatLoginRequest, clientIP string) (*WeChatLoginResponse, error) {
	// 获取配置
	cfg, err := config.Load("config.yaml")
	if err != nil || cfg.WeChat.AppID == "" || cfg.WeChat.AppSecret == "" {
		return nil, errorcode.NewWithMessage(errorcode.ErrInternalServer, "微信登录功能未配置，请检查config.yaml中的wechat配置")
	}

	// 1. 调用微信API，使用code换取openid
	wechatClient := wechat.NewClient(cfg.WeChat.AppID, cfg.WeChat.AppSecret)
	session, err := wechatClient.Code2Session(req.Code)
	if err != nil {
		return nil, errorcode.NewWithMessage(errorcode.ErrInternalServer, "微信登录失败: "+err.Error())
	}

	if session.OpenID == "" {
		return nil, errorcode.NewWithMessage(errorcode.ErrInternalServer, "获取微信OpenID失败")
	}

	// 2. 根据openid查询或创建用户
	user, err := s.userRepo.GetByOpenID(session.OpenID)
	isNew := false

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新用户
			user = &model.User{
				OpenID:  session.OpenID,
				UnionID: session.UnionID,
				Status:  model.StatusEnabled,
			}
			if err := s.userRepo.Create(user); err != nil {
				return nil, errorcode.New(errorcode.ErrDatabase)
			}
			isNew = true
		} else {
			return nil, errorcode.New(errorcode.ErrDatabase)
		}
	}

	// 检查用户状态
	if user.Status == model.StatusDisabled {
		return nil, errorcode.New(errorcode.ErrAccountDisabled)
	}

	// 检查是否被封禁
	if user.IsBlocked() {
		return nil, errorcode.NewWithMessage(
			errorcode.ErrAccountBlocked,
			"您的账号因多次爽约已被暂时封禁，请联系客服",
		)
	}

	// 3. 生成JWT Token
	tokenPair, err := jwt.GenerateTokenPair(user.ID, user.OpenID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrInternalServer)
	}

	// 4. 更新登录信息（记录IP地址和登录时间）
	_ = s.userRepo.UpdateLoginInfo(user.ID, clientIP)

	// 重新查询用户以获取最新信息
	user, _ = s.userRepo.GetByID(user.ID)

	return &WeChatLoginResponse{
		Token:     tokenPair.AccessToken,
		ExpiresIn: int64(time.Until(tokenPair.AccessExpireAt).Seconds()),
		User:      user.ToVO(),
		IsNew:     isNew,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID int64) (*model.UserVO, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrUserNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 检查用户状态
	if user.Status == model.StatusDisabled {
		return nil, errorcode.New(errorcode.ErrAccountDisabled)
	}

	return user.ToVO(), nil
}

// UpdateUserInfo 更新用户信息
func (s *UserService) UpdateUserInfo(userID int64, req *UpdateUserInfoRequest) (*model.UserVO, error) {
	// 查询用户
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrUserNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 检查用户状态
	if user.Status == model.StatusDisabled {
		return nil, errorcode.New(errorcode.ErrAccountDisabled)
	}

	// 如果要绑定手机号，需要验证格式
	if req.Phone != "" {
		if !utils.ValidatePhone(req.Phone) {
			return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "手机号格式错误")
		}

		// 检查手机号是否已被其他用户绑定
		existingUser, err := s.userRepo.GetByPhone(req.Phone)
		if err == nil && existingUser.ID != userID {
			return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "该手机号已被其他用户绑定")
		}
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	user.Gender = req.Gender

	if err := s.userRepo.Update(user); err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 重新查询返回
	user, err = s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	return user.ToVO(), nil
}
