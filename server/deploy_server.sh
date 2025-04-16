#!/bin/bash

# 安装依赖
go mod tidy

# 构建服务端程序
go build -o server main.go

# 启动服务端程序
./server    
