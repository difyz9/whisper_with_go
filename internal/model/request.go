package model

// TranscribeRequest 转录请求
type TranscribeRequest struct {
	Language   string `json:"language" form:"language"`     // 语言代码 (auto, zh, en, etc.)
	OutputType string `json:"output_type" form:"output_type"` // 输出格式: srt, txt, json
	Translate  bool   `json:"translate" form:"translate"`   // 是否翻译为英文
}

// HealthCheckResponse 健康检查响应
type HealthCheckResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Time    string `json:"time"`
}
