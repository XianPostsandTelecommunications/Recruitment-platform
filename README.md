# 实验室招新平台开发文档

## 项目概述

实验室招新平台是一个基于 Go + MySQL + Gin + React 技术栈的现代化Web应用，旨在为高校实验室提供便捷的招新管理解决方案。平台支持学生浏览实验室信息、提交申请，同时为管理员提供完整的后台管理功能。

## 技术栈

### 后端技术栈
- **语言**: Go 1.21+
- **Web框架**: Gin
- **数据库**: MySQL 8.0+
- **ORM**: GORM
- **认证**: JWT-Go
- **密码加密**: bcrypt
- **数据验证**: validator
- **日志**: logrus

### 前端技术栈
- **框架**: React 18 + TypeScript
- **路由**: React Router 6
- **状态管理**: Redux Toolkit
- **UI组件库**: Ant Design
- **样式**: Tailwind CSS
- **构建工具**: Vite
- **HTTP客户端**: Axios
- **图表**: Recharts
- **动画**: Framer Motion

## 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React 前端    │    │   Gin 后端      │    │   MySQL 数据库  │
│                 │    │                 │    │                 │
│ - 用户界面      │◄──►│ - RESTful API   │◄──►│ - 用户数据      │
│ - 状态管理      │    │ - 业务逻辑      │    │ - 实验室数据    │
│ - 路由管理      │    │ - 数据验证      │    │ - 申请数据      │
│ - 组件库        │    │ - 认证授权      │    │ - 通知数据      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 功能模块

### 1. 用户管理模块
- 用户注册/登录
- JWT身份认证
- 角色权限控制（学生/管理员）
- 个人资料管理

### 2. 实验室管理模块
- 实验室信息展示
- 实验室创建/编辑/删除（管理员）
- 实验室搜索和筛选
- 标签分类管理

### 3. 申请管理模块
- 学生提交申请
- 申请状态跟踪
- 管理员审核申请
- 申请历史查看

### 4. 通知管理模块
- 系统通知推送
- 申请状态通知
- 通知已读标记
- 通知历史查看

### 5. 统计管理模块
- 平台数据统计
- 申请数据分析
- 用户活跃度统计
- 可视化图表展示

## API 接口文档

### 认证相关接口

#### 用户注册
```
POST /api/auth/register
Content-Type: application/json

Request Body:
{
  "username": "string",
  "email": "string", 
  "password": "string",
  "role": "student"
}

Response:
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "id": 1,
    "username": "string",
    "email": "string",
    "role": "student",
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

#### 用户登录
```
POST /api/auth/login
Content-Type: application/json

Request Body:
{
  "email": "string",
  "password": "string"
}

Response:
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "jwt_token_string",
    "user": {
      "id": 1,
      "username": "string",
      "email": "string",
      "role": "student"
    }
  }
}
```

### 实验室管理接口

#### 获取实验室列表
```
GET /api/labs?page=1&size=10&search=关键词&tags=标签

Response:
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 1,
        "name": "人工智能实验室",
        "description": "专注于AI技术研究",
        "requirements": "熟悉Python，有机器学习基础",
        "maxMembers": 15,
        "contactEmail": "ai@university.edu",
        "tags": ["AI", "机器学习"],
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 10
  }
}
```

#### 创建实验室
```
POST /api/labs
Authorization: Bearer <token>
Content-Type: application/json

Request Body:
{
  "name": "string",
  "description": "string",
  "requirements": "string",
  "maxMembers": 10,
  "contactEmail": "string",
  "tags": ["tag1", "tag2"]
}
```

### 申请管理接口

#### 提交申请
```
POST /api/applications
Authorization: Bearer <token>
Content-Type: application/json

Request Body:
{
  "labId": 1,
  "motivation": "对AI技术有浓厚兴趣",
  "skills": ["Python", "机器学习"],
  "experience": "参加过相关项目",
  "availableTime": "每周20小时"
}
```

#### 审核申请
```
PUT /api/applications/:id/status
Authorization: Bearer <token>
Content-Type: application/json

Request Body:
{
  "status": "accepted",
  "feedback": "欢迎加入我们的团队！"
}
```

## 数据库设计

### 用户表 (users)
```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL COMMENT '用户名',
    email VARCHAR(100) UNIQUE NOT NULL COMMENT '邮箱',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    role ENUM('student', 'admin') DEFAULT 'student' COMMENT '用户角色',
    avatar VARCHAR(255) COMMENT '头像URL',
    phone VARCHAR(20) COMMENT '手机号',
    student_id VARCHAR(20) COMMENT '学号',
    major VARCHAR(100) COMMENT '专业',
    grade VARCHAR(20) COMMENT '年级',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

### 实验室表 (labs)
```sql
CREATE TABLE labs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '实验室名称',
    description TEXT COMMENT '实验室描述',
    requirements TEXT COMMENT '招新要求',
    max_members INT DEFAULT 10 COMMENT '最大成员数',
    current_members INT DEFAULT 0 COMMENT '当前成员数',
    contact_email VARCHAR(100) COMMENT '联系邮箱',
    contact_phone VARCHAR(20) COMMENT '联系电话',
    location VARCHAR(200) COMMENT '实验室位置',
    tags JSON COMMENT '标签数组',
    cover_image VARCHAR(255) COMMENT '封面图片',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_by BIGINT COMMENT '创建者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id),
    INDEX idx_status (status),
    INDEX idx_created_by (created_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='实验室表';
```

### 申请表 (applications)
```sql
CREATE TABLE applications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '申请人ID',
    lab_id BIGINT NOT NULL COMMENT '实验室ID',
    motivation TEXT COMMENT '申请动机',
    skills JSON COMMENT '技能列表',
    experience TEXT COMMENT '相关经验',
    available_time VARCHAR(100) COMMENT '可用时间',
    resume_url VARCHAR(255) COMMENT '简历文件URL',
    status ENUM('pending', 'accepted', 'rejected') DEFAULT 'pending' COMMENT '申请状态',
    feedback TEXT COMMENT '审核反馈',
    reviewed_by BIGINT COMMENT '审核人ID',
    reviewed_at TIMESTAMP NULL COMMENT '审核时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (lab_id) REFERENCES labs(id),
    FOREIGN KEY (reviewed_by) REFERENCES users(id),
    INDEX idx_user_id (user_id),
    INDEX idx_lab_id (lab_id),
    INDEX idx_status (status),
    UNIQUE KEY uk_user_lab (user_id, lab_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='申请表';
```

### 通知表 (notifications)
```sql
CREATE TABLE notifications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '接收用户ID',
    title VARCHAR(200) NOT NULL COMMENT '通知标题',
    content TEXT COMMENT '通知内容',
    type ENUM('system', 'application', 'lab') DEFAULT 'system' COMMENT '通知类型',
    is_read BOOLEAN DEFAULT FALSE COMMENT '是否已读',
    related_id BIGINT COMMENT '关联ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    INDEX idx_user_id (user_id),
    INDEX idx_is_read (is_read),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知表';
```

## 项目结构

### 后端项目结构
```
├── cmd/
│   └── main.go                 # 程序入口
├── internal/
│   ├── config/
│   │   ├── config.go          # 配置管理
│   │   └── database.go        # 数据库配置
│   ├── models/
│   │   ├── user.go            # 用户模型
│   │   ├── lab.go             # 实验室模型
│   │   ├── application.go     # 申请模型
│   │   └── notification.go    # 通知模型
│   ├── handlers/
│   │   ├── auth.go            # 认证处理器
│   │   ├── lab.go             # 实验室处理器
│   │   ├── application.go     # 申请处理器
│   │   ├── notification.go    # 通知处理器
│   │   └── stats.go           # 统计处理器
│   ├── middleware/
│   │   ├── auth.go            # JWT认证中间件
│   │   ├── cors.go            # CORS中间件
│   │   ├── logger.go          # 日志中间件
│   │   └── validator.go       # 数据验证中间件
│   ├── services/
│   │   ├── auth_service.go    # 认证服务
│   │   ├── lab_service.go     # 实验室服务
│   │   ├── application_service.go # 申请服务
│   │   └── notification_service.go # 通知服务
│   └── utils/
│       ├── jwt.go             # JWT工具
│       ├── password.go        # 密码工具
│       ├── response.go        # 响应工具
│       └── validator.go       # 验证工具
├── pkg/
│   └── logger/
│       └── logger.go          # 日志包
├── migrations/
│   └── *.sql                  # 数据库迁移文件
├── docs/
│   └── api.md                 # API文档
├── go.mod
├── go.sum
└── README.md
```

### 前端项目结构
```
src/
├── components/
│   ├── common/
│   │   ├── Header.tsx         # 头部组件
│   │   ├── Footer.tsx         # 底部组件
│   │   ├── Loading.tsx        # 加载组件
│   │   ├── ErrorBoundary.tsx  # 错误边界
│   │   └── Modal.tsx          # 模态框组件
│   ├── layout/
│   │   ├── MainLayout.tsx     # 主布局
│   │   ├── AdminLayout.tsx    # 管理布局
│   │   └── Sidebar.tsx        # 侧边栏
│   ├── forms/
│   │   ├── LoginForm.tsx      # 登录表单
│   │   ├── RegisterForm.tsx   # 注册表单
│   │   ├── LabForm.tsx        # 实验室表单
│   │   └── ApplicationForm.tsx # 申请表单
│   └── ui/
│       ├── Card.tsx           # 卡片组件
│       ├── Button.tsx         # 按钮组件
│       ├── Tag.tsx            # 标签组件
│       └── StatusBadge.tsx    # 状态徽章
├── pages/
│   ├── auth/
│   │   ├── Login.tsx          # 登录页
│   │   └── Register.tsx       # 注册页
│   ├── labs/
│   │   ├── LabList.tsx        # 实验室列表
│   │   ├── LabDetail.tsx      # 实验室详情
│   │   └── LabForm.tsx        # 实验室表单
│   ├── applications/
│   │   ├── ApplicationList.tsx # 申请列表
│   │   ├── ApplicationDetail.tsx # 申请详情
│   │   └── ApplicationForm.tsx # 申请表单
│   ├── admin/
│   │   ├── Dashboard.tsx      # 管理仪表板
│   │   ├── UserManagement.tsx # 用户管理
│   │   ├── LabManagement.tsx  # 实验室管理
│   │   └── ApplicationReview.tsx # 申请审核
│   └── profile/
│       └── Profile.tsx        # 个人资料
├── hooks/
│   ├── useAuth.ts             # 认证钩子
│   ├── useApi.ts              # API钩子
│   └── useLocalStorage.ts     # 本地存储钩子
├── services/
│   ├── api.ts                 # API配置
│   ├── auth.ts                # 认证服务
│   ├── lab.ts                 # 实验室服务
│   ├── application.ts         # 申请服务
│   └── notification.ts        # 通知服务
├── store/
│   ├── index.ts               # Store配置
│   ├── authSlice.ts           # 认证状态
│   ├── labSlice.ts            # 实验室状态
│   └── applicationSlice.ts    # 申请状态
├── types/
│   ├── user.ts                # 用户类型
│   ├── lab.ts                 # 实验室类型
│   ├── application.ts         # 申请类型
│   └── common.ts              # 通用类型
├── utils/
│   ├── constants.ts           # 常量
│   ├── helpers.ts             # 工具函数
│   └── validation.ts          # 验证函数
├── assets/
│   ├── images/                # 图片资源
│   ├── icons/                 # 图标资源
│   └── styles/                # 样式文件
├── App.tsx
├── main.tsx
└── index.css
```

## 开发环境搭建

### 后端环境要求
- Go 1.21+
- MySQL 8.0+
- Redis (可选，用于缓存)

### 前端环境要求
- Node.js 18+
- npm 或 yarn

### 安装步骤

#### 1. 克隆项目
```bash
git clone <repository-url>
cd lab-recruitment-platform
```

#### 2. 后端设置
```bash
# 进入后端目录
cd backend

# 安装依赖
go mod tidy

# 配置数据库
cp .env.example .env
# 编辑 .env 文件，配置数据库连接信息

# 运行数据库迁移
go run cmd/migrate/main.go

# 启动开发服务器
go run cmd/main.go
```

#### 3. 前端设置
```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置API地址

# 启动开发服务器
npm run dev
```

## 部署指南

### 生产环境部署

#### 1. 后端部署
```bash
# 构建二进制文件
go build -o bin/server cmd/main.go

# 使用 Docker 部署
docker build -t lab-recruitment-backend .
docker run -d -p 8080:8080 lab-recruitment-backend
```

#### 2. 前端部署
```bash
# 构建生产版本
npm run build

# 使用 Nginx 部署
# 将 dist 目录内容复制到 Nginx 静态文件目录
```

### Docker Compose 部署
```yaml
version: '3.8'
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=password
      - DB_NAME=lab_recruitment
    depends_on:
      - mysql

  frontend:
    build: ./frontend
    ports:
      - "3000:80"
    depends_on:
      - backend

  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=password
      - MYSQL_DATABASE=lab_recruitment
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"

volumes:
  mysql_data:
```

## 测试指南

### 后端测试
```bash
# 运行单元测试
go test ./...

# 运行集成测试
go test -tags=integration ./...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 前端测试
```bash
# 运行单元测试
npm test

# 运行 E2E 测试
npm run test:e2e

# 生成测试覆盖率报告
npm run test:coverage
```

## 性能优化

### 后端优化
1. 数据库查询优化
   - 添加适当的索引
   - 使用连接池
   - 实现查询缓存

2. API 性能优化
   - 实现分页查询
   - 添加响应缓存
   - 使用 gzip 压缩

3. 并发处理
   - 使用 goroutines 处理并发请求
   - 实现限流机制

### 前端优化
1. 代码分割
   - 路由级别的代码分割
   - 组件懒加载

2. 资源优化
   - 图片懒加载
   - 静态资源压缩
   - CDN 加速

3. 缓存策略
   - 浏览器缓存
   - 状态管理缓存

## 安全考虑

### 后端安全
1. 身份认证
   - JWT token 过期机制
   - 密码强度验证
   - 多因素认证

2. 数据安全
   - SQL 注入防护
   - XSS 攻击防护
   - CSRF 防护

3. API 安全
   - 请求频率限制
   - 输入数据验证
   - 错误信息脱敏

### 前端安全
1. 数据验证
   - 客户端输入验证
   - XSS 防护
   - 敏感信息保护

2. 网络安全
   - HTTPS 强制
   - CSP 策略
   - 安全头部设置

## 监控和日志

### 日志管理
- 使用结构化日志
- 日志级别分类
- 日志轮转策略

### 监控指标
- API 响应时间
- 错误率统计
- 用户活跃度
- 系统资源使用

## 维护和更新

### 版本管理
- 使用语义化版本号
- 维护更新日志
- 向后兼容性考虑

### 数据备份
- 定期数据库备份
- 文件存储备份
- 灾难恢复计划

## 贡献指南

### 代码规范
- 遵循 Go 官方代码规范
- 使用 ESLint 和 Prettier
- 编写单元测试

### 提交规范
- 使用 Conventional Commits
- 提供详细的提交信息
- 关联 Issue 和 PR

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。

## 联系方式

如有问题或建议，请通过以下方式联系：
- 邮箱: support@lab-recruitment.com
- 项目地址: https://github.com/your-org/lab-recruitment-platform
- 文档地址: https://docs.lab-recruitment.com 