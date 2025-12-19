package service

import (
	"errors"

	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/internal/rbac"
	"huaan-medical/internal/repository"
	"huaan-medical/pkg/errorcode"
)

// PermissionService 权限服务
type PermissionService struct {
	repo      *repository.PermissionRepository
	adminRepo *repository.AdminRepository
}

func NewPermissionService() *PermissionService {
	return &PermissionService{
		repo:      repository.NewPermissionRepository(),
		adminRepo: repository.NewAdminRepository(),
	}
}

func (s *PermissionService) ListAll() ([]model.PermissionVO, error) {
	perms, err := s.repo.ListAll()
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	list := make([]model.PermissionVO, len(perms))
	for i := range perms {
		list[i] = *perms[i].ToVO()
	}
	return list, nil
}

// GetAdminPermissionCodes 获取管理员权限码（包含角色权限汇总）
func (s *PermissionService) GetAdminPermissionCodes(adminID int64) ([]string, bool, error) {
	admin, err := s.adminRepo.GetByIDWithRolesPermissions(adminID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, errorcode.New(errorcode.ErrAdminNotFound)
		}
		return nil, false, errorcode.New(errorcode.ErrDatabase)
	}

	permSet := map[string]struct{}{}
	isSuperAdmin := false
	for _, role := range admin.Roles {
		if role.Code == rbac.RoleSuperAdmin {
			isSuperAdmin = true
		}
		for _, p := range role.Permissions {
			if p.Code == "" {
				continue
			}
			permSet[p.Code] = struct{}{}
		}
	}

	list := make([]string, 0, len(permSet))
	for code := range permSet {
		list = append(list, code)
	}
	return list, isSuperAdmin, nil
}
