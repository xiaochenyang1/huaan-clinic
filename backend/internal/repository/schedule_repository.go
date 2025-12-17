package repository

import (
	"time"

	"huaan-medical/internal/model"
	"huaan-medical/pkg/database"

	"gorm.io/gorm"
)

// ScheduleRepository 排班数据访问层
type ScheduleRepository struct {
	db *gorm.DB
}

// NewScheduleRepository 创建排班仓库实例
func NewScheduleRepository() *ScheduleRepository {
	return &ScheduleRepository{db: database.GetDB()}
}

// Create 创建排班
func (r *ScheduleRepository) Create(schedule *model.Schedule) error {
	return r.db.Create(schedule).Error
}

// BatchCreate 批量创建排班（使用事务）
func (r *ScheduleRepository) BatchCreate(schedules []model.Schedule) error {
	if len(schedules) == 0 {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.CreateInBatches(schedules, 100).Error
	})
}

// Update 更新排班
func (r *ScheduleRepository) Update(schedule *model.Schedule) error {
	return r.db.Save(schedule).Error
}

// Delete 删除排班（软删除）
func (r *ScheduleRepository) Delete(id int64) error {
	return r.db.Delete(&model.Schedule{}, id).Error
}

// GetByID 根据ID查询排班
func (r *ScheduleRepository) GetByID(id int64) (*model.Schedule, error) {
	var schedule model.Schedule
	err := r.db.Preload("Doctor.Department").First(&schedule, id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

// GetByIDSimple 根据ID查询排班（不预加载关联）
func (r *ScheduleRepository) GetByIDSimple(id int64) (*model.Schedule, error) {
	var schedule model.Schedule
	err := r.db.First(&schedule, id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

// Exists 检查排班是否存在
// 用于校验同一医生在同一天同一时段只能有一个排班
func (r *ScheduleRepository) Exists(doctorID int64, scheduleDate time.Time, period string, excludeID ...int64) (bool, error) {
	query := r.db.Model(&model.Schedule{}).
		Where("doctor_id = ? AND schedule_date = ? AND period = ?", doctorID, scheduleDate, period)

	// 排除指定ID（用于更新时的检查）
	if len(excludeID) > 0 && excludeID[0] > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// List 分页查询排班列表（管理后台）
func (r *ScheduleRepository) List(page, pageSize int, doctorID *int64, departmentID *int64, startDate, endDate *time.Time, status *int) ([]model.Schedule, int64, error) {
	var schedules []model.Schedule
	var total int64

	query := r.db.Model(&model.Schedule{}).Preload("Doctor.Department")

	// 医生筛选
	if doctorID != nil && *doctorID > 0 {
		query = query.Where("doctor_id = ?", *doctorID)
	}

	// 科室筛选（通过关联查询）
	if departmentID != nil && *departmentID > 0 {
		query = query.Joins("JOIN doctors ON doctors.id = schedules.doctor_id").
			Where("doctors.department_id = ?", *departmentID)
	}

	// 日期范围筛选
	if startDate != nil {
		query = query.Where("schedule_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("schedule_date <= ?", *endDate)
	}

	// 状态筛选
	if status != nil {
		query = query.Where("schedules.status = ?", *status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("schedule_date DESC, period ASC").
		Offset(offset).Limit(pageSize).
		Find(&schedules).Error

	return schedules, total, err
}

// ListByDoctor 查询医生的排班列表（公开接口）
func (r *ScheduleRepository) ListByDoctor(doctorID int64, startDate, endDate time.Time) ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := r.db.Preload("Doctor.Department").
		Where("doctor_id = ? AND schedule_date >= ? AND schedule_date <= ? AND status = ?",
			doctorID, startDate, endDate, model.StatusEnabled).
		Order("schedule_date ASC, period ASC").
		Find(&schedules).Error
	return schedules, err
}

// ListAvailable 查询可预约的排班列表（公开接口）
func (r *ScheduleRepository) ListAvailable(doctorID *int64, departmentID *int64, startDate, endDate time.Time) ([]model.Schedule, error) {
	var schedules []model.Schedule

	query := r.db.Preload("Doctor.Department").
		Where("schedule_date >= ? AND schedule_date <= ? AND status = ? AND available_slots > 0",
			startDate, endDate, model.StatusEnabled)

	// 医生筛选
	if doctorID != nil && *doctorID > 0 {
		query = query.Where("doctor_id = ?", *doctorID)
	}

	// 科室筛选
	if departmentID != nil && *departmentID > 0 {
		query = query.Joins("JOIN doctors ON doctors.id = schedules.doctor_id").
			Where("doctors.department_id = ?", *departmentID)
	}

	err := query.Order("schedule_date ASC, period ASC").Find(&schedules).Error
	return schedules, err
}

// GetByDoctorAndDate 根据医生ID和日期查询排班
func (r *ScheduleRepository) GetByDoctorAndDate(doctorID int64, scheduleDate time.Time, period string) (*model.Schedule, error) {
	var schedule model.Schedule
	err := r.db.Where("doctor_id = ? AND schedule_date = ? AND period = ?",
		doctorID, scheduleDate, period).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

// HasAppointments 检查排班是否有预约记录
func (r *ScheduleRepository) HasAppointments(id int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Appointment{}).Where("schedule_id = ?", id).Count(&count).Error
	return count > 0, err
}

// UpdateAvailableSlots 更新剩余号源数（预约时使用，需要原子性）
func (r *ScheduleRepository) UpdateAvailableSlots(id int64, delta int) error {
	// 使用 SQL 原子操作，防止并发问题
	return r.db.Model(&model.Schedule{}).
		Where("id = ? AND available_slots + ? >= 0", id, delta).
		Update("available_slots", gorm.Expr("available_slots + ?", delta)).Error
}

// CountByDoctor 统计医生的排班数量
func (r *ScheduleRepository) CountByDoctor(doctorID int64, startDate, endDate *time.Time) (int64, error) {
	query := r.db.Model(&model.Schedule{}).Where("doctor_id = ?", doctorID)

	if startDate != nil {
		query = query.Where("schedule_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("schedule_date <= ?", *endDate)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// UpdateStatus 批量更新排班状态
func (r *ScheduleRepository) UpdateStatus(ids []int64, status int) error {
	return r.db.Model(&model.Schedule{}).
		Where("id IN ?", ids).
		Update("status", status).Error
}

// UpdateStatusByDoctor 更新医生的所有排班状态
func (r *ScheduleRepository) UpdateStatusByDoctor(doctorID int64, status int) error {
	return r.db.Model(&model.Schedule{}).
		Where("doctor_id = ? AND schedule_date >= ?", doctorID, time.Now()).
		Update("status", status).Error
}
