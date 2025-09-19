# API 测试示例 v2.0

这里提供新版本API的完整测试示例。

## 🚀 快速开始

### 1. 初始化数据库

```bash
# 完整初始化（推荐）
go run cmd/db-manager.go -write

# 查看统计信息
go run cmd/db-manager.go -stats
```

### 2. 启动服务

```bash
go run main.go
```

## 📋 基础接口测试

### 健康检查

```bash
curl -X GET http://localhost:8080/health
```

预期响应：
```json
{
  "status": "ok",
  "message": "Image Generate Prompts API is running"
}
```

### 根路径信息

```bash
curl -X GET http://localhost:8080/
```

## 🏷️ 标签管理接口

### 1. 获取所有标签

```bash
curl -X GET http://localhost:8080/api/v1/tags
```

### 2. 创建新标签

```bash
curl -X POST http://localhost:8080/api/v1/tags \
  -H "Content-Type: application/json" \
  -d '{"name": "新标签"}'
```

### 3. 搜索标签

```bash
curl -X GET "http://localhost:8080/api/v1/tags/search?keyword=风景"
```

### 4. 获取标签统计信息

```bash
curl -X GET http://localhost:8080/api/v1/tags/stats
```

### 5. 删除标签

```bash
curl -X DELETE http://localhost:8080/api/v1/tags/1
```

## 📝 提示词管理接口

### 1. 创建提示词（完整示例）

```bash
curl -X POST http://localhost:8080/api/v1/prompts \
  -H "Content-Type: application/json" \
  -d '{
    "prompt_text": "a beautiful sunset over mountains, golden hour, cinematic lighting, high quality",
    "negative_prompt": "ugly, blurry, low quality, pixelated, noise",
    "model_name": "stable-diffusion-v1-5",
    "is_public": true,
    "style_description": "风景摄影风格，温暖的金色调",
    "usage_scenario": "适用于自然风光、旅游宣传、背景图片",
    "atmosphere_description": "宁静、温暖、壮观的黄昏氛围",
    "expressive_intent": "表现大自然的壮美和宁静",
    "structure_analysis": "{\"主体\":\"山峰日落\",\"光照\":\"黄金时刻\",\"质量\":\"高质量\",\"风格\":\"电影感\"}",
    "tag_names": ["风景", "暖色调", "高质量", "4K"]
  }'
```

### 2. 上传图片并创建提示词

使用Postman或支持multipart/form-data的工具：

- URL: `POST http://localhost:8080/api/v1/prompts/upload`
- Content-Type: `multipart/form-data`
- 表单字段：
  - `image`: 选择图片文件
  - `prompt_text`: "a beautiful landscape"
  - `is_public`: true
  - `tag_names`: "风景,测试"
  - 其他描述字段...

或者使用curl：

```bash
curl -X POST http://localhost:8080/api/v1/prompts/upload \
  -F "image=@/path/to/your/image.jpg" \
  -F "prompt_text=beautiful sunset landscape" \
  -F "negative_prompt=ugly, blurry" \
  -F "model_name=stable-diffusion-v1-5" \
  -F "is_public=true" \
  -F "style_description=风景摄影风格" \
  -F "tag_names=风景,测试,上传"
```

### 3. 获取提示词列表

```bash
# 基础查询
curl -X GET "http://localhost:8080/api/v1/prompts?page=1&page_size=10"

# 带搜索条件
curl -X GET "http://localhost:8080/api/v1/prompts?keyword=sunset&page=1&page_size=5"

# 按标签过滤
curl -X GET "http://localhost:8080/api/v1/prompts?tag_names=风景,高质量"

# 按模型过滤
curl -X GET "http://localhost:8080/api/v1/prompts?model_name=stable-diffusion-v1-5"

# 只获取公开的
curl -X GET "http://localhost:8080/api/v1/prompts?is_public=true"

# 排序
curl -X GET "http://localhost:8080/api/v1/prompts?sort_by=created_at&sort_order=desc"
```

### 4. 获取单个提示词

```bash
curl -X GET http://localhost:8080/api/v1/prompts/1
```

### 5. 更新提示词

```bash
curl -X PUT http://localhost:8080/api/v1/prompts/1 \
  -H "Content-Type: application/json" \
  -d '{
    "prompt_text": "updated prompt text",
    "style_description": "更新的风格描述",
    "tag_names": ["更新", "测试"],
    "is_public": false
  }'
```

### 6. 删除提示词

```bash
curl -X DELETE http://localhost:8080/api/v1/prompts/1
```

## 🔍 搜索和过滤功能

### 1. 获取公开提示词

```bash
curl -X GET "http://localhost:8080/api/v1/prompts/public?page=1&page_size=20"
```

### 2. 根据标签搜索提示词

```bash
curl -X GET "http://localhost:8080/api/v1/prompts/search/tags?tags=风景,高质量&page=1&page_size=10"
```

### 3. 检查重复提示词

```bash
curl -X GET "http://localhost:8080/api/v1/prompts/check-duplicate?prompt_text=a+beautiful+sunset"
```

### 4. 获取最近的提示词

```bash
curl -X GET "http://localhost:8080/api/v1/prompts/recent?limit=5"
```

## 📊 统计信息接口

### 1. 获取提示词统计

```bash
curl -X GET http://localhost:8080/api/v1/prompts/stats
```

预期响应：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_prompts": 10,
    "public_prompts": 8,
    "private_prompts": 2,
    "model_stats": [
      {
        "model_name": "stable-diffusion-v1-5",
        "count": 7
      },
      {
        "model_name": "stable-diffusion-xl",
        "count": 3
      }
    ]
  }
}
```

### 2. 获取标签统计

```bash
curl -X GET http://localhost:8080/api/v1/tags/stats
```

预期响应：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_tags": 15,
    "popular_tags": [
      {
        "tag_id": 1,
        "tag_name": "风景",
        "use_count": 5
      },
      {
        "tag_id": 2,
        "tag_name": "高质量",
        "use_count": 4
      }
    ]
  }
}
```

## 🧪 复杂查询示例

### 1. 多条件查询

```bash
curl -X GET "http://localhost:8080/api/v1/prompts?keyword=landscape&tag_names=风景,高质量&model_name=stable-diffusion-v1-5&is_public=true&sort_by=created_at&sort_order=desc&page=1&page_size=10"
```

### 2. 搜索包含特定关键词的提示词

```bash
# 搜索描述中包含"温暖"的提示词
curl -X GET "http://localhost:8080/api/v1/prompts?keyword=温暖"

# 这会搜索以下字段：
# - prompt_text
# - negative_prompt  
# - style_description
# - usage_scenario
```

## 📋 测试流程

### 完整测试流程

```bash
# 1. 初始化数据库
go run cmd/db-manager.go -write

# 2. 启动服务
go run main.go

# 3. 健康检查
curl http://localhost:8080/health

# 4. 查看初始标签
curl http://localhost:8080/api/v1/tags

# 5. 查看示例数据
curl http://localhost:8080/api/v1/prompts

# 6. 创建新标签
curl -X POST http://localhost:8080/api/v1/tags \
  -H "Content-Type: application/json" \
  -d '{"name": "测试标签"}'

# 7. 创建新提示词
curl -X POST http://localhost:8080/api/v1/prompts \
  -H "Content-Type: application/json" \
  -d '{
    "prompt_text": "test prompt for API",
    "tag_names": ["测试标签", "API"],
    "is_public": true,
    "style_description": "测试风格"
  }'

# 8. 搜索测试
curl -X GET "http://localhost:8080/api/v1/prompts?keyword=test"

# 9. 统计信息
curl http://localhost:8080/api/v1/prompts/stats
curl http://localhost:8080/api/v1/tags/stats

# 10. 数据库统计
go run cmd/db-manager.go -stats
```

## 🐛 错误处理示例

### 1. 无效ID

```bash
curl -X GET http://localhost:8080/api/v1/prompts/999999
```

预期响应：
```json
{
  "code": 404,
  "message": "提示词不存在"
}
```

### 2. 参数验证错误

```bash
curl -X POST http://localhost:8080/api/v1/prompts \
  -H "Content-Type: application/json" \
  -d '{"prompt_text": ""}'
```

预期响应：
```json
{
  "code": 400,
  "message": "参数验证失败: Key: 'CreatePromptRequest.PromptText' Error:Field validation for 'PromptText' failed on the 'required' tag"
}
```

### 3. 重复标签

```bash
# 创建已存在的标签（会返回现有标签，不报错）
curl -X POST http://localhost:8080/api/v1/tags \
  -H "Content-Type: application/json" \
  -d '{"name": "风景"}'
```

## 📊 性能测试

### 1. 批量创建测试

```bash
# 创建多个提示词测试性能
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/prompts \
    -H "Content-Type: application/json" \
    -d "{\"prompt_text\": \"test prompt $i\", \"tag_names\": [\"测试\", \"批量$i\"], \"is_public\": true}"
done
```

### 2. 分页测试

```bash
# 测试大量数据的分页
curl -X GET "http://localhost:8080/api/v1/prompts?page=1&page_size=50"
```

## 🛠️ 开发调试

### 查看数据库状态

```bash
# 查看详细统计
go run cmd/db-manager.go -stats

# 验证数据完整性
go run cmd/db-manager.go -validate
```

### 重置开发环境

```bash
# 重置数据库（危险操作）
go run cmd/db-manager.go -reset

# 重新初始化
go run cmd/db-manager.go -write
```

---

**提示**: 所有API都返回统一的JSON格式，包含 `code`、`message` 和 `data` 字段。
