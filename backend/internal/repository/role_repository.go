package repository

import (
	"strings"

	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/pkg/database"
)

// RoleRepository 角色数据访问层
type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{db: database.GetDB()}
}

func NewRoleRepositoryWithDB(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetByID(id int64) (*model.Role, error) {
	var role model.Role
	if err := r.db.Preload("Permissions").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetByCode(code string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Preload("Permissions").Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindByIDs(ids []int64) ([]model.Role, error) {
	var roles []model.Role
	if len(ids) == 0 {
		return roles, nil
	}
	if err := r.db.Where("id IN ?", ids).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) List(page, pageSize int, keyword string, status *int) ([]model.Role, int64, error) {
	var list []model.Role
	var total int64

	query := r.db.Model(&model.Role{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if strings.TrimSpace(keyword) != "" {
		like := "%" + strings.TrimSpace(keyword) + "%"
		query = query.Where(
			r.db.Where("code LIKE ?", like).
				Or("name LIKE ?", like),
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Permissions").Order("sort_order ASC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *RoleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&model.Role{}).Where("id = ?", id).Updates(updates).Error
}

// ReplacePermissions 替换角色权限（permissions 为空则清空）
func (r *RoleRepository) ReplacePermissions(roleID int64, permissions []model.Permission) error {
	var role model.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}
	return r.db.Model(&role).Association("Permissions").Replace(&permissions)
}
