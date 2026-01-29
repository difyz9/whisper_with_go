package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	whisper "github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

func main() {
	// 命令行参数
	var (
		modelPath  = flag.String("model", "models/ggml-base.bin", "Whisper模型文件路径")
		audioPath  = flag.String("audio", "", "MP3音频文件路径")
		language   = flag.String("lang", "auto", "语言代码 (auto, zh, en, etc.)")
		threads    = flag.Int("threads", 4, "使用的线程数")
		outputSRT  = flag.Bool("srt", false, "输出为SRT字幕格式")
		outputFile = flag.String("output", "", "输出文件路径 (默认: 音频文件名.srt)")
	)
	flag.Parse()

	if *audioPath == "" {
		fmt.Println("错误: 必须指定音频文件路径")
		fmt.Println("使用方法: go run main.go -audio <mp3文件路径> [-model <模型路径>] [-lang <语言>] [-threads <线程数>] [-srt] [-output <输出文件>]")
		fmt.Println("\n选项:")
		fmt.Println("  -srt          输出为SRT字幕格式")
		fmt.Println("  -output FILE  指定输出文件（默认: 音频文件名.srt或.txt）")
		os.Exit(1)
	}

	// 检查音频文件是否存在
	if _, err := os.Stat(*audioPath); os.IsNotExist(err) {
		fmt.Printf("错误: 音频文件不存在: %s\n", *audioPath)
		os.Exit(1)
	}

	// 检查模型文件是否存在
	if _, err := os.Stat(*modelPath); os.IsNotExist(err) {
		fmt.Printf("错误: 模型文件不存在: %s\n", *modelPath)
		fmt.Println("请先下载模型文件，例如:")
		fmt.Println("  bash models/download-ggml-model.sh base")
		os.Exit(1)
	}

	// 转换MP3到WAV格式 (Whisper需要16kHz, 单声道, 16位PCM WAV)
	wavPath, err := convertMP3ToWAV(*audioPath)
	if err != nil {
		fmt.Printf("错误: 转换MP3到WAV失败: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(wavPath) // 清理临时文件

	// 确定输出文件路径
	var outPath string
	if *outputFile != "" {
		outPath = *outputFile
	} else {
		// 根据输入文件名生成输出文件名
		baseName := strings.TrimSuffix(*audioPath, filepath.Ext(*audioPath))
		if *outputSRT {
			outPath = baseName + ".srt"
		} else {
			outPath = baseName + ".txt"
		}
	}

	fmt.Printf("音频文件: %s\n", *audioPath)
	fmt.Printf("临时WAV文件: %s\n", wavPath)
	fmt.Printf("模型文件: %s\n", *modelPath)
	fmt.Printf("语言: %s\n", *language)
	fmt.Printf("输出文件: %s\n", outPath)
	if *outputSRT {
		fmt.Println("输出格式: SRT字幕")
	} else {
		fmt.Println("输出格式: 纯文本")
	}
	fmt.Println("开始转录...")
	fmt.Println(strings.Repeat("-", 60))

	// 执行语音识别
	if err := transcribe(*modelPath, wavPath, *language, *threads, *outputSRT, outPath); err != nil {
		fmt.Printf("错误: 转录失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("✓ 转录完成! 结果已保存到: %s\n", outPath)
}

// convertMP3ToWAV 将MP3文件转换为WAV格式 (16kHz, 单声道, 16位PCM)
func convertMP3ToWAV(mp3Path string) (string, error) {
	// 检查ffmpeg是否安装
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return "", fmt.Errorf("ffmpeg未安装，请先安装ffmpeg: brew install ffmpeg (macOS) 或 apt-get install ffmpeg (Linux)")
	}

	// 创建临时WAV文件
	wavPath := strings.TrimSuffix(mp3Path, filepath.Ext(mp3Path)) + "_temp.wav"

	// 使用ffmpeg转换
	// -i: 输入文件
	// -ar 16000: 采样率16kHz
	// -ac 1: 单声道
	// -c:a pcm_s16le: 16位PCM编码
	cmd := exec.Command("ffmpeg", "-i", mp3Path, "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", wavPath, "-y")
	
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg转换失败: %v", err)
	}

	return wavPath, nil
}

// transcribe 执行语音识别
func transcribe(modelPath, wavPath, language string, threads int, outputSRT bool, outputPath string) error {
	// 加载模型
	model, err := whisper.New(modelPath)
	if err != nil {
		return fmt.Errorf("加载模型失败: %v", err)
	}
	defer model.Close()

	// 读取WAV文件
	samples, err := readWAVFile(wavPath)
	if err != nil {
		return fmt.Errorf("读取WAV文件失败: %v", err)
	}

	// 创建处理上下文
	context, err := model.NewContext()
	if err != nil {
		return fmt.Errorf("创建上下文失败: %v", err)
	}
	
	// 设置语言
	if language != "auto" {
		if err := context.SetLanguage(language); err != nil {
			return fmt.Errorf("设置语言失败: %v", err)
		}
	}
	
	// 设置线程数
	context.SetThreads(uint(threads))
	
	// 启用翻译模式（如果需要）
	context.SetTranslate(false)

	// 处理音频
	if err := context.Process(samples, nil, nil, nil); err != nil {
		return fmt.Errorf("处理音频失败: %v", err)
	}

	// 创建输出文件
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer outFile.Close()

	// 收集所有片段
	var segments []whisper.Segment
	for {
		segment, err := context.NextSegment()
		if err != nil {
			break
		}
		segments = append(segments, segment)
	}

	// 根据格式输出
	if outputSRT {
		// 输出SRT格式
		for i, segment := range segments {
			// SRT序号（从1开始）
			fmt.Fprintf(outFile, "%d\n", i+1)
			
			// SRT时间格式: HH:MM:SS,mmm --> HH:MM:SS,mmm
			fmt.Fprintf(outFile, "%s --> %s\n", 
				formatSRTTime(segment.Start),
				formatSRTTime(segment.End))
			
			// 字幕文本
			fmt.Fprintf(outFile, "%s\n\n", strings.TrimSpace(segment.Text))
			
			// 同时输出到控制台
			fmt.Printf("[%6s --> %6s]  %s\n",
				segment.Start.Truncate(time.Millisecond),
				segment.End.Truncate(time.Millisecond),
				segment.Text)
		}
	} else {
		// 输出纯文本格式
		for _, segment := range segments {
			// 带时间戳的文本
			fmt.Fprintf(outFile, "[%s --> %s]  %s\n",
				segment.Start.Truncate(time.Millisecond),
				segment.End.Truncate(time.Millisecond),
				segment.Text)
			
			// 同时输出到控制台
			fmt.Printf("[%6s --> %6s]  %s\n",
				segment.Start.Truncate(time.Millisecond),
				segment.End.Truncate(time.Millisecond),
				segment.Text)
		}
	}

	return nil
}

// readWAVFile 读取WAV文件并返回音频样本
func readWAVFile(wavPath string) ([]float32, error) {
	// 读取文件内容
	data, err := os.ReadFile(wavPath)
	if err != nil {
		return nil, err
	}

	// 跳过WAV文件头（通常44字节）
	const wavHeaderSize = 44
	if len(data) < wavHeaderSize {
		return nil, fmt.Errorf("WAV文件太小")
	}

	// 将PCM数据转换为float32样本
	pcmData := data[wavHeaderSize:]
	samples := make([]float32, len(pcmData)/2)
	
	for i := 0; i < len(samples); i++ {
		// 读取16位PCM样本 (小端序)
		sample := int16(pcmData[i*2]) | int16(pcmData[i*2+1])<<8
		// 转换为-1.0到1.0的浮点数
		samples[i] = float32(sample) / 32768.0
	}

	return samples, nil
}

// formatSRTTime 格式化时间为SRT格式 (HH:MM:SS,mmm)
func formatSRTTime(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000
	
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, milliseconds)
}
