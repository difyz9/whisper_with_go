param(
    [ValidateSet('up', 'down', 'restart', 'logs', 'ps', 'pull')]
    [string]$Action = 'up',

    [string]$Image = $env:WHISPER_BASE_IMAGE,

    [switch]$NoPull,

    [switch]$Foreground
)

$ErrorActionPreference = 'Stop'

function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host "[OK] $Message" -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Fail {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
    exit 1
}

function Test-Command {
    param([string]$Name)
    return $null -ne (Get-Command $Name -ErrorAction SilentlyContinue)
}

function Ensure-Docker {
    if (-not (Test-Command 'docker')) {
        Fail 'Docker 未安装或未加入 PATH。'
    }

    docker compose version | Out-Null
}

function Ensure-ProjectDirectories {
    foreach ($path in @('models', 'uploads', 'outputs')) {
        if (-not (Test-Path $path)) {
            New-Item -ItemType Directory -Path $path | Out-Null
            Write-Info "已创建目录: $path"
        }
    }
}

function Warn-IfModelMissing {
    if (-not (Test-Path 'models/ggml-base.bin')) {
        Write-Warn '未找到 models/ggml-base.bin，服务可以启动，但转录请求会因缺少模型失败。'
    }
}

function Get-ComposeArgs {
    return @('-f', 'docker-compose.base.yml')
}

function Require-Image {
    if ([string]::IsNullOrWhiteSpace($Image)) {
        Fail '请通过 -Image 参数或 WHISPER_BASE_IMAGE 环境变量提供基础镜像地址，例如 difyz9/whisper-go-base:latest。'
    }
}

Ensure-Docker
Ensure-ProjectDirectories
Warn-IfModelMissing

switch ($Action) {
    'pull' {
        Require-Image
        $env:WHISPER_BASE_IMAGE = $Image
        Write-Info "拉取基础镜像: $Image"
        docker pull $Image
        Write-Success '镜像拉取完成。'
    }

    'up' {
        Require-Image
        $env:WHISPER_BASE_IMAGE = $Image

        if (-not $NoPull) {
            Write-Info "拉取基础镜像: $Image"
            docker pull $Image
        }

        $composeArgs = Get-ComposeArgs

        if ($Foreground) {
            Write-Info '以前台模式启动基础镜像服务。'
            docker compose @composeArgs up
        } else {
            Write-Info '以后台模式启动基础镜像服务。'
            docker compose @composeArgs up -d
            Write-Success '服务已启动。'
            Write-Host '访问: http://localhost:8080'
            Write-Host '健康检查: http://localhost:8080/health'
        }
    }

    'down' {
        Write-Info '停止基础镜像服务。'
        docker compose @(Get-ComposeArgs) down
        Write-Success '服务已停止。'
    }

    'restart' {
        Require-Image
        $env:WHISPER_BASE_IMAGE = $Image

        Write-Info '重启基础镜像服务。'
        docker compose @(Get-ComposeArgs) down

        if (-not $NoPull) {
            Write-Info "拉取基础镜像: $Image"
            docker pull $Image
        }

        docker compose @(Get-ComposeArgs) up -d
        Write-Success '服务已重启。'
    }

    'logs' {
        Write-Info '查看服务日志。'
        docker compose @(Get-ComposeArgs) logs -f
    }

    'ps' {
        Write-Info '查看服务状态。'
        docker compose @(Get-ComposeArgs) ps
    }
}