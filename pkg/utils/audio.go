package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ConvertToWAV 将音频文件转换为 WAV 格式 (16kHz, 单声道, 16位PCM)
// Whisper 需要这种特定格式的音频输入
func ConvertToWAV(inputPath string) (string, error) {
	// 检查 ffmpeg 是否安装
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return "", fmt.Errorf("ffmpeg 未安装，请先安装 ffmpeg")
	}

	// 生成临时 WAV 文件路径
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	outputPath := filepath.Join(filepath.Dir(inputPath), baseName+"_temp.wav")

	// 使用 ffmpeg 转换
	// -i: 输入文件
	// -ar 16000: 采样率 16kHz
	// -ac 1: 单声道
	// -c:a pcm_s16le: 16位 PCM 编码
	// -y: 覆盖输出文件
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-ar", "16000",
		"-ac", "1",
		"-c:a", "pcm_s16le",
		"-y",
		outputPath,
	)

	// 执行转换
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg 转换失败: %v", err)
	}

	return outputPath, nil
}

// ReadWAVFile 读取 WAV 文件并返回音频样本
func ReadWAVFile(wavPath string) ([]float32, error) {
	// 读取文件内容
	data, err := os.ReadFile(wavPath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 跳过 WAV 文件头（通常 44 字节）
	const wavHeaderSize = 44
	if len(data) < wavHeaderSize {
		return nil, fmt.Errorf("WAV 文件格式无效")
	}

	// 将 PCM 数据转换为 float32 样本
	pcmData := data[wavHeaderSize:]
	samples := make([]float32, len(pcmData)/2)

	for i := 0; i < len(samples); i++ {
		// 读取 16 位 PCM 样本（小端序）
		sample := int16(pcmData[i*2]) | int16(pcmData[i*2+1])<<8
		// 转换为 -1.0 到 1.0 的浮点数
		samples[i] = float32(sample) / 32768.0
	}

	return samples, nil
}

// IsAudioFile 检查文件是否为支持的音频格式
func IsAudioFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	supportedExts := []string{".mp3", ".wav", ".m4a", ".aac", ".flac", ".ogg"}
	
	for _, supportedExt := range supportedExts {
		if ext == supportedExt {
			return true
		}
	}
	return false
}
