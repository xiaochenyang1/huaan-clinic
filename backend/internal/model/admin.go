package model

import (
	"time"
)

// Admin 管理员模型
type Admin struct {
	BaseModel
	Username    string     `gorm:"type:varchar(64);uniqueIndex;not null;comment:用户名" json:"username"`
	Password    string     `gorm:"type:varchar(128);not null;comment:密码(加密)" json:"-"`
	Nickname    string     `gorm:"type:varchar(64);comment:昵称" json:"nickname"`
	Avatar      string     `gorm:"type:varchar(512);comment:头像" json:"avatar"`
	Phone       string     `gorm:"type:varchar(20);comment:手机号" json:"phone"`
	Email       string     `gorm:"type:varchar(128);comment:邮箱" json:"email"`
	Status      int        `gorm:"type:tinyint;default:1;comment:状态 0禁用 1启用" json:"status"`
	LastLoginAt *time.Time `gorm:"comment:最后登录时间" json:"last_login_at,omitempty"`
	LastLoginIP string     `gorm:"type:varchar(64);comment:最后登录IP" json:"last_login_ip,omitempty"`

	// 关联
	Roles []Role `gorm:"many2many:admin_roles;" json:"roles,omitempty"`
}

// TableName 表名
func (Admin) TableName() string {
	return "admins"
}

// AdminVO 管理员视图对象
type AdminVO struct {
	ID          int64    `json:"id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	Avatar      string   `json:"avatar"`
	Phone       string   `json:"phone"`
	Email       string   `json:"email"`
	Status      int      `json:"status"`
	StatusName  string   `json:"status_name"`
	Roles       []string `json:"roles,omitempty"`
	RoleNames   []string `json:"role_names,omitempty"`
	LastLoginAt string   `json:"last_login_at,omitempty"`
	CreatedAt   string   `json:"created_at"`
}

// ToVO 转换为视图对象
func (a *Admin) ToVO() *AdminVO {
	statusName := "启用"
	if a.Status == StatusDisabled {
		statusName = "禁用"
	}

	vo := &AdminVO{
		ID:         a.ID,
		Username:   a.Username,
		Nickname:   a.Nickname,
		Avatar:     a.Avatar,
		Phone:      a.Phone,
		Email:      a.Email,
		Status:     a.Status,
		StatusName: statusName,
		CreatedAt:  a.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if a.LastLoginAt != nil {
		vo.LastLoginAt = a.LastLoginAt.Format("2006-01-02 15:04:05")
	}

	// 角色信息
	if len(a.Roles) > 0 {
		vo.Roles = make([]string, len(a.Roles))
		vo.RoleNames = make([]string, len(a.Roles))
		for i, role := range a.Roles {
			vo.Roles[i] = role.Code
			vo.RoleNames[i] = role.Name
		}
	}

	return vo
}

// HasRole 判断是否拥有指定角色
func (a *Admin) HasRole(roleCode string) bool {
	for _, role := range a.Roles {
		if role.Code == roleCode {
			return true
		}
	}
	return false
}

// IsSuperAdmin 判断是否为超级管理员
func (a *Admin) IsSuperAdmin() bool {
	return a.HasRole("super_admin")
}
