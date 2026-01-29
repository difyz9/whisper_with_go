# 使用指南速查表

## 🚀 快速开始

### 1. 音频文件识别
```bash
# 使用run.sh（推荐）
./run.sh -audio test.mp3                    # 纯文本输出
./run.sh -audio test.mp3 -lang zh           # 指定中文
./run.sh -audio test.mp3 -srt               # 生成SRT字幕 ⭐
./run.sh -audio test.mp3 -srt -lang zh      # SRT字幕（中文）

# 直接运行
DYLD_LIBRARY_PATH=/usr/local/lib go run main.go -audio test.mp3 -srt
```

### 2. 视频提取音频
```bash
# 提取音频
./extract_audio.sh video.mp4              # 提取为MP3
./extract_audio.sh video.mp4 wav          # 提取为WAV

# 提取并直接转录（一步到位）⭐
./extract_audio.sh video.mp4 -transcribe     # 纯文本
./extract_audio.sh video.mp4 -transcribe zh  # 中文
./extract_audio.sh video.mp4 -srt            # SRT字幕 ⭐
./extract_audio.sh video.mp4 -srt zh         # SRT字幕（中文）
```

### 3. 批量处理
```bash
# 批量提取目录下所有MP4文件
./extract_audio.sh -batch mp3

# 批量转录
for file in *.mp3; do
    ./run.sh -audio "$file" -lang zh > "${file%.mp3}.txt"
done
```

## 📝 常用命令

### 模型管理
```bash
# 下载模型
./download_model.sh base      # 推荐，142MB
./download_model.sh small     # 更准确，466MB
./download_model.sh tiny      # 最快，75MB
./download_model.sh medium    # 高精度，1.5GB
```

### 编译和构建
```bash
# 编译程序
./build.sh

# 运行编译后的程序
DYLD_LIBRARY_PATH=/usr/local/lib ./whisper-transcribe -audio test.mp3
```

## 🎯 实际使用场景

### 场景1：会议录音转文字
```bash
# 如果是MP4录屏（纯文本）
./extract_audio.sh meeting.mp4 -transcribe zh

# 如果已经是音频文件
./run.sh -audio meeting.mp3 -lang zh
```

### 场景2：视频字幕制作 ⭐
```bash
# 一步生成SRT字幕（推荐）
./extract_audio.sh video.mp4 -srt zh

# 或分步操作
./extract_audio.sh video.mp4 mp3
./run.sh -audio video.mp3 -srt -lang zh

# 生成的SRT文件可直接导入视频编辑软件
```

### 场景3：多个视频批量生成字幕
```bash
# 批量生成SRT字幕
for mp4 in *.mp4; do
    echo "处理: $mp4"
    ./extract_audio.sh "$mp4" -srt zh
    echo "已生成: ${mp4%.mp4}_audio.srt"
done

# 或批量转录为文本
for mp4 in *.mp4; do
    ./extract_audio.sh "$mp4" -transcribe zh
done
```

### 场景4：英文音频翻译
```bash
# 英文音频直接转录
./run.sh -audio english.mp3 -lang en

# 自动检测语言
./run.sh -audio audio.mp3 -lang auto
```

## 🔧 参数说明

### run.sh / main.go 参数
| 参数 | 说明 | 示例 |
|------|------|------|
| `-audio` | 音频文件路径（必需） | `-audio test.mp3` |
| `-model` | 模型文件路径 | `-model models/ggml-small.bin` |
| `-lang` | 语言代码 | `-lang zh` 或 `-lang en` |
| `-threads` | CPU线程数 | `-threads 8` |
| `-srt` | 输出SRT字幕格式 | `-srt` |
| `-output` | 指定输出文件 | `-output subtitle.srt` |

### extract_audio.sh 参数
| 参数 | 说明 | 示例 |
|------|------|------|
| `<video.mp4>` | 输入视频文件（必需） | `video.mp4` |
| `mp3/wav` | 输出格式 | `mp3` 或 `wav` |
| `<output>` | 输出文件名（可选） | `audio.mp3` |
| `-transcribe` | 提取后直接转录（文本） | `-transcribe zh` |
| `-srt` | 提取后生成SRT字幕 | `-srt zh` |
| `-batch` | 批量处理模式 | `-batch mp3` |

### 语言代码
- `auto` - 自动检测
- `zh` - 中文
- `en` - 英语
- `ja` - 日语
- `ko` - 韩语
- `es` - 西班牙语
- `fr` - 法语
- `de` - 德语
- `ru` - 俄语
- 更多请查看Whisper文档

## 💡 技巧和窍门

### 提高识别准确度
1. 使用更大的模型（small或medium）
2. 指定正确的语言代码（不要用auto）
3. 确保音频质量清晰
4. 使用WAV格式可能更准确

### 加快处理速度
1. 使用smaller模型（tiny或base）
2. 增加线程数 `-threads 8`
3. 使用MP3格式减少文件大小

### 节省存储空间
1. 视频提取后删除临时文件
2. 使用MP3而不是WAV
3. 定期清理转录结果

## 🐛 常见问题

### 问题1: dyld: Library not loaded
```bash
# 解决：使用run.sh脚本
./run.sh -audio test.mp3
```

### 问题2: 找不到ffmpeg
```bash
# 安装ffmpeg
brew install ffmpeg  # macOS
```

### 问题3: 识别结果为空
- 检查音频文件是否完整
- 尝试指定语言 `-lang zh`
- 使用更大的模型

### 问题4: 内存不足
- 使用较小的模型（tiny或base）
- 分段处理长音频

## 📊 性能参考

在Apple M3 Pro上的性能（1分钟音频）：

| 模型 | 时间 | 内存 | 准确度 | 推荐场景 |
|------|------|------|--------|----------|
| tiny | ~3秒 | 273MB | ⭐⭐⭐ | 快速预览 |
| base | ~5秒 | 388MB | ⭐⭐⭐⭐ | 日常使用 ⭐ |
| small | ~15秒 | 852MB | ⭐⭐⭐⭐⭐ | 高质量转录 |
| medium | ~45秒 | 2.1GB | ⭐⭐⭐⭐⭐⭐ | 专业场景 |

## 📚 更多信息

- **完整文档**: README.md
- **快速参考**: QUICKSTART.md
- **安装指南**: INSTALL_GUIDE.md
- **示例代码**: example.go
