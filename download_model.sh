#!/bin/bash
# 下载Whisper GGML模型的脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Whisper模型下载工具${NC}"
echo "======================================"

# 创建models目录
mkdir -p models

# 检查参数
MODEL_NAME=${1:-base}

# 获取模型URL的函数
get_model_url() {
    case "$1" in
        tiny)
            echo "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-tiny.bin"
            ;;
        tiny.en)
            echo "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-tiny.en.bin"
            ;;
        base)
            echo "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin"
            ;;
        base.en)
            echo "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.en.bin"
            ;;
        small)
            echo "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin"
            ;;
        small.en)
            echo "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.en.bin"
            ;;
        medium)
            echo "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-medium.bin"
            ;;
        medium.en)
            echo "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-medium.en.bin"
            ;;
        large)
            echo "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-large-v3.bin"
            ;;
        *)
            echo ""
            ;;
    esac
}

# 获取模型URL
MODEL_URL=$(get_model_url "$MODEL_NAME")

# 检查模型是否在列表中
if [ -z "$MODEL_URL" ]; then
    echo -e "${RED}错误: 未知模型 '$MODEL_NAME'${NC}"
    echo "可用的模型:"
    echo "  - tiny      (75 MB)"
    echo "  - tiny.en   (75 MB, 仅英文)"
    echo "  - base      (142 MB, 推荐)"
    echo "  - base.en   (142 MB, 仅英文)"
    echo "  - small     (466 MB)"
    echo "  - small.en  (466 MB, 仅英文)"
    echo "  - medium    (1.5 GB)"
    echo "  - medium.en (1.5 GB, 仅英文)"
    echo "  - large     (2.9 GB)"
    exit 1
fi
MODEL_FILE="models/ggml-${MODEL_NAME}.bin"

# 检查模型是否已存在
if [ -f "$MODEL_FILE" ]; then
    echo -e "${YELLOW}模型文件已存在: $MODEL_FILE${NC}"
    read -p "是否重新下载? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "取消下载"
        exit 0
    fi
fi

echo -e "${GREEN}下载模型: $MODEL_NAME${NC}"
echo "URL: $MODEL_URL"
echo "目标: $MODEL_FILE"
echo "======================================"

# 下载模型
if command -v wget &> /dev/null; then
    wget --show-progress -O "$MODEL_FILE" "$MODEL_URL"
elif command -v curl &> /dev/null; then
    curl -L --progress-bar -o "$MODEL_FILE" "$MODEL_URL"
else
    echo -e "${RED}错误: 未找到wget或curl命令${NC}"
    echo "请安装wget或curl后重试"
    exit 1
fi

# 验证下载
if [ -f "$MODEL_FILE" ]; then
    FILE_SIZE=$(du -h "$MODEL_FILE" | cut -f1)
    echo -e "${GREEN}✓ 下载完成!${NC}"
    echo "文件: $MODEL_FILE"
    echo "大小: $FILE_SIZE"
else
    echo -e "${RED}✗ 下载失败${NC}"
    exit 1
fi
