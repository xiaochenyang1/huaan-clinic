package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/internal/repository"
	"huaan-medical/pkg/database"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/utils"
)

// AdminManageService 管理员管理服务
type AdminManageService struct {
	repo     *repository.AdminRepository
	roleRepo *repository.RoleRepository
}

func NewAdminManageService() *AdminManageService {
	return &AdminManageService{
		repo:     repository.NewAdminRepository(),
		roleRepo: repository.NewRoleRepository(),
	}
}

type ListAdminsRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   *int   `form:"status"`
}

func (s *AdminManageService) List(req *ListAdminsRequest) ([]model.AdminVO, int64, error) {
	admins, total, err := s.repo.List(req.Page, req.PageSize, req.Keyword, req.Status)
	if err != nil {
		return nil, 0, errorcode.New(errorcode.ErrDatabase)
	}

	list := make([]model.AdminVO, len(admins))
	for i := range admins {
		list[i] = *admins[i].ToVO()
	}
	return list, total, nil
}

type CreateAdminRequest struct {
	Username string  `json:"username" binding:"required,min=4,max=20"`
	Password string  `json:"password" binding:"required,min=6,max=32"`
	Nickname string  `json:"nickname" binding:"max=64"`
	Phone    string  `json:"phone" binding:"max=20"`
	Email    string  `json:"email" binding:"max=128"`
	Status   *int    `json:"status"`
	RoleIDs  []int64 `json:"role_ids"`
}

func (s *AdminManageService) Create(req *CreateAdminRequest) (*model.AdminVO, error) {
	username := strings.TrimSpace(req.Username)
	if !utils.ValidateUsername(username) {
		return nil, errorcode.New(errorcode.ErrUsernameInvalid)
	}
	if len(strings.TrimSpace(req.Password)) < 6 {
		return nil, errorcode.New(errorcode.ErrPasswordInvalid)
	}

	// 检查用户名是否存在
	if _, err := s.repo.GetByUsername(username); err == nil {
		return nil, errorcode.New(errorcode.ErrUsernameExists)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	status := model.StatusEnabled
	if req.Status != nil {
		if *req.Status != model.StatusEnabled && *req.Status != model.StatusDisabled {
			return nil, errorcode.New(errorcode.ErrInvalidParams)
		}
		status = *req.Status
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrInternalServer)
	}

	var createdID int64
	err = database.GetDB().Transaction(func(tx *gorm.DB) error {
		adminRepo := repository.NewAdminRepositoryWithDB(tx)
		roleRepo := repository.NewRoleRepositoryWithDB(tx)

		admin := &model.Admin{
			Username: username,
			Password: hashedPassword,
			Nickname: strings.TrimSpace(req.Nickname),
			Phone:    strings.TrimSpace(req.Phone),
			Email:    strings.TrimSpace(req.Email),
			Status:   status,
		}
		if err := adminRepo.Create(admin); err != nil {
			return err
		}

		if req.RoleIDs != nil {
			roles, err := roleRepo.FindByIDs(req.RoleIDs)
			if err != nil {
				return err
			}
			if len(req.RoleIDs) > 0 && len(roles) != len(req.RoleIDs) {
				return errorcode.New(errorcode.ErrInvalidParams)
			}
			if err := adminRepo.ReplaceRoles(admin.ID, roles); err != nil {
				return err
			}
		}

		createdID = admin.ID
		return nil
	})
	if err != nil {
		if appErr, ok := err.(*errorcode.AppError); ok {
			return nil, appErr
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	admin, err := s.repo.GetByID(createdID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	return admin.ToVO(), nil
}

type UpdateAdminRequest struct {
	Nickname *string `json:"nickname"`
	Phone    *string `json:"phone"`
	Email    *string `json:"email"`
	Status   *int    `json:"status"`
	RoleIDs  []int64 `json:"role_ids"`
}

func (s *AdminManageService) Update(adminID int64, req *UpdateAdminRequest) (*model.AdminVO, error) {
	_, err := s.repo.GetByID(adminID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrAdminNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	updates := map[string]interface{}{}
	if req.Nickname != nil {
		updates["nickname"] = strings.TrimSpace(*req.Nickname)
	}
	if req.Phone != nil {
		updates["phone"] = strings.TrimSpace(*req.Phone)
	}
	if req.Email != nil {
		updates["email"] = strings.TrimSpace(*req.Email)
	}
	if req.Status != nil {
		if *req.Status != model.StatusEnabled && *req.Status != model.StatusDisabled {
			return nil, errorcode.New(errorcode.ErrInvalidParams)
		}
		updates["status"] = *req.Status
	}

	err = database.GetDB().Transaction(func(tx *gorm.DB) error {
		adminRepo := repository.NewAdminRepositoryWithDB(tx)
		roleRepo := repository.NewRoleRepositoryWithDB(tx)

		if len(updates) > 0 {
			if err := adminRepo.Update(adminID, updates); err != nil {
				return err
			}
		}

		if req.RoleIDs != nil {
			roles, err := roleRepo.FindByIDs(req.RoleIDs)
			if err != nil {
				return err
			}
			if len(req.RoleIDs) > 0 && len(roles) != len(req.RoleIDs) {
				return errorcode.New(errorcode.ErrInvalidParams)
			}
			if err := adminRepo.ReplaceRoles(adminID, roles); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		if appErr, ok := err.(*errorcode.AppError); ok {
			return nil, appErr
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	admin, err := s.repo.GetByID(adminID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	return admin.ToVO(), nil
}

type ResetAdminPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6,max=32"`
}

func (s *AdminManageService) ResetPassword(adminID int64, req *ResetAdminPasswordRequest) error {
	_, err := s.repo.GetByID(adminID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorcode.New(errorcode.ErrAdminNotFound)
		}
		return errorcode.New(errorcode.ErrDatabase)
	}

	if len(strings.TrimSpace(req.Password)) < 6 {
		return errorcode.New(errorcode.ErrPasswordInvalid)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errorcode.New(errorcode.ErrInternalServer)
	}

	if err := s.repo.UpdatePassword(adminID, hashedPassword); err != nil {
		return errorcode.New(errorcode.ErrDatabase)
	}
	return nil
}
