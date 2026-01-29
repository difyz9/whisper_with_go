#!/bin/bash
# 从MP4视频文件中提取音频
# 支持提取为MP3或WAV格式

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 显示使用方法
show_usage() {
    echo -e "${GREEN}MP4音频提取工具${NC}"
    echo "======================================"
    echo "用法:"
    echo "  $0 <input.mp4> [output_format] [output_file]"
    echo ""
    echo "参数:"
    echo "  input.mp4       输入的MP4视频文件 (必需)"
    echo "  output_format   输出格式: mp3 或 wav (默认: mp3)"
    echo "  output_file     输出文件名 (可选，默认自动生成)"
    echo ""
    echo "示例:"
    echo "  $0 video.mp4                    # 提取为video.mp3"
    echo "  $0 video.mp4 wav                # 提取为video.wav"
    echo "  $0 video.mp4 mp3 audio.mp3     # 提取为audio.mp3"
    echo "  $0 video.mp4 wav audio.wav     # 提取为audio.wav"
    echo ""
    echo "直接转录（提取后立即识别）:"
    echo "  $0 video.mp4 -transcribe        # 提取并转录（纯文本）"
    echo "  $0 video.mp4 -transcribe zh     # 提取并转录（中文）"
    echo "  $0 video.mp4 -srt               # 提取并生成SRT字幕"
    echo "  $0 video.mp4 -srt zh            # 提取并生成SRT字幕（中文）"
    echo ""
}

# 检查依赖
check_dependencies() {
    if ! command -v ffmpeg &> /dev/null; then
        echo -e "${RED}错误: ffmpeg未安装${NC}"
        echo "请先安装ffmpeg:"
        echo "  macOS:  brew install ffmpeg"
        echo "  Linux:  sudo apt-get install ffmpeg"
        exit 1
    fi
}

# 提取音频为MP3
extract_to_mp3() {
    local input="$1"
    local output="$2"
    
    echo -e "${BLUE}正在提取音频为MP3格式...${NC}"
    ffmpeg -i "$input" -vn -acodec libmp3lame -q:a 2 "$output" -y 2>&1 | grep -E 'Duration|size=' || true
    
    if [ -f "$output" ]; then
        local size=$(du -h "$output" | cut -f1)
        echo -e "${GREEN}✓ 提取成功: $output (大小: $size)${NC}"
        return 0
    else
        echo -e "${RED}✗ 提取失败${NC}"
        return 1
    fi
}

# 提取音频为WAV (Whisper专用格式)
extract_to_wav() {
    local input="$1"
    local output="$2"
    
    echo -e "${BLUE}正在提取音频为WAV格式 (16kHz, 单声道)...${NC}"
    ffmpeg -i "$input" -vn -ar 16000 -ac 1 -c:a pcm_s16le "$output" -y 2>&1 | grep -E 'Duration|size=' || true
    
    if [ -f "$output" ]; then
        local size=$(du -h "$output" | cut -f1)
        echo -e "${GREEN}✓ 提取成功: $output (大小: $size)${NC}"
        return 0
    else
        echo -e "${RED}✗ 提取失败${NC}"
        return 1
    fi
}

# 批量处理目录中的所有MP4文件
batch_extract() {
    local format="${1:-mp3}"
    local count=0
    
    echo -e "${GREEN}批量提取模式${NC}"
    echo "======================================"
    echo "格式: $format"
    echo ""
    
    for mp4_file in *.mp4; do
        if [ -f "$mp4_file" ]; then
            local output="${mp4_file%.mp4}.$format"
            echo "处理: $mp4_file -> $output"
            
            if [ "$format" == "wav" ]; then
                extract_to_wav "$mp4_file" "$output"
            else
                extract_to_mp3 "$mp4_file" "$output"
            fi
            
            ((count++))
            echo ""
        fi
    done
    
    echo -e "${GREEN}批量提取完成! 共处理 $count 个文件${NC}"
}

# 主函数
main() {
    # 检查参数
    if [ $# -eq 0 ]; then
        show_usage
        exit 0
    fi
    
    # 检查依赖
    check_dependencies
    
    # 批量处理模式
    if [ "$1" == "-batch" ]; then
        batch_extract "${2:-mp3}"
        exit 0
    fi
    
    # 获取输入文件
    INPUT_FILE="$1"
    
    # 检查输入文件是否存在
    if [ ! -f "$INPUT_FILE" ]; then
        echo -e "${RED}错误: 文件不存在: $INPUT_FILE${NC}"
        exit 1
    fi
    
    # 获取输出格式（默认mp3）
    OUTPUT_FORMAT="${2:-mp3}"
    
    # 特殊模式：直接转录
    if [ "$OUTPUT_FORMAT" == "-transcribe" ] || [ "$OUTPUT_FORMAT" == "-srt" ]; then
        if [ "$OUTPUT_FORMAT" == "-srt" ]; then
            echo -e "${GREEN}提取并生成SRT字幕${NC}"
        else
            echo -e "${GREEN}提取并转录模式${NC}"
        fi
        echo "======================================"
        
        # 提取音频为MP3
        TEMP_AUDIO="${INPUT_FILE%.mp4}_audio.mp3"
        extract_to_mp3 "$INPUT_FILE" "$TEMP_AUDIO"
        
        # 获取语言参数
        LANG="${3:-auto}"
        
        # 调用whisper转录
        echo ""
        echo -e "${BLUE}开始语音识别...${NC}"
        
        if [ -f "./run.sh" ]; then
            if [ "$OUTPUT_FORMAT" == "-srt" ]; then
                ./run.sh -audio "$TEMP_AUDIO" -lang "$LANG" -srt
            else
                ./run.sh -audio "$TEMP_AUDIO" -lang "$LANG"
            fi
        elif [ -f "./main.go" ]; then
            if [ "$OUTPUT_FORMAT" == "-srt" ]; then
                DYLD_LIBRARY_PATH=/usr/local/lib go run main.go -audio "$TEMP_AUDIO" -lang "$LANG" -srt
            else
                DYLD_LIBRARY_PATH=/usr/local/lib go run main.go -audio "$TEMP_AUDIO" -lang "$LANG"
            fi
        else
            echo -e "${YELLOW}提示: 找不到whisper程序，仅提取了音频文件${NC}"
            echo "音频文件: $TEMP_AUDIO"
        fi
        
        exit 0
    fi
    
    # 检查输出格式
    if [ "$OUTPUT_FORMAT" != "mp3" ] && [ "$OUTPUT_FORMAT" != "wav" ]; then
        echo -e "${RED}错误: 不支持的格式 '$OUTPUT_FORMAT'${NC}"
        echo "支持的格式: mp3, wav"
        exit 1
    fi
    
    # 获取输出文件名
    if [ -n "$3" ]; then
        OUTPUT_FILE="$3"
    else
        OUTPUT_FILE="${INPUT_FILE%.mp4}.$OUTPUT_FORMAT"
    fi
    
    # 显示信息
    echo -e "${GREEN}MP4音频提取工具${NC}"
    echo "======================================"
    echo "输入文件: $INPUT_FILE"
    echo "输出格式: $OUTPUT_FORMAT"
    echo "输出文件: $OUTPUT_FILE"
    echo ""
    
    # 提取音频
    if [ "$OUTPUT_FORMAT" == "wav" ]; then
        extract_to_wav "$INPUT_FILE" "$OUTPUT_FILE"
    else
        extract_to_mp3 "$INPUT_FILE" "$OUTPUT_FILE"
    fi
    
    # 显示后续操作提示
    echo ""
    echo -e "${YELLOW}提示:${NC}"
    echo "  查看音频信息: ffmpeg -i $OUTPUT_FILE"
    echo "  播放音频:     ffplay $OUTPUT_FILE"
    
    if [ -f "./run.sh" ] || [ -f "./main.go" ]; then
        echo "  转录音频:     ./run.sh -audio $OUTPUT_FILE"
        echo "  转录音频(中文): ./run.sh -audio $OUTPUT_FILE -lang zh"
    fi
}

# 运行主函数
main "$@"
