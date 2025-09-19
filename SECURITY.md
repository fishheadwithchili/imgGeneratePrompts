# 🔒 项目安全指南

## 敏感信息保护

### ✅ 已保护的敏感信息

1. **数据库密码** - 存储在 `apikey/database.env`
   - ✅ 已被 `.gitignore` 忽略
   - ✅ 不会推送到远程仓库

2. **Docker 密码** - 使用环境变量
   - ✅ `docker-compose.yml` 使用环境变量形式
   - ✅ `docker-compose.override.yml` 包含实际密码但被忽略

3. **生产环境配置** - 使用示例文件
   - ✅ `production.env.example` 只包含占位符
   - ✅ 实际配置文件被 `.gitignore` 忽略

### 📂 .gitignore 配置

```gitignore
# API keys and sensitive configuration
apikey/
*.env
.env*

# Docker compose overrides with sensitive data
docker-compose.override.yml
```

## 🚀 部署前检查清单

### 在推送到 GitHub 前，请确认：

- [ ] 检查 `apikey/` 目录是否被忽略
- [ ] 确认没有 `.env` 文件被追踪
- [ ] 验证 `docker-compose.override.yml` 被忽略
- [ ] 检查代码中没有硬编码的密码
- [ ] 确认所有敏感配置都使用环境变量

### 快速检查命令：

```bash
# 检查哪些文件会被推送
git status

# 检查是否有敏感文件被追踪
git ls-files | grep -E "(\.env|apikey|password|secret)"

# 检查代码中是否有硬编码密码
grep -r "password.*=" --include="*.go" . || echo "✅ 没有发现硬编码密码"
```

## 🔧 本地开发环境设置

### 1. 克隆项目后的设置

```bash
# 1. 复制配置文件
cp production.env.example apikey/database.env

# 2. 修改数据库配置
# 编辑 apikey/database.env，设置您的数据库密码

# 3. 复制 Docker 配置
cp docker-compose.yml docker-compose.override.yml
# 编辑 docker-compose.override.yml，设置密码
```

### 2. 团队协作注意事项

- ✅ 切勿将 `apikey/` 目录添加到版本控制
- ✅ 使用 `production.env.example` 作为配置模板
- ✅ 在 README 中说明如何设置本地配置
- ✅ 定期检查 `.gitignore` 是否生效

## 🛡️ 生产环境安全

### 环境变量管理

生产环境建议使用以下方式管理敏感信息：

1. **Docker Secrets**（推荐）
2. **Kubernetes Secrets**
3. **环境变量注入**
4. **专用密钥管理系统**

### 示例：使用环境变量

```bash
# 设置环境变量
export DB_PASSWORD="your_secure_password"
export MYSQL_ROOT_PASSWORD="another_secure_password"

# 运行应用
./app
```

## 🚨 紧急情况处理

### 如果意外推送了敏感信息：

1. **立即更改密码**
2. **使用 git filter-branch 删除历史记录**
3. **强制推送清理后的历史**
4. **通知团队成员重新克隆**

```bash
# 删除敏感文件的历史记录（危险操作）
git filter-branch --force --index-filter \
  'git rm --cached --ignore-unmatch apikey/database.env' \
  --prune-empty --tag-name-filter cat -- --all

# 强制推送
git push origin --force --all
```

## 📋 安全检查清单

### 代码审查时检查：
- [ ] 没有硬编码的密码、API密钥
- [ ] 所有敏感配置都使用环境变量
- [ ] `.gitignore` 正确配置
- [ ] 没有在日志中输出敏感信息
- [ ] 数据库连接字符串不包含明文密码

### 部署时检查：
- [ ] 生产环境使用强密码
- [ ] 启用HTTPS
- [ ] 配置防火墙规则
- [ ] 定期更新依赖包
- [ ] 监控异常访问

---

**记住：安全是一个持续的过程，而不是一次性的任务！** 🔐
