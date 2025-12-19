package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/internal/repository"
	"huaan-medical/pkg/database"
	"huaan-medical/pkg/errorcode"
)

// RoleService 角色服务
type RoleService struct {
	repo     *repository.RoleRepository
	permRepo *repository.PermissionRepository
}

func NewRoleService() *RoleService {
	return &RoleService{
		repo:     repository.NewRoleRepository(),
		permRepo: repository.NewPermissionRepository(),
	}
}

type ListRolesRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   *int   `form:"status"`
}

func (s *RoleService) List(req *ListRolesRequest) ([]model.RoleVO, int64, error) {
	roles, total, err := s.repo.List(req.Page, req.PageSize, req.Keyword, req.Status)
	if err != nil {
		return nil, 0, errorcode.New(errorcode.ErrDatabase)
	}

	list := make([]model.RoleVO, len(roles))
	for i := range roles {
		list[i] = *roles[i].ToVO()
	}
	return list, total, nil
}

type CreateRoleRequest struct {
	Code        string `json:"code" binding:"required,min=2,max=32"`
	Name        string `json:"name" binding:"required,min=2,max=64"`
	Description string `json:"description" binding:"max=256"`
	SortOrder   int    `json:"sort_order"`
	Status      *int   `json:"status"`
}

func (s *RoleService) Create(req *CreateRoleRequest) (*model.RoleVO, error) {
	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	if code == "" || name == "" {
		return nil, errorcode.New(errorcode.ErrInvalidParams)
	}

	// 检查编码是否存在
	if _, err := s.repo.GetByCode(code); err == nil {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "角色编码已存在")
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

	role := &model.Role{
		Code:        code,
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		SortOrder:   req.SortOrder,
		Status:      status,
	}

	if err := s.repo.Create(role); err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	created, err := s.repo.GetByID(role.ID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	return created.ToVO(), nil
}

type UpdateRoleRequest struct {
	Code        *string `json:"code"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
	Status      *int    `json:"status"`
}

func (s *RoleService) Update(roleID int64, req *UpdateRoleRequest) (*model.RoleVO, error) {
	_, err := s.repo.GetByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	updates := map[string]interface{}{}

	if req.Code != nil {
		code := strings.TrimSpace(*req.Code)
		if code == "" {
			return nil, errorcode.New(errorcode.ErrInvalidParams)
		}
		// 编码冲突检查
		if role, err := s.repo.GetByCode(code); err == nil && role.ID != roleID {
			return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "角色编码已存在")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrDatabase)
		}
		updates["code"] = code
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errorcode.New(errorcode.ErrInvalidParams)
		}
		updates["name"] = name
	}
	if req.Description != nil {
		updates["description"] = strings.TrimSpace(*req.Description)
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Status != nil {
		if *req.Status != model.StatusEnabled && *req.Status != model.StatusDisabled {
			return nil, errorcode.New(errorcode.ErrInvalidParams)
		}
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := database.GetDB().Transaction(func(tx *gorm.DB) error {
			roleRepo := repository.NewRoleRepositoryWithDB(tx)
			return roleRepo.Update(roleID, updates)
		}); err != nil {
			return nil, errorcode.New(errorcode.ErrDatabase)
		}
	}

	role, err := s.repo.GetByID(roleID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	return role.ToVO(), nil
}

type UpdateRolePermissionsRequest struct {
	PermissionIDs []int64 `json:"permission_ids"`
}

func (s *RoleService) UpdatePermissions(roleID int64, req *UpdateRolePermissionsRequest) error {
	_, err := s.repo.GetByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorcode.New(errorcode.ErrNotFound)
		}
		return errorcode.New(errorcode.ErrDatabase)
	}

	ids := make([]int64, 0, len(req.PermissionIDs))
	seen := map[int64]struct{}{}
	for _, id := range req.PermissionIDs {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}

	perms, err := s.permRepo.FindByIDs(ids)
	if err != nil {
		return errorcode.New(errorcode.ErrDatabase)
	}
	if len(ids) > 0 && len(perms) != len(ids) {
		return errorcode.New(errorcode.ErrInvalidParams)
	}

	if err := s.repo.ReplacePermissions(roleID, perms); err != nil {
		return errorcode.New(errorcode.ErrDatabase)
	}
	return nil
}
