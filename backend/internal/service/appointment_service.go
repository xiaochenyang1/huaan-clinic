package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/internal/repository"
	"huaan-medical/pkg/database"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/utils"
)

// AppointmentService 预约服务
type AppointmentService struct {
	repo         *repository.AppointmentRepository
	scheduleRepo *repository.ScheduleRepository
	patientRepo  *repository.PatientRepository
	doctorRepo   *repository.DoctorRepository
	userRepo     *repository.UserRepository
	tokenService *TokenService
}

// NewAppointmentService 创建预约服务实例
func NewAppointmentService() *AppointmentService {
	return &AppointmentService{
		repo:         repository.NewAppointmentRepository(),
		scheduleRepo: repository.NewScheduleRepository(),
		patientRepo:  repository.NewPatientRepository(),
		doctorRepo:   repository.NewDoctorRepository(),
		userRepo:     repository.NewUserRepository(),
		tokenService: NewTokenService(),
	}
}

// CreateAppointmentRequest 创建预约请求
type CreateAppointmentRequest struct {
	IdempotentToken string `json:"idempotent_token" binding:"required"`
	ScheduleID      int64  `json:"schedule_id" binding:"required,min=1"`
	PatientID       int64  `json:"patient_id" binding:"required,min=1"`
	Symptom         string `json:"symptom" binding:"max=512"`
}

// ListAppointmentRequest 列表查询请求
type ListAppointmentRequest struct {
	Status *string `form:"status"`
}

// ListAdminAppointmentRequest 管理后台列表查询请求
type ListAdminAppointmentRequest struct {
	Page      int     `form:"page" binding:"required,min=1"`
	PageSize  int     `form:"page_size" binding:"required,min=1,max=100"`
	StartDate string  `form:"start_date"`
	EndDate   string  `form:"end_date"`
	Status    *string `form:"status"`
	Keyword   string  `form:"keyword"`
}

// CancelAppointmentRequest 取消预约请求
type CancelAppointmentRequest struct {
	Reason string `json:"reason" binding:"required,min=2,max=256"`
}

// Create 创建预约
func (s *AppointmentService) Create(userID int64, req *CreateAppointmentRequest) (*model.AppointmentVO, error) {
	// 1. 验证并消费幂等Token
	if err := s.tokenService.ValidateAndConsumeIdempotentToken(userID, req.IdempotentToken); err != nil {
		return nil, err
	}

	// 2. 查询用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrUserNotFound)
	}

	// 检查用户是否被封禁
	if user.IsBlocked() {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "您的账号已被封禁，无法预约")
	}

	// 3. 查询排班信息
	schedule, err := s.scheduleRepo.GetByID(req.ScheduleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrScheduleNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 检查排班状态
	if schedule.Status == model.StatusDisabled {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "该排班已停诊")
	}

	// 检查排班日期是否已过期
	if schedule.ScheduleDate.Before(utils.GetTodayStart()) {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "该排班已过期")
	}

	// 4. 查询就诊人信息（验证就诊人是否存在且属于该用户）
	_, err = s.patientRepo.GetByUserAndID(userID, req.PatientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrPatientNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 5. 检查是否已有同一医生同一时段的预约
	hasAppointment, err := s.repo.CheckUserPendingAppointment(userID, schedule.DoctorID, schedule.ScheduleDate, schedule.Period)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	if hasAppointment {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "您已预约该医生的该时段")
	}

	// 6. 使用事务创建预约并扣减号源
	var appointment *model.Appointment
	err = database.GetDB().Transaction(func(tx *gorm.DB) error {
		// 6.1 扣减号源（使用原子操作）
		result := tx.Model(&model.Schedule{}).
			Where("id = ? AND available_slots > 0", req.ScheduleID).
			Update("available_slots", gorm.Expr("available_slots - 1"))

		if result.Error != nil {
			return result.Error
		}

		// 如果影响行数为0，说明号源不足
		if result.RowsAffected == 0 {
			return errorcode.NewWithMessage(errorcode.ErrInvalidParams, "号源已满")
		}

		// 6.2 获取号序
		slotNumber, err := s.repo.GetNextSlotNumber(req.ScheduleID)
		if err != nil {
			return err
		}

		// 6.3 生成预约编号
		appointmentNo := utils.GenerateAppointmentNo()

		// 6.4 创建预约
		appointment = &model.Appointment{
			AppointmentNo:   appointmentNo,
			UserID:          userID,
			PatientID:       req.PatientID,
			DoctorID:        schedule.DoctorID,
			DepartmentID:    schedule.Doctor.DepartmentID,
			ScheduleID:      req.ScheduleID,
			AppointmentDate: schedule.ScheduleDate,
			Period:          schedule.Period,
			AppointmentTime: schedule.StartTime,
			SlotNumber:      slotNumber,
			Status:          model.AppointmentStatusPending,
			Symptom:         req.Symptom,
		}

		return s.repo.Create(tx, appointment)
	})

	if err != nil {
		return nil, err
	}

	// 7. 重新查询以获取完整数据
	appointment, err = s.repo.GetByID(appointment.ID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	return appointment.ToVO(), nil
}

// Cancel 取消预约
func (s *AppointmentService) Cancel(userID, appointmentID int64, req *CancelAppointmentRequest) error {
	// 1. 查询预约
	appointment, err := s.repo.GetByUserAndID(userID, appointmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorcode.New(errorcode.ErrAppointmentNotFound)
		}
		return errorcode.New(errorcode.ErrDatabase)
	}

	// 2. 检查预约状态
	if appointment.Status != model.AppointmentStatusPending {
		return errorcode.NewWithMessage(errorcode.ErrInvalidParams, "只能取消待就诊的预约")
	}

	// 3. 检查是否允许取消（就诊当天不可取消）
	if !appointment.ToVO().CanCancel {
		return errorcode.NewWithMessage(errorcode.ErrInvalidParams, "就诊当天不可取消")
	}

	// 4. 使用事务取消预约并返还号源
	now := time.Now()
	err = database.GetDB().Transaction(func(tx *gorm.DB) error {
		// 4.1 更新预约状态
		updates := map[string]interface{}{
			"cancel_reason": req.Reason,
			"cancelled_at":  now,
		}
		if err := s.repo.UpdateStatus(appointmentID, model.AppointmentStatusCancelled, updates); err != nil {
			return err
		}

		// 4.2 返还号源
		return s.scheduleRepo.UpdateAvailableSlots(appointment.ScheduleID, 1)
	})

	return err
}

// Checkin 预约签到
func (s *AppointmentService) Checkin(userID, appointmentID int64) error {
	// 1. 查询预约
	appointment, err := s.repo.GetByUserAndID(userID, appointmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorcode.New(errorcode.ErrAppointmentNotFound)
		}
		return errorcode.New(errorcode.ErrDatabase)
	}

	// 2. 检查预约状态
	if appointment.Status != model.AppointmentStatusPending {
		return errorcode.NewWithMessage(errorcode.ErrInvalidParams, "只能签到待就诊的预约")
	}

	// 3. 检查是否允许签到
	if !appointment.ToVO().CanCheckin {
		return errorcode.NewWithMessage(errorcode.ErrInvalidParams, "不在签到时间范围内")
	}

	// 4. 更新预约状态
	now := time.Now()
	updates := map[string]interface{}{
		"checked_in_at": now,
	}
	return s.repo.UpdateStatus(appointmentID, model.AppointmentStatusCheckedIn, updates)
}

// GetByID 获取预约详情
func (s *AppointmentService) GetByID(userID, appointmentID int64) (*model.AppointmentVO, error) {
	appointment, err := s.repo.GetByUserAndID(userID, appointmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrAppointmentNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	return appointment.ToVO(), nil
}

// ListByUser 查询用户的预约列表
func (s *AppointmentService) ListByUser(userID int64, req *ListAppointmentRequest) ([]model.AppointmentListVO, error) {
	appointments, err := s.repo.ListByUser(userID, req.Status)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	voList := make([]model.AppointmentListVO, len(appointments))
	for i, appointment := range appointments {
		voList[i] = *appointment.ToListVO()
	}

	return voList, nil
}

// List 分页查询预约列表（管理后台）
func (s *AppointmentService) List(req *ListAdminAppointmentRequest) ([]model.AppointmentVO, int64, error) {
	var startDate, endDate *time.Time

	// 解析日期范围
	if req.StartDate != "" {
		sd, err := utils.ParseDate(req.StartDate)
		if err != nil {
			return nil, 0, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "开始日期格式错误")
		}
		startDate = &sd
	}

	if req.EndDate != "" {
		ed, err := utils.ParseDate(req.EndDate)
		if err != nil {
			return nil, 0, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "结束日期格式错误")
		}
		endDate = &ed
	}

	appointments, total, err := s.repo.List(req.Page, req.PageSize, startDate, endDate, req.Status, req.Keyword)
	if err != nil {
		return nil, 0, errorcode.New(errorcode.ErrDatabase)
	}

	voList := make([]model.AppointmentVO, len(appointments))
	for i, appointment := range appointments {
		voList[i] = *appointment.ToVO()
	}

	return voList, total, nil
}
