package model

import (
	"time"
)

// OperationLog 操作日志模型
type OperationLog struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminID     int64     `gorm:"index;comment:管理员ID" json:"admin_id"`
	AdminName   string    `gorm:"type:varchar(64);comment:管理员名称" json:"admin_name"`
	Module      string    `gorm:"type:varchar(32);index;comment:操作模块" json:"module"`
	Action      string    `gorm:"type:varchar(32);comment:操作动作" json:"action"`
	Method      string    `gorm:"type:varchar(10);comment:请求方法" json:"method"`
	Path        string    `gorm:"type:varchar(256);comment:请求路径" json:"path"`
	Query       string    `gorm:"type:text;comment:请求参数" json:"query,omitempty"`
	Body        string    `gorm:"type:text;comment:请求体" json:"body,omitempty"`
	Response    string    `gorm:"type:text;comment:响应内容" json:"response,omitempty"`
	IP          string    `gorm:"type:varchar(64);comment:IP地址" json:"ip"`
	UserAgent   string    `gorm:"type:varchar(512);comment:User-Agent" json:"user_agent,omitempty"`
	Status      int       `gorm:"type:int;comment:响应状态码" json:"status"`
	Latency     int64     `gorm:"type:bigint;comment:耗时(ms)" json:"latency"`
	ErrorMsg    string    `gorm:"type:text;comment:错误信息" json:"error_msg,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName 表名
func (OperationLog) TableName() string {
	return "operation_logs"
}

// LoginLog 登录日志模型
type LoginLog struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserType  string    `gorm:"type:varchar(20);index;comment:用户类型 user/admin" json:"user_type"`
	UserID    int64     `gorm:"index;comment:用户ID" json:"user_id"`
	Username  string    `gorm:"type:varchar(64);comment:用户名" json:"username"`
	LoginType string    `gorm:"type:varchar(20);comment:登录方式 wechat/password" json:"login_type"`
	IP        string    `gorm:"type:varchar(64);comment:IP地址" json:"ip"`
	Location  string    `gorm:"type:varchar(128);comment:登录地点" json:"location,omitempty"`
	Device    string    `gorm:"type:varchar(256);comment:设备信息" json:"device,omitempty"`
	OS        string    `gorm:"type:varchar(64);comment:操作系统" json:"os,omitempty"`
	Browser   string    `gorm:"type:varchar(64);comment:浏览器" json:"browser,omitempty"`
	Status    int       `gorm:"type:tinyint;comment:登录状态 0失败 1成功" json:"status"`
	Message   string    `gorm:"type:varchar(256);comment:登录消息" json:"message,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName 表名
func (LoginLog) TableName() string {
	return "login_logs"
}

// 用户类型常量
const (
	UserTypeUser  = "user"  // 普通用户
	UserTypeAdmin = "admin" // 管理员
)

// 登录类型常量
const (
	LoginTypeWeChat   = "wechat"   // 微信登录
	LoginTypePassword = "password" // 密码登录
)

// 登录状态常量
const (
	LoginStatusFailed  = 0 // 登录失败
	LoginStatusSuccess = 1 // 登录成功
)

// OperationLogVO 操作日志视图对象
type OperationLogVO struct {
	ID        int64  `json:"id"`
	AdminID   int64  `json:"admin_id"`
	AdminName string `json:"admin_name"`
	Module    string `json:"module"`
	Action    string `json:"action"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	IP        string `json:"ip"`
	Status    int    `json:"status"`
	Latency   int64  `json:"latency"`
	CreatedAt string `json:"created_at"`
}

// ToVO 转换为视图对象
func (l *OperationLog) ToVO() *OperationLogVO {
	return &OperationLogVO{
		ID:        l.ID,
		AdminID:   l.AdminID,
		AdminName: l.AdminName,
		Module:    l.Module,
		Action:    l.Action,
		Method:    l.Method,
		Path:      l.Path,
		IP:        l.IP,
		Status:    l.Status,
		Latency:   l.Latency,
		CreatedAt: l.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// LoginLogVO 登录日志视图对象
type LoginLogVO struct {
	ID        int64  `json:"id"`
	UserType  string `json:"user_type"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	LoginType string `json:"login_type"`
	IP        string `json:"ip"`
	Location  string `json:"location"`
	Device    string `json:"device"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

// ToVO 转换为视图对象
func (l *LoginLog) ToVO() *LoginLogVO {
	return &LoginLogVO{
		ID:        l.ID,
		UserType:  l.UserType,
		UserID:    l.UserID,
		Username:  l.Username,
		LoginType: l.LoginType,
		IP:        l.IP,
		Location:  l.Location,
		Device:    l.Device,
		Status:    l.Status,
		Message:   l.Message,
		CreatedAt: l.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
