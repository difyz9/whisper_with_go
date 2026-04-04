package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"whisper_with_go/config"
	"whisper_with_go/internal/handler"
	"whisper_with_go/internal/middleware"
	"whisper_with_go/internal/service"
)

// Setup 设置路由
func Setup(cfg *config.Config) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建 Gin 引擎
	r := gin.New()

	// 使用中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 创建服务和处理器
	whisperService := service.NewWhisperService(cfg)
	whisperHandler := handler.NewWhisperHandler(whisperService, cfg)

	// 健康检查
	r.GET("/health", whisperHandler.HealthCheck)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Whisper API Server",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health":     "GET /health",
				"transcribe": "POST /api/v1/transcribe",
				"download":   "GET /api/v1/download/:filename",
				"swagger":    "GET /swagger/index.html",
			},
		})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// Whisper 转录相关
		v1.POST("/transcribe", whisperHandler.Transcribe)
		v1.GET("/download/:filename", whisperHandler.DownloadOutput)
	}

	return r
}
