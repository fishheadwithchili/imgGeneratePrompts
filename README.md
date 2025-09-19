# Image Generate Prompts API v2.0

基于Go语言的AI图片生成提示词管理系统，采用优雅的MVC架构和多对多标签系统。

## 🆕 v2.0 更新内容

- ✅ **简化数据结构**：移除用户系统，专注核心功能
- ✅ **多对多标签系统**：独立的标签表和中间表
- ✅ **丰富描述字段**：新增风格、场景、氛围、意图等描述
- ✅ **数据库管理工具**：命令行工具管理数据库
- ✅ **完善的API接口**：标签管理、搜索、统计等功能

## 📁 项目结构

```
imgGeneratePrompts/
├── 📂 apikey/                 # 敏感配置（已gitignore）
│   └── database.env          # 数据库配置文件
├── 📂 cmd/                   # 命令行工具
│   └── db-manager.go         # 数据库管理工具
├── 📂 config/                # 配置管理
│   ├── config.go            # 应用配置
│   └── database.go          # 数据库连接和迁移
├── 📂 controllers/           # 控制器层（MVC-C）
│   ├── prompt_controller.go  # 提示词控制器
│   └── tag_controller.go     # 标签控制器
├── 📂 models/               # 模型层（MVC-M）
│   └── prompt.go            # 数据模型定义
├── 📂 services/             # 服务层（业务逻辑）
│   ├── prompt_service.go    # 提示词服务
│   └── tag_service.go       # 标签服务
├── 📂 routes/               # 路由配置
│   └── routes.go
├── 📂 utils/                # 工具类
│   ├── database_manager.go  # 数据库管理工具
│   ├── file_utils.go        # 文件处理工具
│   └── response.go          # 响应格式化工具
├── 📂 uploads/              # 文件上传目录
├── 📄 main.go               # 应用入口
├── 📄 go.mod                # Go模块配置
└── 📄 README.md             # 项目文档
```

## 🗄️ 数据库设计

### 1. prompts 表（主表）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 主键，自增 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |
| deleted_at | DATETIME | 软删除时间 |
| prompt_text | TEXT | 正面提示词 |
| negative_prompt | TEXT | 负面提示词 |
| model_name | VARCHAR(100) | AI模型名称 |
| image_url | VARCHAR(500) | 图片URL |
| is_public | TINYINT(1) | 是否公开 |
| **style_description** | VARCHAR(500) | **风格描述** |
| **usage_scenario** | VARCHAR(500) | **适用场景描述** |
| **atmosphere_description** | VARCHAR(500) | **氛围描述** |
| **expressive_intent** | VARCHAR(500) | **表现意图描述** |
| **structure_analysis** | JSON | **提示词结构分析** |

### 2. tags 表（标签表）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 主键，自增 |
| name | VARCHAR(100) | 标签名称（唯一） |
| created_at | DATETIME | 创建时间 |

### 3. prompt_tags 表（中间表）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| prompt_id | BIGINT | 提示词ID |
| tag_id | BIGINT | 标签ID |

## 🚀 快速开始

### 1. 环境准备

```bash
# 确保MySQL运行在3307端口
# 创建数据库
mysql -u root -p12345678 -P 3307 -e "CREATE DATABASE img_generate_prompts;"
```

### 2. 克隆和安装

```bash
cd D:\projects\GolandProjects\imgGeneratePrompts
go mod download
```

### 3. 数据库初始化 ⭐

使用我们提供的数据库管理工具：

```bash
# 方法1：完整初始化（推荐）
go run cmd/db-manager.go -write

# 方法2：分步骤初始化
go run cmd/db-manager.go -init      # 初始化表结构
go run cmd/db-manager.go -sample    # 创建示例数据

# 其他管理命令
go run cmd/db-manager.go -stats     # 查看统计信息
go run cmd/db-manager.go -validate  # 验证数据完整性
go run cmd/db-manager.go -reset     # 重置数据库（危险）
```

### 4. 启动服务

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动！

## 🔧 数据库管理工具

我们提供了强大的命令行数据库管理工具：

```bash
# 🛠️  数据库管理工具
# 
# 用法:
#   go run cmd/db-manager.go [选项]
# 
# 选项:
#   -write     完整写入数据库（推荐：初始化+示例数据）
#   -init      初始化数据库（创建表结构）
#   -sample    创建示例数据
#   -reset     重置数据库（危险操作）
#   -stats     显示数据库统计信息
#   -validate  验证数据完整性
# 
# 示例:
#   go run cmd/db-manager.go -write    # 完整初始化数据库
#   go run cmd/db-manager.go -stats    # 查看统计信息
#   go run cmd/db-manager.go -sample   # 只创建示例数据
```

## 📡 API 接口

### 基础信息

- **基础URL**: `http://localhost:8080/api/v1`
- **响应格式**: JSON

### 提示词接口

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/prompts` | 创建提示词 |
| POST | `/prompts/upload` | 上传图片并创建提示词 |
| GET | `/prompts` | 获取提示词列表（支持搜索和过滤） |
| GET | `/prompts/public` | 获取公开提示词列表 |
| GET | `/prompts/recent` | 获取最近的提示词 |
| GET | `/prompts/stats` | 获取提示词统计信息 |
| GET | `/prompts/search/tags` | 根据标签搜索提示词 |
| GET | `/prompts/check-duplicate` | 检查重复提示词 |
| GET | `/prompts/:id` | 获取单个提示词 |
| PUT | `/prompts/:id` | 更新提示词 |
| DELETE | `/prompts/:id` | 删除提示词 |

### 标签接口 🆕

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/tags` | 创建标签 |
| GET | `/tags` | 获取所有标签 |
| GET | `/tags/search` | 搜索标签 |
| GET | `/tags/stats` | 获取标签统计信息 |
| GET | `/tags/:id` | 获取单个标签 |
| DELETE | `/tags/:id` | 删除标签 |

### 创建提示词示例

```json
POST /api/v1/prompts
{
  "prompt_text": "a beautiful sunset over mountains, golden hour, cinematic lighting",
  "negative_prompt": "ugly, blurry, low quality",
  "model_name": "stable-diffusion-v1-5",
  "is_public": true,
  "style_description": "风景摄影风格，温暖的金色调",
  "usage_scenario": "适用于自然风光、旅游宣传、背景图片",
  "atmosphere_description": "宁静、温暖、壮观的黄昏氛围",
  "expressive_intent": "表现大自然的壮美和宁静",
  "structure_analysis": "{\"主体\":\"山峰日落\",\"光照\":\"黄金时刻\"}",
  "tag_names": ["风景", "暖色调", "高质量"]
}
```

### 查询参数示例

```bash
# 按标签搜索
GET /api/v1/prompts?tag_names=风景,高质量&page=1&page_size=10

# 关键词搜索
GET /api/v1/prompts?keyword=sunset&sort_by=created_at&sort_order=desc

# 获取公开提示词
GET /api/v1/prompts/public?page=1&page_size=20
```

## ✨ 新功能特性

### 🏷️ 多对多标签系统
- 独立的标签管理
- 标签统计和热门排行
- 支持标签搜索和过滤

### 📝 丰富的描述字段
- **风格描述**: 描述图片的艺术风格
- **适用场景**: 说明使用场景和用途
- **氛围描述**: 表达图片营造的氛围
- **表现意图**: 阐述创作意图
- **结构分析**: JSON格式的提示词结构分析

### 🔍 强大的搜索功能
- 关键词搜索（支持多字段）
- 标签过滤
- 模型名称过滤
- 公开/私有过滤

### 📊 统计功能
- 提示词统计（总数、公开、私有、最近）
- 标签使用统计
- 模型使用统计
- 热门标签排行

## 🛠️ 开发指南

### 数据库操作流程

1. **初始化**: `go run cmd/db-manager.go -write`
2. **开发**: 修改模型后重新迁移
3. **测试**: 使用示例数据测试功能
4. **部署**: 生产环境只运行 `-init`

### 添加新字段

1. 在 `models/prompt.go` 中添加字段
2. 运行 `go run cmd/db-manager.go -init` 迁移
3. 更新对应的服务和控制器

### 自定义配置

修改 `apikey/database.env` 文件：

```env
DB_HOST=localhost
DB_PORT=3307
DB_USER=root
DB_PASSWORD=12345678
DB_NAME=img_generate_prompts
DB_CHARSET=utf8mb4
```

## 🧪 测试示例

```bash
# 1. 初始化数据库
go run cmd/db-manager.go -write

# 2. 启动服务
go run main.go

# 3. 测试健康检查
curl http://localhost:8080/health

# 4. 获取所有标签
curl http://localhost:8080/api/v1/tags

# 5. 获取提示词列表
curl http://localhost:8080/api/v1/prompts

# 6. 根据标签搜索
curl "http://localhost:8080/api/v1/prompts?tag_names=风景,高质量"

# 7. 创建新提示词
curl -X POST http://localhost:8080/api/v1/prompts \
  -H "Content-Type: application/json" \
  -d '{"prompt_text":"test prompt","tag_names":["测试"],"is_public":true}'
```

## 📋 TODO 列表

- [ ] 用户认证系统
- [ ] 图片自动生成集成
- [ ] 提示词推荐算法
- [ ] 批量导入/导出功能
- [ ] API文档生成
- [ ] 单元测试

## 🔗 相关资源

- [API测试示例](EXAMPLES.md)
- [数据库结构参考](database_schema.sql)
- [Go官方文档](https://golang.org/doc/)
- [Gin框架文档](https://gin-gonic.com/)
- [GORM文档](https://gorm.io/)

## 📄 许可证

MIT License

---

**开始使用**: `go run cmd/db-manager.go -write && go run main.go` 🚀
