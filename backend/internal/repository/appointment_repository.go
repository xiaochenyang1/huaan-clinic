package repository

import (
	"time"

	"huaan-medical/internal/model"
	"huaan-medical/pkg/database"

	"gorm.io/gorm"
)

// AppointmentRepository 预约数据访问层
type AppointmentRepository struct {
	db *gorm.DB
}

// NewAppointmentRepository 创建预约仓库实例
func NewAppointmentRepository() *AppointmentRepository {
	return &AppointmentRepository{db: database.GetDB()}
}

// Create 创建预约（需要在事务中调用）
func (r *AppointmentRepository) Create(tx *gorm.DB, appointment *model.Appointment) error {
	return tx.Create(appointment).Error
}

// GetByID 根据ID查询预约
func (r *AppointmentRepository) GetByID(id int64) (*model.Appointment, error) {
	var appointment model.Appointment
	err := r.db.Preload("Patient").
		Preload("Doctor").
		Preload("Department").
		Preload("Schedule").
		First(&appointment, id).Error
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

// GetByUserAndID 根据用户ID和预约ID查询（用于权限校验）
func (r *AppointmentRepository) GetByUserAndID(userID, appointmentID int64) (*model.Appointment, error) {
	var appointment model.Appointment
	err := r.db.Preload("Patient").
		Preload("Doctor").
		Preload("Department").
		Preload("Schedule").
		Where("user_id = ? AND id = ?", userID, appointmentID).
		First(&appointment).Error
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

// ListByUser 查询用户的预约列表
func (r *AppointmentRepository) ListByUser(userID int64, status *string) ([]model.Appointment, error) {
	var appointments []model.Appointment

	query := r.db.Preload("Patient").
		Preload("Doctor").
		Preload("Department").
		Where("user_id = ?", userID)

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("appointment_date DESC, appointment_time DESC").
		Find(&appointments).Error
	return appointments, err
}

// List 分页查询预约列表（管理后台）
func (r *AppointmentRepository) List(page, pageSize int, startDate, endDate *time.Time, status *string, keyword string) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var total int64

	query := r.db.Model(&model.Appointment{}).
		Preload("Patient").
		Preload("Doctor").
		Preload("Department")

	// 日期范围筛选
	if startDate != nil {
		query = query.Where("appointment_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("appointment_date <= ?", *endDate)
	}

	// 状态筛选
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	// 关键词搜索（预约编号、患者姓名、医生姓名）
	if keyword != "" {
		query = query.Where(
			r.db.Where("appointment_no LIKE ?", "%"+keyword+"%").
				Or("EXISTS (SELECT 1 FROM patients WHERE patients.id = appointments.patient_id AND patients.name LIKE ?)", "%"+keyword+"%").
				Or("EXISTS (SELECT 1 FROM doctors WHERE doctors.id = appointments.doctor_id AND doctors.name LIKE ?)", "%"+keyword+"%"),
		)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("appointment_date DESC, appointment_time DESC").
		Offset(offset).Limit(pageSize).
		Find(&appointments).Error

	return appointments, total, err
}

// UpdateStatus 更新预约状态
func (r *AppointmentRepository) UpdateStatus(id int64, status string, extraFields map[string]interface{}) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// 合并额外字段
	for k, v := range extraFields {
		updates[k] = v
	}

	return r.db.Model(&model.Appointment{}).Where("id = ?", id).Updates(updates).Error
}

// CheckUserPendingAppointment 检查用户是否有同一医生同一时段的待就诊预约
func (r *AppointmentRepository) CheckUserPendingAppointment(userID, doctorID int64, appointmentDate time.Time, period string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Appointment{}).
		Where("user_id = ? AND doctor_id = ? AND appointment_date = ? AND period = ? AND status = ?",
			userID, doctorID, appointmentDate, period, model.AppointmentStatusPending).
		Count(&count).Error
	return count > 0, err
}

// CountBySchedule 统计排班的预约数量
func (r *AppointmentRepository) CountBySchedule(scheduleID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Appointment{}).
		Where("schedule_id = ?", scheduleID).
		Count(&count).Error
	return count, err
}

// GetNextSlotNumber 获取下一个号序
func (r *AppointmentRepository) GetNextSlotNumber(scheduleID int64) (int, error) {
	var maxSlotNumber int
	err := r.db.Model(&model.Appointment{}).
		Where("schedule_id = ?", scheduleID).
		Select("COALESCE(MAX(slot_number), 0)").
		Scan(&maxSlotNumber).Error
	return maxSlotNumber + 1, err
}
