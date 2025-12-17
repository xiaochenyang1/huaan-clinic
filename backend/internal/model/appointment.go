package model

import (
	"time"
)

// Appointment 预约模型
type Appointment struct {
	BaseModel
	AppointmentNo   string     `gorm:"type:varchar(32);uniqueIndex;not null;comment:预约编号" json:"appointment_no"`
	UserID          int64      `gorm:"index;not null;comment:用户ID" json:"user_id"`
	PatientID       int64      `gorm:"index;not null;comment:就诊人ID" json:"patient_id"`
	DoctorID        int64      `gorm:"index;not null;comment:医生ID" json:"doctor_id"`
	DepartmentID    int64      `gorm:"index;not null;comment:科室ID" json:"department_id"`
	ScheduleID      int64      `gorm:"index;not null;comment:排班ID" json:"schedule_id"`
	AppointmentDate time.Time  `gorm:"type:date;index;not null;comment:预约日期" json:"appointment_date"`
	Period          string     `gorm:"type:varchar(20);not null;comment:时段" json:"period"`
	AppointmentTime string     `gorm:"type:varchar(10);not null;comment:预约时间 HH:mm" json:"appointment_time"`
	SlotNumber      int        `gorm:"type:int;not null;comment:号序" json:"slot_number"`
	Status          string     `gorm:"type:varchar(20);default:'pending';index;comment:状态" json:"status"`
	Symptom         string     `gorm:"type:varchar(512);comment:症状描述" json:"symptom"`
	CancelReason    string     `gorm:"type:varchar(256);comment:取消原因" json:"cancel_reason"`
	CancelledAt     *time.Time `gorm:"comment:取消时间" json:"cancelled_at,omitempty"`
	CheckedInAt     *time.Time `gorm:"comment:签到时间" json:"checked_in_at,omitempty"`
	CompletedAt     *time.Time `gorm:"comment:完成时间" json:"completed_at,omitempty"`

	// 关联
	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Patient    *Patient    `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	Doctor     *Doctor     `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Schedule   *Schedule   `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
}

// TableName 表名
func (Appointment) TableName() string {
	return "appointments"
}

// AppointmentVO 预约视图对象
type AppointmentVO struct {
	ID              int64  `json:"id"`
	AppointmentNo   string `json:"appointment_no"`
	UserID          int64  `json:"user_id"`
	PatientID       int64  `json:"patient_id"`
	PatientName     string `json:"patient_name"`
	DoctorID        int64  `json:"doctor_id"`
	DoctorName      string `json:"doctor_name"`
	DoctorTitle     string `json:"doctor_title"`
	DoctorAvatar    string `json:"doctor_avatar,omitempty"`
	DepartmentID    int64  `json:"department_id"`
	DepartmentName  string `json:"department_name"`
	AppointmentDate string `json:"appointment_date"`
	Period          string `json:"period"`
	PeriodName      string `json:"period_name"`
	AppointmentTime string `json:"appointment_time"`
	SlotNumber      int    `json:"slot_number"`
	Status          string `json:"status"`
	StatusName      string `json:"status_name"`
	Symptom         string `json:"symptom,omitempty"`
	CancelReason    string `json:"cancel_reason,omitempty"`
	CancelledAt     string `json:"cancelled_at,omitempty"`
	CheckedInAt     string `json:"checked_in_at,omitempty"`
	CompletedAt     string `json:"completed_at,omitempty"`
	CreatedAt       string `json:"created_at"`
	CanCancel       bool   `json:"can_cancel"`   // 是否可取消
	CanCheckin      bool   `json:"can_checkin"`  // 是否可签到
}

// ToVO 转换为视图对象
func (a *Appointment) ToVO() *AppointmentVO {
	vo := &AppointmentVO{
		ID:              a.ID,
		AppointmentNo:   a.AppointmentNo,
		UserID:          a.UserID,
		PatientID:       a.PatientID,
		DoctorID:        a.DoctorID,
		DepartmentID:    a.DepartmentID,
		AppointmentDate: a.AppointmentDate.Format("2006-01-02"),
		Period:          a.Period,
		PeriodName:      GetPeriodName(a.Period),
		AppointmentTime: a.AppointmentTime,
		SlotNumber:      a.SlotNumber,
		Status:          a.Status,
		StatusName:      GetAppointmentStatusName(a.Status),
		Symptom:         a.Symptom,
		CancelReason:    a.CancelReason,
		CreatedAt:       a.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 关联信息
	if a.Patient != nil {
		vo.PatientName = maskName(a.Patient.Name)
	}
	if a.Doctor != nil {
		vo.DoctorName = a.Doctor.Name
		vo.DoctorTitle = GetTitleName(a.Doctor.Title)
		vo.DoctorAvatar = a.Doctor.Avatar
	}
	if a.Department != nil {
		vo.DepartmentName = a.Department.Name
	}

	// 时间格式化
	if a.CancelledAt != nil {
		vo.CancelledAt = a.CancelledAt.Format("2006-01-02 15:04:05")
	}
	if a.CheckedInAt != nil {
		vo.CheckedInAt = a.CheckedInAt.Format("2006-01-02 15:04:05")
	}
	if a.CompletedAt != nil {
		vo.CompletedAt = a.CompletedAt.Format("2006-01-02 15:04:05")
	}

	// 计算是否可取消、可签到
	vo.CanCancel = a.canCancel()
	vo.CanCheckin = a.canCheckin()

	return vo
}

// canCancel 判断是否可取消
func (a *Appointment) canCancel() bool {
	if a.Status != AppointmentStatusPending {
		return false
	}
	// 就诊当天不可取消
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return a.AppointmentDate.After(today)
}

// canCheckin 判断是否可签到
func (a *Appointment) canCheckin() bool {
	if a.Status != AppointmentStatusPending {
		return false
	}

	now := time.Now()

	// 解析预约时间
	appointmentDateTime, err := time.ParseInLocation(
		"2006-01-02 15:04",
		a.AppointmentDate.Format("2006-01-02")+" "+a.AppointmentTime,
		now.Location(),
	)
	if err != nil {
		return false
	}

	// 可提前30分钟签到
	earliestCheckin := appointmentDateTime.Add(-30 * time.Minute)
	// 迟到15分钟后不可签到
	latestCheckin := appointmentDateTime.Add(15 * time.Minute)

	return now.After(earliestCheckin) && now.Before(latestCheckin)
}

// AppointmentListVO 预约列表视图对象（简化版）
type AppointmentListVO struct {
	ID              int64  `json:"id"`
	AppointmentNo   string `json:"appointment_no"`
	PatientName     string `json:"patient_name"`
	DoctorName      string `json:"doctor_name"`
	DoctorAvatar    string `json:"doctor_avatar"`
	DepartmentName  string `json:"department_name"`
	AppointmentDate string `json:"appointment_date"`
	PeriodName      string `json:"period_name"`
	AppointmentTime string `json:"appointment_time"`
	Status          string `json:"status"`
	StatusName      string `json:"status_name"`
	CanCancel       bool   `json:"can_cancel"`
	CanCheckin      bool   `json:"can_checkin"`
}

// ToListVO 转换为列表视图对象
func (a *Appointment) ToListVO() *AppointmentListVO {
	vo := &AppointmentListVO{
		ID:              a.ID,
		AppointmentNo:   a.AppointmentNo,
		AppointmentDate: a.AppointmentDate.Format("2006-01-02"),
		PeriodName:      GetPeriodName(a.Period),
		AppointmentTime: a.AppointmentTime,
		Status:          a.Status,
		StatusName:      GetAppointmentStatusName(a.Status),
		CanCancel:       a.canCancel(),
		CanCheckin:      a.canCheckin(),
	}

	if a.Patient != nil {
		vo.PatientName = maskName(a.Patient.Name)
	}
	if a.Doctor != nil {
		vo.DoctorName = a.Doctor.Name
		vo.DoctorAvatar = a.Doctor.Avatar
	}
	if a.Department != nil {
		vo.DepartmentName = a.Department.Name
	}

	return vo
}
