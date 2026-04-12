package handler

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"whisper_with_go/config"
	"whisper_with_go/internal/model"
	"whisper_with_go/internal/service"
	"whisper_with_go/pkg/utils"
)

// WhisperHandler Whisper HTTP 处理器
type WhisperHandler struct {
	service service.WhisperService
	config  *config.Config
}

// NewWhisperHandler 创建新的 Whisper 处理器
func NewWhisperHandler(svc service.WhisperService, cfg *config.Config) *WhisperHandler {
	return &WhisperHandler{
		service: svc,
		config:  cfg,
	}
}

// HealthCheck 健康检查
// @Summary 健康检查
// @Description 检查服务是否正常运行
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} model.HealthCheckResponse
// @Router /health [get]
func (h *WhisperHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, model.HealthCheckResponse{
		Status:  "ok",
		Version: "1.0.0",
		Time:    time.Now().Format(time.RFC3339),
	})
}

// Transcribe 转录音频文件
// @Summary 转录音频
// @Description 上传音频文件并创建异步转录任务
// @Tags Whisper
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "音频文件"
// @Param language formData string false "语言代码 (auto, zh, en, etc.)" default(auto)
// @Param output_type formData string false "输出格式 (json, srt, txt)" default(json)
// @Param translate formData boolean false "是否翻译为英文" default(false)
// @Success 202 {object} model.Response
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/transcribe [post]
func (h *WhisperHandler) Transcribe(c *gin.Context) {
	// 1. 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("未找到上传文件", err))
		return
	}

	// 2. 验证文件大小
	if file.Size > h.config.Upload.MaxFileSize {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			fmt.Sprintf("文件过大，最大允许 %d MB", h.config.Upload.MaxFileSize/1024/1024),
			nil,
		))
		return
	}

	// 3. 验证文件类型
	if !utils.IsAudioFile(file.Filename) {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			"不支持的文件格式，仅支持: mp3, wav, m4a, aac, flac, ogg",
			nil,
		))
		return
	}

	// 4. 保存上传文件
	filename := utils.GenerateUniqueFilename(file.Filename)
	uploadPath := filepath.Join(h.config.Upload.UploadDir, filename)

	if err := utils.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("保存文件失败", err))
		return
	}

	// 5. 获取请求参数
	var req model.TranscribeRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.RemoveFile(uploadPath)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("参数错误", err))
		return
	}

	// 设置默认值
	if req.Language == "" {
		req.Language = "auto"
	}
	if req.OutputType == "" {
		req.OutputType = "json"
	}

	// 验证输出格式
	if !isValidOutputType(req.OutputType) {
		utils.RemoveFile(uploadPath)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			"不支持的输出格式，仅支持: json, srt, txt",
			nil,
		))
		return
	}

	// 6. 创建异步任务
	result, err := h.service.CreateTask(uploadPath, file.Filename, req)
	if err != nil {
		utils.RemoveFile(uploadPath)
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("创建转录任务失败", err))
		return
	}

	// 7. 返回结果
	c.JSON(http.StatusAccepted, model.NewSuccessResponse("任务已创建", result))
}

// GetTaskStatus 查询转录任务状态
// @Summary 查询转录任务
// @Description 通过任务 ID 查询转录状态和结果
// @Tags Whisper
// @Produce json
// @Param task_id path string true "任务 ID"
// @Success 200 {object} model.Response
// @Failure 404 {object} model.ErrorResponse
// @Router /api/v1/tasks/{task_id} [get]
func (h *WhisperHandler) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("任务 ID 不能为空", nil))
		return
	}

	result, err := h.service.GetTask(taskID)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, model.NewErrorResponse("任务不存在", nil))
			return
		}

		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("查询任务失败", err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse("查询成功", result))
}

// DownloadOutput 下载输出文件
// @Summary 下载输出文件
// @Description 下载转录生成的字幕或文本文件
// @Tags Whisper
// @Produce octet-stream
// @Param filename path string true "文件名"
// @Success 200 {file} file
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/v1/download/{filename} [get]
func (h *WhisperHandler) DownloadOutput(c *gin.Context) {
	filename := c.Param("filename")

	// 防止路径遍历攻击
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("非法文件名", nil))
		return
	}

	filepath := filepath.Join(h.config.Upload.OutputDir, filename)

	// 检查文件是否存在
	if !utils.FileExists(filepath) {
		c.JSON(http.StatusNotFound, model.NewErrorResponse("文件不存在", nil))
		return
	}

	// 设置下载响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")

	// 发送文件
	c.File(filepath)
}

// isValidOutputType 验证输出格式是否有效
func isValidOutputType(outputType string) bool {
	validTypes := []string{"json", "srt", "txt"}
	for _, t := range validTypes {
		if outputType == t {
			return true
		}
	}
	return false
}
