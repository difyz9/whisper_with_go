# Legacy Files

本目录包含项目重构前的旧版本文件，保留作为参考。

## CLI 工具文件

### 代码文件
- **main.go** - 原始的 CLI 主程序（使用 Whisper Go 绑定）
- **main_cli.go** - CLI 封装工具（调用 whisper-cli 命令）
- **example.go** - Whisper Go 绑定的示例代码

### 脚本文件
- **build.sh** - 原构建脚本
- **run.sh** - 原运行脚本
- **setup.sh** - 原环境设置脚本
- **extract_audio.sh** - 音频提取脚本
- **Makefile.old** - 原 Makefile

### 文档文件
- **INSTALL_GUIDE.md** - 原安装指南
- **QUICKSTART.md** - 原快速开始指南
- **USAGE.md** - 原使用指南
- **SRT_GUIDE.md** - 原 SRT 字幕指南
- **dev.md** - 原开发文档

## 如何使用旧版 CLI 工具

如果需要使用旧版本的 CLI 工具：

```bash
# 方式 1: 使用 Whisper Go 绑定版本
cd legacy
go run main.go -audio ../test_data/test.mp3 -lang auto

# 方式 2: 使用 CLI 封装版本
go run main_cli.go -audio ../test_data/test.mp3

# 使用旧的 Makefile
make -f Makefile.old help
```

## 为什么保留这些文件？

1. **参考价值** - 展示了从 CLI 工具到 Web API 的演进过程
2. **学习资源** - 包含了 Whisper Go 绑定的使用示例
3. **备份** - 如果需要 CLI 功能，可以随时参考

## 新项目

新的 Web API 项目位于父目录，采用标准的 Gin 框架和分层架构。

请查看主目录的 README.md 了解新项目的使用方法。
