#!/bin/bash
export CGO_ENABLED=1
export CGO_LDFLAGS="-L/usr/local/lib"
export CGO_CFLAGS="-I/usr/local/include"
export DYLD_LIBRARY_PATH=/usr/local/lib

go run cmd/server/main.go

<!-- 下载模型 -->

 bash download_model.sh

 