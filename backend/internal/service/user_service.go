package service

import (
	"errors"

	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/internal/repository"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/utils"
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
// TODO: 需要配置微信小程序 AppID 和 AppSecret
// TODO: 调用微信API: https://api.weixin.qq.com/sns/jscode2session
func (s *UserService) WeChatLogin(req *WeChatLoginRequest, clientIP string) (*WeChatLoginResponse, error) {
	// 这里需要实现：
	// 1. 调用微信API，使用code换取openid
	// 2. 根据openid查询或创建用户
	// 3. 生成JWT Token
	// 4. 更新登录信息

	// 暂时返回错误，提示需要配置
	return nil, errorcode.NewWithMessage(errorcode.ErrInternalServer, "微信登录功能需要配置小程序AppID和AppSecret")
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
