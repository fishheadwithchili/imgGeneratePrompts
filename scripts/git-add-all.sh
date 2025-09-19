#!/bin/bash
echo "===================================="
echo "     添加项目文件到 Git"
echo "===================================="
echo

echo "正在添加源代码文件..."
git add *.go
git add go.mod go.sum

echo "正在添加配置文件..."
git add config/
git add models/
git add controllers/
git add services/
git add routes/
git add utils/
git add cmd/

echo "正在添加文档文件..."
git add *.md
git add *.sql

echo "正在添加Docker配置..."
git add Dockerfile
git add docker-compose.yml
git add nginx.conf

echo "正在添加脚本文件..."
git add scripts/
git add *.sh
git add Makefile

echo "正在添加示例配置..."
git add production.env.example

echo "正在添加上传目录结构..."
git add uploads/.gitkeep

echo
echo "===================================="
echo "检查将要提交的文件："
echo "===================================="
git status

echo
echo "===================================="
echo "如果上面的文件列表正确，请运行："
echo "git commit -m \"初始提交：添加项目源代码\""
echo "git push origin main"
echo "===================================="
