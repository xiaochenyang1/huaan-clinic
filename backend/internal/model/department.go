package model

// Department 科室模型
type Department struct {
	BaseModel
	Name        string `gorm:"type:varchar(64);not null;comment:科室名称" json:"name"`
	Description string `gorm:"type:varchar(512);comment:科室描述" json:"description"`
	Icon        string `gorm:"type:varchar(256);comment:科室图标" json:"icon"`
	SortOrder   int    `gorm:"type:int;default:0;comment:排序序号" json:"sort_order"`
	Status      int    `gorm:"type:tinyint;default:1;comment:状态 0停用 1启用" json:"status"`

	// 关联
	Doctors []Doctor `gorm:"foreignKey:DepartmentID" json:"doctors,omitempty"`
}

// TableName 表名
func (Department) TableName() string {
	return "departments"
}

// DepartmentVO 科室视图对象
type DepartmentVO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
	Status      int    `json:"status"`
	StatusName  string `json:"status_name"`
	DoctorCount int    `json:"doctor_count,omitempty"` // 医生数量
}

// ToVO 转换为视图对象
func (d *Department) ToVO() *DepartmentVO {
	statusName := "启用"
	if d.Status == StatusDisabled {
		statusName = "停用"
	}
	return &DepartmentVO{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Icon:        d.Icon,
		SortOrder:   d.SortOrder,
		Status:      d.Status,
		StatusName:  statusName,
		DoctorCount: len(d.Doctors),
	}
}
