#!/bin/bash
# Linux/macOS 开发脚本
# 用法: ./scripts/dev.sh [命令]

set -e  # 遇到错误时退出

show_help() {
    echo ""
    echo "🛠️  开发脚本 - Image Generate Prompts"
    echo ""
    echo "用法: ./scripts/dev.sh [命令]"
    echo ""
    echo "命令:"
    echo "  init     - 初始化项目（数据库 + 依赖）"
    echo "  start    - 启动开发服务器"
    echo "  migrate  - 执行数据库迁移"
    echo "  clean    - 清理依赖和缓存"
    echo "  reset    - 重置数据库"
    echo "  test     - 运行测试"
    echo "  build    - 构建项目"
    echo "  help     - 显示此帮助信息"
    echo ""
}

init_project() {
    echo "🚀 初始化项目..."
    echo ""
    echo "1. 下载依赖..."
    go mod download
    go mod tidy
    echo ""
    echo "2. 初始化数据库..."
    go run cmd/db-manager.go -write
    echo ""
    echo "✅ 项目初始化完成！"
    echo ""
    echo "现在可以运行：./scripts/dev.sh start"
}

start_server() {
    echo "🚀 启动开发服务器..."
    echo ""
    echo "📊 数据库状态："
    go run cmd/db-manager.go -stats
    echo ""
    echo "💡 提示：如果是首次运行或修改了数据表结构，请先运行："
    echo "    ./scripts/dev.sh migrate"
    echo ""
    echo "🌐 启动服务器在 http://localhost:8080"
    go run main.go
}

migrate_database() {
    echo "🔄 执行数据库迁移..."
    echo ""
    go run main.go -migrate
    echo ""
    echo "✅ 迁移完成！现在可以运行： ./scripts/dev.sh start"
}

clean_project() {
    echo "🧹 清理项目..."
    [ -f go.sum ] && rm go.sum
    go clean -cache
    go mod download
    go mod tidy
    echo "✅ 清理完成！"
}

reset_database() {
    echo "⚠️  重置数据库"
    read -p "确定要重置数据库吗？这将删除所有数据！(y/N): " confirm
    if [[ $confirm == [yY] ]]; then
        go run cmd/db-manager.go -reset
        echo "✅ 数据库重置完成！"
    else
        echo "❌ 操作已取消"
    fi
}

run_tests() {
    echo "🧪 运行测试..."
    echo ""
    echo "1. 检查服务是否启动..."
    if ! curl -s http://localhost:8080/health > /dev/null; then
        echo "请先启动服务: ./scripts/dev.sh start"
        return 1
    fi
    echo ""
    echo "2. 运行API测试..."
    echo "测试健康检查..."
    curl -s http://localhost:8080/health | jq .
    echo ""
    echo "测试标签接口..."
    curl -s http://localhost:8080/api/v1/tags | jq .
    echo ""
    echo "测试提示词接口..."
    curl -s "http://localhost:8080/api/v1/prompts?page=1&page_size=5" | jq .
    echo ""
    echo "✅ 测试完成！"
}

build_project() {
    echo "🔨 构建项目..."
    mkdir -p bin
    go build -o bin/imgGeneratePrompts main.go
    echo "✅ 构建完成！可执行文件: bin/imgGeneratePrompts"
}

# 主逻辑
case "${1:-help}" in
    "init")
        init_project
        ;;
    "start")
        start_server
        ;;
    "migrate")
        migrate_database
        ;;
    "clean")
        clean_project
        ;;
    "reset")
        reset_database
        ;;
    "test")
        run_tests
        ;;
    "build")
        build_project
        ;;
    "help"|*)
        show_help
        ;;
esac

echo ""
