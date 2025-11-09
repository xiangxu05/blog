#!/bin/sh
set -e

echo "确保 data 目录存在（支持挂载卷）..."
mkdir -p data

echo "设置 data 目录权限（如果需要）..."
# 确保 data 目录可写（如果挂载的卷没有正确权限）
chmod 755 data 2>/dev/null || true

echo "检查数据库文件是否存在..."
if [ ! -f "data/blogData.db" ]; then
    echo "数据库文件不存在，创建空的数据库文件..."
    # 先创建空的数据库文件，确保目录权限正确
    touch data/blogData.db
    chmod 644 data/blogData.db 2>/dev/null || true
    
    echo "初始化数据库..."
    ./blog regen || {
        echo "数据库初始化失败！"
        exit 1
    }
    echo "数据库初始化成功！"
else
    echo "数据库文件已存在，跳过初始化"
fi

echo "启动应用..."
exec ./blog blogsys

