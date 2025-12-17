package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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

	// API路由组
	api := r.Group("/api")
	{
		// 公开接口（无需认证）
		setupPublicRoutes(api)

		// 用户接口（需要用户认证）
		setupUserRoutes(api)

		// 管理后台接口（需要管理员认证）
		setupAdminRoutes(api)
	}

	return r
}

// setupPublicRoutes 设置公开路由（无需认证）
func setupPublicRoutes(rg *gin.RouterGroup) {
	// 用户登录
	rg.POST("/user/login", placeholder("微信登录"))

	// Token刷新
	rg.POST("/auth/refresh", placeholder("刷新Token"))

	// 科室列表（公开）
	rg.GET("/departments", placeholder("获取科室列表"))

	// 医生列表（公开）
	rg.GET("/doctors", placeholder("获取医生列表"))
	rg.GET("/doctors/:id", placeholder("获取医生详情"))

	// 排班查询（公开）
	rg.GET("/schedule", placeholder("查询医生排班"))
	rg.GET("/schedule/available", placeholder("获取可预约时段"))
}

// setupUserRoutes 设置用户路由（需要用户认证）
func setupUserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("")
	user.Use(middleware.JWTAuth())
	{
		// 用户信息
		user.GET("/user/info", placeholder("获取用户信息"))
		user.PUT("/user/info", placeholder("更新用户信息"))

		// 就诊人管理
		user.GET("/user/patients", placeholder("获取就诊人列表"))
		user.POST("/user/patients", placeholder("添加就诊人"))
		user.PUT("/user/patients/:id", placeholder("编辑就诊人"))
		user.DELETE("/user/patients/:id", placeholder("删除就诊人"))

		// 幂等Token
		user.GET("/token/idempotent", placeholder("获取幂等Token"))

		// 预约管理
		user.POST("/appointments", placeholder("创建预约"))
		user.GET("/appointments", placeholder("获取预约列表"))
		user.GET("/appointments/:id", placeholder("获取预约详情"))
		user.PUT("/appointments/:id/cancel", placeholder("取消预约"))
		user.POST("/appointments/:id/checkin", placeholder("预约签到"))

		// 就诊记录
		user.GET("/records", placeholder("获取就诊记录列表"))
		user.GET("/records/:id", placeholder("获取就诊记录详情"))
	}
}

// setupAdminRoutes 设置管理后台路由（需要管理员认证）
func setupAdminRoutes(rg *gin.RouterGroup) {
	// 管理员登录（公开）
	rg.POST("/admin/login", placeholder("管理员登录"))

	admin := rg.Group("/admin")
	admin.Use(middleware.JWTAdminAuth())
	{
		// 管理员信息
		admin.GET("/info", placeholder("获取管理员信息"))
		admin.PUT("/password", placeholder("修改密码"))

		// 仪表盘
		admin.GET("/dashboard", placeholder("仪表盘数据"))

		// 预约管理
		admin.GET("/appointments", placeholder("预约列表"))
		admin.PUT("/appointments/:id", placeholder("更新预约状态"))
		admin.GET("/appointments/export", placeholder("导出预约数据"))

		// 患者管理
		admin.GET("/patients", placeholder("患者列表"))
		admin.GET("/patients/:id", placeholder("患者详情"))

		// 科室管理
		admin.GET("/departments", placeholder("科室列表"))
		admin.POST("/departments", placeholder("添加科室"))
		admin.PUT("/departments/:id", placeholder("编辑科室"))
		admin.DELETE("/departments/:id", placeholder("删除科室"))

		// 医生管理
		admin.GET("/doctors", placeholder("医生列表"))
		admin.POST("/doctors", placeholder("添加医生"))
		admin.PUT("/doctors/:id", placeholder("编辑医生"))
		admin.DELETE("/doctors/:id", placeholder("删除医生"))
		admin.POST("/upload/avatar", placeholder("上传医生头像"))

		// 排班管理
		admin.GET("/schedules", placeholder("排班列表"))
		admin.POST("/schedules", placeholder("创建排班"))
		admin.POST("/schedules/batch", placeholder("批量创建排班"))
		admin.PUT("/schedules/:id", placeholder("编辑排班"))
		admin.DELETE("/schedules/:id", placeholder("删除排班"))

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
