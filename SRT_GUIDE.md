# SRT字幕生成指南

## 🎬 什么是SRT字幕？

SRT（SubRip Subtitle）是最常用的字幕格式，几乎所有视频播放器和编辑软件都支持。

## 📝 SRT文件格式

```srt
1
00:00:00,000 --> 00:00:05,600
Hello, this is the test of the Whisper Speech Recognition System using Go.

2
00:00:05,600 --> 00:00:10,200
This is the second subtitle line.
```

格式说明：
- 第1行：序号（从1开始）
- 第2行：时间轴 `开始时间 --> 结束时间`（格式：HH:MM:SS,mmm）
- 第3行：字幕文本
- 第4行：空行分隔

## 🚀 快速生成字幕

### 从音频文件生成
```bash
# 基本用法
./run.sh -audio audio.mp3 -srt

# 指定语言（中文）
./run.sh -audio audio.mp3 -srt -lang zh

# 自定义输出文件名
./run.sh -audio audio.mp3 -srt -output my_subtitle.srt

# 使用更准确的模型
./run.sh -audio audio.mp3 -srt -lang zh -model models/ggml-small.bin
```

### 从视频文件生成（一键完成）⭐
```bash
# 自动提取音频并生成字幕
./extract_audio.sh video.mp4 -srt zh

# 会生成：video_audio.mp3 和 video_audio.srt
```

## 📺 使用生成的字幕

### 1. VLC播放器
```bash
# 将SRT文件与视频放在同一目录，且文件名相同
video.mp4
video.srt

# VLC会自动加载字幕
```

### 2. macOS QuickTime
```bash
# 打开视频
open video.mp4

# 在"窗口" -> "字幕轨道" 中加载SRT文件
```

### 3. 视频编辑软件

**Final Cut Pro:**
1. 导入视频和SRT文件
2. 将SRT拖到时间线
3. 编辑字幕样式

**Adobe Premiere:**
1. 文件 -> 导入 -> 选择SRT文件
2. 拖到字幕轨道
3. 调整样式和位置

**DaVinci Resolve:**
1. 媒体池中右键 -> 导入字幕
2. 将字幕拖到时间线
3. 在字幕面板中编辑

### 4. YouTube / 视频网站
```bash
# 上传视频时，可以直接上传SRT字幕文件
# 支持多语言字幕
```

## 🎯 实际应用场景

### 场景1：会议录像加字幕
```bash
# 1. 生成字幕
./extract_audio.sh meeting.mp4 -srt zh

# 2. 在视频编辑软件中导入字幕
# 3. 调整字幕样式
# 4. 导出带字幕的视频
```

### 场景2：教学视频制作
```bash
# 1. 录制教学视频
# 2. 自动生成字幕
./extract_audio.sh lecture.mp4 -srt zh

# 3. 手动校对字幕文本（用文本编辑器打开SRT）
# 4. 导入视频编辑软件
# 5. 美化字幕样式
```

### 场景3：多语言字幕
```bash
# 生成中文字幕
./extract_audio.sh video.mp4 -srt zh
mv video_audio.srt video_zh.srt

# 生成英文字幕
./extract_audio.sh video.mp4 -srt en
mv video_audio.srt video_en.srt

# 在播放器中可以切换不同语言
```

### 场景4：批量处理
```bash
# 为所有MP4视频生成字幕
for video in *.mp4; do
    echo "处理: $video"
    ./extract_audio.sh "$video" -srt zh
done

# 生成的SRT文件命名为：video_audio.srt
# 重命名以匹配视频文件
for video in *.mp4; do
    audio_srt="${video%.mp4}_audio.srt"
    if [ -f "$audio_srt" ]; then
        mv "$audio_srt" "${video%.mp4}.srt"
        echo "重命名: $audio_srt -> ${video%.mp4}.srt"
    fi
done
```

## ✏️ 编辑SRT字幕

### 使用文本编辑器
```bash
# 用任何文本编辑器打开
vim video.srt
code video.srt
open -a TextEdit video.srt
```

### 常见调整

**调整时间轴：**
```srt
# 原始
00:00:05,600 --> 00:00:10,200

# 延后2秒
00:00:07,600 --> 00:00:12,200
```

**分割长字幕：**
```srt
# 原始
1
00:00:00,000 --> 00:00:10,000
This is a very long subtitle that should be split into multiple lines.

# 分割后
1
00:00:00,000 --> 00:00:05,000
This is a very long subtitle

2
00:00:05,000 --> 00:00:10,000
that should be split into multiple lines.
```

**合并短字幕：**
```srt
# 原始
1
00:00:00,000 --> 00:00:02,000
Hello

2
00:00:02,000 --> 00:00:04,000
World

# 合并后
1
00:00:00,000 --> 00:00:04,000
Hello World
```

## 🔧 高级选项

### 提高字幕准确度
```bash
# 1. 使用更大的模型
./run.sh -audio video.mp3 -srt -model models/ggml-medium.bin -lang zh

# 2. 确保音频质量
./extract_audio.sh video.mp4 wav  # WAV质量更好
./run.sh -audio video.wav -srt -lang zh

# 3. 指定正确的语言（不要用auto）
./run.sh -audio video.mp3 -srt -lang zh
```

### 生成双语字幕
```bash
# 1. 生成中文字幕
./run.sh -audio video.mp3 -srt -lang zh -output video_zh.srt

# 2. 生成英文翻译（如果音频是中文，Whisper不能直接翻译）
# 需要使用翻译工具处理SRT文件

# 3. 合并为双语字幕（手动编辑）
# 在文本编辑器中将两个SRT合并
```

## 📊 字幕质量对比

| 模型 | 准确度 | 速度 | 适用场景 |
|------|--------|------|----------|
| tiny | ⭐⭐⭐ | 最快 | 快速预览 |
| base | ⭐⭐⭐⭐ | 快 | 一般视频 |
| small | ⭐⭐⭐⭐⭐ | 中等 | 高质量字幕 ⭐ |
| medium | ⭐⭐⭐⭐⭐⭐ | 慢 | 专业制作 |

推荐：
- **教学视频**: small或medium
- **会议记录**: base或small
- **快速预览**: base
- **专业制作**: medium

## 💡 技巧和建议

1. **清晰的音频**：字幕质量取决于音频质量
2. **指定语言**：手动指定语言比auto更准确
3. **后期校对**：AI生成的字幕需要人工检查
4. **时间轴调整**：可能需要微调时间以匹配画面
5. **字幕长度**：每行不超过42个字符（英文）或20个字（中文）
6. **显示时长**：至少1秒，最多7秒

## 🐛 常见问题

### 字幕不同步
- 检查视频和字幕的起始点
- 使用字幕编辑软件调整时间轴

### 识别不准确
- 使用更大的模型（small或medium）
- 确保音频清晰
- 指定正确的语言

### 字幕太长
- 手动编辑分割长字幕
- 每条字幕不超过2行

### 乱码问题
- 确保SRT文件编码为UTF-8
- 使用支持UTF-8的编辑器

## 📚 更多资源

- **SRT格式规范**: https://en.wikipedia.org/wiki/SubRip
- **字幕编辑工具**: Aegisub, Subtitle Edit
- **在线字幕编辑**: https://www.happyscribe.com/subtitle-tools
