#!/bin/bash
# scripts/check-files.sh - 检查文件挂载和访问情况

set -e

echo "=========================================="
echo "检查文件挂载和访问情况"
echo "=========================================="

# 检查容器是否运行
if ! docker ps | grep -q blog; then
    echo "错误: 容器未运行"
    exit 1
fi

echo ""
echo "1. 检查本地 files 目录..."
if [ -d "data/files" ]; then
    echo "✓ 本地 files 目录存在"
    echo "  文件数量: $(ls -1 data/files/ 2>/dev/null | wc -l)"
    echo "  目录大小: $(du -sh data/files/ | cut -f1)"
else
    echo "✗ 本地 files 目录不存在"
fi

echo ""
echo "2. 检查容器内 files 目录..."
if docker exec blog test -d /app/data/files; then
    echo "✓ 容器内 files 目录存在"
    FILE_COUNT=$(docker exec blog sh -c "ls -1 /app/data/files/ 2>/dev/null | wc -l")
    echo "  文件数量: $FILE_COUNT"
    echo "  目录大小: $(docker exec blog du -sh /app/data/files/ | cut -f1)"
else
    echo "✗ 容器内 files 目录不存在"
fi

echo ""
echo "3. 检查数据库中的文件路径..."
docker exec blog sh -c "cd /app && sqlite3 data/blogData.db 'SELECT id, filename, store_path FROM tb_file LIMIT 5;'" 2>/dev/null || echo "无法查询数据库"

echo ""
echo "4. 测试文件访问..."
TEST_FILE=$(ls data/files/*.jpg 2>/dev/null | head -1)
if [ -n "$TEST_FILE" ]; then
    BASENAME=$(basename "$TEST_FILE")
    echo "测试文件: $BASENAME"
    
    if docker exec blog test -f "/app/data/files/$BASENAME"; then
        echo "✓ 容器内可以访问文件"
        
        # 检查文件路径在数据库中的存储
        DB_PATH=$(docker exec blog sh -c "cd /app && sqlite3 data/blogData.db \"SELECT store_path FROM tb_file WHERE store_path LIKE '%$BASENAME%' LIMIT 1;\"" 2>/dev/null || echo "")
        if [ -n "$DB_PATH" ]; then
            echo "  数据库中的路径: $DB_PATH"
            
            # 检查应用是否能访问
            if docker exec blog test -f "/app/$DB_PATH"; then
                echo "✓ 应用可以访问数据库中的路径"
            else
                echo "✗ 应用无法访问数据库中的路径: /app/$DB_PATH"
                echo "  提示: 检查路径是相对路径还是绝对路径"
            fi
        else
            echo "  提示: 数据库中未找到该文件的记录"
        fi
    else
        echo "✗ 容器内无法访问文件"
    fi
else
    echo "提示: 没有找到测试文件"
fi

echo ""
echo "=========================================="
echo "检查完成！"
echo "=========================================="


