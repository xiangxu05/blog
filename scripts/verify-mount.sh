#!/bin/bash
# scripts/verify-mount.sh - 验证 Docker 挂载是否正常工作

set -e

echo "=========================================="
echo "验证 Docker 挂载配置"
echo "=========================================="

# 检查容器是否运行
if ! docker ps | grep -q blog; then
    echo "错误: 容器未运行，请先启动容器"
    echo "运行: ./scripts/docker.sh up"
    exit 1
fi

echo ""
echo "1. 检查 docker-compose.yml 挂载配置..."
docker-compose config | grep -A 5 volumes || echo "未找到挂载配置"

echo ""
echo "2. 检查容器内的挂载信息..."
docker inspect blog | grep -A 10 '"Mounts"' || echo "未找到挂载信息"

echo ""
echo "3. 检查本地 data 目录..."
if [ -d "data" ]; then
    echo "本地 data 目录存在"
    ls -lah data/ | head -10
else
    echo "警告: 本地 data 目录不存在"
fi

echo ""
echo "4. 检查容器内的 data 目录..."
docker exec blog ls -lah /app/data/ 2>/dev/null || echo "容器内 /app/data 目录不存在或无法访问"

echo ""
echo "5. 创建测试文件验证挂载..."
TEST_FILE="data/.mount_test_$(date +%s)"
echo "测试内容 $(date)" > "$TEST_FILE"
echo "在本地创建测试文件: $TEST_FILE"

sleep 1

if docker exec blog test -f "/app/data/.mount_test_$(basename $TEST_FILE)"; then
    echo "✓ 挂载正常！容器内可以看到本地文件"
    rm -f "$TEST_FILE"
    docker exec blog rm -f "/app/data/.mount_test_$(basename $TEST_FILE)" 2>/dev/null || true
else
    echo "✗ 挂载失败！容器内看不到本地文件"
    rm -f "$TEST_FILE"
fi

echo ""
echo "6. 检查数据库文件..."
if [ -f "data/blogData.db" ]; then
    echo "✓ 本地数据库文件存在: data/blogData.db"
    echo "文件大小: $(du -h data/blogData.db | cut -f1)"
    
    if docker exec blog test -f "/app/data/blogData.db"; then
        echo "✓ 容器内数据库文件存在"
        echo "容器内文件大小: $(docker exec blog du -h /app/data/blogData.db | cut -f1)"
    else
        echo "✗ 容器内数据库文件不存在"
    fi
else
    echo "提示: 本地数据库文件不存在（首次启动时会自动创建）"
fi

echo ""
echo "=========================================="
echo "验证完成！"
echo "=========================================="


