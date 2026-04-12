# Test Data

本目录包含用于测试的音频文件和输出样本。

## 文件说明

### 音频文件
- **test.mp3** - 英文测试音频
- **test.aiff** - AIFF 格式测试音频
- **001.mp3** - 音频提取样本
- **001.mp4** - 视频文件（用于音频提取测试）

### 字幕文件
- **test.srt** - 测试生成的 SRT 字幕
- **test_output.srt** - 输出示例
- **001.srt** - 字幕样本

## 使用这些测试文件

### 测试 API

```bash
# 从项目根目录运行
curl -X POST http://localhost:8080/api/v1/transcribe \
  -F "file=@test_data/test.mp3" \
  -F "language=auto" \
  -F "output_type=json"

# 再通过返回的 task_id 查询结果
curl http://localhost:8080/api/v1/tasks/<task_id>
```

### 测试旧版 CLI

```bash
# 使用 legacy 中的旧版工具
cd legacy
go run main.go -audio ../test_data/test.mp3
```

## 添加更多测试文件

你可以添加更多测试音频文件到此目录。支持的格式：
- MP3
- WAV
- M4A
- AAC
- FLAC
- OGG

## 注意

测试文件不会被提交到 Git 仓库（已在 .gitignore 中配置）。
