package config

import (
	"os"
	"strconv"
)

// Config 应用配置
type Config struct {
	Server  ServerConfig
	Whisper WhisperConfig
	Upload  UploadConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string // debug, release, test
}

// WhisperConfig Whisper 模型配置
type WhisperConfig struct {
	ModelPath string
	Threads   int
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxFileSize   int64  // 最大文件大小（字节）
	UploadDir     string // 上传目录
	OutputDir     string // 输出目录
	AllowedExts   []string
	AllowedMimes  []string
}

// Load 加载配置
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8084"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Whisper: WhisperConfig{
			ModelPath: getEnv("WHISPER_MODEL_PATH", "models/ggml-base.bin"),
			Threads:   getEnvAsInt("WHISPER_THREADS", 4),
		},
		Upload: UploadConfig{
			MaxFileSize: getEnvAsInt64("MAX_FILE_SIZE", 100*1024*1024), // 默认100MB
			UploadDir:   getEnv("UPLOAD_DIR", "uploads"),
			OutputDir:   getEnv("OUTPUT_DIR", "outputs"),
			AllowedExts: []string{".mp3", ".wav", ".m4a", ".aac", ".flac", ".ogg"},
			AllowedMimes: []string{
				"audio/mpeg",
				"audio/wav",
				"audio/x-wav",
				"audio/mp4",
				"audio/aac",
				"audio/flac",
				"audio/ogg",
			},
		},
	}
}

// getEnv 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsInt64 获取环境变量并转换为 int64
func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}
