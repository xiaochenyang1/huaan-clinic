package service

import (
	"errors"

	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/internal/repository"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/jwt"
	"huaan-medical/pkg/utils"
)

// AdminService 管理员服务
type AdminService struct {
	repo *repository.AdminRepository
	logRepo *repository.LogRepository
}

// NewAdminService 创建管理员服务实例
func NewAdminService() *AdminService {
	return &AdminService{
		repo:    repository.NewAdminRepository(),
		logRepo: repository.NewLogRepository(),
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=2,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string         `json:"token"`
	ExpiresIn int64          `json:"expires_in"`
	Admin     *model.AdminVO `json:"admin"`
}

// Login 管理员登录
func (s *AdminService) Login(req *LoginRequest, clientIP string, userAgent string) (*LoginResponse, error) {
	// 查询管理员
	admin, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.writeLoginLog(req.Username, 0, clientIP, userAgent, model.LoginStatusFailed, "管理员不存在")
			return nil, errorcode.New(errorcode.ErrAdminNotFound)
		}
		s.writeLoginLog(req.Username, 0, clientIP, userAgent, model.LoginStatusFailed, "数据库错误")
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 检查状态
	if admin.Status == model.StatusDisabled {
		s.writeLoginLog(admin.Username, admin.ID, clientIP, userAgent, model.LoginStatusFailed, "账号已禁用")
		return nil, errorcode.New(errorcode.ErrAccountDisabled)
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, admin.Password) {
		s.writeLoginLog(admin.Username, admin.ID, clientIP, userAgent, model.LoginStatusFailed, "密码错误")
		return nil, errorcode.New(errorcode.ErrPasswordWrong)
	}

	// 获取角色
	role := ""
	if len(admin.Roles) > 0 {
		role = admin.Roles[0].Code
	}

	// 生成Token
	token, err := jwt.GenerateAdminToken(admin.ID, admin.Username, role)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrInternalServer)
	}

	// 更新登录信息
	_ = s.repo.UpdateLoginInfo(admin.ID, clientIP)
	s.writeLoginLog(admin.Username, admin.ID, clientIP, userAgent, model.LoginStatusSuccess, "登录成功")

	return &LoginResponse{
		Token:     token,
		ExpiresIn: 7200, // 2小时
		Admin:     admin.ToVO(),
	}, nil
}

func (s *AdminService) writeLoginLog(username string, adminID int64, ip string, userAgent string, status int, msg string) {
	if s.logRepo == nil {
		return
	}

	device := userAgent
	if len(device) > 256 {
		device = device[:256]
	}

	_ = s.logRepo.CreateLoginLog(&model.LoginLog{
		UserType:  model.UserTypeAdmin,
		UserID:    adminID,
		Username:  username,
		LoginType: model.LoginTypePassword,
		IP:        ip,
		Device:    device,
		Status:    status,
		Message:   msg,
	})
}

// GetAdminInfo 获取管理员信息
func (s *AdminService) GetAdminInfo(adminID int64) (*model.AdminVO, error) {
	admin, err := s.repo.GetByID(adminID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrAdminNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	return admin.ToVO(), nil
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=32"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
}

// ChangePassword 修改密码
func (s *AdminService) ChangePassword(adminID int64, req *ChangePasswordRequest) error {
	admin, err := s.repo.GetByID(adminID)
	if err != nil {
		return errorcode.New(errorcode.ErrAdminNotFound)
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, admin.Password) {
		return errorcode.New(errorcode.ErrPasswordWrong)
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errorcode.New(errorcode.ErrInternalServer)
	}

	return s.repo.UpdatePassword(adminID, hashedPassword)
}
