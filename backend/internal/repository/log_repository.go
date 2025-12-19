package repository

import (
	"time"

	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/pkg/database"
)

// LogRepository 日志仓库
type LogRepository struct {
	db *gorm.DB
}

func NewLogRepository() *LogRepository {
	return &LogRepository{db: database.GetDB()}
}

func (r *LogRepository) CreateOperationLog(log *model.OperationLog) error {
	return r.db.Create(log).Error
}

func (r *LogRepository) ListOperationLogs(page, pageSize int, adminID *int64, module, keyword string, startDate, endDate *time.Time) ([]model.OperationLog, int64, error) {
	var list []model.OperationLog
	var total int64

	query := r.db.Model(&model.OperationLog{})
	if adminID != nil && *adminID > 0 {
		query = query.Where("admin_id = ?", *adminID)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			r.db.Where("admin_name LIKE ?", like).
				Or("path LIKE ?", like).
				Or("module LIKE ?", like).
				Or("action LIKE ?", like),
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *LogRepository) CreateLoginLog(log *model.LoginLog) error {
	return r.db.Create(log).Error
}

func (r *LogRepository) ListLoginLogs(page, pageSize int, userType, keyword string, status *int, startDate, endDate *time.Time) ([]model.LoginLog, int64, error) {
	var list []model.LoginLog
	var total int64

	query := r.db.Model(&model.LoginLog{})
	if userType != "" {
		query = query.Where("user_type = ?", userType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			r.db.Where("username LIKE ?", like).
				Or("ip LIKE ?", like).
				Or("device LIKE ?", like),
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

