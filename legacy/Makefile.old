.PHONY: help setup build run clean test download-model

# 默认目标
help:
	@echo "Whisper.cpp Go 绑定 - 可用命令:"
	@echo ""
	@echo "  make setup          - 设置环境和下载依赖"
	@echo "  make download-model - 下载base模型"
	@echo "  make build          - 编译可执行文件"
	@echo "  make run            - 运行示例（需设置AUDIO变量）"
	@echo "  make test           - 生成并测试音频"
	@echo "  make clean          - 清理生成的文件"
	@echo ""
	@echo "示例:"
	@echo "  make run AUDIO=test.mp3"
	@echo "  make run AUDIO=test.mp3 LANG=zh"
	@echo ""

# 设置环境
setup:
	@echo "检查环境..."
	@which go > /dev/null || (echo "错误: Go未安装" && exit 1)
	@which ffmpeg > /dev/null || (echo "错误: FFmpeg未安装" && exit 1)
	@echo "下载Go依赖..."
	go mod tidy
	@echo "✓ 环境设置完成"

# 下载模型
download-model:
	@echo "下载Whisper模型..."
	@mkdir -p models
	./download_model.sh base

# 编译
build:
	@echo "编译程序..."
	go build -o whisper-transcribe main.go
	@echo "✓ 编译完成: ./whisper-transcribe"

# 运行（需要AUDIO参数）
run:
ifndef AUDIO
	@echo "错误: 请指定AUDIO参数"
	@echo "示例: make run AUDIO=test.mp3"
	@exit 1
endif
	@echo "处理音频文件: $(AUDIO)"
ifdef LANG
	go run main.go -audio $(AUDIO) -lang $(LANG)
else
	go run main.go -audio $(AUDIO)
endif

# 生成测试音频并运行（仅macOS）
test:
	@echo "生成测试音频..."
	@say "Hello, this is a test of the Whisper speech recognition system in Go" -o test.aiff
	@ffmpeg -i test.aiff -ar 16000 -ac 1 test.mp3 -y 2>/dev/null
	@rm test.aiff
	@echo "✓ 测试音频已生成: test.mp3"
	@echo ""
	@echo "运行识别..."
	@go run main.go -audio test.mp3

# 中文测试
test-zh:
	@echo "生成中文测试音频..."
	@say -v "Ting-Ting" "你好，这是一个语音识别系统的测试" -o test_zh.aiff
	@ffmpeg -i test_zh.aiff -ar 16000 -ac 1 test_zh.mp3 -y 2>/dev/null
	@rm test_zh.aiff
	@echo "✓ 测试音频已生成: test_zh.mp3"
	@echo ""
	@echo "运行识别..."
	@go run main.go -audio test_zh.mp3 -lang zh

# 清理
clean:
	@echo "清理文件..."
	@rm -f whisper-transcribe
	@rm -f test.mp3 test_zh.mp3
	@rm -f *_temp.wav
	@echo "✓ 清理完成"

# 完整清理（包括模型）
clean-all: clean
	@echo "清理模型文件..."
	@rm -rf models/*.bin
	@echo "✓ 完整清理完成"

# 安装（复制到系统路径）
install: build
	@echo "安装到 /usr/local/bin ..."
	@sudo cp whisper-transcribe /usr/local/bin/
	@echo "✓ 安装完成，现在可以使用: whisper-transcribe"

# 卸载
uninstall:
	@echo "卸载..."
	@sudo rm -f /usr/local/bin/whisper-transcribe
	@echo "✓ 卸载完成"

# 显示项目信息
info:
	@echo "项目信息:"
	@echo "  Go版本: $$(go version)"
	@echo "  FFmpeg: $$(which ffmpeg || echo '未安装')"
	@echo "  模型目录: ./models"
	@echo "  已下载模型:"
	@ls -lh models/*.bin 2>/dev/null || echo "    无"
