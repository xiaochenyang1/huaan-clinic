package service

import (
	"time"

	"huaan-medical/internal/model"
	"huaan-medical/internal/repository"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/utils"
)

// LogService 日志服务
type LogService struct {
	repo *repository.LogRepository
}

func NewLogService() *LogService {
	return &LogService{repo: repository.NewLogRepository()}
}

type ListOperationLogsRequest struct {
	Page      int    `form:"page" binding:"required,min=1"`
	PageSize  int    `form:"page_size" binding:"required,min=1,max=100"`
	AdminID   int64  `form:"admin_id"`
	Module    string `form:"module"`
	Keyword   string `form:"keyword"`
	StartDate string `form:"start_date"` // YYYY-MM-DD
	EndDate   string `form:"end_date"`   // YYYY-MM-DD
}

func (s *LogService) ListOperationLogs(req *ListOperationLogsRequest) ([]model.OperationLogVO, int64, error) {
	var adminID *int64
	if req.AdminID > 0 {
		adminID = &req.AdminID
	}

	var startDate *time.Time
	var endDate *time.Time
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
		// 结束日期包含当天
		tmp := ed.Add(24*time.Hour - time.Nanosecond)
		endDate = &tmp
	}

	logs, total, err := s.repo.ListOperationLogs(req.Page, req.PageSize, adminID, req.Module, req.Keyword, startDate, endDate)
	if err != nil {
		return nil, 0, errorcode.New(errorcode.ErrDatabase)
	}

	voList := make([]model.OperationLogVO, len(logs))
	for i := range logs {
		voList[i] = *logs[i].ToVO()
	}
	return voList, total, nil
}

type ListLoginLogsRequest struct {
	Page      int    `form:"page" binding:"required,min=1"`
	PageSize  int    `form:"page_size" binding:"required,min=1,max=100"`
	UserType  string `form:"user_type"` // user/admin
	Keyword   string `form:"keyword"`
	Status    *int   `form:"status"` // 0/1
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

func (s *LogService) ListLoginLogs(req *ListLoginLogsRequest) ([]model.LoginLogVO, int64, error) {
	var startDate *time.Time
	var endDate *time.Time
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
		tmp := ed.Add(24*time.Hour - time.Nanosecond)
		endDate = &tmp
	}

	logs, total, err := s.repo.ListLoginLogs(req.Page, req.PageSize, req.UserType, req.Keyword, req.Status, startDate, endDate)
	if err != nil {
		return nil, 0, errorcode.New(errorcode.ErrDatabase)
	}

	voList := make([]model.LoginLogVO, len(logs))
	for i := range logs {
		voList[i] = *logs[i].ToVO()
	}
	return voList, total, nil
}

