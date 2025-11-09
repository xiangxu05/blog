#!/bin/bash
# scripts/dev.sh - 本地开发环境启动脚本

set -e

echo "=========================================="
echo "启动本地开发环境"
echo "=========================================="

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo "错误: 未找到 Go 环境，请先安装 Go"
    exit 1
fi

echo "Go 版本: $(go version)"

# 检查数据库目录
if [ ! -d "data" ]; then
    echo "创建 data 目录..."
    mkdir -p data
fi

# 检查数据库文件是否存在
if [ ! -f "data/blogData.db" ]; then
    echo "数据库文件不存在，初始化数据库..."
    go run main.go regen || {
        echo "数据库初始化失败！"
        exit 1
    }
    echo "数据库初始化成功！"
else
    echo "数据库文件已存在，跳过初始化"
fi

# 设置环境变量
export GIN_MODE=debug
export TZ=Asia/Shanghai

# 启动应用
echo "=========================================="
echo "启动应用服务器..."
echo "访问地址: http://localhost:80"
echo "按 Ctrl+C 停止服务"
echo "=========================================="

go run main.go blogsys

