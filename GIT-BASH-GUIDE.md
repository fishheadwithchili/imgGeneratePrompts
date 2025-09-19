# 🚀 GoLand Git Bash 快速操作指南

## 在 GoLand 的 Git Bash 中运行：

### 1. 📋 检查 Git 状态
```bash
# 检查哪些文件没有被添加到Git
./scripts/git-check.sh
```

### 2. 📂 添加所有项目文件
```bash
# 一键添加所有需要的文件
./scripts/git-add-all.sh
```

### 3. 💾 提交并推送
```bash
# 提交代码
git commit -m "初始提交：添加项目源代码和文档"

# 推送到远程仓库
git push origin main
```

### 4. 🔄 数据库迁移
```bash
# 执行数据库迁移（首次运行或修改数据表结构时）
./scripts/dev.sh migrate
# 或者
./scripts/migrate.sh
```

### 5. 🚀 启动服务
```bash
# 正常启动服务
./scripts/dev.sh start
```

## 📁 脚本说明

### Git 相关脚本：
- `./scripts/git-check.sh` - 检查Git状态，查看哪些文件没有被跟踪
- `./scripts/git-add-all.sh` - 批量添加所有项目文件到Git

### 开发相关脚本：
- `./scripts/dev.sh migrate` - 执行数据库迁移
- `./scripts/dev.sh start` - 启动开发服务器
- `./scripts/dev.sh help` - 查看所有可用命令

## 🔍 常用命令组合

### 首次设置项目：
```bash
# 1. 检查文件状态
./scripts/git-check.sh

# 2. 添加所有文件
./scripts/git-add-all.sh

# 3. 提交
git commit -m "初始提交"

# 4. 推送
git push origin main

# 5. 执行数据库迁移
./scripts/dev.sh migrate

# 6. 启动服务
./scripts/dev.sh start
```

### 日常开发：
```bash
# 直接启动服务（数据库已迁移的情况下）
./scripts/dev.sh start
```

### 修改数据表结构后：
```bash
# 执行迁移
./scripts/dev.sh migrate

# 启动服务
./scripts/dev.sh start
```

## ⚠️ 注意事项

1. **确保脚本有执行权限**：
   ```bash
   chmod +x scripts/*.sh
   ```

2. **在项目根目录执行**：
   ```bash
   # 确保在正确的目录
   pwd
   # 应该显示：/d/projects/GolandProjects/imgGeneratePrompts
   ```

3. **Git Bash 中的路径**：
   - 使用 `./scripts/` 而不是 `scripts\`
   - 使用正斜杠 `/` 而不是反斜杠 `\`

## 🎯 如果脚本执行出错

### 权限问题：
```bash
# 给脚本添加执行权限
chmod +x scripts/*.sh
```

### 路径问题：
```bash
# 确保在项目根目录
cd /d/projects/GolandProjects/imgGeneratePrompts

# 查看当前目录内容
ls -la
```

### Git 问题：
```bash
# 检查Git状态
git status

# 手动添加文件
git add main.go go.mod go.sum
git add config/ models/ controllers/ services/ routes/ utils/
git add *.md Dockerfile docker-compose.yml
```

## 🔧 GoLand 中的 Git Bash 设置

1. **打开终端**：`Alt + F12`
2. **选择 Git Bash**：点击终端右上角的下拉箭头 → 选择 "Git Bash"
3. **设为默认**：右键点击终端标签 → "Set as Default"

现在您可以愉快地使用 bash 脚本了！🎉
