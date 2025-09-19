#!/bin/bash
echo "===================================="
echo "     Git 状态检查"
echo "===================================="
echo

echo "1. 当前 Git 状态："
echo "------------------------------------"
git status

echo
echo "2. 已跟踪的文件："
echo "------------------------------------"
git ls-files

echo
echo "3. 未跟踪的文件（应该添加的）："
echo "------------------------------------"
git ls-files --others --exclude-standard

echo
echo "4. 被忽略的文件："
echo "------------------------------------"
git ls-files --others --ignored --exclude-standard

echo
echo "===================================="
echo "如果看到很多未跟踪的文件，请运行："
echo "./scripts/git-add-all.sh"
echo "===================================="
