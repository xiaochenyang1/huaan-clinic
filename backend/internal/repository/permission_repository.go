package repository

import (
	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/pkg/database"
)

// PermissionRepository 权限数据访问层
type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{db: database.GetDB()}
}

func NewPermissionRepositoryWithDB(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) FindByIDs(ids []int64) ([]model.Permission, error) {
	var list []model.Permission
	if len(ids) == 0 {
		return list, nil
	}
	if err := r.db.Where("id IN ?", ids).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PermissionRepository) ListAll() ([]model.Permission, error) {
	var list []model.Permission
	if err := r.db.Order("module ASC, sort_order ASC, id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
