# Docker 快速开始

## 🚀 三步启动

### 1. 下载模型

```bash
bash download_model.sh
```

### 2. 启动服务

**CPU 版本（推荐初学者）:**
```bash
./scripts/docker.sh build
./scripts/docker.sh up
```

**GPU 版本（需要 NVIDIA GPU）:**
```bash
./scripts/docker.sh build-gpu
./scripts/docker.sh up-gpu
```

### 3. 测试 API

```bash
# 健康检查
curl http://localhost:8080/health

# 转录测试
curl -X POST http://localhost:8080/api/v1/transcribe \
  -F "file=@test_data/test.mp3" \
  -F "output_type=json"
```

## ✅ 就这么简单！

服务运行在：`http://localhost:8080`

## 📚 更多信息

- 详细部署指南：[docker/README.md](docker/README.md)
- API 文档：[docs/API_REFERENCE.md](docs/API_REFERENCE.md)
- 主文档：[README.md](README.md)

## 🛑 停止服务

```bash
./scripts/docker.sh down
```

## 🔄 重启服务

```bash
./scripts/docker.sh restart
```

## 📊 查看日志

```bash
./scripts/docker.sh logs
```
