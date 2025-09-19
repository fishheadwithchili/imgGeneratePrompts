#!/bin/bash
echo "🚀 一键快速设置 - Image Generate Prompts"
echo "========================================"
echo

# 1. 运行项目初始化
echo "📁 初始化项目文件..."
./scripts/init-project.sh

echo
echo "⚠️  重要：请先编辑 apikey/database.env 文件，设置您的数据库密码！"
echo
read -p "按回车继续（请确保已编辑密码）..."
echo

# 2. 启动数据库
echo
echo "🐳 启动 MySQL 数据库..."
docker-compose up mysql -d

# 3. 等待数据库启动
echo "⏳ 等待数据库启动（30秒）..."
sleep 30

# 4. 执行数据库迁移
echo "🔄 执行数据库迁移..."
./scripts/dev.sh migrate

# 5. 启动服务
echo
echo "🌐 启动开发服务器..."
echo "服务将在 http://localhost:8080 启动"
echo
echo "按 Ctrl+C 停止服务"
sleep 3
./scripts/dev.sh start
