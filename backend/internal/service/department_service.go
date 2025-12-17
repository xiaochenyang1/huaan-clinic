package service

import (
	"errors"

	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/internal/repository"
	"huaan-medical/pkg/errorcode"
)

// DepartmentService 科室服务
type DepartmentService struct {
	repo *repository.DepartmentRepository
}

// NewDepartmentService 创建科室服务实例
func NewDepartmentService() *DepartmentService {
	return &DepartmentService{
		repo: repository.NewDepartmentRepository(),
	}
}

// CreateRequest 创建科室请求
type CreateDepartmentRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=64"`
	Description string `json:"description" binding:"max=512"`
	Icon        string `json:"icon" binding:"max=256"`
	SortOrder   int    `json:"sort_order"`
	Status      int    `json:"status"`
}

// UpdateRequest 更新科室请求
type UpdateDepartmentRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=64"`
	Description string `json:"description" binding:"max=512"`
	Icon        string `json:"icon" binding:"max=256"`
	SortOrder   int    `json:"sort_order"`
	Status      int    `json:"status"`
}

// ListRequest 列表查询请求
type ListDepartmentRequest struct {
	Page     int  `form:"page" binding:"required,min=1"`
	PageSize int  `form:"page_size" binding:"required,min=1,max=100"`
	Status   *int `form:"status"`
}

// Create 创建科室
func (s *DepartmentService) Create(req *CreateDepartmentRequest) (*model.DepartmentVO, error) {
	// 检查名称是否重复
	existing, _ := s.repo.GetByName(req.Name)
	if existing != nil {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "科室名称已存在")
	}

	dept := &model.Department{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	}

	if err := s.repo.Create(dept); err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	return dept.ToVO(), nil
}

// Update 更新科室
func (s *DepartmentService) Update(id int64, req *UpdateDepartmentRequest) (*model.DepartmentVO, error) {
	dept, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrDepartmentNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	// 检查名称是否重复（排除自己）
	existing, _ := s.repo.GetByName(req.Name)
	if existing != nil && existing.ID != id {
		return nil, errorcode.NewWithMessage(errorcode.ErrInvalidParams, "科室名称已存在")
	}

	dept.Name = req.Name
	dept.Description = req.Description
	dept.Icon = req.Icon
	dept.SortOrder = req.SortOrder
	dept.Status = req.Status

	if err := s.repo.Update(dept); err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	return dept.ToVO(), nil
}

// Delete 删除科室
func (s *DepartmentService) Delete(id int64) error {
	// 检查科室是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorcode.New(errorcode.ErrDepartmentNotFound)
		}
		return errorcode.New(errorcode.ErrDatabase)
	}

	// 检查是否有医生
	hasDoctors, err := s.repo.HasDoctors(id)
	if err != nil {
		return errorcode.New(errorcode.ErrDatabase)
	}
	if hasDoctors {
		return errorcode.New(errorcode.ErrDepartmentHasDoctor)
	}

	return s.repo.Delete(id)
}

// GetByID 获取科室详情
func (s *DepartmentService) GetByID(id int64) (*model.DepartmentVO, error) {
	dept, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.New(errorcode.ErrDepartmentNotFound)
		}
		return nil, errorcode.New(errorcode.ErrDatabase)
	}
	return dept.ToVO(), nil
}

// List 分页查询科室列表（管理后台）
func (s *DepartmentService) List(req *ListDepartmentRequest) ([]model.DepartmentVO, int64, error) {
	departments, total, err := s.repo.List(req.Page, req.PageSize, req.Status)
	if err != nil {
		return nil, 0, errorcode.New(errorcode.ErrDatabase)
	}

	voList := make([]model.DepartmentVO, len(departments))
	for i, dept := range departments {
		voList[i] = *dept.ToVO()
	}

	return voList, total, nil
}

// ListAll 查询所有启用的科室（公开接口）
func (s *DepartmentService) ListAll() ([]model.DepartmentVO, error) {
	departments, err := s.repo.ListAll()
	if err != nil {
		return nil, errorcode.New(errorcode.ErrDatabase)
	}

	voList := make([]model.DepartmentVO, len(departments))
	for i, dept := range departments {
		voList[i] = *dept.ToVO()
	}

	return voList, nil
}
