@echo off
:: 一键开始脚本 - Image Generate Prompts API (Windows)
:: 用法: quick-start.bat

setlocal enabledelayedexpansion

echo 🚀 Image Generate Prompts API - 一键开始
echo ========================================
echo.

:: 检查参数
if "%1"=="--help" goto help
if "%1"=="-h" goto help
if "%1"=="/?" goto help

:: 检查系统要求
echo 🔍 检查系统要求...

:: 检查 Go
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go 未安装。请先安装 Go 1.21 或更高版本。
    echo    下载地址: https://golang.org/dl/
    pause
    exit /b 1
)
for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
echo ✅ Go 版本: %GO_VERSION%

:: 检查 MySQL 客户端
mysql --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  MySQL 客户端未安装。请确保 MySQL 服务器运行在 3307 端口。
) else (
    echo ✅ MySQL 客户端已安装
)

:: 检查 curl
curl --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  curl 未安装，部分测试功能可能无法使用。
) else (
    echo ✅ curl 已安装
)
echo.

:: 设置项目
echo 📦 设置项目...

echo   - 下载 Go 依赖...
go mod download
go mod tidy

echo   - 创建目录...
if not exist uploads mkdir uploads
if not exist logs mkdir logs
if not exist bin mkdir bin

echo ✅ 项目设置完成
echo.

:: 初始化数据库
if "%1"=="--skip-db" goto skip_db

echo 🗄️  初始化数据库...

:: 检查数据库配置文件
if not exist "apikey\database.env" (
    echo ❌ 数据库配置文件不存在: apikey\database.env
    echo    请确保文件存在并包含正确的数据库配置。
    pause
    exit /b 1
)

echo   - 创建表结构和示例数据...
go run cmd/db-manager.go -write
if %errorlevel% neq 0 (
    echo ❌ 数据库初始化失败
    echo    请检查:
    echo    1. MySQL 服务是否在 3307 端口运行
    echo    2. 数据库配置是否正确
    echo    3. 数据库 'img_generate_prompts' 是否已创建
    pause
    exit /b 1
)
echo ✅ 数据库初始化成功
echo.

:skip_db

:: 检查是否只初始化
if "%1"=="--init-only" goto init_only

:: 启动服务
echo 🌐 启动开发服务器...
echo.
echo 📊 数据库统计信息:
go run cmd/db-manager.go -stats
echo.
echo 🚀 服务器启动中...
echo    访问地址: http://localhost:8080
echo    API文档: http://localhost:8080/api/v1
echo    健康检查: http://localhost:8080/health
echo.
echo 按 Ctrl+C 停止服务器
echo ========================================
echo.
echo 💡 提示: 在另一个终端中运行 'scripts\test-api.bat' 来测试API
echo.

:: 启动服务器
go run main.go
goto end

:init_only
echo 🎉 初始化完成!
echo.
echo 现在可以运行以下命令启动服务器:
echo   go run main.go
echo   # 或者
echo   scripts\dev.bat start
echo.
goto show_help

:help
echo 用法: %0 [选项]
echo.
echo 选项:
echo   --help, -h, /?    显示此帮助信息
echo   --init-only       只初始化，不启动服务
echo   --skip-db         跳过数据库初始化
echo.
goto end

:show_help
echo 🛠️  开发命令:
echo.
echo # 查看项目信息
echo curl http://localhost:8080/
echo.
echo # 健康检查
echo curl http://localhost:8080/health
echo.
echo # 获取标签列表
echo curl http://localhost:8080/api/v1/tags
echo.
echo # 获取提示词列表
echo curl http://localhost:8080/api/v1/prompts
echo.
echo # 运行完整API测试
echo scripts\test-api.bat
echo.
echo # 数据库管理
echo go run cmd/db-manager.go -stats    :: 查看统计
echo go run cmd/db-manager.go -sample   :: 添加示例数据
echo go run cmd/db-manager.go -reset    :: 重置数据库
echo.
echo # 使用开发脚本
echo scripts\dev.bat help               :: 查看所有命令
echo scripts\dev.bat start              :: 启动服务
echo scripts\dev.bat test               :: 运行测试
echo.

:end
echo.
pause
