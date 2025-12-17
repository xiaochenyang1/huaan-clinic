package service

import (
	"errors"

	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/internal/repository"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/utils"
)

// PatientService 就诊人服务
type PatientService struct {
	repo *repository.PatientRepository
}

// NewPatientService 创建就诊人服务实例
func NewPatientService() *PatientService {
	return &PatientService{
		repo: repository.NewPatientRepository(),
	}
}

// CreatePatientRequest 创建就诊人请求
type CreatePatientRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=32"`
	IDCard    string `json:"id_card" binding:"required,len=18"`
	Phone     string `json:"phone" binding:"required,len=11"`
	Relation  string `json:"relation" binding:"required,oneof=self parent child spouse other"`
	IsDefault int    `json:"is_default" binding:"oneof=0 1"`
}

// UpdatePatientRequest 更新就诊人请求
type UpdatePatientRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=32"`
	IDCard    string `json:"id_card" binding:"required,len=18"`
	Phone     string `json:"phone" binding:"required,len=11"`
	Relation  string `json:"relation" binding:"required,oneof=self parent child spouse other"`
	IsDefault int    `json:"is_default" binding:"oneof=0 1"`
}

// Create 创建就诊人
func (s *PatientService) Create(userID int64, req *CreatePatientRequest) (*model.PatientVO, error) {
	// 验证身份证号格式
	if !utils.ValidateIDCard(req.IDCard) {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "身份证号格式错误")
	}

	// 验证手机号格式
	if !utils.ValidatePhone(req.Phone) {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "手机号格式错误")
	}

	// 检查就诊人数量限制（每个用户最多10个就诊人）
	count, err := s.repo.CountByUser(userID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	if count >= 10 {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "最多只能添加10个就诊人")
	}

	// 检查身份证号是否重复（同一用户下）
	exists, err := s.repo.ExistsByIDCard(userID, req.IDCard)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	if exists {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "该身份证号已添加")
	}

	// 从身份证号提取信息
	gender := utils.GetGenderFromIDCard(req.IDCard)
	birthDate := utils.GetBirthDateFromIDCard(req.IDCard)

	// 如果是第一个就诊人，自动设为默认
	if count == 0 {
		req.IsDefault = 1
	}

	// 如果要设置为默认，先清除其他默认标记
	if req.IsDefault == 1 {
		if err := s.repo.ClearDefaultByUser(userID); err != nil {
			return nil, errorcode.New(errorcode.ErrDatabase)
		}
	}

	// 创建就诊人
	patient := &model.Patient{
		UserID:    userID,
		Name:      req.Name,
		IDCard:    req.IDCard,
		Phone:     req.Phone,
		Gender:    gender,
		BirthDate: birthDate,
		Relation:  req.Relation,
		IsDefault: req.IsDefault,
	}

	if err := s.repo.Create(patient); err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 重新查询以获取完整数据
	patient, err = s.repo.GetByID(patient.ID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	return patient.ToFullVO(), nil
}

// Update 更新就诊人
func (s *PatientService) Update(userID, patientID int64, req *UpdatePatientRequest) (*model.PatientVO, error) {
	// 验证身份证号格式
	if !utils.ValidateIDCard(req.IDCard) {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "身份证号格式错误")
	}

	// 验证手机号格式
	if !utils.ValidatePhone(req.Phone) {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "手机号格式错误")
	}

	// 检查就诊人是否存在且属于该用户
	patient, err := s.repo.GetByUserAndID(userID, patientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrPatientNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 检查身份证号是否重复（同一用户下，排除自己）
	exists, err := s.repo.ExistsByIDCard(userID, req.IDCard, patientID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	if exists {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "该身份证号已添加")
	}

	// 从身份证号提取信息
	gender := utils.GetGenderFromIDCard(req.IDCard)
	birthDate := utils.GetBirthDateFromIDCard(req.IDCard)

	// 如果要设置为默认，先清除其他默认标记
	if req.IsDefault == 1 && patient.IsDefault == 0 {
		if err := s.repo.ClearDefaultByUser(userID); err != nil {
			return nil, errorcode.New(errorcode.ErrDatabase)
		}
	}

	// 更新就诊人信息
	patient.Name = req.Name
	patient.IDCard = req.IDCard
	patient.Phone = req.Phone
	patient.Gender = gender
	patient.BirthDate = birthDate
	patient.Relation = req.Relation
	patient.IsDefault = req.IsDefault

	if err := s.repo.Update(patient); err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 重新查询以获取完整数据
	patient, err = s.repo.GetByID(patientID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	return patient.ToFullVO(), nil
}

// Delete 删除就诊人
func (s *PatientService) Delete(userID, patientID int64) error {
	// 检查就诊人是否存在且属于该用户
	patient, err := s.repo.GetByUserAndID(userID, patientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorcode.New(errorcode.ErrPatientNotFound)
		}
		return errorcode.New(errorcode.ErrDatabase)
	}

	// 检查是否有预约记录
	hasAppointments, err := s.repo.HasAppointments(patientID)
	if err != nil {
		return errorcode.New(errorcode.ErrDatabase)
	}
	if hasAppointments {
		return errorcode.NewWithMessage(errorcode.ErrInvalidParams, "该就诊人存在预约记录，无法删除")
	}

	// 删除就诊人
	if err := s.repo.Delete(patientID); err != nil {
		return errorcode.New(errorcode.ErrDatabase)
	}

	// 如果删除的是默认就诊人，自动设置第一个为默认
	if patient.IsDefault == 1 {
		patients, err := s.repo.ListByUser(userID)
		if err == nil && len(patients) > 0 {
			_ = s.repo.SetDefault(userID, patients[0].ID)
		}
	}

	return nil
}

// GetByID 获取就诊人详情
func (s *PatientService) GetByID(userID, patientID int64) (*model.PatientVO, error) {
	// 检查就诊人是否存在且属于该用户
	patient, err := s.repo.GetByUserAndID(userID, patientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrPatientNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	return patient.ToFullVO(), nil
}

// List 查询用户的就诊人列表
func (s *PatientService) List(userID int64) ([]model.PatientVO, error) {
	patients, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	voList := make([]model.PatientVO, len(patients))
	for i, patient := range patients {
		voList[i] = *patient.ToVO()
	}

	return voList, nil
}

// SetDefault 设置默认就诊人
func (s *PatientService) SetDefault(userID, patientID int64) error {
	// 检查就诊人是否存在且属于该用户
	_, err := s.repo.GetByUserAndID(userID, patientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorcode.New(errorcode.ErrPatientNotFound)
		}
		return errorcode.New(errorcode.ErrDatabase)
	}

	return s.repo.SetDefault(userID, patientID)
}

// ListAdminPatientsRequest 管理后台患者列表查询请求
type ListAdminPatientsRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Keyword  string `form:"keyword"`
}

// ListAdmin 分页查询患者列表（管理后台）
func (s *PatientService) ListAdmin(req *ListAdminPatientsRequest) ([]model.PatientVO, int64, error) {
	patients, total, err := s.repo.List(req.Page, req.PageSize, req.Keyword)
	if err != nil {
		return nil, 0, errorcode.New(errorcode.ErrDatabase)
	}

	voList := make([]model.PatientVO, len(patients))
	for i, patient := range patients {
		vo := patient.ToVO()
		// 填充用户的最后登录IP
		if patient.User != nil {
			vo.LastLoginIP = patient.User.LastLoginIP
		}
		voList[i] = *vo
	}

	return voList, total, nil
}

// GetByIDAdmin 获取患者详情（管理后台，无需权限校验）
func (s *PatientService) GetByIDAdmin(patientID int64) (*model.PatientVO, error) {
	patient, err := s.repo.GetByID(patientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrPatientNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	return patient.ToVO(), nil
}
