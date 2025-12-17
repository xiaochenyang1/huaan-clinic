package model

import (
	"time"
)

// Schedule 排班模型
type Schedule struct {
	BaseModel
	DoctorID       int64     `gorm:"index;not null;comment:医生ID" json:"doctor_id"`
	ScheduleDate   time.Time `gorm:"type:date;index;not null;comment:排班日期" json:"schedule_date"`
	Period         string    `gorm:"type:varchar(20);not null;comment:时段 morning/afternoon" json:"period"`
	StartTime      string    `gorm:"type:varchar(10);not null;comment:开始时间 HH:mm" json:"start_time"`
	EndTime        string    `gorm:"type:varchar(10);not null;comment:结束时间 HH:mm" json:"end_time"`
	TotalSlots     int       `gorm:"type:int;not null;comment:总号源数" json:"total_slots"`
	AvailableSlots int       `gorm:"type:int;not null;comment:剩余号源数" json:"available_slots"`
	Status         int       `gorm:"type:tinyint;default:1;comment:状态 0停诊 1正常" json:"status"`

	// 关联
	Doctor *Doctor `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
}

// TableName 表名
func (Schedule) TableName() string {
	return "schedules"
}

// ScheduleVO 排班视图对象
type ScheduleVO struct {
	ID             int64  `json:"id"`
	DoctorID       int64  `json:"doctor_id"`
	DoctorName     string `json:"doctor_name,omitempty"`
	DepartmentID   int64  `json:"department_id,omitempty"`
	DepartmentName string `json:"department_name,omitempty"`
	ScheduleDate   string `json:"schedule_date"`
	Period         string `json:"period"`
	PeriodName     string `json:"period_name"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	TotalSlots     int    `json:"total_slots"`
	AvailableSlots int    `json:"available_slots"`
	Status         int    `json:"status"`
	StatusName     string `json:"status_name"`
	IsAvailable    bool   `json:"is_available"` // 是否可预约
}

// ToVO 转换为视图对象
func (s *Schedule) ToVO() *ScheduleVO {
	statusName := "正常"
	if s.Status == StatusDisabled {
		statusName = "停诊"
	}

	vo := &ScheduleVO{
		ID:             s.ID,
		DoctorID:       s.DoctorID,
		ScheduleDate:   s.ScheduleDate.Format("2006-01-02"),
		Period:         s.Period,
		PeriodName:     GetPeriodName(s.Period),
		StartTime:      s.StartTime,
		EndTime:        s.EndTime,
		TotalSlots:     s.TotalSlots,
		AvailableSlots: s.AvailableSlots,
		Status:         s.Status,
		StatusName:     statusName,
		IsAvailable:    s.Status == StatusEnabled && s.AvailableSlots > 0,
	}

	if s.Doctor != nil {
		vo.DoctorName = s.Doctor.Name
		vo.DepartmentID = s.Doctor.DepartmentID
		if s.Doctor.Department != nil {
			vo.DepartmentName = s.Doctor.Department.Name
		}
	}

	return vo
}

// TimeSlot 时间段
type TimeSlot struct {
	StartTime   string `json:"start_time"`   // 开始时间 HH:mm
	EndTime     string `json:"end_time"`     // 结束时间 HH:mm
	SlotNumber  int    `json:"slot_number"`  // 号序
	IsAvailable bool   `json:"is_available"` // 是否可预约
}

// ScheduleWithSlots 带时间段的排班
type ScheduleWithSlots struct {
	*ScheduleVO
	TimeSlots []TimeSlot `json:"time_slots"`
}
