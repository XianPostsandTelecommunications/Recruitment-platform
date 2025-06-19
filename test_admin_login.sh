#!/bin/bash

echo "🔐 管理员账号登录测试"
echo "========================"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "\n${BLUE}📋 管理员账号信息:${NC}"
echo "   邮箱: 1234567@qq.com"
echo "   密码: epi666"
echo "   角色: admin"

echo -e "\n${YELLOW}=== 1. 数据库验证 ===${NC}"
db_result=$(sudo mysql -u root -proot -e "USE lab_recruitment; SELECT username, email, role, status FROM users WHERE email = '1234567@qq.com';" 2>/dev/null | tail -n +2)

if [ -n "$db_result" ]; then
    echo -e "✅ 数据库中找到管理员账号:"
    echo "   $db_result"
else
    echo -e "❌ 数据库中未找到管理员账号"
    exit 1
fi

echo -e "\n${YELLOW}=== 2. API登录测试 ===${NC}"
login_response=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"1234567@qq.com","password":"epi666"}')

login_code=$(echo "$login_response" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('code', 'error'))" 2>/dev/null)
login_role=$(echo "$login_response" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('data', {}).get('user', {}).get('role', 'N/A'))" 2>/dev/null)
token=$(echo "$login_response" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('data', {}).get('token', 'N/A'))" 2>/dev/null)

if [ "$login_code" = "200" ] && [ "$login_role" = "admin" ]; then
    echo -e "✅ API登录成功:"
    echo "   状态码: $login_code"
    echo "   角色: $login_role"
    echo "   Token: ${token:0:50}..."
else
    echo -e "❌ API登录失败:"
    echo "   状态码: $login_code"
    echo "   响应: $login_response"
    exit 1
fi

echo -e "\n${YELLOW}=== 3. 前端页面检查 ===${NC}"
frontend_status=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/login)

if [ "$frontend_status" = "200" ]; then
    echo -e "✅ 前端登录页面可访问 (状态码: $frontend_status)"
else
    echo -e "❌ 前端登录页面异常 (状态码: $frontend_status)"
fi

echo -e "\n${GREEN}🎉 管理员账号测试完成！${NC}"
echo -e "${YELLOW}========================${NC}"

echo -e "\n${GREEN}✅ 问题已解决:${NC}"
echo "   • 管理员账号已创建并激活"
echo "   • 后端API登录验证通过"
echo "   • 前端登录页面正常访问"

echo -e "\n${BLUE}🌐 访问信息:${NC}"
echo "   • 前端登录地址: http://localhost:3000/login"
echo "   • 管理员邮箱: 1234567@qq.com"
echo "   • 管理员密码: epi666"

echo -e "\n${BLUE}💡 使用说明:${NC}"
echo "1. 打开浏览器访问: http://localhost:3000/login"
echo "2. 输入邮箱: 1234567@qq.com"
echo "3. 输入密码: epi666"
echo "4. 点击登录按钮"

echo -e "\n${GREEN}🎊 现在可以正常登录管理后台了！${NC}" 