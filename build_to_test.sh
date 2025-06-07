#!/bin/bash

# 清理之前的构建
rm -f yop

# 构建服务器
export CGO_ENABLED=0
go build -o yop cmd/server/main.go

# 检查构建结果
if [ -f "yop" ]; then
  echo "构建成功"
  ./yop
else
  echo "构建失败"
  exit 1
fi