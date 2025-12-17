package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"huaan-medical/internal/middleware"
	"huaan-medical/internal/model"
	"huaan-medical/internal/router"
	"huaan-medical/internal/scheduler"
	"huaan-medical/pkg/config"
	"huaan-medical/pkg/database"
	"huaan-medical/pkg/jwt"
	"huaan-medical/pkg/logger"
	"huaan-medical/pkg/redis"

	_ "huaan-medical/docs" // 导入Swagger文档
)

// @title 华安医疗预约系统 API
// @version 1.0
// @description 华安医疗门诊预约系统后端API接口文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@huaan-medical.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 请在Token前加上"Bearer "前缀

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.Init(&cfg.Log); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("配置加载成功",
		zap.String("mode", cfg.Server.Mode),
		zap.Int("port", cfg.Server.Port),
	)

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		logger.Fatal("初始化数据库失败", zap.Error(err))
	}
	defer database.Close()
	logger.Info("数据库连接成功")

	// 自动迁移数据库表
	if err := model.AutoMigrate(database.GetDB()); err != nil {
		logger.Fatal("数据库迁移失败", zap.Error(err))
	}
	logger.Info("数据库迁移成功")

	// 初始化Redis（可选）
	redis.TryInit(&cfg.Redis)
	if redis.IsEnabled() {
		defer redis.Close()
		logger.Info("Redis连接成功")
	} else {
		logger.Warn("Redis未启用或连接失败，部分缓存功能将不可用")
	}

	// 初始化JWT
	jwt.Init(&cfg.JWT)
	logger.Info("JWT初始化成功")

	// 初始化限流器
	middleware.InitRateLimiter(&cfg.RateLimit)
	logger.Info("限流器初始化成功")

	// 初始化定时任务
	scheduler.Init()
	defer scheduler.Stop()
	logger.Info("定时任务初始化成功")

	// 设置路由
	r := router.Setup(cfg.Server.Mode)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动服务器（非阻塞）
	go func() {
		logger.Info("服务器启动中...",
			zap.Int("port", cfg.Server.Port),
			zap.String("mode", cfg.Server.Mode),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	// 优雅关闭，等待10秒处理完现有请求
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务器强制关闭", zap.Error(err))
	}

	logger.Info("服务器已关闭")
}
