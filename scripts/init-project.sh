#!/bin/bash
echo "🚀 项目初始化脚本"
echo "========================================"
echo "正在为您设置开发环境..."
echo

# 1. 创建 apikey 目录
if [ ! -d "apikey" ]; then
    echo "📁 创建 apikey 目录..."
    mkdir -p apikey
fi

# 2. 复制数据库配置文件
if [ ! -f "apikey/database.env" ]; then
    echo "📄 创建数据库配置文件..."
    cp apikey/database.env.example apikey/database.env
    echo "✅ 已创建 apikey/database.env"
    echo "⚠️  请编辑 apikey/database.env 文件，将 YOUR_PASSWORD_HERE 替换为您的实际数据库密码"
else
    echo "✅ apikey/database.env 已存在"
fi

# 3. 复制 Docker 配置文件
if [ ! -f "docker-compose.override.yml" ]; then
    echo "🐳 创建 Docker 开发配置..."
    if [ -f "docker-compose.override.yml.template" ]; then
        cp docker-compose.override.yml.template docker-compose.override.yml
        echo "✅ 已从模板创建 docker-compose.override.yml"
    else
        cat > docker-compose.override.yml << 'EOF'
# docker-compose.override.yml
# 本地开发环境覆盖配置
# 此文件包含敏感信息，已被 .gitignore 忽略

version: '3.8'

services:
  mysql:
    environment:
      MYSQL_ROOT_PASSWORD: YOUR_MYSQL_ROOT_PASSWORD
      MYSQL_PASSWORD: YOUR_MYSQL_PASSWORD
EOF
        echo "✅ 已创建 docker-compose.override.yml"
    fi
else
    echo "✅ docker-compose.override.yml 已存在"
fi

# 4. 创建日志目录
if [ ! -d "logs" ]; then
    echo "📝 创建日志目录..."
    mkdir -p logs
    touch logs/.gitkeep
fi

# 5. 确保上传目录存在
if [ ! -d "uploads" ]; then
    echo "📤 创建上传目录..."
    mkdir -p uploads
    touch uploads/.gitkeep
fi

# 6. 下载 Go 依赖
echo "📦 下载 Go 依赖..."
go mod download
go mod tidy

echo
echo "✅ 项目初始化完成！"
echo
echo "📋 下一步操作："
echo "1. 编辑 apikey/database.env 文件，将 YOUR_PASSWORD_HERE 替换为您的数据库密码"
echo "2. 启动数据库：docker-compose up mysql -d"
echo "3. 执行数据库迁移：./scripts/dev.sh migrate"
echo "4. 启动服务：./scripts/dev.sh start"
echo
echo "🔧 或者使用一键启动："
echo "   ./scripts/quick-setup.sh"
