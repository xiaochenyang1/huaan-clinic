package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// 状态常量
const (
	StatusDisabled = 0 // 禁用
	StatusEnabled  = 1 // 启用
)

// 性别常量
const (
	GenderUnknown = 0 // 未知
	GenderMale    = 1 // 男
	GenderFemale  = 2 // 女
)

// 预约状态常量
const (
	AppointmentStatusPending   = "pending"   // 待就诊
	AppointmentStatusCheckedIn = "checked_in" // 已签到
	AppointmentStatusCompleted = "completed" // 已完成
	AppointmentStatusCancelled = "cancelled" // 已取消
	AppointmentStatusMissed    = "missed"    // 已爽约
)

// 排班时段常量
const (
	PeriodMorning   = "morning"   // 上午
	PeriodAfternoon = "afternoon" // 下午
)

// 与用户关系常量
const (
	RelationSelf   = "self"   // 本人
	RelationParent = "parent" // 父母
	RelationChild  = "child"  // 子女
	RelationSpouse = "spouse" // 配偶
	RelationOther  = "other"  // 其他
)

// 医生职称常量
const (
	TitleChiefPhysician          = "chief_physician"           // 主任医师
	TitleAssociateChiefPhysician = "associate_chief_physician" // 副主任医师
	TitleAttendingPhysician      = "attending_physician"       // 主治医师
	TitleResidentPhysician       = "resident_physician"        // 住院医师
)

// GetTitleName 获取职称名称
func GetTitleName(title string) string {
	titles := map[string]string{
		TitleChiefPhysician:          "主任医师",
		TitleAssociateChiefPhysician: "副主任医师",
		TitleAttendingPhysician:      "主治医师",
		TitleResidentPhysician:       "住院医师",
	}
	if name, ok := titles[title]; ok {
		return name
	}
	return title
}

// GetRelationName 获取关系名称
func GetRelationName(relation string) string {
	relations := map[string]string{
		RelationSelf:   "本人",
		RelationParent: "父母",
		RelationChild:  "子女",
		RelationSpouse: "配偶",
		RelationOther:  "其他",
	}
	if name, ok := relations[relation]; ok {
		return name
	}
	return relation
}

// GetAppointmentStatusName 获取预约状态名称
func GetAppointmentStatusName(status string) string {
	statuses := map[string]string{
		AppointmentStatusPending:   "待就诊",
		AppointmentStatusCheckedIn: "已签到",
		AppointmentStatusCompleted: "已完成",
		AppointmentStatusCancelled: "已取消",
		AppointmentStatusMissed:    "已爽约",
	}
	if name, ok := statuses[status]; ok {
		return name
	}
	return status
}

// GetPeriodName 获取时段名称
func GetPeriodName(period string) string {
	periods := map[string]string{
		PeriodMorning:   "上午",
		PeriodAfternoon: "下午",
	}
	if name, ok := periods[period]; ok {
		return name
	}
	return period
}
