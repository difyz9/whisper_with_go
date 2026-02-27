package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 简单的Go封装whisper.cpp CLI工具
// 这个版本不需要CGO，直接调用whisper-cli命令

func main() {
	// 解析命令行参数
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	var audioPath, modelPath, language string
	threads := "4"

	// 简单的参数解析
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-audio":
			if i+1 < len(os.Args) {
				audioPath = os.Args[i+1]
				i++
			}
		case "-model":
			if i+1 < len(os.Args) {
				modelPath = os.Args[i+1]
				i++
			}
		case "-lang":
			if i+1 < len(os.Args) {
				language = os.Args[i+1]
				i++
			}
		case "-threads":
			if i+1 < len(os.Args) {
				threads = os.Args[i+1]
				i++
			}
		}
	}

	if audioPath == "" {
		fmt.Println("错误: 必须指定 -audio 参数")
		printUsage()
		os.Exit(1)
	}

	if modelPath == "" {
		modelPath = "models/ggml-base.bin"
	}

	if language == "" {
		language = "auto"
	}

	// 检查文件是否存在
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		fmt.Printf("错误: 音频文件不存在: %s\n", audioPath)
		os.Exit(1)
	}

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		fmt.Printf("错误: 模型文件不存在: %s\n", modelPath)
		fmt.Println("请先下载模型: ./download_model.sh base")
		os.Exit(1)
	}

	// 检查whisper-cli是否安装
	if _, err := exec.LookPath("whisper-cli"); err != nil {
		fmt.Println("错误: whisper-cli未找到")
		fmt.Println("\n请先安装whisper.cpp:")
		fmt.Println("  git clone https://github.com/ggerganov/whisper.cpp.git")
		fmt.Println("  cd whisper.cpp")
		fmt.Println("  cmake -B build && cmake --build build -j")
		fmt.Println("  sudo cp build/bin/whisper-cli /usr/local/bin/")
		fmt.Println("\n或者查看 INSTALL_GUIDE.md 了解详细步骤")
		os.Exit(1)
	}

	fmt.Printf("音频文件: %s\n", audioPath)
	fmt.Printf("模型: %s\n", modelPath)
	fmt.Printf("语言: %s\n", language)
	fmt.Println(strings.Repeat("-", 60))

	// 如果是MP3，先转换为WAV
	wavPath := audioPath
	needCleanup := false

	if strings.HasSuffix(strings.ToLower(audioPath), ".mp3") {
		fmt.Println("检测到MP3格式，正在转换为WAV...")
		wavPath = strings.TrimSuffix(audioPath, filepath.Ext(audioPath)) + "_temp.wav"
		
		if err := convertToWAV(audioPath, wavPath); err != nil {
			fmt.Printf("转换失败: %v\n", err)
			os.Exit(1)
		}
		needCleanup = true
		defer func() {
			if needCleanup {
				os.Remove(wavPath)
			}
		}()
		fmt.Println("转换完成")
	}

	// 调用whisper-cli
	fmt.Println("开始识别...")
	fmt.Println(strings.Repeat("-", 60))

	args := []string{
		"-m", modelPath,
		"-f", wavPath,
		"-t", threads,
	}

	if language != "auto" {
		args = append(args, "-l", language)
	}

	cmd := exec.Command("whisper-cli", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("\n识别失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("识别完成!")
}

func convertToWAV(inputPath, outputPath string) error {
	// 检查ffmpeg
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg未安装，请先安装: brew install ffmpeg")
	}

	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-ar", "16000",
		"-ac", "1",
		"-c:a", "pcm_s16le",
		outputPath,
		"-y",
	)

	// 隐藏ffmpeg的输出
	cmd.Stderr = nil
	cmd.Stdout = nil

	return cmd.Run()
}

func printUsage() {
	fmt.Println("Go Whisper 语音识别工具 (CLI封装版本)")
	fmt.Println("\n用法:")
	fmt.Println("  go run main_cli.go -audio <音频文件> [选项]")
	fmt.Println("\n选项:")
	fmt.Println("  -audio <path>   音频文件路径 (必需)")
	fmt.Println("  -model <path>   模型文件路径 (默认: models/ggml-base.bin)")
	fmt.Println("  -lang <code>    语言代码 (默认: auto)")
	fmt.Println("  -threads <n>    线程数 (默认: 4)")
	fmt.Println("\n示例:")
	fmt.Println("  go run main_cli.go -audio test.mp3")
	fmt.Println("  go run main_cli.go -audio test.mp3 -lang zh")
	fmt.Println("  go run main_cli.go -audio test.mp3 -model models/ggml-small.bin")
}
