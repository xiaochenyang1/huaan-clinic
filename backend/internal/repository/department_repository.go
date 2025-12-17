package repository

import (
	"huaan-medical/internal/model"
	"huaan-medical/pkg/database"

	"gorm.io/gorm"
)

// DepartmentRepository 科室数据访问层
type DepartmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository 创建科室仓库实例
func NewDepartmentRepository() *DepartmentRepository {
	return &DepartmentRepository{db: database.GetDB()}
}

// Create 创建科室
func (r *DepartmentRepository) Create(dept *model.Department) error {
	return r.db.Create(dept).Error
}

// Update 更新科室
func (r *DepartmentRepository) Update(dept *model.Department) error {
	return r.db.Save(dept).Error
}

// Delete 删除科室（软删除）
func (r *DepartmentRepository) Delete(id int64) error {
	return r.db.Delete(&model.Department{}, id).Error
}

// GetByID 根据ID查询科室
func (r *DepartmentRepository) GetByID(id int64) (*model.Department, error) {
	var dept model.Department
	err := r.db.First(&dept, id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

// GetByIDWithDoctors 根据ID查询科室（含医生）
func (r *DepartmentRepository) GetByIDWithDoctors(id int64) (*model.Department, error) {
	var dept model.Department
	err := r.db.Preload("Doctors").First(&dept, id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

// List 查询科室列表
func (r *DepartmentRepository) List(page, pageSize int, status *int) ([]model.Department, int64, error) {
	var departments []model.Department
	var total int64

	query := r.db.Model(&model.Department{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("sort_order ASC, id ASC").
		Offset(offset).Limit(pageSize).
		Find(&departments).Error

	return departments, total, err
}

// ListAll 查询所有启用的科室（公开接口用）
func (r *DepartmentRepository) ListAll() ([]model.Department, error) {
	var departments []model.Department
	err := r.db.Where("status = ?", model.StatusEnabled).
		Order("sort_order ASC, id ASC").
		Find(&departments).Error
	return departments, err
}

// HasDoctors 检查科室下是否有医生
func (r *DepartmentRepository) HasDoctors(id int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Doctor{}).Where("department_id = ?", id).Count(&count).Error
	return count > 0, err
}

// GetByName 根据名称查询科室
func (r *DepartmentRepository) GetByName(name string) (*model.Department, error) {
	var dept model.Department
	err := r.db.Where("name = ?", name).First(&dept).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}
