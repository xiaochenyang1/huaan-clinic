package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/internal/repository"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/utils"
)

// ScheduleService 排班服务
type ScheduleService struct {
	repo       *repository.ScheduleRepository
	doctorRepo *repository.DoctorRepository
}

// NewScheduleService 创建排班服务实例
func NewScheduleService() *ScheduleService {
	return &ScheduleService{
		repo:       repository.NewScheduleRepository(),
		doctorRepo: repository.NewDoctorRepository(),
	}
}

// CreateScheduleRequest 创建排班请求
type CreateScheduleRequest struct {
	DoctorID     int64  `json:"doctor_id" binding:"required,min=1"`
	ScheduleDate string `json:"schedule_date" binding:"required"` // YYYY-MM-DD
	Period       string `json:"period" binding:"required,oneof=morning afternoon"`
	StartTime    string `json:"start_time" binding:"required"` // HH:mm
	EndTime      string `json:"end_time" binding:"required"`   // HH:mm
	TotalSlots   int    `json:"total_slots" binding:"required,min=1,max=999"`
	Status       int    `json:"status" binding:"oneof=0 1"`
}

// UpdateScheduleRequest 更新排班请求
type UpdateScheduleRequest struct {
	StartTime  string `json:"start_time" binding:"required"` // HH:mm
	EndTime    string `json:"end_time" binding:"required"`   // HH:mm
	TotalSlots int    `json:"total_slots" binding:"required,min=1,max=999"`
	Status     int    `json:"status" binding:"oneof=0 1"`
}

// BatchCreateScheduleRequest 批量创建排班请求
type BatchCreateScheduleRequest struct {
	DoctorID   int64    `json:"doctor_id" binding:"required,min=1"`
	StartDate  string   `json:"start_date" binding:"required"`  // YYYY-MM-DD
	EndDate    string   `json:"end_date" binding:"required"`    // YYYY-MM-DD
	Periods    []string `json:"periods" binding:"required,min=1,dive,oneof=morning afternoon"`
	WeekDays   []int    `json:"week_days" binding:"required,min=1,dive,min=0,max=6"` // 0=周日, 1=周一...6=周六
	StartTimes []string `json:"start_times" binding:"required,len=2"`                // [上午开始时间, 下午开始时间]
	EndTimes   []string `json:"end_times" binding:"required,len=2"`                  // [上午结束时间, 下午结束时间]
	TotalSlots int      `json:"total_slots" binding:"required,min=1,max=999"`
}

// ListScheduleRequest 列表查询请求
type ListScheduleRequest struct {
	Page         int    `form:"page" binding:"required,min=1"`
	PageSize     int    `form:"page_size" binding:"required,min=1,max=100"`
	DoctorID     *int64 `form:"doctor_id"`
	DepartmentID *int64 `form:"department_id"`
	StartDate    string `form:"start_date"` // YYYY-MM-DD
	EndDate      string `form:"end_date"`   // YYYY-MM-DD
	Status       *int   `form:"status"`
}

// ListAvailableScheduleRequest 查询可预约排班请求
type ListAvailableScheduleRequest struct {
	DoctorID     *int64 `form:"doctor_id"`
	DepartmentID *int64 `form:"department_id"`
	StartDate    string `form:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate      string `form:"end_date" binding:"required"`   // YYYY-MM-DD
}

// Create 创建排班
func (s *ScheduleService) Create(req *CreateScheduleRequest) (*model.ScheduleVO, error) {
	// 解析日期
	scheduleDate, err := utils.ParseDate(req.ScheduleDate)
	if err != nil {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "排班日期格式错误")
	}

	// 日期不能早于今天
	if scheduleDate.Before(utils.GetTodayStart()) {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "排班日期不能早于今天")
	}

	// 检查医生是否存在
	doctor, err := s.doctorRepo.GetByIDSimple(req.DoctorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrDoctorNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 检查医生是否启用
	if doctor.Status == model.StatusDisabled {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "该医生已停诊")
	}

	// 检查排班是否已存在
	exists, err := s.repo.Exists(req.DoctorID, scheduleDate, req.Period)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	if exists {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "该时段排班已存在")
	}

	// 创建排班
	schedule := &model.Schedule{
		DoctorID:       req.DoctorID,
		ScheduleDate:   scheduleDate,
		Period:         req.Period,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		TotalSlots:     req.TotalSlots,
		AvailableSlots: req.TotalSlots, // 初始剩余号源等于总号源
		Status:         req.Status,
	}

	if err := s.repo.Create(schedule); err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 重新查询以获取关联数据
	schedule, err = s.repo.GetByID(schedule.ID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	return schedule.ToVO(), nil
}

// BatchCreate 批量创建排班
func (s *ScheduleService) BatchCreate(req *BatchCreateScheduleRequest) (int, error) {
	// 解析日期
	startDate, err := utils.ParseDate(req.StartDate)
	if err != nil {
		return 0, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "开始日期格式错误")
	}

	endDate, err := utils.ParseDate(req.EndDate)
	if err != nil {
		return 0, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "结束日期格式错误")
	}

	// 日期验证
	if startDate.Before(utils.GetTodayStart()) {
		return 0, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "开始日期不能早于今天")
	}

	if endDate.Before(startDate) {
		return 0, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "结束日期不能早于开始日期")
	}

	// 日期跨度限制（最多90天）
	if endDate.Sub(startDate).Hours() > 90*24 {
		return 0, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "日期跨度不能超过90天")
	}

	// 检查医生是否存在
	doctor, err := s.doctorRepo.GetByIDSimple(req.DoctorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errorcode.New(errorcode.ErrDoctorNotFound)
		}
		return 0, errorcode.New(errorcode.ErrDatabase)
	}

	// 检查医生是否启用
	if doctor.Status == model.StatusDisabled {
		return 0, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "该医生已停诊")
	}

	// 构建星期几的映射
	weekDayMap := make(map[int]bool)
	for _, wd := range req.WeekDays {
		weekDayMap[wd] = true
	}

	// 构建时段时间映射
	timeMap := map[string]struct {
		StartTime string
		EndTime   string
	}{
		model.PeriodMorning:   {StartTime: req.StartTimes[0], EndTime: req.EndTimes[0]},
		model.PeriodAfternoon: {StartTime: req.StartTimes[1], EndTime: req.EndTimes[1]},
	}

	// 生成排班列表
	var schedules []model.Schedule
	currentDate := startDate

	for !currentDate.After(endDate) {
		// 检查是否在指定的星期几
		weekday := int(currentDate.Weekday())
		if !weekDayMap[weekday] {
			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}

		// 为每个时段创建排班
		for _, period := range req.Periods {
			// 检查是否已存在
			exists, err := s.repo.Exists(req.DoctorID, currentDate, period)
			if err != nil {
				return 0, errorcode.New(errorcode.ErrDatabase)
			}

			// 跳过已存在的排班
			if exists {
				currentDate = currentDate.AddDate(0, 0, 1)
				continue
			}

			timeInfo := timeMap[period]
			schedules = append(schedules, model.Schedule{
				DoctorID:       req.DoctorID,
				ScheduleDate:   currentDate,
				Period:         period,
				StartTime:      timeInfo.StartTime,
				EndTime:        timeInfo.EndTime,
				TotalSlots:     req.TotalSlots,
				AvailableSlots: req.TotalSlots,
				Status:         model.StatusEnabled,
			})
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// 批量创建
	if len(schedules) == 0 {
		return 0, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "没有可创建的排班")
	}

	if err := s.repo.BatchCreate(schedules); err != nil {
		return 0, errorcode.New(errorcode.ErrDatabase)
	}

	return len(schedules), nil
}

// Update 更新排班
func (s *ScheduleService) Update(id int64, req *UpdateScheduleRequest) (*model.ScheduleVO, error) {
	// 检查排班是否存在
	schedule, err := s.repo.GetByIDSimple(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrScheduleNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 如果减少总号源数，需要检查是否小于已预约数
	bookedSlots := schedule.TotalSlots - schedule.AvailableSlots
	if req.TotalSlots < bookedSlots {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "总号源数不能少于已预约数")
	}

	// 更新排班信息
	schedule.StartTime = req.StartTime
	schedule.EndTime = req.EndTime
	schedule.TotalSlots = req.TotalSlots
	schedule.AvailableSlots = req.TotalSlots - bookedSlots // 重新计算剩余号源
	schedule.Status = req.Status

	if err := s.repo.Update(schedule); err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 重新查询以获取关联数据
	schedule, err = s.repo.GetByID(id)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	return schedule.ToVO(), nil
}

// Delete 删除排班
func (s *ScheduleService) Delete(id int64) error {
	// 检查排班是否存在
	_, err := s.repo.GetByIDSimple(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorcode.New(errorcode.ErrScheduleNotFound)
		}
		return errorcode.New(errorcode.ErrDatabase)
	}

	// 检查是否有预约记录
	hasAppointments, err := s.repo.HasAppointments(id)
	if err != nil {
		return errorcode.New(errorcode.ErrDatabase)
	}
	if hasAppointments {
		return errorcode.NewWithMessage(errorcode.ErrInvalidParams, "该排班存在预约记录，无法删除")
	}

	return s.repo.Delete(id)
}

// GetByID 获取排班详情
func (s *ScheduleService) GetByID(id int64) (*model.ScheduleVO, error) {
	schedule, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrScheduleNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	return schedule.ToVO(), nil
}

// List 分页查询排班列表（管理后台）
func (s *ScheduleService) List(req *ListScheduleRequest) ([]model.ScheduleVO, int64, error) {
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

	schedules, total, err := s.repo.List(req.Page, req.PageSize, req.DoctorID, req.DepartmentID, startDate, endDate, req.Status)
	if err != nil {
		return nil, 0, errorcode.New(errorcode.ErrDatabase)
	}

	voList := make([]model.ScheduleVO, len(schedules))
	for i, schedule := range schedules {
		voList[i] = *schedule.ToVO()
	}

	return voList, total, nil
}

// ListByDoctor 查询医生的排班列表（公开接口）
func (s *ScheduleService) ListByDoctor(doctorID int64, startDate, endDate string) ([]model.ScheduleVO, error) {
	// 解析日期
	sd, err := utils.ParseDate(startDate)
	if err != nil {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "开始日期格式错误")
	}

	ed, err := utils.ParseDate(endDate)
	if err != nil {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "结束日期格式错误")
	}

	// 检查医生是否存在
	_, err = s.doctorRepo.GetByIDSimple(doctorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrDoctorNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	schedules, err := s.repo.ListByDoctor(doctorID, sd, ed)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	voList := make([]model.ScheduleVO, len(schedules))
	for i, schedule := range schedules {
		voList[i] = *schedule.ToVO()
	}

	return voList, nil
}

// ListAvailable 查询可预约的排班列表（公开接口）
func (s *ScheduleService) ListAvailable(req *ListAvailableScheduleRequest) ([]model.ScheduleVO, error) {
	// 解析日期
	startDate, err := utils.ParseDate(req.StartDate)
	if err != nil {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "开始日期格式错误")
	}

	endDate, err := utils.ParseDate(req.EndDate)
	if err != nil {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "结束日期格式错误")
	}

	schedules, err := s.repo.ListAvailable(req.DoctorID, req.DepartmentID, startDate, endDate)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	voList := make([]model.ScheduleVO, len(schedules))
	for i, schedule := range schedules {
		voList[i] = *schedule.ToVO()
	}

	return voList, nil
}
