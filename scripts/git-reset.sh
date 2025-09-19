#!/bin/bash
echo "🔧 Git 仓库完全重置脚本"
echo "========================================"
echo "⚠️  警告：这将重置您的 Git 历史！"
echo "========================================"
echo

read -p "确定要重置 Git 仓库吗？(y/N): " confirm
if [[ $confirm != [yY] ]]; then
    echo "操作已取消"
    exit 0
fi

echo "🗑️  删除 .git 目录..."
rm -rf .git

echo "🆕 重新初始化 Git 仓库..."
git init

echo "📝 检查 .gitignore..."
if [[ ! -f .gitignore ]]; then
    echo "⚠️  .gitignore 文件不存在！"
    exit 1
fi

echo "📂 添加所有项目文件..."
git add main.go
git add go.mod go.sum
git add config/
git add models/
git add controllers/
git add services/
git add routes/
git add utils/
git add cmd/
git add *.md
git add Dockerfile
git add docker-compose.yml
git add nginx.conf
git add production.env.example
git add Makefile
git add scripts/
git add *.sql
git add uploads/.gitkeep

echo "📊 检查将要添加的文件："
git status

echo
echo "✅ 文件添加完成！"
echo "现在运行以下命令完成设置："
echo
echo "git commit -m \"初始提交：完整项目文件\""
echo "git branch -M main"
echo "git remote add origin https://github.com/fishheadwithchili/imgGeneratePrompts.git"
echo "git push -u origin main"
