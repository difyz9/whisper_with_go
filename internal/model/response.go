package model

import "time"

// Response 通用响应结构
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"参数错误"`
	Error   string `json:"error,omitempty" example:"invalid request"`
}

// TranscribeSuccessResponse 转录成功响应
type TranscribeSuccessResponse struct {
	Success bool               `json:"success" example:"true"`
	Message string             `json:"message" example:"转录成功"`
	Data    TranscribeResponse `json:"data"`
}

// TranscribeResponse 转录响应
type TranscribeResponse struct {
	TaskID      string    `json:"task_id"`
	Filename    string    `json:"filename"`
	Language    string    `json:"language"`
	Duration    float64   `json:"duration_seconds"`
	Segments    []Segment `json:"segments,omitempty"`
	Text        string    `json:"text,omitempty"`
	OutputFile  string    `json:"output_file,omitempty"`
	ProcessTime float64   `json:"process_time_seconds"`
}

// Segment 转录片段
type Segment struct {
	Index int     `json:"index"`
	Start float64 `json:"start"` // 开始时间（秒）
	End   float64 `json:"end"`   // 结束时间（秒）
	Text  string  `json:"text"`  // 文本内容
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(message string, data interface{}) *Response {
	return &Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(message string, err error) *Response {
	resp := &Response{
		Success: false,
		Message: message,
	}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp
}

// DurationToSeconds 将 time.Duration 转换为浮点数秒
func DurationToSeconds(d time.Duration) float64 {
	return d.Seconds()
}
