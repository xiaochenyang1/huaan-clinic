package repository

import (
	"huaan-medical/internal/model"
	"huaan-medical/pkg/database"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.GetDB()}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// GetByID 根据ID查询用户
func (r *UserRepository) GetByID(id int64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByOpenID 根据OpenID查询用户
func (r *UserRepository) GetByOpenID(openID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("open_id = ?", openID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPhone 根据手机号查询用户
func (r *UserRepository) GetByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateLoginInfo 更新登录信息
func (r *UserRepository) UpdateLoginInfo(id int64, ip string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_at": gorm.Expr("NOW()"),
			"last_login_ip": ip,
		}).Error
}

// GetByUsername 根据用户名查询用户
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsUsername 检查用户名是否已存在
func (r *UserRepository) ExistsUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsPhone 检查手机号是否已存在
func (r *UserRepository) ExistsPhone(phone string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("phone = ?", phone).Count(&count).Error
	return count > 0, err
}

// IncrementMissedCount 增加爽约次数
func (r *UserRepository) IncrementMissedCount(id int64) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		Update("missed_count", gorm.Expr("missed_count + 1")).Error
}

// BlockUser 封禁用户
func (r *UserRepository) BlockUser(id int64, until string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		Update("blocked_until", until).Error
}

// UnblockUser 解封用户
func (r *UserRepository) UnblockUser(id int64) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		Update("blocked_until", nil).Error
}
