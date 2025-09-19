# Makefile for Image Generate Prompts API

.PHONY: help init start build clean test reset deps install run dev

# 默认目标
.DEFAULT_GOAL := help

# 变量定义
BINARY_NAME=imgGeneratePrompts
BUILD_DIR=bin
MAIN_FILE=main.go
DB_MANAGER=cmd/db-manager.go

# 帮助信息
help: ## 显示帮助信息
	@echo ""
	@echo "🛠️  Image Generate Prompts API - Makefile"
	@echo ""
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
	@echo ""

# 项目初始化
init: deps db-init ## 完整初始化项目（依赖+数据库）
	@echo "✅ 项目初始化完成！"

# 安装依赖
deps: ## 下载和整理依赖
	@echo "📦 下载依赖..."
	@go mod download
	@go mod tidy
	@echo "✅ 依赖安装完成！"

# 数据库初始化
db-init: ## 初始化数据库
	@echo "🗄️  初始化数据库..."
	@go run $(DB_MANAGER) -write
	@echo "✅ 数据库初始化完成！"

# 启动开发服务器
start: ## 启动开发服务器
dev: start
run: start
start:
	@echo "🚀 启动开发服务器..."
	@echo "📊 数据库状态："
	@go run $(DB_MANAGER) -stats
	@echo ""
	@echo "🌐 服务器启动在 http://localhost:8080"
	@go run $(MAIN_FILE)

# 构建项目
build: ## 构建可执行文件
	@echo "🔨 构建项目..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ 构建完成！文件位置: $(BUILD_DIR)/$(BINARY_NAME)"

# 构建多平台版本
build-all: ## 构建多平台版本
	@echo "🔨 构建多平台版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	@echo "✅ 多平台构建完成！"

# 清理项目
clean: ## 清理构建文件和缓存
	@echo "🧹 清理项目..."
	@rm -rf $(BUILD_DIR)
	@go clean -cache
	@go clean -modcache
	@echo "✅ 清理完成！"

# 运行测试
test: ## 运行API测试
	@echo "🧪 运行测试..."
	@./scripts/test-api.sh || (echo "❌ 请确保服务正在运行" && exit 1)

# 数据库管理
db-stats: ## 显示数据库统计信息
	@go run $(DB_MANAGER) -stats

db-reset: ## 重置数据库（危险操作）
	@go run $(DB_MANAGER) -reset

db-sample: ## 创建示例数据
	@go run $(DB_MANAGER) -sample

db-validate: ## 验证数据完整性
	@go run $(DB_MANAGER) -validate

# 代码质量
fmt: ## 格式化代码
	@echo "🎨 格式化代码..."
	@go fmt ./...
	@echo "✅ 代码格式化完成！"

vet: ## 代码静态检查
	@echo "🔍 代码静态检查..."
	@go vet ./...
	@echo "✅ 静态检查完成！"

# 开发工具
install: ## 安装到系统路径
	@echo "📦 安装到系统..."
	@go install .
	@echo "✅ 安装完成！"

# 版本信息
version: ## 显示版本信息
	@echo "📋 版本信息："
	@echo "Go version: $(shell go version)"
	@echo "Git commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "Build time: $(shell date)"

# 快速开发流程
quick: deps db-init start ## 快速开始（依赖+数据库+启动）

# 生产部署准备
deploy: clean deps build ## 准备生产部署
	@echo "🚀 准备生产部署..."
	@echo "✅ 部署文件已准备就绪！"

# 监控日志（如果在后台运行）
logs: ## 查看日志（如果有的话）
	@echo "📝 查看日志..."
	@tail -f logs/*.log 2>/dev/null || echo "暂无日志文件"

# 健康检查
health: ## 检查服务健康状态
	@echo "🏥 健康检查..."
	@curl -s http://localhost:8080/health | jq . || echo "服务未启动或不可访问"

# 全面检查
check: fmt vet health ## 全面检查（格式+静态检查+健康检查）
	@echo "✅ 全面检查完成！"
