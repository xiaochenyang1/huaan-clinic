package model

// Doctor 医生模型
type Doctor struct {
	BaseModel
	DepartmentID int64  `gorm:"index;not null;comment:所属科室ID" json:"department_id"`
	Name         string `gorm:"type:varchar(32);not null;comment:姓名" json:"name"`
	Avatar       string `gorm:"type:varchar(512);comment:头像URL" json:"avatar"`
	Title        string `gorm:"type:varchar(32);not null;comment:职称" json:"title"`
	Specialty    string `gorm:"type:varchar(256);comment:擅长领域" json:"specialty"`
	Introduction string `gorm:"type:text;comment:个人简介" json:"introduction"`
	SortOrder    int    `gorm:"type:int;default:0;comment:排序序号" json:"sort_order"`
	Status       int    `gorm:"type:tinyint;default:1;comment:状态 0停诊 1正常" json:"status"`

	// 关联
	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
}

// TableName 表名
func (Doctor) TableName() string {
	return "doctors"
}

// DoctorVO 医生视图对象
type DoctorVO struct {
	ID             int64  `json:"id"`
	DepartmentID   int64  `json:"department_id"`
	DepartmentName string `json:"department_name,omitempty"`
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
	Title          string `json:"title"`
	TitleName      string `json:"title_name"`
	Specialty      string `json:"specialty"`
	Introduction   string `json:"introduction,omitempty"`
	SortOrder      int    `json:"sort_order"`
	Status         int    `json:"status"`
	StatusName     string `json:"status_name"`
}

// ToVO 转换为视图对象
func (d *Doctor) ToVO() *DoctorVO {
	statusName := "正常"
	if d.Status == StatusDisabled {
		statusName = "停诊"
	}

	vo := &DoctorVO{
		ID:           d.ID,
		DepartmentID: d.DepartmentID,
		Name:         d.Name,
		Avatar:       d.Avatar,
		Title:        d.Title,
		TitleName:    GetTitleName(d.Title),
		Specialty:    d.Specialty,
		Introduction: d.Introduction,
		SortOrder:    d.SortOrder,
		Status:       d.Status,
		StatusName:   statusName,
	}

	if d.Department != nil {
		vo.DepartmentName = d.Department.Name
	}

	return vo
}

// DoctorListVO 医生列表视图对象（简化版）
type DoctorListVO struct {
	ID             int64  `json:"id"`
	DepartmentID   int64  `json:"department_id"`
	DepartmentName string `json:"department_name,omitempty"`
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
	Title          string `json:"title"`
	TitleName      string `json:"title_name"`
	Specialty      string `json:"specialty"`
	Status         int    `json:"status"`
	StatusName     string `json:"status_name"`
}

// ToListVO 转换为列表视图对象
func (d *Doctor) ToListVO() *DoctorListVO {
	statusName := "正常"
	if d.Status == StatusDisabled {
		statusName = "停诊"
	}

	vo := &DoctorListVO{
		ID:           d.ID,
		DepartmentID: d.DepartmentID,
		Name:         d.Name,
		Avatar:       d.Avatar,
		Title:        d.Title,
		TitleName:    GetTitleName(d.Title),
		Specialty:    d.Specialty,
		Status:       d.Status,
		StatusName:   statusName,
	}

	if d.Department != nil {
		vo.DepartmentName = d.Department.Name
	}

	return vo
}
