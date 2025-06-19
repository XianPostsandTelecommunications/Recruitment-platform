package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"lab-recruitment-platform/internal/config"
	"lab-recruitment-platform/internal/handlers"
	"lab-recruitment-platform/internal/middleware"
	"lab-recruitment-platform/internal/models"
	"lab-recruitment-platform/pkg/logger"
	"lab-recruitment-platform/pkg/validator"
)

// @title 实验室招新平台 API
// @version 1.0
// @description 实验室招新平台后端API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 请输入JWT令牌，格式：Bearer <token>

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.InitLogger(&logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		Output:     cfg.Log.Output,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
		Compress:   cfg.Log.Compress,
	}); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	logger.Info("开始启动实验室招新平台...")

	// 初始化验证器
	validator.InitValidator()

	// 初始化数据库
	if err := config.InitDatabase(&cfg.Database); err != nil {
		logger.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化Redis
	if err := config.InitRedis(&cfg.Redis); err != nil {
		logger.Fatalf("初始化Redis失败: %v", err)
	}

	// 自动迁移数据库表
	if err := config.AutoMigrate(
		&models.User{},
		&models.Lab{},
		&models.Application{},
		&models.Notification{},
	); err != nil {
		logger.Fatalf("数据库迁移失败: %v", err)
	}

	// 设置Gin模式
	if cfg.Server.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	r := gin.New()

	// 添加中间件
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestLoggerMiddleware())
	r.Use(gin.Recovery())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "实验室招新平台运行正常",
			"time":    time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证相关路由
		authHandler := handlers.NewAuthHandler()
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
			auth.POST("/refresh", middleware.AuthMiddleware(), authHandler.RefreshToken)
			auth.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
			auth.PUT("/profile", middleware.AuthMiddleware(), authHandler.UpdateProfile)
			auth.POST("/change-password", middleware.AuthMiddleware(), authHandler.ChangePassword)
		}

		// 申请相关路由
		applicationHandler := handlers.NewApplicationHandler()
		api.POST("/send-code", applicationHandler.SendCode)
		api.POST("/apply", applicationHandler.Apply)

		// 用户管理路由（需要管理员权限）
		// userHandler := handlers.NewUserHandler()
		// users := api.Group("/users")
		// users.Use(middleware.AdminMiddleware())
		// {
		// 	users.GET("", userHandler.ListUsers)
		// 	users.GET("/:id", userHandler.GetUser)
		// 	users.PUT("/:id", userHandler.UpdateUser)
		// 	users.DELETE("/:id", userHandler.DeleteUser)
		// 	users.POST("/:id/reset-password", userHandler.ResetPassword)
		// 	users.PUT("/:id/status", userHandler.UpdateUserStatus)
		// 	users.GET("/stats", userHandler.GetUserStats)
		// }

		// 实验室管理路由
		// labHandler := handlers.NewLabHandler()
		// labs := api.Group("/labs")
		// {
		// 	labs.GET("", labHandler.ListLabs)
		// 	labs.GET("/:id", labHandler.GetLab)
		// 	labs.POST("", middleware.AuthMiddleware(), labHandler.CreateLab)
		// 	labs.PUT("/:id", middleware.AuthMiddleware(), labHandler.UpdateLab)
		// 	labs.DELETE("/:id", middleware.AdminMiddleware(), labHandler.DeleteLab)
		// 	labs.GET("/:id/applications", middleware.AuthMiddleware(), labHandler.GetLabApplications)
		// }

		// 申请管理路由
		// applicationHandler := handlers.NewApplicationHandler()
		// applications := api.Group("/applications")
		// applications.Use(middleware.AuthMiddleware())
		// {
		// 	applications.GET("", applicationHandler.ListApplications)
		// 	applications.GET("/:id", applicationHandler.GetApplication)
		// 	applications.POST("", applicationHandler.CreateApplication)
		// 	applications.PUT("/:id", applicationHandler.UpdateApplication)
		// 	applications.DELETE("/:id", applicationHandler.DeleteApplication)
		// 	applications.POST("/:id/review", middleware.AdminMiddleware(), applicationHandler.ReviewApplication)
		// 	applications.GET("/stats", middleware.AdminMiddleware(), applicationHandler.GetApplicationStats)
		// }

		// 通知管理路由
		// notificationHandler := handlers.NewNotificationHandler()
		// notifications := api.Group("/notifications")
		// notifications.Use(middleware.AuthMiddleware())
		// {
		// 	notifications.GET("", notificationHandler.ListNotifications)
		// 	notifications.GET("/:id", notificationHandler.GetNotification)
		// 	notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
		// 	notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
		// 	notifications.GET("/stats", notificationHandler.GetNotificationStats)
		// }

		// 文件上传路由
		// uploadHandler := handlers.NewUploadHandler()
		// upload := api.Group("/upload")
		// upload.Use(middleware.AuthMiddleware())
		// {
		// 	upload.POST("/image", uploadHandler.UploadImage)
		// 	upload.POST("/file", uploadHandler.UploadFile)
		// }
	}

	// Swagger文档
	if cfg.Server.IsDevelopment() {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		logger.Info("Swagger文档地址: http://localhost:" + cfg.Server.Port + "/swagger/index.html")
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 启动服务器
	go func() {
		logger.Infof("服务器启动在端口: %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("服务器关闭失败: %v", err)
	}

	// 关闭数据库连接
	if err := config.CloseDatabase(); err != nil {
		logger.Errorf("关闭数据库连接失败: %v", err)
	}

	// 关闭Redis连接
	if err := config.CloseRedis(); err != nil {
		logger.Errorf("关闭Redis连接失败: %v", err)
	}

	logger.Info("服务器已关闭")
} 