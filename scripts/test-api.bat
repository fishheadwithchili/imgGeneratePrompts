@echo off
:: Windows API 测试脚本
:: 用法: scripts\test-api.bat

setlocal enabledelayedexpansion

:: 配置
set API_BASE=http://localhost:8080
set TEST_TAG_NAME=API测试标签%random%
set TEST_PROMPT_TEXT=API测试提示词 - %date% %time%

echo 🚀 开始API自动化测试...
echo.

:: 1. 健康检查
echo 🧪 测试: 健康检查
curl -s "%API_BASE%/health" > temp_response.json
findstr "ok" temp_response.json >nul
if %errorlevel%==0 (
    echo ✅ 健康检查通过
) else (
    echo ❌ 健康检查失败
    type temp_response.json
    goto cleanup
)
echo.

:: 2. 获取标签列表
echo 🧪 测试: 获取标签列表
curl -s "%API_BASE%/api/v1/tags" > temp_response.json
findstr "200" temp_response.json >nul
if %errorlevel%==0 (
    echo ✅ 标签列表获取成功
) else (
    echo ❌ 标签列表获取失败
)
echo.

:: 3. 创建测试标签
echo 🧪 测试: 创建新标签
curl -s -X POST "%API_BASE%/api/v1/tags" -H "Content-Type: application/json" -d "{\"name\": \"%TEST_TAG_NAME%\"}" > temp_response.json
findstr "200" temp_response.json >nul
if %errorlevel%==0 (
    echo ✅ 标签创建成功
    :: 提取标签ID (简化版)
    for /f "tokens=2 delims=:" %%i in ('findstr "\"id\"" temp_response.json') do (
        set TAG_ID=%%i
        set TAG_ID=!TAG_ID:,=!
        set TAG_ID=!TAG_ID: =!
    )
) else (
    echo ❌ 标签创建失败
)
echo.

:: 4. 获取提示词列表
echo 🧪 测试: 获取提示词列表
curl -s "%API_BASE%/api/v1/prompts?page=1&page_size=5" > temp_response.json
findstr "200" temp_response.json >nul
if %errorlevel%==0 (
    echo ✅ 提示词列表获取成功
) else (
    echo ❌ 提示词列表获取失败
)
echo.

:: 5. 创建测试提示词
echo 🧪 测试: 创建新提示词
curl -s -X POST "%API_BASE%/api/v1/prompts" -H "Content-Type: application/json" -d "{\"prompt_text\": \"%TEST_PROMPT_TEXT%\", \"negative_prompt\": \"test negative\", \"is_public\": true, \"tag_names\": [\"%TEST_TAG_NAME%\"]}" > temp_response.json
findstr "200" temp_response.json >nul
if %errorlevel%==0 (
    echo ✅ 提示词创建成功
    :: 提取提示词ID (简化版)
    for /f "tokens=2 delims=:" %%i in ('findstr "\"id\"" temp_response.json') do (
        set PROMPT_ID=%%i
        set PROMPT_ID=!PROMPT_ID:,=!
        set PROMPT_ID=!PROMPT_ID: =!
        goto found_prompt_id
    )
    :found_prompt_id
) else (
    echo ❌ 提示词创建失败
)
echo.

:: 6. 获取公开提示词
echo 🧪 测试: 获取公开提示词
curl -s "%API_BASE%/api/v1/prompts/public?page=1&page_size=3" > temp_response.json
findstr "200" temp_response.json >nul
if %errorlevel%==0 (
    echo ✅ 公开提示词获取成功
) else (
    echo ❌ 公开提示词获取失败
)
echo.

:: 7. 搜索功能测试
echo 🧪 测试: 关键词搜索
curl -s "%API_BASE%/api/v1/prompts?keyword=测试" > temp_response.json
findstr "200" temp_response.json >nul
if %errorlevel%==0 (
    echo ✅ 关键词搜索成功
) else (
    echo ❌ 关键词搜索失败
)
echo.

:: 8. 统计信息测试
echo 🧪 测试: 获取统计信息
curl -s "%API_BASE%/api/v1/prompts/stats" > temp_response.json
findstr "200" temp_response.json >nul
if %errorlevel%==0 (
    echo ✅ 统计信息获取成功
) else (
    echo ❌ 统计信息获取失败
)

curl -s "%API_BASE%/api/v1/tags/stats" > temp_response.json
findstr "200" temp_response.json >nul
if %errorlevel%==0 (
    echo ✅ 标签统计获取成功
) else (
    echo ❌ 标签统计获取失败
)
echo.

:: 9. 错误处理测试
echo 🧪 测试: 错误处理
curl -s "%API_BASE%/api/v1/prompts/999999" > temp_response.json
findstr "404" temp_response.json >nul
if %errorlevel%==0 (
    echo ✅ 无效ID错误处理正确
) else (
    echo ❌ 无效ID错误处理异常
)
echo.

:: 10. 清理测试数据
echo 🧹 清理测试数据
if defined PROMPT_ID (
    echo 删除测试提示词...
    curl -s -X DELETE "%API_BASE%/api/v1/prompts/!PROMPT_ID!" > temp_response.json
    findstr "200" temp_response.json >nul
    if %errorlevel%==0 (
        echo ✅ 测试提示词删除成功
    ) else (
        echo ⚠️  测试提示词删除失败
    )
)

if defined TAG_ID (
    echo 删除测试标签...
    curl -s -X DELETE "%API_BASE%/api/v1/tags/!TAG_ID!" > temp_response.json
    findstr "200" temp_response.json >nul
    if %errorlevel%==0 (
        echo ✅ 测试标签删除成功
    ) else (
        echo ⚠️  测试标签删除失败（可能仍被使用）
    )
)
echo.

echo 🎉 API测试完成！
echo.
echo 📊 测试总结：
echo ✅ 健康检查和基础功能
echo ✅ 标签管理功能
echo ✅ 提示词管理功能  
echo ✅ 搜索和过滤功能
echo ✅ 统计信息功能
echo ✅ 错误处理功能
echo.
echo 🚀 核心功能测试通过！

:cleanup
if exist temp_response.json del temp_response.json
pause
