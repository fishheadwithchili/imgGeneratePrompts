# 部署指南

本文档详细介绍了如何在不同环境中部署 Image Generate Prompts API。

## 🚀 快速开始

### 本地开发环境

```bash
# 1. 克隆项目
git clone <repository-url>
cd imgGeneratePrompts

# 2. 使用开发脚本快速开始
make quick
# 或者
./scripts/dev.sh init && ./scripts/dev.sh start
# Windows: scripts\dev.bat init && scripts\dev.bat start
```

## 🐳 Docker 部署

### 使用 Docker Compose（推荐）

```bash
# 1. 启动所有服务
docker-compose up -d

# 2. 查看服务状态
docker-compose ps

# 3. 查看日志
docker-compose logs -f api

# 4. 初始化数据库（第一次启动后）
docker-compose exec api ./main -db-init

# 5. 停止服务
docker-compose down
```

### 单独使用 Docker

```bash
# 1. 构建镜像
docker build -t img-prompts-api .

# 2. 启动 MySQL（可选，如果已有数据库）
docker run -d \
  --name img-prompts-mysql \
  -e MYSQL_ROOT_PASSWORD=12345678 \
  -e MYSQL_DATABASE=img_generate_prompts \
  -p 3307:3306 \
  mysql:8.0

# 3. 启动应用
docker run -d \
  --name img-prompts-api \
  --link img-prompts-mysql:mysql \
  -p 8080:8080 \
  -v $(pwd)/uploads:/root/uploads \
  -v $(pwd)/apikey:/root/apikey:ro \
  img-prompts-api
```

## 🖥️ 生产环境部署

### 1. 服务器准备

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y git curl wget mysql-server nginx

# CentOS/RHEL
sudo yum update
sudo yum install -y git curl wget mysql-server nginx
```

### 2. MySQL 配置

```bash
# 启动 MySQL
sudo systemctl start mysql
sudo systemctl enable mysql

# 安全配置
sudo mysql_secure_installation

# 创建数据库和用户
sudo mysql -u root -p
```

```sql
CREATE DATABASE img_generate_prompts CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'img_prompts_user'@'localhost' IDENTIFIED BY 'your_secure_password';
GRANT ALL PRIVILEGES ON img_generate_prompts.* TO 'img_prompts_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 3. 应用部署

```bash
# 1. 克隆代码
git clone <repository-url> /opt/img-prompts
cd /opt/img-prompts

# 2. 构建应用
make build

# 3. 复制配置文件
cp production.env.example apikey/production.env
# 编辑 apikey/production.env 文件

# 4. 创建系统用户
sudo useradd -r -s /bin/false imgprompts

# 5. 设置权限
sudo chown -R imgprompts:imgprompts /opt/img-prompts
sudo chmod +x /opt/img-prompts/bin/imgGeneratePrompts

# 6. 创建上传目录
sudo mkdir -p /var/uploads/img-prompts
sudo chown imgprompts:imgprompts /var/uploads/img-prompts

# 7. 初始化数据库
sudo -u imgprompts ./bin/imgGeneratePrompts -db-init
```

### 4. Systemd 服务配置

创建 `/etc/systemd/system/img-prompts.service`：

```ini
[Unit]
Description=Image Generate Prompts API
After=network.target mysql.service
Requires=mysql.service

[Service]
Type=simple
User=imgprompts
Group=imgprompts
WorkingDirectory=/opt/img-prompts
ExecStart=/opt/img-prompts/bin/imgGeneratePrompts
Restart=always
RestartSec=5
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=img-prompts

# 环境变量
Environment=GIN_MODE=release

# 安全配置
NoNewPrivileges=true
PrivateTmp=true
PrivateDevices=true
ProtectHome=true
ProtectSystem=strict
ReadWritePaths=/var/uploads/img-prompts /var/log

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable img-prompts
sudo systemctl start img-prompts
sudo systemctl status img-prompts
```

### 5. Nginx 配置

创建 `/etc/nginx/sites-available/img-prompts`：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    client_max_body_size 10M;

    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /uploads/ {
        proxy_pass http://127.0.0.1:8080;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    location /health {
        proxy_pass http://127.0.0.1:8080;
    }

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

启用站点：

```bash
sudo ln -s /etc/nginx/sites-available/img-prompts /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 6. SSL 配置（使用 Let's Encrypt）

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo crontab -e
# 添加: 0 12 * * * /usr/bin/certbot renew --quiet
```

## 📊 监控和日志

### 1. 日志配置

```bash
# 创建日志目录
sudo mkdir -p /var/log/img-prompts
sudo chown imgprompts:imgprompts /var/log/img-prompts

# 配置 logrotate
sudo tee /etc/logrotate.d/img-prompts << EOF
/var/log/img-prompts/*.log {
    daily
    missingok
    rotate 52
    compress
    delaycompress
    notifempty
    create 644 imgprompts imgprompts
    postrotate
        systemctl reload img-prompts
    endscript
}
EOF
```

### 2. 监控脚本

创建 `/opt/img-prompts/scripts/monitor.sh`：

```bash
#!/bin/bash
# 健康检查脚本

URL="http://localhost:8080/health"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" $URL)

if [ $STATUS -eq 200 ]; then
    echo "$(date): Service is healthy"
    exit 0
else
    echo "$(date): Service is unhealthy (HTTP $STATUS)"
    # 可以在这里添加重启逻辑或发送告警
    exit 1
fi
```

添加到 crontab：

```bash
# 每分钟检查一次
* * * * * /opt/img-prompts/scripts/monitor.sh >> /var/log/img-prompts/monitor.log 2>&1
```

## 🔒 安全配置

### 1. 防火墙设置

```bash
# UFW (Ubuntu)
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# 或者 iptables
sudo iptables -A INPUT -p tcp --dport 22 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 80 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 443 -j ACCEPT
sudo iptables -A INPUT -j DROP
```

### 2. 数据库安全

```bash
# MySQL 配置优化
sudo vim /etc/mysql/mysql.conf.d/mysqld.cnf
```

```ini
[mysqld]
bind-address = 127.0.0.1
max_connections = 100
innodb_buffer_pool_size = 256M
```

### 3. 应用安全

- 定期更新依赖：`go mod update`
- 使用强密码
- 定期备份数据库
- 限制文件上传大小
- 实现 API 限流

## 📈 性能优化

### 1. 数据库优化

```sql
-- 添加索引
CREATE INDEX idx_prompts_public_created ON prompts(is_public, created_at);
CREATE INDEX idx_prompts_model_public ON prompts(model_name, is_public);

-- 配置优化
SET GLOBAL innodb_buffer_pool_size = 256*1024*1024;
```

### 2. 应用优化

- 使用连接池
- 实现缓存（Redis）
- 图片 CDN
- Gzip 压缩

### 3. 服务器优化

```bash
# 增加文件描述符限制
echo "* soft nofile 65536" >> /etc/security/limits.conf
echo "* hard nofile 65536" >> /etc/security/limits.conf
```

## 🔄 备份和恢复

### 1. 数据库备份

```bash
#!/bin/bash
# 创建备份脚本 /opt/img-prompts/scripts/backup.sh

BACKUP_DIR="/opt/backups/img-prompts"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="img_generate_prompts"

mkdir -p $BACKUP_DIR

# 数据库备份
mysqldump -u img_prompts_user -p$DB_PASSWORD $DB_NAME > $BACKUP_DIR/db_$DATE.sql

# 文件备份
tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz /var/uploads/img-prompts

# 删除30天前的备份
find $BACKUP_DIR -name "*.sql" -mtime +30 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +30 -delete

echo "Backup completed: $DATE"
```

### 2. 自动备份

```bash
# 添加到 crontab
0 2 * * * /opt/img-prompts/scripts/backup.sh >> /var/log/img-prompts/backup.log 2>&1
```

## 🚨 故障排除

### 常见问题

1. **服务无法启动**
   ```bash
   sudo journalctl -u img-prompts -f
   ```

2. **数据库连接失败**
   ```bash
   mysql -u img_prompts_user -p
   ```

3. **文件上传失败**
   ```bash
   ls -la /var/uploads/img-prompts
   sudo chown imgprompts:imgprompts /var/uploads/img-prompts
   ```

4. **Nginx 502 错误**
   ```bash
   sudo nginx -t
   sudo systemctl status img-prompts
   ```

### 维护命令

```bash
# 重启服务
sudo systemctl restart img-prompts

# 重新加载配置
sudo systemctl reload img-prompts

# 查看日志
sudo journalctl -u img-prompts -f

# 查看系统状态
sudo systemctl status img-prompts mysql nginx
```

---

有任何部署问题，请查看日志文件或联系技术支持。
