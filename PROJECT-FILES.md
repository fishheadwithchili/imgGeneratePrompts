# 📁 项目文件说明

## ✅ 应该存在的文件（会被 Git 跟踪）

### 🔧 配置文件
- `apikey/database.env.example` - 数据库配置示例
- `apikey/production.env.example` - 生产环境配置示例
- `apikey/README.md` - 配置文件说明
- `docker-compose.yml` - Docker 基础配置
- `docker-compose.override.yml.template` - Docker 覆盖配置模板
- `.gitignore` - Git 忽略配置

### 📜 脚本文件
- `scripts/init-project.sh` - 项目初始化脚本
- `scripts/quick-setup.sh` - 一键设置脚本
- `scripts/dev.sh` - 开发环境管理
- `scripts/migrate.sh` - 数据库迁移
- `scripts/git-*.sh` - Git 相关脚本

### 📚 源代码
- `main.go` - 应用入口
- `go.mod`, `go.sum` - Go 模块文件
- `config/`, `models/`, `controllers/`, `services/`, `routes/`, `utils/` - 源代码目录
- `cmd/` - 命令行工具

### 📖 文档
- `README.md` - 项目说明（包含首次设置指南）
- `MIGRATION.md` - 迁移使用说明
- `SECURITY.md` - 安全指南
- `GIT-BASH-GUIDE.md` - Git Bash 使用指南
- 其他 `.md` 文档文件

## ❌ 不应该存在的文件（被 Git 忽略）

### 🔒 敏感配置
- `apikey/database.env` - 包含真实数据库密码
- `docker-compose.override.yml` - 包含开发环境密码
- 任何 `.env` 文件

### 🔧 IDE 和系统文件
- `.idea/` - IntelliJ IDEA 配置
- `.vscode/` - VS Code 配置
- `.DS_Store` - macOS 系统文件

### 📝 运行时文件
- `logs/` - 应用日志（但 `logs/.gitkeep` 会被跟踪）
- `uploads/` - 上传的文件（但 `uploads/.gitkeep` 会被跟踪）
- `bin/` - 编译后的二进制文件

## 🚀 首次克隆项目后的操作

### 方法 1：一键设置（推荐）
```bash
chmod +x scripts/*.sh
./scripts/init-project.sh
```

### 方法 2：手动设置
```bash
# 1. 创建数据库配置
cp apikey/database.env.example apikey/database.env
# 编辑 apikey/database.env 设置密码

# 2. 创建 Docker 配置
cp docker-compose.override.yml.template docker-compose.override.yml
# 如需要，编辑 docker-compose.override.yml

# 3. 下载依赖
go mod download

# 4. 启动数据库
docker-compose up mysql -d

# 5. 执行迁移
./scripts/dev.sh migrate

# 6. 启动服务
./scripts/dev.sh start
```

## 🔍 检查项目状态

```bash
# 检查必需文件是否存在
ls -la apikey/database.env          # 应该存在
ls -la docker-compose.override.yml  # 可选，但推荐存在

# 检查 Git 状态
git status  # 应该显示 "working tree clean"

# 检查服务是否运行
curl http://localhost:8080/health   # 应该返回健康状态
```

## ⚠️ 注意事项

1. **首次运行**：必须先创建 `apikey/database.env` 文件
2. **安全性**：永远不要将包含真实密码的文件提交到 Git
3. **团队协作**：团队成员需要各自创建自己的配置文件
4. **生产部署**：使用环境变量而不是配置文件来管理敏感信息

## 🆘 故障排除

### 错误：无法打开数据库配置文件
```bash
# 解决方案：创建配置文件
cp apikey/database.env.example apikey/database.env
```

### 错误：Docker 启动失败
```bash
# 解决方案：创建 Docker 配置
cp docker-compose.override.yml.template docker-compose.override.yml
```

### 错误：数据库连接失败
```bash
# 解决方案：检查数据库配置和密码
cat apikey/database.env
docker-compose ps mysql
```
