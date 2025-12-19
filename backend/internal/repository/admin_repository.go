package repository

import (
	"strings"

	"huaan-medical/internal/model"
	"huaan-medical/pkg/database"

	"gorm.io/gorm"
)

// AdminRepository 管理员数据访问层
type AdminRepository struct {
	db *gorm.DB
}

// NewAdminRepository 创建管理员仓库实例
func NewAdminRepository() *AdminRepository {
	return &AdminRepository{db: database.GetDB()}
}

func NewAdminRepositoryWithDB(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// GetByUsername 根据用户名查询管理员
func (r *AdminRepository) GetByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.Preload("Roles").Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// GetByID 根据ID查询管理员
func (r *AdminRepository) GetByID(id int64) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.Preload("Roles").First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// GetByIDWithRolesPermissions 根据ID查询管理员（预加载角色与权限）
func (r *AdminRepository) GetByIDWithRolesPermissions(id int64) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.Preload("Roles.Permissions").First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// UpdateLoginInfo 更新登录信息
func (r *AdminRepository) UpdateLoginInfo(id int64, ip string) error {
	return r.db.Model(&model.Admin{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_at": gorm.Expr("NOW()"),
			"last_login_ip": ip,
		}).Error
}

// UpdatePassword 更新密码
func (r *AdminRepository) UpdatePassword(id int64, password string) error {
	return r.db.Model(&model.Admin{}).Where("id = ?", id).
		Update("password", password).Error
}

// Create 创建管理员
func (r *AdminRepository) Create(admin *model.Admin) error {
	return r.db.Create(admin).Error
}

// List 分页查询管理员列表
func (r *AdminRepository) List(page, pageSize int, keyword string, status *int) ([]model.Admin, int64, error) {
	var list []model.Admin
	var total int64

	query := r.db.Model(&model.Admin{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if strings.TrimSpace(keyword) != "" {
		like := "%" + strings.TrimSpace(keyword) + "%"
		query = query.Where(
			r.db.Where("username LIKE ?", like).
				Or("nickname LIKE ?", like).
				Or("phone LIKE ?", like).
				Or("email LIKE ?", like),
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Roles").Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// Update 更新管理员信息
func (r *AdminRepository) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&model.Admin{}).Where("id = ?", id).Updates(updates).Error
}

// ReplaceRoles 替换管理员角色（roleIDs 为空则清空）
func (r *AdminRepository) ReplaceRoles(adminID int64, roles []model.Role) error {
	var admin model.Admin
	if err := r.db.First(&admin, adminID).Error; err != nil {
		return err
	}
	return r.db.Model(&admin).Association("Roles").Replace(&roles)
}
