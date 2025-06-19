# 🚀 EPI实验室面试申请系统 - 生产环境部署指南

## 📋 系统概述

这是一个完整的实验室面试申请系统，包含：
- 📝 学生面试申请表单
- 👨‍💼 管理员后台管理
- 📧 邮箱验证功能
- 🎉 成功/失败弹窗提示
- 🔍 申请搜索和过滤
- 📊 申请数据统计

## 🛠 技术栈

### 前端
- **框架**: React 18 + TypeScript + Vite
- **UI库**: Ant Design 5.x
- **状态管理**: Redux Toolkit
- **路由**: React Router 6

### 后端
- **语言**: Go 1.21+
- **框架**: Gin Web Framework
- **数据库**: MySQL 8.0
- **ORM**: GORM
- **认证**: JWT Token
- **邮件**: SMTP (支持QQ邮箱)

## 📦 部署环境要求

### 服务器配置
- **CPU**: 2核心以上
- **内存**: 4GB以上
- **存储**: 20GB以上
- **操作系统**: Ubuntu 20.04+ / CentOS 8+

### 软件依赖
- **Go**: 1.21+
- **Node.js**: 18+
- **MySQL**: 8.0+
- **Nginx**: 1.18+ (可选，用于反向代理)

## 🔧 部署步骤

### 1. 克隆代码
```bash
git clone https://github.com/XianPostsandTelecommunications/Recruitment-platform.git
cd Recruitment-platform
```

### 2. 数据库配置
```sql
-- 创建数据库
CREATE DATABASE lab_recruitment CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户（可选）
CREATE USER 'lab_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON lab_recruitment.* TO 'lab_user'@'localhost';
FLUSH PRIVILEGES;
```

### 3. 环境配置
编辑 `config.yaml`：
```yaml
server:
  port: 8080
  mode: release  # 生产环境设置为 release

database:
  host: localhost
  port: 3306
  database: lab_recruitment
  username: lab_user
  password: your_password

jwt:
  secret: your_jwt_secret_key_here
  expire_hours: 24

email:
  smtp_host: smtp.qq.com
  smtp_port: 587
  username: your_email@qq.com
  password: your_email_password
  from_name: EPI实验室
```

### 4. 后端部署
```bash
# 编译后端
go mod tidy
go build -o main cmd/main/main.go

# 运行后端（推荐使用进程管理器）
nohup ./main > app.log 2>&1 &

# 或使用 systemd (推荐)
sudo cp deployment/lab-recruitment.service /etc/systemd/system/
sudo systemctl enable lab-recruitment
sudo systemctl start lab-recruitment
```

### 5. 前端部署
```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 构建生产版本
npm run build

# 部署到web服务器
sudo cp -r dist/* /var/www/html/
```

### 6. Nginx配置（推荐）
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    # 前端静态文件
    location / {
        root /var/www/html;
        try_files $uri $uri/ /index.html;
    }
    
    # 后端API代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## 🔐 安全配置

### 1. 防火墙设置
```bash
# 只开放必要端口
sudo ufw allow 22    # SSH
sudo ufw allow 80    # HTTP
sudo ufw allow 443   # HTTPS (如果使用SSL)
sudo ufw enable
```

### 2. SSL证书（推荐）
```bash
# 使用 Let's Encrypt
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

### 3. 管理员账号
系统默认管理员账号：
- **邮箱**: `1234567@qq.com`
- **密码**: `epi666`

**⚠️ 部署后请立即修改默认密码！**

## 🎯 系统功能

### 用户端功能
- **申请表单**: http://your-domain.com/apply
- **邮箱验证**: 6位数字验证码
- **实时验证**: 表单字段实时验证
- **年级选择**: 大一、大二、大三

### 管理端功能
- **管理后台**: http://your-domain.com/admin
- **申请列表**: 分页显示，状态过滤
- **姓名搜索**: 模糊匹配搜索
- **状态管理**: pending → interviewed → passed/rejected
- **数据统计**: 各状态申请数量

## 📧 邮件配置

### QQ邮箱配置
1. 登录QQ邮箱 → 设置 → 账户
2. 开启SMTP服务
3. 获取授权码
4. 在 `config.yaml` 中配置邮箱信息

### 其他邮箱
支持任何标准SMTP服务，修改对应配置即可。

## 🔍 监控和维护

### 1. 日志监控
```bash
# 查看应用日志
tail -f app.log

# 查看系统状态
sudo systemctl status lab-recruitment
```

### 2. 数据库备份
```bash
# 每日备份脚本
mysqldump -u lab_user -p lab_recruitment > backup_$(date +%Y%m%d).sql
```

### 3. 性能监控
- 监控服务器CPU、内存使用率
- 监控数据库连接数
- 监控申请提交成功率

## 🚨 故障排除

### 常见问题

1. **邮件发送失败**
   - 检查SMTP配置
   - 确认邮箱授权码正确
   - 检查防火墙端口587

2. **数据库连接失败**
   - 检查MySQL服务状态
   - 确认数据库配置正确
   - 检查用户权限

3. **前端页面空白**
   - 检查Nginx配置
   - 确认静态文件路径正确
   - 检查API代理配置

### 测试功能
部署完成后，使用测试验证码 `999999` 进行功能测试。

## 📞 技术支持

如有部署问题，请查看：
1. 应用日志文件
2. 系统日志: `journalctl -u lab-recruitment`
3. Nginx日志: `/var/log/nginx/error.log`

## 🎉 部署完成

部署成功后，访问：
- **用户申请页面**: http://your-domain.com/apply
- **管理后台**: http://your-domain.com/admin

祝贺！EPI实验室面试申请系统已成功上线！ 🚀 