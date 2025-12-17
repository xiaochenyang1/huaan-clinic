package repository

import (
	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/pkg/database"
)

// MedicalRecordRepository 就诊记录仓储
type MedicalRecordRepository struct {
	db *gorm.DB
}

// NewMedicalRecordRepository 创建就诊记录仓储实例
func NewMedicalRecordRepository() *MedicalRecordRepository {
	return &MedicalRecordRepository{
		db: database.GetDB(),
	}
}

// GetByID 根据ID获取就诊记录
func (r *MedicalRecordRepository) GetByID(id int64) (*model.MedicalRecord, error) {
	var record model.MedicalRecord
	err := r.db.
		Preload("Appointment").
		Preload("Patient").
		Preload("Doctor").
		Preload("Department").
		First(&record, id).Error
	return &record, err
}

// GetByUserAndID 根据用户ID和记录ID获取就诊记录（权限校验）
func (r *MedicalRecordRepository) GetByUserAndID(userID, recordID int64) (*model.MedicalRecord, error) {
	var record model.MedicalRecord
	err := r.db.
		Preload("Appointment").
		Preload("Patient").
		Preload("Doctor").
		Preload("Department").
		Joins("JOIN patients ON patients.id = medical_records.patient_id").
		Where("medical_records.id = ? AND patients.user_id = ?", recordID, userID).
		First(&record).Error
	return &record, err
}

// ListByUser 查询用户的就诊记录列表
func (r *MedicalRecordRepository) ListByUser(userID int64) ([]*model.MedicalRecord, error) {
	var records []*model.MedicalRecord
	err := r.db.
		Preload("Patient").
		Preload("Doctor").
		Preload("Department").
		Joins("JOIN patients ON patients.id = medical_records.patient_id").
		Where("patients.user_id = ?", userID).
		Order("medical_records.visit_date DESC, medical_records.created_at DESC").
		Find(&records).Error
	return records, err
}

// GetByAppointmentID 根据预约ID获取就诊记录
func (r *MedicalRecordRepository) GetByAppointmentID(appointmentID int64) (*model.MedicalRecord, error) {
	var record model.MedicalRecord
	err := r.db.
		Preload("Appointment").
		Preload("Patient").
		Preload("Doctor").
		Preload("Department").
		Where("appointment_id = ?", appointmentID).
		First(&record).Error
	return &record, err
}

// Create 创建就诊记录
func (r *MedicalRecordRepository) Create(record *model.MedicalRecord) error {
	return r.db.Create(record).Error
}

// Update 更新就诊记录
func (r *MedicalRecordRepository) Update(record *model.MedicalRecord) error {
	return r.db.Save(record).Error
}
