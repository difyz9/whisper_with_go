#!/bin/bash

# API 测试脚本

BASE_URL="http://localhost:8084"

echo "================================"
echo "Whisper API 测试脚本"
echo "================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. 测试健康检查
printf "${YELLOW}1. 测试健康检查${NC}\n"
echo "GET $BASE_URL/health"
curl -s $BASE_URL/health | jq .
echo ""
echo ""

# 2. 测试根路径
printf "${YELLOW}2. 测试根路径${NC}\n"
echo "GET $BASE_URL/"
curl -s $BASE_URL/ | jq .
echo ""
echo ""

# 3. 测试转录 API（需要音频文件）
if [ -f "test.mp3" ] || [ -f "test.wav" ]; then
    printf "${YELLOW}3. 测试转录 API${NC}\n"
    
    # 查找测试文件
    TEST_FILE=""
    if [ -f "001.mp3" ]; then
        TEST_FILE="001.mp3"
    elif [ -f "001.wav" ]; then
        TEST_FILE="001.wav"
    fi
    
    echo "POST $BASE_URL/api/v1/transcribe"
    echo "文件: $TEST_FILE"
    echo ""
    
    # JSON 格式输出
    printf "${GREEN}测试 JSON 输出:${NC}\n"
    curl -X POST $BASE_URL/api/v1/transcribe \
        -F "file=@$TEST_FILE" \
        -F "language=auto" \
        -F "output_type=json" | jq .
    echo ""
    echo ""
    
    # SRT 格式输出
    printf "${GREEN}测试 SRT 输出:${NC}\n"
    curl -X POST $BASE_URL/api/v1/transcribe \
        -F "file=@$TEST_FILE" \
        -F "language=auto" \
        -F "output_type=srt" | jq .
    echo ""
    echo ""
    
else
    printf "${RED}3. 跳过转录测试（未找到 test.mp3 或 test.wav 文件）${NC}\n"
    echo ""
fi

# 4. 测试错误处理
printf "${YELLOW}4. 测试错误处理${NC}\n"
echo "POST $BASE_URL/api/v1/transcribe (无文件)"
curl -X POST $BASE_URL/api/v1/transcribe | jq .
echo ""
echo ""

echo "================================"
echo "测试完成"
echo "================================"
