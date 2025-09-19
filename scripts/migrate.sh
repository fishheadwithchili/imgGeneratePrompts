#!/bin/bash
echo "=================================="
echo "     执行数据库迁移"
echo "=================================="
echo

go run main.go -migrate

echo
echo "迁移完成！现在可以使用 ./scripts/dev.sh start 启动服务"
