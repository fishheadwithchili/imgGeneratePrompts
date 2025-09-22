# 图像生成提示词管理系统后端 V3.1

这是一个用于管理AI图像生成提示词的后端API系统，基于Golang、Gin框架和GORM ORM。

## 🆕 V3.1 更新

### 1. 数据库字段优化
- **输入图片字段**：`input_image_url` 替代原有的 `reference_images`，存储格式优化为逗号分隔的字符串
- **输出图片字段**：`output_image_url` 替代原有的 `output_image`
- **存储格式**：多个图片URL以逗号分隔存储，如：`/uploads/180151.jpg,/uploads/180150.jpg`
- **向后兼容**：API响应仍返回数组格式，确保前端兼容性

### 2. 字段映射关系
- 数据库：`input_image_url`（varchar 500，逗号分隔）→ API：`input_image_urls`（数组）
- 数据库：`output_image_url`（varchar 500）→ API：`output_image_url`（字符串）

## 主要特性

- 🎨 **提示词管理**：完整的CRUD操作，支持创建、查看、更新和删除提示词
- 🏷️ **标签系统**：灵活的标签管理，支持多对多关联
- 📸 **多图片上传**：支持参考图和输出图的分别上传
- 🤖 **AI智能生成**：基于图片和基础提示词自动生成完整描述
- 🔍 **高级搜索**：支持关键词、模型、标签等多维度搜索
- 📊 **统计分析**：提供提示词和标签的统计信息
- ✅ **重复检测**：自动检测重复的提示词内容

## 技术栈

- **语言**: Go 1.19+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 5.7+
- **API格式**: RESTful JSON

## 项目结构

```
imgGeneratePrompts/
├── config/               # 配置管理
│   ├── config.go        # 应用配置
│   └── database.go      # 数据库配置
├── controllers/         # 控制器层
│   ├── prompt_controller.go  # 提示词控制器
│   └── tag_controller.go     # 标签控制器
├── models/              # 数据模型
│   └── prompt.go        # 提示词和标签模型
├── services/            # 业务逻辑层
│   ├── prompt_service.go     # 提示词服务
│   └── tag_service.go        # 标签服务
├── routes/              # 路由定义
│   └── routes.go        # API路由配置
├── utils/               # 工具函数
│   ├── file_utils.go    # 文件处理
│   └── response.go      # 响应格式化
├── scripts/             # 脚本文件
│   ├── init.sql         # 数据库初始化
│   └── migrate_v3.sql   # V3.0数据库迁移
└── uploads/             # 图片上传目录
```

## 数据库设计

### prompts表（V3.1更新）
```sql
- id                     # 主键
- created_at            # 创建时间
- updated_at            # 更新时间
- deleted_at            # 软删除时间
- prompt_text           # 正面提示词
- negative_prompt       # 负面提示词
- model_name           # 模型名称
- input_image_url      # 输入参照图片URL（逗号分隔）
- output_image_url     # 输出参照图片URL
- is_public            # 是否公开
- style_description    # 风格描述
- usage_scenario       # 使用场景
- atmosphere_description # 氛围描述
- expressive_intent    # 表现意图
- structure_analysis   # 结构分析（JSON）
```

### tags表
```sql
- id         # 主键
- name       # 标签名称
- created_at # 创建时间
```

### prompt_tags表（关联表）
```sql
- prompt_id  # 提示词ID
- tag_id     # 标签ID
```

## API接口详情

### 提示词接口

#### 创建提示词
```http
POST /api/v1/prompts/
Content-Type: application/json

{
  "prompt_text": "提示词内容",
  "negative_prompt": "负面提示词",
  "model_name": "模型名称",
  "input_image_urls": ["/uploads/image1.jpg", "/uploads/image2.jpg"],
  "output_image_url": "/uploads/output.jpg",
  "is_public": true,
  "style_description": "风格描述",
  "usage_scenario": "使用场景",
  "atmosphere_description": "氛围描述",
  "expressive_intent": "表现意图",
  "structure_analysis": "{\"主体\":\"描述\"}",
  "tag_names": ["标签1", "标签2"]
}
```

#### 上传图片并创建提示词
```http
POST /api/v1/prompts/upload
Content-Type: multipart/form-data

- input_images: 输入参考图片文件（支持多个）
- output_image: 输出图片文件（单个）
- prompt_text: 提示词内容
- negative_prompt: 负面提示词
- model_name: 模型名称
- is_public: 是否公开
- style_description: 风格描述
- usage_scenario: 使用场景
- atmosphere_description: 氛围描述
- expressive_intent: 表现意图
- structure_analysis: 结构分析
- tag_names: 标签名称（逗号分隔）
```

#### AI智能分析
```http
POST /api/v1/prompts/analyze
Content-Type: multipart/form-data

- output_image: 输出图片文件（必需）
- input_images: 输入参考图片文件（可选，支持多个）
- prompt_text: 基础提示词
- model_name: 模型名称
```

#### 响应格式示例
```json
{
  "success": true,
  "message": "操作成功",
  "data": {
    "id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "prompt_text": "提示词内容",
    "negative_prompt": "负面提示词",
    "model_name": "模型名称",
    "input_image_urls": ["/uploads/image1.jpg", "/uploads/image2.jpg"],
    "output_image_url": "/uploads/output.jpg",
    "is_public": true,
    "style_description": "风格描述",
    "usage_scenario": "使用场景",
    "atmosphere_description": "氛围描述",
    "expressive_intent": "表现意图",
    "structure_analysis": "{\"主体\":\"描述\"}",
    "tags": [
      {"id": 1, "name": "标签1", "created_at": "2024-01-01T00:00:00Z"},
      {"id": 2, "name": "标签2", "created_at": "2024-01-01T00:00:00Z"}
    ]
  }
}
```

### 完整API接口列表

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/prompts/ | 创建提示词 |
| POST | /api/v1/prompts/upload | 上传图片并创建提示词 |
| POST | /api/v1/prompts/analyze | AI智能分析 |
| GET | /api/v1/prompts/ | 获取提示词列表 |
| GET | /api/v1/prompts/:id | 获取单个提示词 |
| PUT | /api/v1/prompts/:id | 更新提示词 |
| DELETE | /api/v1/prompts/:id | 删除提示词 |
| GET | /api/v1/prompts/public | 获取公开提示词 |
| GET | /api/v1/prompts/recent | 获取最近提示词 |
| GET | /api/v1/prompts/stats | 获取统计信息 |
| GET | /api/v1/prompts/search/tags | 按标签搜索 |
| GET | /api/v1/prompts/check-duplicate | 检查重复 |

### 标签接口

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/tags/ | 创建标签 |
| GET | /api/v1/tags/ | 获取所有标签 |
| GET | /api/v1/tags/:id | 获取单个标签 |
| DELETE | /api/v1/tags/:id | 删除标签 |
| GET | /api/v1/tags/search | 搜索标签 |
| GET | /api/v1/tags/stats | 获取标签统计 |

### 系统接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /health | 健康检查 |
| GET | /db-status | 数据库状态检查 |
| GET | / | API信息 |

## 快速开始

### 1. 环境准备

确保已安装：
- Go 1.19+
- MySQL 5.7+

### 2. 克隆项目

```bash
git clone <repository-url>
cd imgGeneratePrompts
```

### 3. 配置数据库

创建数据库配置文件：
```bash
cp apikey/database.env.example apikey/database.env
```

编辑 `apikey/database.env` 文件，配置数据库连接信息：
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=img_prompts
```

### 4. 初始化数据库

执行数据库初始化脚本：
```bash
mysql -u your_username -p < scripts/init.sql
```

### 5. 安装依赖

```bash
go mod download
```

### 6. 运行项目

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动

## V3.1 升级指南

如果您从V3.0升级到V3.1，数据库结构已更新。请按以下步骤操作：

1. **备份数据库**
   ```bash
   mysqldump -u username -p img_prompts > backup_v3.sql
   ```

2. **执行字段重命名（如果需要）**
   ```sql
   -- 如果你的表中仍使用旧字段名，请执行以下SQL
   ALTER TABLE prompts 
   CHANGE COLUMN reference_images input_image_url VARCHAR(500) COMMENT '输入的参照图片的存储路径或URL；可能多个图片';
   
   ALTER TABLE prompts 
   CHANGE COLUMN output_image output_image_url VARCHAR(500) COMMENT '输出的参照图片的存储路径或URL';
   ```

3. **数据格式转换（如果需要）**
   ```sql
   -- 如果原来存储的是JSON格式，转换为逗号分隔格式
   -- 这个步骤需要根据具体数据情况编写转换脚本
   ```

4. **更新代码**
   ```bash
   git pull origin master
   go mod download
   ```

5. **重启服务**
   ```bash
   go run main.go
   ```

## 存储格式说明

### 输入图片URL存储

- **数据库格式**：`/uploads/180151.jpg,/uploads/180150.jpg`（逗号分隔字符串）
- **API响应格式**：`["uploads/180151.jpg", "/uploads/180150.jpg"]`（字符串数组）
- **最大长度**：500字符
- **分隔符**：英文逗号（`,`）

### 示例数据
```sql
INSERT INTO prompts (
  prompt_text, 
  input_image_url, 
  output_image_url,
  ...
) VALUES (
  'Create a professional e-commerce fashion photo...',
  '/uploads/180151.jpg,/uploads/180150.jpg',
  '/uploads/180152.jpg',
  ...
);
```

## AI集成指南

当前版本的AI分析功能使用模拟数据。要集成真实的AI服务，请修改 `services/prompt_service.go` 中的 `AnalyzePromptData` 方法：

```go
func (s *PromptService) AnalyzePromptData(promptText, modelName, outputImageBase64 string, inputImagesBase64 []string) (*models.AnalyzePromptResponse, error) {
    // 替换为您的AI API调用
    // 例如：Google Gemini, OpenAI Vision API等
}
```

## 开发指南

### 运行测试

```bash
go test ./...
```

### 构建二进制文件

```bash
go build -o bin/img-prompts main.go
```

### 代码格式化

```bash
go fmt ./...
```

## 贡献指南

欢迎提交Pull Request或Issue！

## 许可证

MIT License

## 更新日志

### V3.1.0 (2024-09)
- 🔄 优化数据库字段结构
- 📝 更新字段命名：`input_image_url` 替代 `reference_images`
- 💾 优化存储格式：逗号分隔字符串替代JSON
- 🔄 保持API兼容性
- 📚 更新文档和示例

### V3.0.0 (2024-01)
- ✨ 新增多图片上传支持
- ✨ 新增AI智能生成功能
- ✨ 向后兼容旧版本
- 🔥 移除Docker相关文件
- 🐛 修复已知问题
- 📝 更新文档

### V2.0.0
- 初始版本发布
- 基础CRUD功能
- 标签系统
- 搜索功能
