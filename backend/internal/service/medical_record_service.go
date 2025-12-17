package service

import (
	"errors"

	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/internal/repository"
	"huaan-medical/pkg/errorcode"
)

// MedicalRecordService 就诊记录服务
type MedicalRecordService struct {
	repo *repository.MedicalRecordRepository
}

// NewMedicalRecordService 创建就诊记录服务实例
func NewMedicalRecordService() *MedicalRecordService {
	return &MedicalRecordService{
		repo: repository.NewMedicalRecordRepository(),
	}
}

// GetByID 获取就诊记录详情（需要权限校验）
func (s *MedicalRecordService) GetByID(userID, recordID int64) (*model.MedicalRecordVO, error) {
	record, err := s.repo.GetByUserAndID(userID, recordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrRecordNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	return record.ToVO(), nil
}

// ListByUser 查询用户的就诊记录列表
func (s *MedicalRecordService) ListByUser(userID int64) ([]model.MedicalRecordListVO, error) {
	records, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	voList := make([]model.MedicalRecordListVO, len(records))
	for i, record := range records {
		voList[i] = *record.ToListVO()
	}

	return voList, nil
}
