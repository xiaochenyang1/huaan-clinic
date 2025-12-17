package repository

import (
	"huaan-medical/internal/model"
	"huaan-medical/pkg/database"

	"gorm.io/gorm"
)

// PatientRepository 就诊人数据访问层
type PatientRepository struct {
	db *gorm.DB
}

// NewPatientRepository 创建就诊人仓库实例
func NewPatientRepository() *PatientRepository {
	return &PatientRepository{db: database.GetDB()}
}

// Create 创建就诊人
func (r *PatientRepository) Create(patient *model.Patient) error {
	return r.db.Create(patient).Error
}

// Update 更新就诊人
func (r *PatientRepository) Update(patient *model.Patient) error {
	return r.db.Save(patient).Error
}

// Delete 删除就诊人（软删除）
func (r *PatientRepository) Delete(id int64) error {
	return r.db.Delete(&model.Patient{}, id).Error
}

// GetByID 根据ID查询就诊人
func (r *PatientRepository) GetByID(id int64) (*model.Patient, error) {
	var patient model.Patient
	err := r.db.First(&patient, id).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

// GetByUserAndID 根据用户ID和就诊人ID查询（用于权限校验）
func (r *PatientRepository) GetByUserAndID(userID, patientID int64) (*model.Patient, error) {
	var patient model.Patient
	err := r.db.Where("user_id = ? AND id = ?", userID, patientID).First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

// ListByUser 根据用户ID查询就诊人列表
func (r *PatientRepository) ListByUser(userID int64) ([]model.Patient, error) {
	var patients []model.Patient
	err := r.db.Where("user_id = ?", userID).
		Order("is_default DESC, id DESC").
		Find(&patients).Error
	return patients, err
}

// GetDefaultByUser 获取用户的默认就诊人
func (r *PatientRepository) GetDefaultByUser(userID int64) (*model.Patient, error) {
	var patient model.Patient
	err := r.db.Where("user_id = ? AND is_default = 1", userID).First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

// ExistsByIDCard 检查身份证号是否已存在（用于同一用户下去重）
func (r *PatientRepository) ExistsByIDCard(userID int64, idCard string, excludeID ...int64) (bool, error) {
	query := r.db.Model(&model.Patient{}).
		Where("user_id = ? AND id_card = ?", userID, idCard)

	// 排除指定ID（用于更新时的检查）
	if len(excludeID) > 0 && excludeID[0] > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// CountByUser 统计用户的就诊人数量
func (r *PatientRepository) CountByUser(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Patient{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// ClearDefaultByUser 清除用户的所有默认就诊人标记（设置新默认时使用）
func (r *PatientRepository) ClearDefaultByUser(userID int64) error {
	return r.db.Model(&model.Patient{}).
		Where("user_id = ? AND is_default = 1", userID).
		Update("is_default", 0).Error
}

// SetDefault 设置默认就诊人（先清除旧的，再设置新的）
func (r *PatientRepository) SetDefault(userID, patientID int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 清除该用户的所有默认标记
		if err := tx.Model(&model.Patient{}).
			Where("user_id = ? AND is_default = 1", userID).
			Update("is_default", 0).Error; err != nil {
			return err
		}

		// 设置新的默认就诊人
		return tx.Model(&model.Patient{}).
			Where("id = ? AND user_id = ?", patientID, userID).
			Update("is_default", 1).Error
	})
}

// HasAppointments 检查就诊人是否有预约记录
func (r *PatientRepository) HasAppointments(id int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Appointment{}).
		Where("patient_id = ?", id).
		Count(&count).Error
	return count > 0, err
}

// List 分页查询患者列表（管理后台）
func (r *PatientRepository) List(page, pageSize int, keyword string) ([]model.Patient, int64, error) {
	var patients []model.Patient
	var total int64

	query := r.db.Model(&model.Patient{})

	// 关键词搜索（姓名、手机号）
	if keyword != "" {
		query = query.Where("name LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&patients).Error

	return patients, total, err
}
