#!/bin/bash

# 安装依赖
go mod tidy

# 构建 Agent 端程序
go build -o agent main.go

# 启动 Agent 端程序
./agent    
