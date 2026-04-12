package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	whisper "github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"whisper_with_go/config"
	"whisper_with_go/internal/model"
	"whisper_with_go/pkg/utils"
)

var ErrTaskNotFound = errors.New("task not found")

// WhisperService Whisper 服务接口
type WhisperService interface {
	CreateTask(audioPath, originalFilename string, req model.TranscribeRequest) (*model.TaskSubmissionResponse, error)
	GetTask(taskID string) (*model.TaskStatusResponse, error)
}

type whisperService struct {
	config *config.Config
	mu     sync.RWMutex
	tasks  map[string]*taskRecord
}

type taskRecord struct {
	ID               string
	AudioPath        string
	OriginalFilename string
	Request          model.TranscribeRequest
	Status           model.TaskStatus
	Result           *model.TranscribeResponse
	Error            string
	CreatedAt        time.Time
	StartedAt        *time.Time
	CompletedAt      *time.Time
}

// NewWhisperService 创建新的 Whisper 服务实例
func NewWhisperService(cfg *config.Config) WhisperService {
	return &whisperService{
		config: cfg,
		tasks:  make(map[string]*taskRecord),
	}
}

// CreateTask 创建异步转录任务
func (s *whisperService) CreateTask(audioPath, originalFilename string, req model.TranscribeRequest) (*model.TaskSubmissionResponse, error) {
	taskID := generateTaskID()
	task := &taskRecord{
		ID:               taskID,
		AudioPath:        audioPath,
		OriginalFilename: originalFilename,
		Request:          req,
		Status:           model.TaskStatusPending,
		CreatedAt:        time.Now(),
	}

	s.mu.Lock()
	s.tasks[taskID] = task
	s.mu.Unlock()

	go s.processTask(taskID)

	return &model.TaskSubmissionResponse{
		TaskID:    taskID,
		Status:    model.TaskStatusPending,
		StatusURL: fmt.Sprintf("/api/v1/tasks/%s", taskID),
	}, nil
}

// GetTask 查询任务状态和结果
func (s *whisperService) GetTask(taskID string) (*model.TaskStatusResponse, error) {
	s.mu.RLock()
	task, ok := s.tasks[taskID]
	if !ok {
		s.mu.RUnlock()
		return nil, ErrTaskNotFound
	}

	response := &model.TaskStatusResponse{
		TaskID:      task.ID,
		Status:      task.Status,
		Error:       task.Error,
		Result:      task.Result,
		CreatedAt:   task.CreatedAt,
		StartedAt:   task.StartedAt,
		CompletedAt: task.CompletedAt,
	}
	s.mu.RUnlock()

	return response, nil
}

func (s *whisperService) processTask(taskID string) {
	startedAt := time.Now()

	s.mu.Lock()
	task, ok := s.tasks[taskID]
	if !ok {
		s.mu.Unlock()
		return
	}
	task.Status = model.TaskStatusProcessing
	task.StartedAt = &startedAt
	audioPath := task.AudioPath
	originalFilename := task.OriginalFilename
	req := task.Request
	s.mu.Unlock()

	defer utils.RemoveFile(audioPath)

	result, err := s.transcribeAudio(taskID, audioPath, originalFilename, req.Language, req.OutputType, req.Translate)
	completedAt := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok = s.tasks[taskID]
	if !ok {
		return
	}
	task.CompletedAt = &completedAt
	if err != nil {
		task.Status = model.TaskStatusFailed
		task.Error = err.Error()
		return
	}

	task.Status = model.TaskStatusCompleted
	task.Result = result
}

// transcribeAudio 执行语音转录
func (s *whisperService) transcribeAudio(taskID, audioPath, originalFilename, language, outputType string, translate bool) (*model.TranscribeResponse, error) {
	startTime := time.Now()

	// 1. 转换音频格式
	wavPath, err := utils.ConvertToWAV(audioPath)
	if err != nil {
		return nil, fmt.Errorf("音频转换失败: %v", err)
	}
	defer os.Remove(wavPath) // 清理临时文件

	// 2. 加载 Whisper 模型
	whisperModel, err := whisper.New(s.config.Whisper.ModelPath)
	if err != nil {
		return nil, fmt.Errorf("加载模型失败: %v", err)
	}
	defer whisperModel.Close()

	// 3. 读取音频样本
	samples, err := utils.ReadWAVFile(wavPath)
	if err != nil {
		return nil, fmt.Errorf("读取音频失败: %v", err)
	}

	// 4. 创建处理上下文
	context, err := whisperModel.NewContext()
	if err != nil {
		return nil, fmt.Errorf("创建上下文失败: %v", err)
	}

	// 5. 配置参数
	if language != "" && language != "auto" {
		if err := context.SetLanguage(language); err != nil {
			return nil, fmt.Errorf("设置语言失败: %v", err)
		}
	}
	context.SetThreads(uint(s.config.Whisper.Threads))
	context.SetTranslate(translate)
	context.SetTokenTimestamps(true)
	context.SetSplitOnWord(true)

	// 6. 处理音频
	if err := context.Process(samples, nil, nil, nil); err != nil {
		return nil, fmt.Errorf("处理音频失败: %v", err)
	}

	// 7. 收集转录结果
	var segments []model.Segment
	var fullText strings.Builder

	for {
		segment, err := context.NextSegment()
		if err != nil {
			break
		}

		text := strings.TrimSpace(segment.Text)
		if text == "" {
			continue
		}

		start, end := refinedSegmentTiming(context, segment)

		segments = append(segments, model.Segment{
			Index: len(segments) + 1,
			Start: durationToSeconds(start),
			End:   durationToSeconds(end),
			Text:  text,
		})

		fullText.WriteString(text)
		fullText.WriteString(" ")
	}

	// 8. 生成输出文件（如果需要）
	var outputFile string
	if outputType == "srt" || outputType == "txt" {
		outputFile, err = s.saveOutput(audioPath, segments, outputType)
		if err != nil {
			return nil, fmt.Errorf("保存输出文件失败: %v", err)
		}
	}

	// 9. 构建响应
	response := &model.TranscribeResponse{
		TaskID:      taskID,
		Filename:    originalFilename,
		Language:    language,
		Text:        strings.TrimSpace(fullText.String()),
		ProcessTime: time.Since(startTime).Seconds(),
		OutputFile:  outputFile,
	}

	// 根据输出类型决定是否包含详细的片段信息
	if outputType == "json" {
		response.Segments = segments
	}

	// 计算音频时长
	if len(segments) > 0 {
		response.Duration = segments[len(segments)-1].End
	}

	return response, nil
}

// saveOutput 保存输出文件
func (s *whisperService) saveOutput(audioPath string, segments []model.Segment, outputType string) (string, error) {
	// 确保输出目录存在
	if err := utils.EnsureDirExists(s.config.Upload.OutputDir); err != nil {
		return "", err
	}

	// 生成输出文件路径 - 只提取文件名，不包含目录
	filename := filepath.Base(audioPath)
	// 去掉扩展名
	baseName := strings.TrimSuffix(filename, filepath.Ext(filename))
	// 去掉 _temp 后缀（如果有）
	baseName = strings.TrimSuffix(baseName, "_temp")
	outputPath := filepath.Join(s.config.Upload.OutputDir, fmt.Sprintf("%s.%s", baseName, outputType))

	// 创建输出文件
	file, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 根据格式写入内容
	switch outputType {
	case "srt":
		for _, segment := range segments {
			fmt.Fprintf(file, "%d\n", segment.Index)
			fmt.Fprintf(file, "%s --> %s\n",
				utils.FormatSRTTime(time.Duration(segment.Start*float64(time.Second))),
				utils.FormatSRTTime(time.Duration(segment.End*float64(time.Second))),
			)
			fmt.Fprintf(file, "%s\n\n", segment.Text)
		}
	case "txt":
		for _, segment := range segments {
			fmt.Fprintf(file, "[%s --> %s] %s\n",
				utils.FormatTimestamp(time.Duration(segment.Start*float64(time.Second))),
				utils.FormatTimestamp(time.Duration(segment.End*float64(time.Second))),
				segment.Text,
			)
		}
	}

	return filepath.Base(outputPath), nil
}

// generateTaskID 生成任务ID
func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}

func refinedSegmentTiming(context whisper.Context, segment whisper.Segment) (time.Duration, time.Duration) {
	start := segment.Start
	end := segment.End
	hasTokenTiming := false

	for _, token := range segment.Tokens {
		if !context.IsText(token) || strings.TrimSpace(token.Text) == "" {
			continue
		}
		if token.End <= token.Start {
			continue
		}

		if !hasTokenTiming {
			start = token.Start
			end = token.End
			hasTokenTiming = true
			continue
		}

		if token.Start < start {
			start = token.Start
		}
		if token.End > end {
			end = token.End
		}
	}

	if !hasTokenTiming || end <= start {
		return segment.Start, segment.End
	}

	return start, end
}

func durationToSeconds(value time.Duration) float64 {
	return float64(value.Milliseconds()) / 1000
}
