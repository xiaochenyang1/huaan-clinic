package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"huaan-medical/internal/handler"
	"huaan-medical/internal/middleware"
)

// Setup 初始化路由
func Setup(mode string) *gin.Engine {
	// 设置运行模式
	gin.SetMode(mode)

	r := gin.New()

	// 全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "华安医疗预约系统 API 服务运行中",
		})
	})

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 静态文件服务（用于访问上传的文件）
	r.Static("/uploads", "./uploads")

	// 初始化Handler
	adminHandler := handler.NewAdminHandler()
	deptHandler := handler.NewDepartmentHandler()
	doctorHandler := handler.NewDoctorHandler()
	scheduleHandler := handler.NewScheduleHandler()
	uploadHandler := handler.NewUploadHandler()
	userHandler := handler.NewUserHandler()
	patientHandler := handler.NewPatientHandler()
	tokenHandler := handler.NewTokenHandler()
	appointmentHandler := handler.NewAppointmentHandler()

	// API路由组
	api := r.Group("/api")
	{
		// 公开接口（无需认证）
		setupPublicRoutes(api, deptHandler, doctorHandler, scheduleHandler, userHandler)

		// 用户接口（需要用户认证）
		setupUserRoutes(api, userHandler, patientHandler, tokenHandler, appointmentHandler)

		// 管理后台接口（需要管理员认证）
		setupAdminRoutes(api, adminHandler, deptHandler, doctorHandler, scheduleHandler, uploadHandler, appointmentHandler)
	}

	return r
}

// setupPublicRoutes 设置公开路由（无需认证）
func setupPublicRoutes(rg *gin.RouterGroup, deptHandler *handler.DepartmentHandler, doctorHandler *handler.DoctorHandler, scheduleHandler *handler.ScheduleHandler, userHandler *handler.UserHandler) {
	// 用户登录
	rg.POST("/user/login", userHandler.WeChatLogin)

	// Token刷新
	rg.POST("/auth/refresh", placeholder("刷新Token"))

	// 科室列表（公开）
	rg.GET("/departments", deptHandler.ListAll)

	// 医生列表（公开）
	rg.GET("/doctors", doctorHandler.ListPublic)
	rg.GET("/doctors/:id", doctorHandler.GetByIDPublic)

	// 排班查询（公开）
	rg.GET("/schedule", scheduleHandler.ListByDoctor)
	rg.GET("/schedule/available", scheduleHandler.ListAvailable)
}

// setupUserRoutes 设置用户路由（需要用户认证）
func setupUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler, patientHandler *handler.PatientHandler, tokenHandler *handler.TokenHandler, appointmentHandler *handler.AppointmentHandler) {
	user := rg.Group("")
	user.Use(middleware.JWTAuth())
	{
		// 用户信息
		user.GET("/user/info", userHandler.GetInfo)
		user.PUT("/user/info", userHandler.UpdateInfo)

		// 就诊人管理
		user.GET("/user/patients", patientHandler.List)
		user.GET("/user/patients/:id", patientHandler.GetByID)
		user.POST("/user/patients", patientHandler.Create)
		user.PUT("/user/patients/:id", patientHandler.Update)
		user.DELETE("/user/patients/:id", patientHandler.Delete)

		// 幂等Token
		user.GET("/token/idempotent", tokenHandler.GetIdempotentToken)

		// 预约管理
		user.POST("/appointments", appointmentHandler.Create)
		user.GET("/appointments", appointmentHandler.List)
		user.GET("/appointments/:id", appointmentHandler.GetByID)
		user.PUT("/appointments/:id/cancel", appointmentHandler.Cancel)
		user.POST("/appointments/:id/checkin", appointmentHandler.Checkin)

		// 就诊记录
		user.GET("/records", placeholder("获取就诊记录列表"))
		user.GET("/records/:id", placeholder("获取就诊记录详情"))
	}
}

// setupAdminRoutes 设置管理后台路由（需要管理员认证）
func setupAdminRoutes(rg *gin.RouterGroup, adminHandler *handler.AdminHandler, deptHandler *handler.DepartmentHandler, doctorHandler *handler.DoctorHandler, scheduleHandler *handler.ScheduleHandler, uploadHandler *handler.UploadHandler, appointmentHandler *handler.AppointmentHandler) {
	// 管理员登录（公开）
	rg.POST("/admin/login", adminHandler.Login)

	admin := rg.Group("/admin")
	admin.Use(middleware.JWTAdminAuth())
	{
		// 管理员信息
		admin.GET("/info", adminHandler.GetInfo)
		admin.PUT("/password", adminHandler.ChangePassword)

		// 仪表盘
		admin.GET("/dashboard", placeholder("仪表盘数据"))

		// 预约管理
		admin.GET("/appointments", appointmentHandler.ListAdmin)
		admin.PUT("/appointments/:id", placeholder("更新预约状态"))
		admin.GET("/appointments/export", placeholder("导出预约数据"))

		// 患者管理
		admin.GET("/patients", placeholder("患者列表"))
		admin.GET("/patients/:id", placeholder("患者详情"))

		// 科室管理
		admin.GET("/departments", deptHandler.List)
		admin.GET("/departments/:id", deptHandler.GetByID)
		admin.POST("/departments", deptHandler.Create)
		admin.PUT("/departments/:id", deptHandler.Update)
		admin.DELETE("/departments/:id", deptHandler.Delete)

		// 医生管理
		admin.GET("/doctors", doctorHandler.List)
		admin.GET("/doctors/:id", doctorHandler.GetByID)
		admin.POST("/doctors", doctorHandler.Create)
		admin.PUT("/doctors/:id", doctorHandler.Update)
		admin.DELETE("/doctors/:id", doctorHandler.Delete)

		// 文件上传
		admin.POST("/upload/avatar", uploadHandler.UploadAvatar)
		admin.POST("/upload/image", uploadHandler.UploadImage)

		// 排班管理
		admin.GET("/schedules", scheduleHandler.List)
		admin.GET("/schedules/:id", scheduleHandler.GetByID)
		admin.POST("/schedules", scheduleHandler.Create)
		admin.POST("/schedules/batch", scheduleHandler.BatchCreate)
		admin.PUT("/schedules/:id", scheduleHandler.Update)
		admin.DELETE("/schedules/:id", scheduleHandler.Delete)

		// 数据统计
		admin.GET("/statistics", placeholder("数据统计"))
	}
}

// placeholder 占位处理函数（后续会被实际Handler替换）
func placeholder(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    200000,
			"message": "接口开发中",
			"data": gin.H{
				"api": name,
			},
		})
	}
}
