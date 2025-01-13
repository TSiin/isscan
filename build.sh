#!/bin/bash

# 设置环境变量
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo "开始构建 Server Manager..."

# 检查依赖
echo "检查依赖..."
go mod tidy
if [ $? -ne 0 ]; then
    echo -e "${RED}依赖检查失败${NC}"
    exit 1
fi

# 运行测试
echo "运行测试..."
go test ./...
if [ $? -ne 0 ]; then
    echo -e "${RED}测试失败${NC}"
    exit 1
fi

# 构建
echo "编译程序..."
go build -o isscan main.go
if [ $? -ne 0 ]; then
    echo -e "${RED}编译失败${NC}"
    exit 1
fi

echo -e "${GREEN}构建成功!${NC}" 