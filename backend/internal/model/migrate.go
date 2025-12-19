package model

import (
	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// 用户相关
		&User{},
		&Patient{},

		// 医院相关
		&Department{},
		&Doctor{},
		&Schedule{},

		// 预约相关
		&Appointment{},
		&MedicalRecord{},

		// 管理员相关
		&Admin{},
		&Role{},
		&Permission{},
		&AdminRole{},
		&RolePermission{},

		// 日志相关
		&OperationLog{},
		&LoginLog{},
	)
}

// GetAllModels 获取所有模型（用于文档生成等）
func GetAllModels() []interface{} {
	return []interface{}{
		&User{},
		&Patient{},
		&Department{},
		&Doctor{},
		&Schedule{},
		&Appointment{},
		&MedicalRecord{},
		&Admin{},
		&Role{},
		&Permission{},
		&AdminRole{},
		&RolePermission{},
		&OperationLog{},
		&LoginLog{},
	}
}
