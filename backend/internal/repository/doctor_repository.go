package repository

import (
	"huaan-medical/internal/model"
	"huaan-medical/pkg/database"

	"gorm.io/gorm"
)

// DoctorRepository 医生数据访问层
type DoctorRepository struct {
	db *gorm.DB
}

// NewDoctorRepository 创建医生仓库实例
func NewDoctorRepository() *DoctorRepository {
	return &DoctorRepository{db: database.GetDB()}
}

// Create 创建医生
func (r *DoctorRepository) Create(doctor *model.Doctor) error {
	return r.db.Create(doctor).Error
}

// Update 更新医生
func (r *DoctorRepository) Update(doctor *model.Doctor) error {
	return r.db.Save(doctor).Error
}

// Delete 删除医生（软删除）
func (r *DoctorRepository) Delete(id int64) error {
	return r.db.Delete(&model.Doctor{}, id).Error
}

// GetByID 根据ID查询医生
func (r *DoctorRepository) GetByID(id int64) (*model.Doctor, error) {
	var doctor model.Doctor
	err := r.db.Preload("Department").First(&doctor, id).Error
	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

// GetByIDSimple 根据ID查询医生（不预加载关联）
func (r *DoctorRepository) GetByIDSimple(id int64) (*model.Doctor, error) {
	var doctor model.Doctor
	err := r.db.First(&doctor, id).Error
	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

// List 分页查询医生列表（管理后台）
func (r *DoctorRepository) List(page, pageSize int, departmentID *int64, status *int, keyword string) ([]model.Doctor, int64, error) {
	var doctors []model.Doctor
	var total int64

	query := r.db.Model(&model.Doctor{}).Preload("Department")

	// 科室筛选
	if departmentID != nil && *departmentID > 0 {
		query = query.Where("department_id = ?", *departmentID)
	}

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 关键词搜索（姓名、职称、擅长领域）
	if keyword != "" {
		query = query.Where("name LIKE ? OR title LIKE ? OR specialty LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("sort_order ASC, id DESC").
		Offset(offset).Limit(pageSize).
		Find(&doctors).Error

	return doctors, total, err
}

// ListPublic 查询医生列表（公开接口）
func (r *DoctorRepository) ListPublic(departmentID *int64, keyword string) ([]model.Doctor, error) {
	var doctors []model.Doctor

	query := r.db.Model(&model.Doctor{}).Preload("Department").
		Where("status = ?", model.StatusEnabled)

	// 科室筛选
	if departmentID != nil && *departmentID > 0 {
		query = query.Where("department_id = ?", *departmentID)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR specialty LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Order("sort_order ASC, id DESC").Find(&doctors).Error
	return doctors, err
}

// HasAppointments 检查医生是否有预约记录
func (r *DoctorRepository) HasAppointments(id int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Appointment{}).
		Where("doctor_id = ?", id).
		Count(&count).Error
	return count > 0, err
}

// HasSchedules 检查医生是否有排班记录
func (r *DoctorRepository) HasSchedules(id int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Schedule{}).
		Where("doctor_id = ?", id).
		Count(&count).Error
	return count > 0, err
}

// CountByDepartment 统计科室下的医生数量
func (r *DoctorRepository) CountByDepartment(departmentID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Doctor{}).
		Where("department_id = ? AND status = ?", departmentID, model.StatusEnabled).
		Count(&count).Error
	return count, err
}

// GetByDepartment 根据科室ID获取医生列表
func (r *DoctorRepository) GetByDepartment(departmentID int64) ([]model.Doctor, error) {
	var doctors []model.Doctor
	err := r.db.Where("department_id = ? AND status = ?", departmentID, model.StatusEnabled).
		Order("sort_order ASC, id DESC").
		Find(&doctors).Error
	return doctors, err
}

// ExistsByName 检查医生姓名是否存在（同一科室下）
func (r *DoctorRepository) ExistsByName(name string, departmentID int64, excludeID ...int64) (bool, error) {
	query := r.db.Model(&model.Doctor{}).
		Where("name = ? AND department_id = ?", name, departmentID)

	// 排除指定ID（用于更新时的重复检查）
	if len(excludeID) > 0 && excludeID[0] > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// UpdateStatus 批量更新医生状态
func (r *DoctorRepository) UpdateStatus(ids []int64, status int) error {
	return r.db.Model(&model.Doctor{}).
		Where("id IN ?", ids).
		Update("status", status).Error
}
