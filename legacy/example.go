package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	whisper "github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

// 这个示例展示了如何以编程方式使用Whisper Go绑定

func main() {
	// 简单示例: 转录一个音频文件
	err := simpleTranscription()
	if err != nil {
		log.Fatal(err)
	}
}

// simpleTranscription 演示基本的语音转录
func simpleTranscription() error {
	modelPath := "models/ggml-base.bin"
	audioPath := "test.mp3"
	
	fmt.Println("=== 简单转录示例 ===")
	fmt.Printf("模型: %s\n", modelPath)
	fmt.Printf("音频: %s\n", audioPath)
	fmt.Println()
	
	// 步骤1: 转换音频格式
	wavPath, err := convertToWAV(audioPath)
	if err != nil {
		return fmt.Errorf("音频转换失败: %v", err)
	}
	defer cleanupTempFile(wavPath)
	
	// 步骤2: 加载Whisper模型
	model, err := whisper.New(modelPath)
	if err != nil {
		return fmt.Errorf("加载模型失败: %v", err)
	}
	defer model.Close()
	
	fmt.Println("✓ 模型加载成功")
	
	// 步骤3: 读取音频数据
	samples, err := readAudioSamples(wavPath)
	if err != nil {
		return fmt.Errorf("读取音频失败: %v", err)
	}
	
	fmt.Printf("✓ 音频加载成功 (样本数: %d)\n", len(samples))
	
	// 步骤4: 创建处理上下文并设置参数
	ctx := model.NewContext()
	
	// 设置参数
	ctx.SetLanguage("auto") // 自动检测语言
	ctx.SetThreads(4)       // 使用4个线程
	ctx.SetTranslate(false) // 不翻译
	ctx.SetPrintTimestamps(true) // 打印时间戳
	
	fmt.Println("✓ 参数配置完成")
	fmt.Println()
	fmt.Println("开始转录...")
	fmt.Println(strings.Repeat("-", 60))
	
	// 步骤5: 处理音频
	processCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	if err := ctx.Process(processCtx, samples); err != nil {
		return fmt.Errorf("音频处理失败: %v", err)
	}
	
	// 步骤6: 获取并显示结果
	for {
		segment, err := ctx.NextSegment()
		if err != nil {
			break
		}
		
		// 输出带时间戳的文本
		fmt.Printf("[%s --> %s]  %s\n",
			formatTime(segment.Start),
			formatTime(segment.End),
			segment.Text)
	}
	
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("✓ 转录完成")
	
	return nil
}

// convertToWAV 将MP3转换为Whisper所需的WAV格式
func convertToWAV(inputPath string) (string, error) {
	outputPath := strings.TrimSuffix(inputPath, ".mp3") + "_temp.wav"
	
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-ar", "16000",  // 16kHz采样率
		"-ac", "1",       // 单声道
		"-c:a", "pcm_s16le", // 16位PCM
		outputPath,
		"-y", // 覆盖已存在文件
	)
	
	if err := cmd.Run(); err != nil {
		return "", err
	}
	
	return outputPath, nil
}

// readAudioSamples 从WAV文件读取音频样本
func readAudioSamples(wavPath string) ([]float32, error) {
	// 使用FFmpeg读取音频并转换为float32数组
	// 这里简化处理，实际项目中可以使用专门的音频库
	
	// 读取文件内容（跳过WAV头部）
	data, err := os.ReadFile(wavPath)
	if err != nil {
		return nil, err
	}
	
	// WAV文件头通常是44字节
	const headerSize = 44
	if len(data) < headerSize {
		return nil, fmt.Errorf("无效的WAV文件")
	}
	
	pcmData := data[headerSize:]
	samples := make([]float32, len(pcmData)/2)
	
	// 将16位PCM转换为float32
	for i := 0; i < len(samples); i++ {
		// 小端序读取
		sample := int16(pcmData[i*2]) | int16(pcmData[i*2+1])<<8
		samples[i] = float32(sample) / 32768.0 // 归一化到[-1, 1]
	}
	
	return samples, nil
}

// formatTime 格式化时间戳
func formatTime(ms int64) string {
	seconds := ms / 1000
	milliseconds := ms % 1000
	
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60
	
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, secs, milliseconds)
}

// cleanupTempFile 清理临时文件
func cleanupTempFile(path string) {
	os.Remove(path)
}

/* 高级用法示例:

// 1. 批量处理多个音频文件
func batchTranscription(audioPaths []string) {
	model, _ := whisper.New("models/ggml-base.bin")
	defer model.Close()
	
	for _, path := range audioPaths {
		// 处理每个文件
		transcribeFile(model, path)
	}
}

// 2. 实时音频流处理
func realtimeTranscription() {
	model, _ := whisper.New("models/ggml-base.bin")
	defer model.Close()
	
	// 从麦克风或音频流读取数据
	// audioStream := startAudioCapture()
	// 
	// for samples := range audioStream {
	//     ctx := model.NewContext()
	//     ctx.Process(context.Background(), samples)
	//     // 处理结果
	// }
}

// 3. 自定义回调函数
func transcriptionWithCallback() {
	model, _ := whisper.New("models/ggml-base.bin")
	defer model.Close()
	
	ctx := model.NewContext()
	
	// 设置进度回调
	ctx.SetProgressCallback(func(progress float32) {
		fmt.Printf("进度: %.1f%%\n", progress*100)
	})
	
	// 处理音频
	samples, _ := readAudioSamples("audio.wav")
	ctx.Process(context.Background(), samples)
}

// 4. 翻译模式
func translateToEnglish(audioPath string) {
	model, _ := whisper.New("models/ggml-base.bin")
	defer model.Close()
	
	ctx := model.NewContext()
	ctx.SetTranslate(true)  // 启用翻译到英文
	ctx.SetLanguage("zh")   // 源语言是中文
	
	samples, _ := readAudioSamples(audioPath)
	ctx.Process(context.Background(), samples)
	
	// 输出英文翻译
	for {
		segment, err := ctx.NextSegment()
		if err != nil {
			break
		}
		fmt.Println(segment.Text)
	}
}

*/
