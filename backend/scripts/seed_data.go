// 数据库初始化数据脚本
//go:build ignore
// +build ignore

package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"huaan-medical/internal/rbac"
)

// 定义临时的模型结构（避免依赖主项目）
type Admin struct {
	ID        int64     `gorm:"primarykey"`
	Username  string    `gorm:"size:50;not null;uniqueIndex"`
	Password  string    `gorm:"size:255;not null"`
	Nickname  string    `gorm:"size:50"`
	Avatar    string    `gorm:"size:255"`
	Phone     string    `gorm:"size:20"`
	Email     string    `gorm:"size:100"`
	Status    int       `gorm:"type:tinyint;not null;default:1"`
	LastLoginAt *time.Time
	LastLoginIP string `gorm:"size:64"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Admin) TableName() string {
	return "admins"
}

type Department struct {
	ID          int64  `gorm:"primarykey"`
	Name        string `gorm:"size:50;not null;uniqueIndex"`
	Description string `gorm:"type:text"`
	Icon        string `gorm:"size:255"`
	Status      int    `gorm:"type:tinyint;not null;default:1"`
	SortOrder   int    `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Department) TableName() string {
	return "departments"
}

type Role struct {
	ID          int64  `gorm:"primarykey"`
	Name        string `gorm:"size:50;not null"`
	Code        string `gorm:"size:50;not null;uniqueIndex"`
	Description string `gorm:"size:255"`
	SortOrder   int    `gorm:"default:0"`
	Status      int    `gorm:"type:tinyint;not null;default:1"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Role) TableName() string {
	return "roles"
}

type Permission struct {
	ID          int64  `gorm:"primarykey"`
	Name        string `gorm:"size:50;not null"`
	Code        string `gorm:"size:100;not null;uniqueIndex"`
	Module      string `gorm:"size:50;not null"`
	Description string `gorm:"size:255"`
	SortOrder   int    `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Permission) TableName() string {
	return "permissions"
}

type AdminRole struct {
	AdminID int64 `gorm:"primarykey"`
	RoleID  int64 `gorm:"primarykey"`
}

func (AdminRole) TableName() string {
	return "admin_roles"
}

type RolePermission struct {
	RoleID       int64 `gorm:"primarykey"`
	PermissionID int64 `gorm:"primarykey"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func main() {
	// 数据库连接配置（从 config.yaml 读取）
	dsn := "root:12345678@tcp(localhost:3306)/huaan_medical?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	fmt.Println("开始初始化数据...")

	// 1. 初始化管理员
	if err := seedAdmins(db); err != nil {
		log.Fatalf("初始化管理员失败: %v", err)
	}

	// 2. 初始化科室
	if err := seedDepartments(db); err != nil {
		log.Fatalf("初始化科室失败: %v", err)
	}

	// 3. 初始化角色
	if err := seedRoles(db); err != nil {
		log.Fatalf("初始化角色失败: %v", err)
	}

	// 4. 初始化权限
	if err := seedPermissions(db); err != nil {
		log.Fatalf("初始化权限失败: %v", err)
	}

	fmt.Println("✅ 数据初始化完成!")
	fmt.Println("\n默认管理员账号:")
	fmt.Println("  用户名: admin")
	fmt.Println("  密码: admin123")
}

// seedAdmins 初始化管理员
func seedAdmins(db *gorm.DB) error {
	fmt.Println("→ 初始化管理员...")

	// 检查默认管理员是否存在
	var existing Admin
	if err := db.Where("username = ?", "admin").First(&existing).Error; err == nil {
		fmt.Println("  默认管理员已存在，跳过创建")
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 生成密码哈希
	password := "admin123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admins := []Admin{
		{
			Username: "admin",
			Password: string(hash),
			Nickname: "系统管理员",
			Phone:    "13800138000",
			Email:    "admin@huaan-medical.com",
			Status:   1,
		},
	}

	result := db.Create(&admins)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("  ✅ 成功创建 %d 个管理员\n", result.RowsAffected)
	return nil
}

// seedDepartments 初始化科室
func seedDepartments(db *gorm.DB) error {
	fmt.Println("→ 初始化科室...")

	// 检查是否已存在科室
	var count int64
	db.Model(&Department{}).Count(&count)
	if count > 0 {
		fmt.Println("  科室已存在，跳过")
		return nil
	}

	departments := []Department{
		{
			Name:        "内科",
			Description: "内科门诊，诊治各种内科疾病",
			Status:      1,
			SortOrder:   1,
		},
		{
			Name:        "外科",
			Description: "外科门诊，诊治各种外科疾病",
			Status:      1,
			SortOrder:   2,
		},
		{
			Name:        "儿科",
			Description: "儿科门诊，专注儿童疾病诊疗",
			Status:      1,
			SortOrder:   3,
		},
		{
			Name:        "妇产科",
			Description: "妇产科门诊，提供妇科和产科医疗服务",
			Status:      1,
			SortOrder:   4,
		},
		{
			Name:        "骨科",
			Description: "骨科门诊，诊治骨骼肌肉系统疾病",
			Status:      1,
			SortOrder:   5,
		},
		{
			Name:        "皮肤科",
			Description: "皮肤科门诊，诊治各种皮肤疾病",
			Status:      1,
			SortOrder:   6,
		},
		{
			Name:        "眼科",
			Description: "眼科门诊，诊治眼部疾病",
			Status:      1,
			SortOrder:   7,
		},
		{
			Name:        "耳鼻喉科",
			Description: "耳鼻喉科门诊，诊治耳鼻喉疾病",
			Status:      1,
			SortOrder:   8,
		},
		{
			Name:        "口腔科",
			Description: "口腔科门诊，提供口腔诊疗服务",
			Status:      1,
			SortOrder:   9,
		},
		{
			Name:        "中医科",
			Description: "中医科门诊，提供中医诊疗服务",
			Status:      1,
			SortOrder:   10,
		},
	}

	result := db.Create(&departments)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("  ✅ 成功创建 %d 个科室\n", result.RowsAffected)
	return nil
}

// seedRoles 初始化角色
func seedRoles(db *gorm.DB) error {
	fmt.Println("→ 初始化角色...")

	// 检查是否已存在角色
	var count int64
	db.Model(&Role{}).Count(&count)
	if count > 0 {
		fmt.Println("  角色已存在，跳过")
		return nil
	}

	roles := []Role{
		{
			Name:        "超级管理员",
			Code:        "super_admin",
			Description: "拥有所有权限",
			Status:      "active",
		},
		{
			Name:        "医院管理员",
			Code:        "hospital_admin",
			Description: "管理医院业务",
			Status:      "active",
		},
		{
			Name:        "科室管理员",
			Code:        "department_admin",
			Description: "管理科室事务",
			Status:      "active",
		},
	}

	result := db.Create(&roles)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("  ✅ 成功创建 %d 个角色\n", result.RowsAffected)
	return nil
}

// seedPermissions 初始化权限
func seedPermissions(db *gorm.DB) error {
	fmt.Println("→ 初始化权限...")

	// 检查是否已存在权限
	var count int64
	db.Model(&Permission{}).Count(&count)
	if count > 0 {
		fmt.Println("  权限已存在，跳过")
		return nil
	}

	permissions := []Permission{
		// 科室管理
		{Name: "查看科室", Code: "department:view", Category: "科室管理", Description: "查看科室信息"},
		{Name: "创建科室", Code: "department:create", Category: "科室管理", Description: "创建新科室"},
		{Name: "编辑科室", Code: "department:edit", Category: "科室管理", Description: "编辑科室信息"},
		{Name: "删除科室", Code: "department:delete", Category: "科室管理", Description: "删除科室"},

		// 医生管理
		{Name: "查看医生", Code: "doctor:view", Category: "医生管理", Description: "查看医生信息"},
		{Name: "创建医生", Code: "doctor:create", Category: "医生管理", Description: "添加新医生"},
		{Name: "编辑医生", Code: "doctor:edit", Category: "医生管理", Description: "编辑医生信息"},
		{Name: "删除医生", Code: "doctor:delete", Category: "医生管理", Description: "删除医生"},

		// 排班管理
		{Name: "查看排班", Code: "schedule:view", Category: "排班管理", Description: "查看排班信息"},
		{Name: "创建排班", Code: "schedule:create", Category: "排班管理", Description: "创建排班"},
		{Name: "编辑排班", Code: "schedule:edit", Category: "排班管理", Description: "编辑排班"},
		{Name: "删除排班", Code: "schedule:delete", Category: "排班管理", Description: "删除排班"},

		// 预约管理
		{Name: "查看预约", Code: "appointment:view", Category: "预约管理", Description: "查看预约信息"},
		{Name: "处理预约", Code: "appointment:handle", Category: "预约管理", Description: "处理预约状态"},
		{Name: "取消预约", Code: "appointment:cancel", Category: "预约管理", Description: "取消预约"},
		{Name: "导出预约", Code: "appointment:export", Category: "预约管理", Description: "导出预约数据"},

		// 患者管理
		{Name: "查看患者", Code: "patient:view", Category: "患者管理", Description: "查看患者信息"},
		{Name: "编辑患者", Code: "patient:edit", Category: "患者管理", Description: "编辑患者信息"},

		// 系统管理
		{Name: "查看管理员", Code: "admin:view", Category: "系统管理", Description: "查看管理员"},
		{Name: "创建管理员", Code: "admin:create", Category: "系统管理", Description: "创建管理员"},
		{Name: "编辑管理员", Code: "admin:edit", Category: "系统管理", Description: "编辑管理员"},
		{Name: "删除管理员", Code: "admin:delete", Category: "系统管理", Description: "删除管理员"},

		// 数据统计
		{Name: "查看统计", Code: "statistics:view", Category: "数据统计", Description: "查看数据统计"},
		{Name: "导出数据", Code: "data:export", Category: "数据统计", Description: "导出数据"},
	}

	result := db.Create(&permissions)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("  ✅ 成功创建 %d 个权限\n", result.RowsAffected)
	return nil
}
