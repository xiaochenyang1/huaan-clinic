package model

import (
	"time"
)

// User 用户模型
type User struct {
	BaseModel
	OpenID       string     `gorm:"type:varchar(64);uniqueIndex;not null;comment:微信OpenID" json:"open_id"`
	UnionID      string     `gorm:"type:varchar(64);index;comment:微信UnionID" json:"union_id,omitempty"`
	Nickname     string     `gorm:"type:varchar(64);comment:昵称" json:"nickname"`
	Avatar       string     `gorm:"type:varchar(512);comment:头像URL" json:"avatar"`
	Phone        string     `gorm:"type:varchar(20);index;comment:手机号" json:"phone"`
	Gender       int        `gorm:"type:tinyint;default:0;comment:性别 0未知 1男 2女" json:"gender"`
	Status       int        `gorm:"type:tinyint;default:1;comment:状态 0禁用 1启用" json:"status"`
	BlockedUntil *time.Time `gorm:"comment:封禁截止时间(爽约惩罚)" json:"blocked_until,omitempty"`
	MissedCount  int        `gorm:"type:int;default:0;comment:累计爽约次数" json:"missed_count"`
	LastLoginAt  *time.Time `gorm:"comment:最后登录时间" json:"last_login_at,omitempty"`
	LastLoginIP  string     `gorm:"type:varchar(64);comment:最后登录IP" json:"last_login_ip,omitempty"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// IsBlocked 判断用户是否被封禁
func (u *User) IsBlocked() bool {
	if u.BlockedUntil == nil {
		return false
	}
	return u.BlockedUntil.After(time.Now())
}

// UserVO 用户信息视图对象（用于响应）
type UserVO struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Phone     string `json:"phone"`
	Gender    int    `json:"gender"`
	HasPhone  bool   `json:"has_phone"`  // 是否绑定手机
	IsBlocked bool   `json:"is_blocked"` // 是否被封禁
}

// ToVO 转换为视图对象
func (u *User) ToVO() *UserVO {
	return &UserVO{
		ID:        u.ID,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Phone:     maskPhone(u.Phone),
		Gender:    u.Gender,
		HasPhone:  u.Phone != "",
		IsBlocked: u.IsBlocked(),
	}
}

// maskPhone 手机号脱敏
func maskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}
