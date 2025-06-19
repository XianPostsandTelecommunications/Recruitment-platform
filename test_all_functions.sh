#!/bin/bash

# 招聘平台功能全链路测试脚本
echo "🚀 开始执行招聘平台全功能测试..."
echo "================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 测试结果统计
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试函数
test_api() {
    local test_name="$1"
    local method="$2"
    local url="$3"
    local data="$4"
    local expected_code="$5"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "\n${BLUE}[测试 $TOTAL_TESTS] $test_name${NC}"
    echo "请求: $method $url"
    
    if [ -n "$data" ]; then
        echo "数据: $data"
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$url")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url")
    fi
    
    # 分离响应体和状态码
    http_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | sed '$d')
    
    echo "状态码: $http_code"
    echo "响应: $response_body"
    
    if [ "$http_code" = "$expected_code" ]; then
        echo -e "${GREEN}✅ 测试通过${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        echo -e "${RED}❌ 测试失败 (期望: $expected_code, 实际: $http_code)${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 3

# 测试基础服务
echo -e "\n${YELLOW}=== 基础服务测试 ===${NC}"

# 1. 健康检查
test_api "健康检查" "GET" "http://localhost:8080/health" "" "200"

# 2. 前端服务检查
test_api "前端服务检查" "GET" "http://localhost:3000" "" "200"

# 测试API接口
echo -e "\n${YELLOW}=== API接口测试 ===${NC}"

# 3. 发送验证码接口
TEST_EMAIL="test@example.com"
test_api "发送验证码" "POST" "http://localhost:8080/api/v1/send-code" \
    "{\"email\":\"$TEST_EMAIL\"}" "200"

# 等待邮件发送
echo "⏳ 等待邮件发送..."
sleep 2

# 4. 测试无效邮箱格式
test_api "无效邮箱格式" "POST" "http://localhost:8080/api/v1/send-code" \
    "{\"email\":\"invalid-email\"}" "400"

# 5. 测试空邮箱
test_api "空邮箱参数" "POST" "http://localhost:8080/api/v1/send-code" \
    "{}" "400"

# 测试申请接口（使用模拟验证码）
echo -e "\n${YELLOW}=== 申请功能测试 ===${NC}"

# 6. 测试申请提交（错误验证码）
test_api "错误验证码申请" "POST" "http://localhost:8080/api/v1/apply" \
    "{
        \"name\":\"张三\",
        \"email\":\"$TEST_EMAIL\",
        \"phone\":\"13800138000\",
        \"student_id\":\"2024001\",
        \"major\":\"计算机科学与技术\",
        \"grade\":\"2024\",
        \"interview_time\":\"2024-12-25 14:00\",
        \"verification_code\":\"123456\"
    }" "400"

# 7. 测试缺少参数的申请
test_api "缺少参数申请" "POST" "http://localhost:8080/api/v1/apply" \
    "{
        \"name\":\"张三\",
        \"email\":\"$TEST_EMAIL\"
    }" "400"

# 测试认证接口
echo -e "\n${YELLOW}=== 认证功能测试 ===${NC}"

# 8. 管理员账号验证
echo -e "\n${BLUE}[测试 $((TOTAL_TESTS + 1))] 管理员账号验证${NC}"
TOTAL_TESTS=$((TOTAL_TESTS + 1))
echo "检查管理员账号是否存在并可用"
admin_exists=$(sudo mysql -u root -proot -e "USE lab_recruitment; SELECT COUNT(*) FROM users WHERE email='1234567@qq.com' AND role='admin';" 2>/dev/null | tail -1)
if [ "$admin_exists" -gt 0 ]; then
    echo -e "${GREEN}✅ 管理员账号存在${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    echo -e "${RED}❌ 管理员账号不存在${NC}"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

# 9. 管理员登录
LOGIN_DATA="{
    \"email\":\"1234567@qq.com\",
    \"password\":\"epi666\"
}"
login_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "$LOGIN_DATA" \
    "http://localhost:8080/api/v1/auth/login")

login_code=$(curl -s -w "%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -d "$LOGIN_DATA" \
    "http://localhost:8080/api/v1/auth/login" | tail -c 3)

echo -e "\n${BLUE}[测试 $((TOTAL_TESTS + 1))] 管理员登录${NC}"
TOTAL_TESTS=$((TOTAL_TESTS + 1))
echo "请求: POST http://localhost:8080/api/v1/auth/login"
echo "状态码: $login_code"

if [ "$login_code" = "200" ]; then
    echo -e "${GREEN}✅ 测试通过${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
    
    # 提取token用于后续测试
    TOKEN=$(echo "$login_response" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('data', {}).get('token', ''))" 2>/dev/null || echo "")
    echo "Token: ${TOKEN:0:20}..."
else
    echo -e "${RED}❌ 测试失败${NC}"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

# 10. 获取用户信息（需要token）
if [ -n "$TOKEN" ]; then
    profile_response=$(curl -s -w "\n%{http_code}" -X GET \
        -H "Authorization: Bearer $TOKEN" \
        "http://localhost:8080/api/v1/auth/profile")
    
    profile_code=$(echo "$profile_response" | tail -n1)
    
    echo -e "\n${BLUE}[测试 $((TOTAL_TESTS + 1))] 获取用户信息${NC}"
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo "状态码: $profile_code"
    
    if [ "$profile_code" = "200" ]; then
        echo -e "${GREEN}✅ 测试通过${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ 测试失败${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
fi

# 数据库连接测试
echo -e "\n${YELLOW}=== 数据库连接测试 ===${NC}"

# 11. 检查数据库表
echo -e "\n${BLUE}[测试 $((TOTAL_TESTS + 1))] 数据库表检查${NC}"
TOTAL_TESTS=$((TOTAL_TESTS + 1))

db_tables=$(sudo mysql -u root -proot -e "USE lab_recruitment; SHOW TABLES;" 2>/dev/null | grep -v "Tables_in_lab_recruitment" | wc -l)

if [ "$db_tables" -gt 0 ]; then
    echo "数据库表数量: $db_tables"
    echo -e "${GREEN}✅ 数据库连接正常${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    echo -e "${RED}❌ 数据库连接异常${NC}"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

# 服务进程检查
echo -e "\n${YELLOW}=== 服务进程检查 ===${NC}"

# 12. 后端进程检查
echo -e "\n${BLUE}[测试 $((TOTAL_TESTS + 1))] 后端进程检查${NC}"
TOTAL_TESTS=$((TOTAL_TESTS + 1))

backend_process=$(ps aux | grep "go run cmd/main/main.go" | grep -v grep | wc -l)
if [ "$backend_process" -gt 0 ]; then
    echo -e "${GREEN}✅ 后端服务进程正常${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    echo -e "${RED}❌ 后端服务进程异常${NC}"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

# 13. 前端进程检查
echo -e "\n${BLUE}[测试 $((TOTAL_TESTS + 1))] 前端进程检查${NC}"
TOTAL_TESTS=$((TOTAL_TESTS + 1))

frontend_process=$(ps aux | grep "npm run dev" | grep -v grep | wc -l)
if [ "$frontend_process" -gt 0 ]; then
    echo -e "${GREEN}✅ 前端服务进程正常${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    echo -e "${RED}❌ 前端服务进程异常${NC}"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

# 输出测试总结
echo -e "\n${YELLOW}=================================${NC}"
echo -e "${YELLOW}📊 测试结果总结${NC}"
echo -e "${YELLOW}=================================${NC}"
echo -e "总测试数: ${BLUE}$TOTAL_TESTS${NC}"
echo -e "通过测试: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失败测试: ${RED}$FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 所有测试通过！系统运行正常！${NC}"
    echo -e "\n${GREEN}📋 系统访问地址：${NC}"
    echo -e "🌐 前端界面: http://localhost:3000"
    echo -e "🔧 后端API: http://localhost:8080"
    echo -e "📚 API文档: http://localhost:8080/swagger/index.html"
    echo -e "❤️ 健康检查: http://localhost:8080/health"
    exit 0
else
    echo -e "\n${RED}⚠️ 有 $FAILED_TESTS 项测试失败，请检查系统状态${NC}"
    exit 1
fi 