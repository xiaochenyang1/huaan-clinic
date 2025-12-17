package service

import (
	"errors"

	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/internal/repository"
	"huaan-medical/pkg/errorcode"
)

// DoctorService 医生服务
type DoctorService struct {
	repo     *repository.DoctorRepository
	deptRepo *repository.DepartmentRepository
}

// NewDoctorService 创建医生服务实例
func NewDoctorService() *DoctorService {
	return &DoctorService{
		repo:     repository.NewDoctorRepository(),
		deptRepo: repository.NewDepartmentRepository(),
	}
}

// CreateDoctorRequest 创建医生请求
type CreateDoctorRequest struct {
	DepartmentID int64  `json:"department_id" binding:"required,min=1"`
	Name         string `json:"name" binding:"required,min=2,max=32"`
	Avatar       string `json:"avatar" binding:"max=512"`
	Title        string `json:"title" binding:"required,oneof=chief_physician associate_chief_physician attending_physician resident_physician"`
	Specialty    string `json:"specialty" binding:"max=256"`
	Introduction string `json:"introduction" binding:"max=2000"`
	SortOrder    int    `json:"sort_order"`
	Status       int    `json:"status" binding:"oneof=0 1"`
}

// UpdateDoctorRequest 更新医生请求
type UpdateDoctorRequest struct {
	DepartmentID int64  `json:"department_id" binding:"required,min=1"`
	Name         string `json:"name" binding:"required,min=2,max=32"`
	Avatar       string `json:"avatar" binding:"max=512"`
	Title        string `json:"title" binding:"required,oneof=chief_physician associate_chief_physician attending_physician resident_physician"`
	Specialty    string `json:"specialty" binding:"max=256"`
	Introduction string `json:"introduction" binding:"max=2000"`
	SortOrder    int    `json:"sort_order"`
	Status       int    `json:"status" binding:"oneof=0 1"`
}

// ListDoctorRequest 列表查询请求
type ListDoctorRequest struct {
	Page         int    `form:"page" binding:"required,min=1"`
	PageSize     int    `form:"page_size" binding:"required,min=1,max=100"`
	DepartmentID *int64 `form:"department_id"`
	Status       *int   `form:"status"`
	Keyword      string `form:"keyword"`
}

// ListPublicDoctorRequest 公开接口列表查询请求
type ListPublicDoctorRequest struct {
	DepartmentID *int64 `form:"department_id"`
	Keyword      string `form:"keyword"`
}

// Create 创建医生
func (s *DoctorService) Create(req *CreateDoctorRequest) (*model.DoctorVO, error) {
	// 检查科室是否存在
	dept, err := s.deptRepo.GetByID(req.DepartmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrDepartmentNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 检查科室是否启用
	if dept.Status == model.StatusDisabled {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "所属科室已停用")
	}

	// 检查医生姓名是否重复（同一科室下）
	exists, err := s.repo.ExistsByName(req.Name, req.DepartmentID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	if exists {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "该科室下已存在同名医生")
	}

	// 创建医生
	doctor := &model.Doctor{
		DepartmentID: req.DepartmentID,
		Name:         req.Name,
		Avatar:       req.Avatar,
		Title:        req.Title,
		Specialty:    req.Specialty,
		Introduction: req.Introduction,
		SortOrder:    req.SortOrder,
		Status:       req.Status,
	}

	if err := s.repo.Create(doctor); err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 重新查询以获取关联数据
	doctor, err = s.repo.GetByID(doctor.ID)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	return doctor.ToVO(), nil
}

// Update 更新医生
func (s *DoctorService) Update(id int64, req *UpdateDoctorRequest) (*model.DoctorVO, error) {
	// 检查医生是否存在
	doctor, err := s.repo.GetByIDSimple(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrDoctorNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 检查科室是否存在
	dept, err := s.deptRepo.GetByID(req.DepartmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrDepartmentNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 检查科室是否启用
	if dept.Status == model.StatusDisabled {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "所属科室已停用")
	}

	// 检查医生姓名是否重复（同一科室下，排除自己）
	exists, err := s.repo.ExistsByName(req.Name, req.DepartmentID, id)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	if exists {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "该科室下已存在同名医生")
	}

	// 更新医生信息
	doctor.DepartmentID = req.DepartmentID
	doctor.Name = req.Name
	doctor.Avatar = req.Avatar
	doctor.Title = req.Title
	doctor.Specialty = req.Specialty
	doctor.Introduction = req.Introduction
	doctor.SortOrder = req.SortOrder
	doctor.Status = req.Status

	if err := s.repo.Update(doctor); err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 重新查询以获取关联数据
	doctor, err = s.repo.GetByID(id)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	return doctor.ToVO(), nil
}

// Delete 删除医生
func (s *DoctorService) Delete(id int64) error {
	// 检查医生是否存在
	_, err := s.repo.GetByIDSimple(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorcode.New(errorcode.ErrDoctorNotFound)
		}
		return errorcode.New(errorcode.ErrDatabase)
	}

	// 检查是否有预约记录
	hasAppointments, err := s.repo.HasAppointments(id)
	if err != nil {
		return errorcode.New(errorcode.ErrDatabase)
	}
	if hasAppointments {
		return errorcode.NewWithMessage(errorcode.ErrInvalidParams, "该医生存在预约记录，无法删除")
	}

	// 检查是否有排班记录
	hasSchedules, err := s.repo.HasSchedules(id)
	if err != nil {
		return errorcode.New(errorcode.ErrDatabase)
	}
	if hasSchedules {
		return errorcode.NewWithMessage(errorcode.ErrInvalidParams, "该医生存在排班记录，请先删除排班")
	}

	return s.repo.Delete(id)
}

// GetByID 获取医生详情
func (s *DoctorService) GetByID(id int64) (*model.DoctorVO, error) {
	doctor, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrDoctorNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	return doctor.ToVO(), nil
}

// List 分页查询医生列表（管理后台）
func (s *DoctorService) List(req *ListDoctorRequest) ([]model.DoctorListVO, int64, error) {
	doctors, total, err := s.repo.List(req.Page, req.PageSize, req.DepartmentID, req.Status, req.Keyword)
	if err != nil {
		return nil, 0, errorcode.New(errorcode.ErrDatabase)
	}

	voList := make([]model.DoctorListVO, len(doctors))
	for i, doctor := range doctors {
		voList[i] = *doctor.ToListVO()
	}

	return voList, total, nil
}

// ListPublic 查询医生列表（公开接口）
func (s *DoctorService) ListPublic(req *ListPublicDoctorRequest) ([]model.DoctorListVO, error) {
	doctors, err := s.repo.ListPublic(req.DepartmentID, req.Keyword)
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	voList := make([]model.DoctorListVO, len(doctors))
	for i, doctor := range doctors {
		voList[i] = *doctor.ToListVO()
	}

	return voList, nil
}

// UpdateStatus 批量更新医生状态
func (s *DoctorService) UpdateStatus(ids []int64, status int) error {
	if len(ids) == 0 {
		return errorcode.NewWithMessage(errorcode.ErrInvalidParams, "请选择要操作的医生")
	}

	if status != model.StatusEnabled && status != model.StatusDisabled {
		return errorcode.NewWithMessage(errorcode.ErrInvalidParams, "状态值无效")
	}

	return s.repo.UpdateStatus(ids, status)
}
