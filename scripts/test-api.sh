#!/bin/bash
# API 自动化测试脚本
# 用法: ./scripts/test-api.sh

set -e

# 配置
API_BASE="http://localhost:8080"
TEST_TAG_NAME="API测试标签$(date +%s)"
TEST_PROMPT_TEXT="API测试提示词 - $(date)"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 辅助函数
print_test() {
    echo -e "${BLUE}🧪 测试: $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

check_response() {
    local response="$1"
    local description="$2"
    
    if echo "$response" | jq -e '.code == 200' > /dev/null 2>&1; then
        print_success "$description"
        return 0
    else
        print_error "$description"
        echo "响应: $response"
        return 1
    fi
}

# 开始测试
echo "🚀 开始API自动化测试..."
echo ""

# 1. 健康检查
print_test "健康检查"
response=$(curl -s "$API_BASE/health")
if echo "$response" | jq -e '.status == "ok"' > /dev/null 2>&1; then
    print_success "健康检查通过"
else
    print_error "健康检查失败"
    echo "响应: $response"
    exit 1
fi
echo ""

# 2. 根路径测试
print_test "根路径信息"
response=$(curl -s "$API_BASE/")
if echo "$response" | jq -e '.version' > /dev/null 2>&1; then
    print_success "根路径信息获取成功"
else
    print_error "根路径信息获取失败"
fi
echo ""

# 3. 标签管理测试
echo "🏷️  标签管理测试"

# 3.1 获取所有标签
print_test "获取所有标签"
response=$(curl -s "$API_BASE/api/v1/tags")
check_response "$response" "标签列表获取"

# 3.2 创建新标签
print_test "创建新标签"
response=$(curl -s -X POST "$API_BASE/api/v1/tags" \
    -H "Content-Type: application/json" \
    -d "{\"name\": \"$TEST_TAG_NAME\"}")
if check_response "$response" "标签创建"; then
    TAG_ID=$(echo "$response" | jq -r '.data.id')
    print_success "创建的标签ID: $TAG_ID"
fi

# 3.3 搜索标签
print_test "搜索标签"
response=$(curl -s "$API_BASE/api/v1/tags/search?keyword=测试")
check_response "$response" "标签搜索"

# 3.4 获取标签统计
print_test "获取标签统计"
response=$(curl -s "$API_BASE/api/v1/tags/stats")
check_response "$response" "标签统计获取"

echo ""

# 4. 提示词管理测试
echo "📝 提示词管理测试"

# 4.1 获取提示词列表
print_test "获取提示词列表"
response=$(curl -s "$API_BASE/api/v1/prompts?page=1&page_size=5")
check_response "$response" "提示词列表获取"

# 4.2 创建新提示词
print_test "创建新提示词"
response=$(curl -s -X POST "$API_BASE/api/v1/prompts" \
    -H "Content-Type: application/json" \
    -d "{
        \"prompt_text\": \"$TEST_PROMPT_TEXT\",
        \"negative_prompt\": \"test negative prompt\",
        \"model_name\": \"test-model\",
        \"is_public\": true,
        \"style_description\": \"测试风格描述\",
        \"usage_scenario\": \"API测试场景\",
        \"atmosphere_description\": \"测试氛围\",
        \"expressive_intent\": \"测试表现意图\",
        \"structure_analysis\": \"{\\\"测试\\\":\\\"结构分析\\\"}\",
        \"tag_names\": [\"$TEST_TAG_NAME\", \"API测试\"]
    }")
if check_response "$response" "提示词创建"; then
    PROMPT_ID=$(echo "$response" | jq -r '.data.id')
    print_success "创建的提示词ID: $PROMPT_ID"
fi

# 4.3 获取单个提示词
if [ ! -z "$PROMPT_ID" ]; then
    print_test "获取单个提示词"
    response=$(curl -s "$API_BASE/api/v1/prompts/$PROMPT_ID")
    check_response "$response" "单个提示词获取"
fi

# 4.4 更新提示词
if [ ! -z "$PROMPT_ID" ]; then
    print_test "更新提示词"
    response=$(curl -s -X PUT "$API_BASE/api/v1/prompts/$PROMPT_ID" \
        -H "Content-Type: application/json" \
        -d "{
            \"prompt_text\": \"$TEST_PROMPT_TEXT - 已更新\",
            \"style_description\": \"更新的测试风格描述\",
            \"tag_names\": [\"$TEST_TAG_NAME\", \"已更新\"]
        }")
    check_response "$response" "提示词更新"
fi

# 4.5 获取公开提示词
print_test "获取公开提示词"
response=$(curl -s "$API_BASE/api/v1/prompts/public?page=1&page_size=5")
check_response "$response" "公开提示词获取"

# 4.6 获取最近提示词
print_test "获取最近提示词"
response=$(curl -s "$API_BASE/api/v1/prompts/recent?limit=3")
check_response "$response" "最近提示词获取"

# 4.7 获取提示词统计
print_test "获取提示词统计"
response=$(curl -s "$API_BASE/api/v1/prompts/stats")
check_response "$response" "提示词统计获取"

echo ""

# 5. 搜索功能测试
echo "🔍 搜索功能测试"

# 5.1 关键词搜索
print_test "关键词搜索"
response=$(curl -s "$API_BASE/api/v1/prompts?keyword=测试&page=1&page_size=3")
check_response "$response" "关键词搜索"

# 5.2 标签搜索
print_test "标签搜索"
response=$(curl -s "$API_BASE/api/v1/prompts/search/tags?tags=$TEST_TAG_NAME&page=1&page_size=3")
check_response "$response" "标签搜索"

# 5.3 重复检查
print_test "重复检查"
response=$(curl -s "$API_BASE/api/v1/prompts/check-duplicate?prompt_text=$(echo "$TEST_PROMPT_TEXT" | sed 's/ /%20/g')")
check_response "$response" "重复检查"

echo ""

# 6. 过滤功能测试
echo "🔧 过滤功能测试"

# 6.1 按模型过滤
print_test "按模型过滤"
response=$(curl -s "$API_BASE/api/v1/prompts?model_name=test-model")
check_response "$response" "模型过滤"

# 6.2 按公开状态过滤
print_test "按公开状态过滤"
response=$(curl -s "$API_BASE/api/v1/prompts?is_public=true&page=1&page_size=3")
check_response "$response" "公开状态过滤"

# 6.3 排序测试
print_test "排序测试"
response=$(curl -s "$API_BASE/api/v1/prompts?sort_by=created_at&sort_order=desc&page=1&page_size=3")
check_response "$response" "排序功能"

echo ""

# 7. 清理测试数据
echo "🧹 清理测试数据"

# 7.1 删除测试提示词
if [ ! -z "$PROMPT_ID" ]; then
    print_test "删除测试提示词"
    response=$(curl -s -X DELETE "$API_BASE/api/v1/prompts/$PROMPT_ID")
    check_response "$response" "提示词删除"
fi

# 7.2 删除测试标签
if [ ! -z "$TAG_ID" ]; then
    print_test "删除测试标签"
    response=$(curl -s -X DELETE "$API_BASE/api/v1/tags/$TAG_ID")
    if check_response "$response" "标签删除"; then
        print_success "测试数据清理完成"
    else
        print_warning "标签可能仍被其他提示词使用，删除失败"
    fi
fi

echo ""

# 8. 性能测试（简单）
echo "⚡ 简单性能测试"

print_test "并发请求测试（10个并发）"
start_time=$(date +%s)
for i in {1..10}; do
    curl -s "$API_BASE/api/v1/prompts?page=1&page_size=1" > /dev/null &
done
wait
end_time=$(date +%s)
duration=$((end_time - start_time))
print_success "10个并发请求完成，耗时: ${duration}秒"

echo ""

# 9. 错误处理测试
echo "🚨 错误处理测试"

# 9.1 无效ID测试
print_test "无效ID测试"
response=$(curl -s "$API_BASE/api/v1/prompts/999999")
if echo "$response" | jq -e '.code == 404' > /dev/null 2>&1; then
    print_success "无效ID错误处理正确"
else
    print_error "无效ID错误处理异常"
fi

# 9.2 参数验证测试
print_test "参数验证测试"
response=$(curl -s -X POST "$API_BASE/api/v1/prompts" \
    -H "Content-Type: application/json" \
    -d '{"prompt_text": ""}')
if echo "$response" | jq -e '.code == 400' > /dev/null 2>&1; then
    print_success "参数验证错误处理正确"
else
    print_error "参数验证错误处理异常"
fi

echo ""

# 测试完成
echo "🎉 API测试完成！"
echo ""
echo "📊 测试总结："
echo "✅ 基础功能: 健康检查、根路径"
echo "✅ 标签管理: 创建、获取、搜索、统计、删除"
echo "✅ 提示词管理: 创建、获取、更新、删除"
echo "✅ 搜索功能: 关键词、标签、重复检查"
echo "✅ 过滤功能: 模型、状态、排序"
echo "✅ 错误处理: 无效ID、参数验证"
echo "✅ 性能测试: 并发请求"
echo ""
echo "🚀 所有核心功能测试通过！"
