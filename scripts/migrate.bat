@echo off
echo ==================================
echo     执行数据库迁移
echo ==================================
echo.

go run main.go -migrate

echo.
echo 迁移完成！现在可以使用 dev.bat 启动服务
pause
