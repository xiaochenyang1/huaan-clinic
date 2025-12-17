package service

import (
	"time"

	"huaan-medical/internal/model"
	"huaan-medical/pkg/database"
)

// StatisticsService 统计服务
type StatisticsService struct{}

// NewStatisticsService 创建统计服务实例
func NewStatisticsService() *StatisticsService {
	return &StatisticsService{}
}

// DashboardData 仪表盘数据
type DashboardData struct {
	TodayAppointments   int64                    `json:"today_appointments"`    // 今日预约数
	TodayCheckins       int64                    `json:"today_checkins"`        // 今日签到数
	TotalAppointments   int64                    `json:"total_appointments"`    // 总预约数
	TotalUsers          int64                    `json:"total_users"`           // 总用户数
	TotalDoctors        int64                    `json:"total_doctors"`         // 总医生数
	TotalDepartments    int64                    `json:"total_departments"`     // 总科室数
	RecentAppointments  []model.AppointmentListVO `json:"recent_appointments"`   // 最近预约
	AppointmentTrend    []TrendData              `json:"appointment_trend"`     // 预约趋势（最近7天）
	DepartmentStats     []DepartmentStat         `json:"department_stats"`      // 科室统计
}

// TrendData 趋势数据
type TrendData struct {
	Date  string `json:"date"`  // 日期 YYYY-MM-DD
	Count int64  `json:"count"` // 数量
}

// DepartmentStat 科室统计
type DepartmentStat struct {
	DepartmentName string `json:"department_name"` // 科室名称
	DoctorCount    int64  `json:"doctor_count"`    // 医生数
	AppointmentCount int64  `json:"appointment_count"` // 预约数
}

// StatisticsData 统计数据
type StatisticsData struct {
	AppointmentStats    AppointmentStats    `json:"appointment_stats"`    // 预约统计
	DoctorStats         []DoctorStat        `json:"doctor_stats"`         // 医生统计
	DepartmentRanking   []DepartmentRanking `json:"department_ranking"`   // 科室排行
	TimeSlotDistribution []TimeSlotStat      `json:"time_slot_distribution"` // 时段分布
}

// AppointmentStats 预约统计
type AppointmentStats struct {
	Total      int64 `json:"total"`       // 总预约数
	Pending    int64 `json:"pending"`     // 待就诊
	CheckedIn  int64 `json:"checked_in"`  // 已签到
	Completed  int64 `json:"completed"`   // 已完成
	Cancelled  int64 `json:"cancelled"`   // 已取消
	Missed     int64 `json:"missed"`      // 已爽约
}

// DoctorStat 医生统计
type DoctorStat struct {
	DoctorID   int64  `json:"doctor_id"`
	DoctorName string `json:"doctor_name"`
	DepartmentName string `json:"department_name"`
	AppointmentCount int64  `json:"appointment_count"`
	CompletedCount int64  `json:"completed_count"`
	Rating float64 `json:"rating"` // 评分（预留）
}

// DepartmentRanking 科室排行
type DepartmentRanking struct {
	DepartmentID   int64  `json:"department_id"`
	DepartmentName string `json:"department_name"`
	AppointmentCount int64  `json:"appointment_count"`
}

// TimeSlotStat 时段统计
type TimeSlotStat struct {
	Period string `json:"period"` // morning/afternoon
	Count  int64  `json:"count"`
}

// GetDashboard 获取仪表盘数据
func (s *StatisticsService) GetDashboard() (*DashboardData, error) {
	db := database.GetDB()
	today := time.Now().Format("2006-01-02")

	var data DashboardData

	// 今日预约数
	db.Model(&model.Appointment{}).
		Where("DATE(appointment_date) = ?", today).
		Count(&data.TodayAppointments)

	// 今日签到数
	db.Model(&model.Appointment{}).
		Where("DATE(appointment_date) = ? AND status IN ?", today, []string{model.AppointmentStatusCheckedIn, model.AppointmentStatusCompleted}).
		Count(&data.TodayCheckins)

	// 总预约数
	db.Model(&model.Appointment{}).Count(&data.TotalAppointments)

	// 总用户数
	db.Model(&model.User{}).Count(&data.TotalUsers)

	// 总医生数
	db.Model(&model.Doctor{}).Count(&data.TotalDoctors)

	// 总科室数
	db.Model(&model.Department{}).Count(&data.TotalDepartments)

	// 最近10条预约
	var appointments []model.Appointment
	db.Model(&model.Appointment{}).
		Preload("Patient").
		Preload("Doctor").
		Preload("Department").
		Order("created_at DESC").
		Limit(10).
		Find(&appointments)

	data.RecentAppointments = make([]model.AppointmentListVO, len(appointments))
	for i, apt := range appointments {
		data.RecentAppointments[i] = *apt.ToListVO()
	}

	// 最近7天预约趋势
	data.AppointmentTrend = make([]TrendData, 7)
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		db.Model(&model.Appointment{}).
			Where("DATE(appointment_date) = ?", date).
			Count(&count)
		data.AppointmentTrend[6-i] = TrendData{
			Date:  date,
			Count: count,
		}
	}

	// 科室统计
	type DeptStat struct {
		DepartmentID   int64
		DepartmentName string
		DoctorCount    int64
		AppointmentCount int64
	}
	var deptStats []DeptStat
	db.Model(&model.Department{}).
		Select("departments.id as department_id, departments.name as department_name, "+
			"COUNT(DISTINCT doctors.id) as doctor_count, "+
			"COUNT(appointments.id) as appointment_count").
		Joins("LEFT JOIN doctors ON doctors.department_id = departments.id").
		Joins("LEFT JOIN appointments ON appointments.department_id = departments.id").
		Group("departments.id").
		Find(&deptStats)

	data.DepartmentStats = make([]DepartmentStat, len(deptStats))
	for i, stat := range deptStats {
		data.DepartmentStats[i] = DepartmentStat{
			DepartmentName:   stat.DepartmentName,
			DoctorCount:      stat.DoctorCount,
			AppointmentCount: stat.AppointmentCount,
		}
	}

	return &data, nil
}

// GetStatistics 获取统计数据
func (s *StatisticsService) GetStatistics(startDate, endDate string) (*StatisticsData, error) {
	db := database.GetDB()

	var data StatisticsData

	// 预约统计
	query := db.Model(&model.Appointment{})
	if startDate != "" && endDate != "" {
		query = query.Where("DATE(appointment_date) BETWEEN ? AND ?", startDate, endDate)
	}

	query.Count(&data.AppointmentStats.Total)
	query.Where("status = ?", model.AppointmentStatusPending).Count(&data.AppointmentStats.Pending)
	query.Where("status = ?", model.AppointmentStatusCheckedIn).Count(&data.AppointmentStats.CheckedIn)
	query.Where("status = ?", model.AppointmentStatusCompleted).Count(&data.AppointmentStats.Completed)
	query.Where("status = ?", model.AppointmentStatusCancelled).Count(&data.AppointmentStats.Cancelled)
	query.Where("status = ?", model.AppointmentStatusMissed).Count(&data.AppointmentStats.Missed)

	// 医生统计（TOP 10）
	type DocStat struct {
		DoctorID       int64
		DoctorName     string
		DepartmentName string
		AppointmentCount int64
		CompletedCount int64
	}
	var docStats []DocStat
	docQuery := db.Model(&model.Appointment{}).
		Select("doctors.id as doctor_id, doctors.name as doctor_name, "+
			"departments.name as department_name, "+
			"COUNT(appointments.id) as appointment_count, "+
			"SUM(CASE WHEN appointments.status = ? THEN 1 ELSE 0 END) as completed_count", model.AppointmentStatusCompleted).
		Joins("JOIN doctors ON doctors.id = appointments.doctor_id").
		Joins("JOIN departments ON departments.id = appointments.department_id")

	if startDate != "" && endDate != "" {
		docQuery = docQuery.Where("DATE(appointments.appointment_date) BETWEEN ? AND ?", startDate, endDate)
	}

	docQuery.Group("doctors.id").
		Order("appointment_count DESC").
		Limit(10).
		Find(&docStats)

	data.DoctorStats = make([]DoctorStat, len(docStats))
	for i, stat := range docStats {
		data.DoctorStats[i] = DoctorStat{
			DoctorID:         stat.DoctorID,
			DoctorName:       stat.DoctorName,
			DepartmentName:   stat.DepartmentName,
			AppointmentCount: stat.AppointmentCount,
			CompletedCount:   stat.CompletedCount,
		}
	}

	// 科室排行（TOP 10）
	type DeptRank struct {
		DepartmentID   int64
		DepartmentName string
		AppointmentCount int64
	}
	var deptRanks []DeptRank
	deptQuery := db.Model(&model.Appointment{}).
		Select("departments.id as department_id, departments.name as department_name, "+
			"COUNT(appointments.id) as appointment_count").
		Joins("JOIN departments ON departments.id = appointments.department_id")

	if startDate != "" && endDate != "" {
		deptQuery = deptQuery.Where("DATE(appointments.appointment_date) BETWEEN ? AND ?", startDate, endDate)
	}

	deptQuery.Group("departments.id").
		Order("appointment_count DESC").
		Limit(10).
		Find(&deptRanks)

	data.DepartmentRanking = make([]DepartmentRanking, len(deptRanks))
	for i, rank := range deptRanks {
		data.DepartmentRanking[i] = DepartmentRanking{
			DepartmentID:     rank.DepartmentID,
			DepartmentName:   rank.DepartmentName,
			AppointmentCount: rank.AppointmentCount,
		}
	}

	// 时段分布
	type TimeSlot struct {
		Period string
		Count  int64
	}
	var timeSlots []TimeSlot
	timeQuery := db.Model(&model.Appointment{}).
		Select("period, COUNT(*) as count")

	if startDate != "" && endDate != "" {
		timeQuery = timeQuery.Where("DATE(appointment_date) BETWEEN ? AND ?", startDate, endDate)
	}

	timeQuery.Group("period").Find(&timeSlots)

	data.TimeSlotDistribution = make([]TimeSlotStat, len(timeSlots))
	for i, slot := range timeSlots {
		data.TimeSlotDistribution[i] = TimeSlotStat{
			Period: slot.Period,
			Count:  slot.Count,
		}
	}

	return &data, nil
}
