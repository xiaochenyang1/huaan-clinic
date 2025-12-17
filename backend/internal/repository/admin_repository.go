package repository

import (
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
