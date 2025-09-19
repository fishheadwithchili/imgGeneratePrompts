#!/bin/bash
# 一键开始脚本 - Image Generate Prompts API
# 用法: ./quick-start.sh

set -e

echo "🚀 Image Generate Prompts API - 一键开始"
echo "========================================"
echo ""

# 检查是否安装了必要的工具
check_requirements() {
    echo "🔍 检查系统要求..."
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        echo "❌ Go 未安装。请先安装 Go 1.21 或更高版本。"
        echo "   下载地址: https://golang.org/dl/"
        exit 1
    fi
    
    local go_version=$(go version | grep -oE '[0-9]+\.[0-9]+' | head -1)
    echo "✅ Go 版本: $go_version"
    
    # 检查 MySQL
    if ! command -v mysql &> /dev/null; then
        echo "⚠️  MySQL 客户端未安装。请确保 MySQL 服务器运行在 3307 端口。"
    else
        echo "✅ MySQL 客户端已安装"
    fi
    
    # 检查 curl
    if ! command -v curl &> /dev/null; then
        echo "⚠️  curl 未安装，部分测试功能可能无法使用。"
    else
        echo "✅ curl 已安装"
    fi
    
    echo ""
}

# 设置项目
setup_project() {
    echo "📦 设置项目..."
    
    # 下载依赖
    echo "  - 下载 Go 依赖..."
    go mod download
    go mod tidy
    
    # 创建必要目录
    echo "  - 创建目录..."
    mkdir -p uploads logs bin
    
    # 设置脚本权限
    if [ -f "scripts/dev.sh" ]; then
        chmod +x scripts/dev.sh
    fi
    if [ -f "scripts/test-api.sh" ]; then
        chmod +x scripts/test-api.sh
    fi
    
    echo "✅ 项目设置完成"
    echo ""
}

# 初始化数据库
init_database() {
    echo "🗄️  初始化数据库..."
    
    # 检查数据库配置文件
    if [ ! -f "apikey/database.env" ]; then
        echo "❌ 数据库配置文件不存在: apikey/database.env"
        echo "   请确保文件存在并包含正确的数据库配置。"
        exit 1
    fi
    
    # 运行数据库初始化
    echo "  - 创建表结构和示例数据..."
    if go run cmd/db-manager.go -write; then
        echo "✅ 数据库初始化成功"
    else
        echo "❌ 数据库初始化失败"
        echo "   请检查:"
        echo "   1. MySQL 服务是否在 3307 端口运行"
        echo "   2. 数据库配置是否正确"
        echo "   3. 数据库 'img_generate_prompts' 是否已创建"
        exit 1
    fi
    echo ""
}

# 启动服务
start_service() {
    echo "🌐 启动开发服务器..."
    echo ""
    echo "📊 数据库统计信息:"
    go run cmd/db-manager.go -stats
    echo ""
    echo "🚀 服务器启动中..."
    echo "   访问地址: http://localhost:8080"
    echo "   API文档: http://localhost:8080/api/v1"
    echo "   健康检查: http://localhost:8080/health"
    echo ""
    echo "按 Ctrl+C 停止服务器"
    echo "========================================"
    
    # 启动服务器
    go run main.go
}

# 显示使用帮助
show_help() {
    echo "🛠️  开发命令:"
    echo ""
    echo "# 查看项目信息"
    echo "curl http://localhost:8080/"
    echo ""
    echo "# 健康检查"
    echo "curl http://localhost:8080/health"
    echo ""
    echo "# 获取标签列表"
    echo "curl http://localhost:8080/api/v1/tags"
    echo ""
    echo "# 获取提示词列表"
    echo "curl http://localhost:8080/api/v1/prompts"
    echo ""
    echo "# 运行完整API测试"
    echo "./scripts/test-api.sh"
    echo ""
    echo "# 数据库管理"
    echo "go run cmd/db-manager.go -stats    # 查看统计"
    echo "go run cmd/db-manager.go -sample   # 添加示例数据"
    echo "go run cmd/db-manager.go -reset    # 重置数据库"
    echo ""
    echo "# 使用 Makefile"
    echo "make help                         # 查看所有命令"
    echo "make start                        # 启动服务"
    echo "make test                         # 运行测试"
    echo ""
}

# 主函数
main() {
    # 检查参数
    if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
        echo "用法: $0 [选项]"
        echo ""
        echo "选项:"
        echo "  --help, -h     显示此帮助信息"
        echo "  --init-only    只初始化，不启动服务"
        echo "  --skip-db      跳过数据库初始化"
        echo ""
        exit 0
    fi
    
    # 执行步骤
    check_requirements
    setup_project
    
    if [ "$1" != "--skip-db" ]; then
        init_database
    fi
    
    if [ "$1" = "--init-only" ]; then
        echo "🎉 初始化完成!"
        echo ""
        echo "现在可以运行以下命令启动服务器:"
        echo "  go run main.go"
        echo "  # 或者"
        echo "  make start"
        echo ""
        show_help
        exit 0
    fi
    
    # 显示帮助信息
    echo "💡 提示: 在另一个终端中运行 './scripts/test-api.sh' 来测试API"
    echo ""
    
    # 启动服务
    start_service
}

# 运行主函数
main "$@"
