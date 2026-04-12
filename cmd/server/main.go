package main

import (
	"fmt"
	"log"
	"os"

	"whisper_with_go/config"
	docs "whisper_with_go/docs"
	"whisper_with_go/internal/router"
	"whisper_with_go/pkg/utils"
)

// @title Whisper API Server
// @version 1.0.0
// @description 基于 Gin 和 whisper.cpp 的语音转文字 API 服务。
// @BasePath /
// @schemes http https
// @contact.name API Support
// @contact.url https://github.com/difyz9/whisper_with_go
// @license.name MIT

func main() {
	// 打印启动横幅
	printBanner()

	// 加载配置
	cfg := config.Load()

	// 确保必要的目录存在
	if err := utils.EnsureDirExists(cfg.Upload.UploadDir); err != nil {
		log.Fatalf("创建上传目录失败: %v", err)
	}
	if err := utils.EnsureDirExists(cfg.Upload.OutputDir); err != nil {
		log.Fatalf("创建输出目录失败: %v", err)
	}

	// 检查模型文件是否存在
	if _, err := os.Stat(cfg.Whisper.ModelPath); os.IsNotExist(err) {
		log.Printf("警告: 模型文件不存在: %s", cfg.Whisper.ModelPath)
		log.Println("请先下载模型文件，例如:")
		log.Println("  bash download_model.sh")
	}

	addr := fmt.Sprintf(":%s", cfg.Server.Port)

	// 设置路由
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost%s", addr)
	docs.SwaggerInfo.BasePath = "/"
	r := router.Setup(cfg)

	// 启动服务器
	log.Printf("服务器启动在 http://localhost%s", addr)
	log.Printf("模式: %s", cfg.Server.Mode)
	log.Printf("模型: %s", cfg.Whisper.ModelPath)
	log.Printf("线程数: %d", cfg.Whisper.Threads)
	log.Println("--------------------------------------")
	log.Println("API 端点:")
	log.Printf("  健康检查: GET  http://localhost%s/health", addr)
	log.Printf("  转录音频: POST http://localhost%s/api/v1/transcribe", addr)
	log.Printf("  任务查询: GET  http://localhost%s/api/v1/tasks/:task_id", addr)
	log.Printf("  下载文件: GET  http://localhost%s/api/v1/download/:filename", addr)
	log.Printf("  Swagger : GET  http://localhost%s/swagger/index.html", addr)
	log.Println("--------------------------------------")

	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

func printBanner() {
	banner := `
╦ ╦┬ ┬┬┌─┐┌─┐┌─┐┬─┐  ╔═╗╔═╗╦  
║║║├─┤│└─┐├─┘├┤ ├┬┘  ╠═╣╠═╝║  
╚╩╝┴ ┴┴└─┘┴  └─┘┴└─  ╩ ╩╩  ╩  
语音转文字 API 服务
Version: 1.0.0
	`
	fmt.Println(banner)
}
