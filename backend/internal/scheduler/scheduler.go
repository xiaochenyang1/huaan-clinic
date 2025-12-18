package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"huaan-medical/internal/model"
	"huaan-medical/pkg/config"
	"huaan-medical/pkg/database"
	"huaan-medical/pkg/logger"
	"huaan-medical/pkg/redis"
	"huaan-medical/pkg/wechat"
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

	cfg, err := config.Load("config.yaml")
	if err != nil {
		logger.Error("加载配置失败", zap.Error(err))
		return
	}

	tplID := cfg.WeChat.Subscribe.AppointmentReminderTemplateID
	if cfg.WeChat.AppID == "" || cfg.WeChat.AppSecret == "" || tplID == "" {
		logger.Info("微信订阅消息未配置，跳过推送")
		return
	}

	// 查询明天所有待就诊的预约
	var appointments []model.Appointment
	err = db.Model(&model.Appointment{}).
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

	wechatClient := wechat.NewClient(cfg.WeChat.AppID, cfg.WeChat.AppSecret)
	page := cfg.WeChat.Subscribe.AppointmentReminderPage

	var okCount int
	var failCount int
	for _, a := range appointments {
		openID := ""
		if a.User != nil {
			openID = a.User.OpenID
		}
		if openID == "" {
			failCount++
			continue
		}

		// 注意：订阅消息模板字段由你在微信后台创建的模板决定。
		// 这里使用常见字段名示例（thing1/time2/thing3），如不匹配会发送失败。
		data := map[string]interface{}{
			"thing1": map[string]string{"value": a.Patient.Name},                                                               // 就诊人
			"time2":  map[string]string{"value": a.AppointmentDate.Format("2006-01-02") + " " + model.GetPeriodName(a.Period)}, // 时间
			"thing3": map[string]string{"value": a.Department.Name + " " + a.Doctor.Name},                                      // 科室/医生
		}

		req := &wechat.SubscribeMessageRequest{
			ToUser:     openID,
			TemplateID: tplID,
			Page:       page,
			Data:       data,
		}

		if err := wechatClient.SendSubscribeMessage(req); err != nil {
			failCount++
			logger.Warn("发送就诊提醒失败", zap.Error(err), zap.Int64("appointment_id", a.ID))
			continue
		}
		okCount++
	}

	logger.Info("发送就诊提醒完成", zap.Int("total", len(appointments)), zap.Int("ok", okCount), zap.Int("fail", failCount))
}

// cleanExpiredTokens 清理过期Token
// 每小时执行一次（预留功能）
func cleanExpiredTokens() {
	logger.Info("开始清理过期Token")

	if !redis.IsEnabled() {
		logger.Info("Redis未启用，跳过清理")
		return
	}

	// Redis会自动清理带 TTL 的 key；这里额外处理“无过期时间”的黑名单 key（通常属于异常数据）。
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := redis.GetClient()
	pattern := fmt.Sprintf(redis.KeyTokenBlacklist, "*")
	var cursor uint64

	var scanned int
	var deleted int
	var noExpire int

	for {
		keys, nextCursor, err := client.Scan(ctx, cursor, pattern, 200).Result()
		if err != nil {
			logger.Warn("扫描Token黑名单失败", zap.Error(err))
			break
		}

		scanned += len(keys)
		for _, key := range keys {
			ttl, err := client.TTL(ctx, key).Result()
			if err != nil {
				continue
			}
			// -1: 永不过期（异常）；-2: 不存在
			if ttl == -1 {
				noExpire++
				_ = client.Del(ctx, key).Err()
				deleted++
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	logger.Info("清理过期Token完成", zap.Int("scanned", scanned), zap.Int("deleted", deleted), zap.Int("no_expire", noExpire))
}
