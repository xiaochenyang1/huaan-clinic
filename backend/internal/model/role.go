package model

// Role 角色模型
type Role struct {
	BaseModel
	Code        string `gorm:"type:varchar(32);uniqueIndex;not null;comment:角色编码" json:"code"`
	Name        string `gorm:"type:varchar(64);not null;comment:角色名称" json:"name"`
	Description string `gorm:"type:varchar(256);comment:角色描述" json:"description"`
	SortOrder   int    `gorm:"type:int;default:0;comment:排序序号" json:"sort_order"`
	Status      int    `gorm:"type:tinyint;default:1;comment:状态 0禁用 1启用" json:"status"`

	// 关联
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Admins      []Admin      `gorm:"many2many:admin_roles;" json:"admins,omitempty"`
}

// TableName 表名
func (Role) TableName() string {
	return "roles"
}

// Permission 权限模型
type Permission struct {
	BaseModel
	Code        string `gorm:"type:varchar(64);uniqueIndex;not null;comment:权限编码" json:"code"`
	Name        string `gorm:"type:varchar(64);not null;comment:权限名称" json:"name"`
	Module      string `gorm:"type:varchar(32);index;not null;comment:所属模块" json:"module"`
	Description string `gorm:"type:varchar(256);comment:权限描述" json:"description"`
	SortOrder   int    `gorm:"type:int;default:0;comment:排序序号" json:"sort_order"`

	// 关联
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// TableName 表名
func (Permission) TableName() string {
	return "permissions"
}

// AdminRole 管理员角色关联表
type AdminRole struct {
	AdminID int64 `gorm:"primaryKey" json:"admin_id"`
	RoleID  int64 `gorm:"primaryKey" json:"role_id"`
}

// TableName 表名
func (AdminRole) TableName() string {
	return "admin_roles"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID       int64 `gorm:"primaryKey" json:"role_id"`
	PermissionID int64 `gorm:"primaryKey" json:"permission_id"`
}

// TableName 表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// RoleVO 角色视图对象
type RoleVO struct {
	ID          int64    `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	SortOrder   int      `json:"sort_order"`
	Status      int      `json:"status"`
	StatusName  string   `json:"status_name"`
	Permissions []string `json:"permissions,omitempty"`
}

// ToVO 转换为视图对象
func (r *Role) ToVO() *RoleVO {
	statusName := "启用"
	if r.Status == StatusDisabled {
		statusName = "禁用"
	}

	vo := &RoleVO{
		ID:          r.ID,
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		SortOrder:   r.SortOrder,
		Status:      r.Status,
		StatusName:  statusName,
	}

	if len(r.Permissions) > 0 {
		vo.Permissions = make([]string, len(r.Permissions))
		for i, p := range r.Permissions {
			vo.Permissions[i] = p.Code
		}
	}

	return vo
}

// PermissionVO 权限视图对象
type PermissionVO struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Module      string `json:"module"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// ToVO 转换为视图对象
func (p *Permission) ToVO() *PermissionVO {
	return &PermissionVO{
		ID:          p.ID,
		Code:        p.Code,
		Name:        p.Name,
		Module:      p.Module,
		Description: p.Description,
		SortOrder:   p.SortOrder,
	}
}

// PermissionGroup 权限分组
type PermissionGroup struct {
	Module      string          `json:"module"`
	ModuleName  string          `json:"module_name"`
	Permissions []*PermissionVO `json:"permissions"`
}
