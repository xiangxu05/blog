#!/bin/bash
# scripts/docker.sh - Docker Compose 启动脚本

set -e

echo "=========================================="
echo "Docker Compose 启动脚本"
echo "=========================================="

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: 未找到 Docker，请先安装 Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "错误: 未找到 docker-compose，请先安装 docker-compose"
    exit 1
fi

echo "Docker 版本: $(docker --version)"
echo "Docker Compose 版本: $(docker-compose --version)"

# 检查 docker-compose.yml 是否存在
if [ ! -f "docker-compose.yml" ]; then
    echo "错误: 未找到 docker-compose.yml 文件"
    exit 1
fi

# 检查数据目录
if [ ! -d "data" ]; then
    echo "创建 data 目录..."
    mkdir -p data
    chmod 755 data
fi

# 解析命令行参数
ACTION="${1:-up}"

case "$ACTION" in
    build)
        echo "=========================================="
        echo "构建 Docker 镜像..."
        echo "=========================================="
        docker-compose build
        ;;
    up)
        echo "=========================================="
        echo "构建并启动服务（后台运行）..."
        echo "=========================================="
        docker-compose up -d --build
        echo ""
        echo "服务已启动！"
        echo "访问地址: http://localhost:80"
        echo ""
        echo "查看日志: docker-compose logs -f"
        echo "停止服务: docker-compose down"
        ;;
    down)
        echo "=========================================="
        echo "停止并删除服务..."
        echo "=========================================="
        docker-compose down
        ;;
    restart)
        echo "=========================================="
        echo "重启服务..."
        echo "=========================================="
        docker-compose restart
        ;;
    logs)
        echo "=========================================="
        echo "查看服务日志..."
        echo "=========================================="
        docker-compose logs -f
        ;;
    ps)
        echo "=========================================="
        echo "查看服务状态..."
        echo "=========================================="
        docker-compose ps
        ;;
    shell)
        echo "=========================================="
        echo "进入容器 Shell..."
        echo "=========================================="
        docker-compose exec blog sh
        ;;
    clean)
        echo "=========================================="
        echo "清理 Docker 资源..."
        echo "=========================================="
        docker-compose down -v
        docker system prune -f
        echo "清理完成！"
        ;;
    *)
        echo "用法: $0 [build|up|down|restart|logs|ps|shell|clean]"
        echo ""
        echo "命令说明:"
        echo "  build   - 构建 Docker 镜像"
        echo "  up      - 构建并启动服务（默认）"
        echo "  down    - 停止并删除服务"
        echo "  restart - 重启服务"
        echo "  logs    - 查看服务日志"
        echo "  ps      - 查看服务状态"
        echo "  shell   - 进入容器 Shell"
        echo "  clean   - 清理 Docker 资源"
        exit 1
        ;;
esac

