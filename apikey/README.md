# 🔑 API Keys & 配置文件目录

此目录包含应用程序的敏感配置文件。

## 📁 文件说明

### ✅ 已提供的示例文件
- `database.env.example` - 数据库配置示例
- `production.env.example` - 生产环境配置示例

### ❌ 需要您创建的文件（被Git忽略）
- `database.env` - 实际的数据库配置（包含真实密码）
- `production.env` - 生产环境配置（包含真实密码）

## 🚀 快速设置

### 开发环境
```bash
# 复制示例文件
cp apikey/database.env.example apikey/database.env

# 编辑配置文件，设置您的数据库密码
nano apikey/database.env  # 或使用您喜欢的编辑器
```

### 生产环境
```bash
# 复制示例文件
cp apikey/production.env.example apikey/production.env

# 编辑配置文件，设置生产环境的安全配置
nano apikey/production.env
```

## ⚠️ 安全提醒

1. **永远不要提交真实的配置文件**到Git仓库
2. **使用强密码**用于生产环境
3. **定期更换密码**
4. **不要在日志中输出敏感信息**

## 🔧 配置示例

### 开发环境配置
```env
DB_HOST=localhost
DB_PORT=3307
DB_USER=root
DB_PASSWORD=your_dev_password
DB_NAME=img_generate_prompts
```

### 生产环境配置
```env
DB_HOST=your-production-host
DB_PORT=3306
DB_USER=your_prod_user
DB_PASSWORD=your_very_secure_password
DB_NAME=img_generate_prompts_prod
```

## 🆘 故障排除

**错误：无法打开数据库配置文件**
```bash
# 确保文件存在
ls -la apikey/database.env

# 如果不存在，从示例文件复制
cp apikey/database.env.example apikey/database.env
```

**错误：权限被拒绝**
```bash
# 检查文件权限
chmod 600 apikey/database.env
```
