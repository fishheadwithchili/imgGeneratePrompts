@echo off
:: Windows 开发脚本
:: 用法: scripts\dev.bat [命令]

if "%1"=="" goto help
if "%1"=="help" goto help
if "%1"=="init" goto init
if "%1"=="start" goto start
if "%1"=="migrate" goto migrate
if "%1"=="clean" goto clean
if "%1"=="reset" goto reset
if "%1"=="test" goto test
if "%1"=="build" goto build
goto help

:help
echo.
echo 🛠️  开发脚本 - Image Generate Prompts
echo.
echo 用法: scripts\dev.bat [命令]
echo.
echo 命令:
echo   init     - 初始化项目（数据库 + 依赖）
echo   start    - 启动开发服务器
echo   migrate  - 执行数据库迁移
echo   clean    - 清理依赖和缓存
echo   reset    - 重置数据库
echo   test     - 运行测试
echo   build    - 构建项目
echo   help     - 显示此帮助信息
echo.
goto end

:init
echo 🚀 初始化项目...
echo.
echo 1. 下载依赖...
go mod download
go mod tidy
echo.
echo 2. 初始化数据库...
go run cmd/db-manager.go -write
echo.
echo ✅ 项目初始化完成！
echo.
echo 现在可以运行：scripts\dev.bat start
goto end

:start
echo 🚀 启动开发服务器...
echo.
echo 📊 数据库状态：
go run cmd/db-manager.go -stats
echo.
echo 💡 提示：如果是首次运行或修改了数据表结构，请先运行：
echo    scripts\migrate.bat
echo.
echo 🌐 启动服务器在 http://localhost:8080
go run main.go
goto end

:migrate
echo 🔄 执行数据库迁移...
echo.
go run main.go -migrate
echo.
echo ✅ 迁移完成！现在可以运行： scripts\dev.bat start
goto end

:clean
echo 🧹 清理项目...
if exist go.sum del go.sum
go clean -cache
go mod download
go mod tidy
echo ✅ 清理完成！
goto end

:reset
echo ⚠️  重置数据库
set /p confirm="确定要重置数据库吗？这将删除所有数据！(y/N): "
if /i "%confirm%"=="y" (
    go run cmd/db-manager.go -reset
    echo ✅ 数据库重置完成！
) else (
    echo ❌ 操作已取消
)
goto end

:test
echo 🧪 运行测试...
echo.
echo 1. 检查服务是否启动...
curl -s http://localhost:8080/health || echo 请先启动服务: scripts\dev.bat start
echo.
echo 2. 运行API测试...
echo 测试健康检查...
curl -s http://localhost:8080/health
echo.
echo 测试标签接口...
curl -s http://localhost:8080/api/v1/tags
echo.
echo 测试提示词接口...
curl -s "http://localhost:8080/api/v1/prompts?page=1&page_size=5"
echo.
echo ✅ 测试完成！
goto end

:build
echo 🔨 构建项目...
if not exist bin mkdir bin
go build -o bin\imgGeneratePrompts.exe main.go
echo ✅ 构建完成！可执行文件: bin\imgGeneratePrompts.exe
goto end

:end
echo.
