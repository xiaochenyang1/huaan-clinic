package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"huaan-medical/internal/model"
	"huaan-medical/pkg/database"
	"huaan-medical/pkg/logger"
)

var cronJob *cron.Cron

// Init 初始化定时任务
func Init() {
	cronJob = cron.New(cron.WithSeconds())

	// 每天22:00处理爽约
	cronJob.AddFunc("0 0 22 * * *", handleMissedAppointments)

	// 每天20:00推送明日就诊提醒（预留功能）
	cronJob.AddFunc("0 0 20 * * *", sendAppointmentReminders)

	// 每小时清理过期Token（预留功能）
	cronJob.AddFunc("0 0 * * * *", cleanExpiredTokens)

	cronJob.Start()
	logger.Info("定时任务已启动")
}

// Stop 停止定时任务
func Stop() {
	if cronJob != nil {
		cronJob.Stop()
		logger.Info("定时任务已停止")
	}
}

// handleMissedAppointments 处理爽约预约
// 每天22:00执行，将今天所有状态为"待就诊"的预约标记为"爽约"
func handleMissedAppointments() {
	logger.Info("开始处理爽约预约")

	db := database.GetDB()
	today := time.Now().Format("2006-01-02")

	// 更新今天所有待就诊的预约为爽约状态
	result := db.Model(&model.Appointment{}).
		Where("DATE(appointment_date) = ? AND status = ?", today, model.AppointmentStatusPending).
		Updates(map[string]interface{}{
			"status": model.AppointmentStatusMissed,
		})

	if result.Error != nil {
		logger.Error("处理爽约预约失败", zap.Error(result.Error))
		return
	}

	logger.Info("处理爽约预约完成", zap.Int64("count", result.RowsAffected))
}

// sendAppointmentReminders 发送就诊提醒
// 每天20:00执行，向明天有预约的用户发送提醒（预留功能）
func sendAppointmentReminders() {
	logger.Info("开始发送就诊提醒")

	db := database.GetDB()
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	// 查询明天所有待就诊的预约
	var appointments []model.Appointment
	err := db.Model(&model.Appointment{}).
		Preload("User").
		Preload("Patient").
		Preload("Doctor").
		Preload("Department").
		Where("DATE(appointment_date) = ? AND status = ?", tomorrow, model.AppointmentStatusPending).
		Find(&appointments).Error

	if err != nil {
		logger.Error("查询明日预约失败", zap.Error(err))
		return
	}

	// TODO: 实现微信模板消息推送
	// 这里需要调用微信API发送模板消息
	// 每个预约发送一条提醒消息给用户

	logger.Info("发送就诊提醒完成", zap.Int("count", len(appointments)))
}

// cleanExpiredTokens 清理过期Token
// 每小时执行一次（预留功能）
func cleanExpiredTokens() {
	logger.Info("开始清理过期Token")

	// TODO: 清理Redis中的过期Token黑名单等
	// Redis会自动清理过期key，这里可以做一些额外的清理工作

	logger.Info("清理过期Token完成")
}
