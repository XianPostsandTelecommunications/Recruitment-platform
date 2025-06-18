# 开发指南

## 开发环境搭建

### 环境要求

#### 后端环境
- **Go**: 1.21+
- **MySQL**: 8.0+
- **Git**: 最新版本
- **IDE**: GoLand, VS Code 或 Vim

#### 前端环境
- **Node.js**: 18+
- **npm**: 9+ 或 **yarn**: 1.22+
- **Git**: 最新版本
- **IDE**: VS Code, WebStorm 或 Vim

### 快速开始

#### 1. 克隆项目
```bash
git clone https://github.com/your-org/lab-recruitment-platform.git
cd lab-recruitment-platform
```

#### 2. 后端开发环境

##### 2.1 安装 Go
```bash
# macOS (使用 Homebrew)
brew install go

# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# Windows
# 下载并安装 Go 官方安装包
```

##### 2.2 配置 Go 环境
```bash
# 设置 GOPATH 和 GOROOT
export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# 验证安装
go version
```

##### 2.3 安装 MySQL
```bash
# macOS
brew install mysql
brew services start mysql

# Ubuntu/Debian
sudo apt install mysql-server
sudo systemctl start mysql
sudo systemctl enable mysql

# Windows
# 下载并安装 MySQL 官方安装包
```

##### 2.4 配置数据库
```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE lab_recruitment CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'dev_user'@'localhost' IDENTIFIED BY 'dev_password';
GRANT ALL PRIVILEGES ON lab_recruitment.* TO 'dev_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

##### 2.5 运行数据库迁移
```bash
cd backend
mysql -u dev_user -p lab_recruitment < ../migrations/001_initial_schema.sql
```

##### 2.6 配置环境变量
```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件
DB_HOST=localhost
DB_PORT=3306
DB_USER=dev_user
DB_PASSWORD=dev_password
DB_NAME=lab_recruitment

JWT_SECRET=dev_jwt_secret_key
JWT_EXPIRE_HOURS=24

SERVER_PORT=8080
SERVER_MODE=debug

LOG_LEVEL=debug
LOG_FILE=logs/app.log

UPLOAD_PATH=uploads
MAX_FILE_SIZE=10485760
```

##### 2.7 安装依赖并启动
```bash
# 安装依赖
go mod tidy

# 启动开发服务器
go run cmd/main.go
```

#### 3. 前端开发环境

##### 3.1 安装 Node.js
```bash
# macOS (使用 Homebrew)
brew install node

# Ubuntu/Debian
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Windows
# 下载并安装 Node.js 官方安装包
```

##### 3.2 配置前端项目
```bash
cd frontend

# 安装依赖
npm install

# 配置环境变量
cp .env.example .env

# 编辑 .env 文件
VITE_API_BASE_URL=http://localhost:8080/api
VITE_APP_TITLE=实验室招新平台
```

##### 3.3 启动开发服务器
```bash
npm run dev
```

## 项目结构说明

### 后端项目结构
```
backend/
├── cmd/                    # 程序入口
│   ├── main.go            # 主程序
│   └── migrate/           # 数据库迁移工具
├── internal/              # 内部包
│   ├── config/            # 配置管理
│   ├── models/            # 数据模型
│   ├── handlers/          # HTTP 处理器
│   ├── middleware/        # 中间件
│   ├── services/          # 业务逻辑
│   └── utils/             # 工具函数
├── pkg/                   # 公共包
├── migrations/            # 数据库迁移文件
├── uploads/               # 文件上传目录
├── logs/                  # 日志文件
├── go.mod                 # Go 模块文件
├── go.sum                 # 依赖校验文件
└── .env                   # 环境变量
```

### 前端项目结构
```
frontend/
├── src/
│   ├── components/        # 组件
│   ├── pages/            # 页面
│   ├── hooks/            # 自定义钩子
│   ├── services/         # API 服务
│   ├── store/            # 状态管理
│   ├── types/            # TypeScript 类型
│   ├── utils/            # 工具函数
│   └── assets/           # 静态资源
├── public/               # 公共资源
├── dist/                 # 构建输出
├── package.json          # 项目配置
├── vite.config.ts        # Vite 配置
└── .env                  # 环境变量
```

## 开发规范

### 代码规范

#### Go 代码规范
1. **命名规范**
   ```go
   // 包名使用小写
   package config
   
   // 变量和函数使用驼峰命名
   var userName string
   func getUserInfo() {}
   
   // 常量使用大写
   const MAX_RETRY_COUNT = 3
   
   // 结构体字段使用驼峰命名
   type User struct {
       ID       int64  `json:"id"`
       Username string `json:"username"`
       Email    string `json:"email"`
   }
   ```

2. **错误处理**
   ```go
   // 始终检查错误
   if err != nil {
       return fmt.Errorf("failed to get user: %w", err)
   }
   
   // 使用自定义错误类型
   var (
       ErrUserNotFound = errors.New("user not found")
       ErrInvalidInput = errors.New("invalid input")
   )
   ```

3. **注释规范**
   ```go
   // Package config provides configuration management
   package config
   
   // User represents a user in the system
   type User struct {
       ID int64 `json:"id"`
   }
   
   // GetUser retrieves a user by ID
   func GetUser(id int64) (*User, error) {
       // implementation
   }
   ```

#### React 代码规范
1. **组件命名**
   ```tsx
   // 使用 PascalCase
   const UserProfile = () => {
       return <div>User Profile</div>
   }
   
   // 文件名与组件名一致
   // UserProfile.tsx
   ```

2. **Props 类型定义**
   ```tsx
   interface UserProfileProps {
       userId: number;
       onUpdate?: (user: User) => void;
   }
   
   const UserProfile: React.FC<UserProfileProps> = ({ userId, onUpdate }) => {
       // implementation
   }
   ```

3. **Hook 使用规范**
   ```tsx
   // 自定义 Hook 以 use 开头
   const useUser = (userId: number) => {
       const [user, setUser] = useState<User | null>(null);
       const [loading, setLoading] = useState(true);
       
       useEffect(() => {
           // fetch user data
       }, [userId]);
       
       return { user, loading };
   }
   ```

### Git 工作流

#### 分支策略
```
main                    # 主分支，生产环境
├── develop            # 开发分支
├── feature/user-auth  # 功能分支
├── bugfix/login-issue # 修复分支
└── hotfix/security    # 热修复分支
```

#### 提交规范
```bash
# 使用 Conventional Commits
feat: add user authentication
fix: resolve login issue
docs: update API documentation
style: format code
refactor: restructure user service
test: add unit tests for auth
chore: update dependencies
```

#### 提交示例
```bash
# 功能开发
git checkout -b feature/user-auth
# 开发完成后
git add .
git commit -m "feat: add user authentication with JWT"
git push origin feature/user-auth

# 创建 Pull Request
# 代码审查通过后合并到 develop
```

### API 开发规范

#### 响应格式
```go
// 统一响应结构
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Errors  []Error     `json:"errors,omitempty"`
}

type Error struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

// 成功响应
func SuccessResponse(data interface{}) Response {
    return Response{
        Code:    200,
        Message: "success",
        Data:    data,
    }
}

// 错误响应
func ErrorResponse(code int, message string, errors []Error) Response {
    return Response{
        Code:    code,
        Message: message,
        Errors:  errors,
    }
}
```

#### 错误码定义
```go
const (
    // 成功
    CodeSuccess = 200
    
    // 客户端错误
    CodeBadRequest     = 400
    CodeUnauthorized   = 401
    CodeForbidden      = 403
    CodeNotFound       = 404
    CodeConflict       = 409
    CodeValidationFail = 422
    
    // 服务器错误
    CodeInternalError = 500
    CodeNotImplemented = 501
    CodeServiceUnavailable = 503
)
```

## 测试指南

### 后端测试

#### 单元测试
```go
// handlers/auth_test.go
package handlers

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "bytes"
    "encoding/json"
)

func TestRegister(t *testing.T) {
    // 准备测试数据
    reqBody := map[string]interface{}{
        "username": "testuser",
        "email":    "test@example.com",
        "password": "password123",
    }
    
    body, _ := json.Marshal(reqBody)
    
    // 创建请求
    req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    
    // 创建响应记录器
    w := httptest.NewRecorder()
    
    // 执行请求
    Register(w, req)
    
    // 验证响应
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
    
    // 验证响应体
    var response Response
    json.Unmarshal(w.Body.Bytes(), &response)
    
    if response.Code != CodeSuccess {
        t.Errorf("Expected code %d, got %d", CodeSuccess, response.Code)
    }
}
```

#### 集成测试
```go
// tests/integration_test.go
package tests

import (
    "testing"
    "net/http"
    "net/http/httptest"
)

func TestAPIIntegration(t *testing.T) {
    // 设置测试数据库
    setupTestDB()
    defer cleanupTestDB()
    
    // 创建测试服务器
    server := httptest.NewServer(setupRouter())
    defer server.Close()
    
    // 测试用户注册
    resp, err := http.Post(server.URL+"/api/auth/register", "application/json", bytes.NewBuffer(registerData))
    if err != nil {
        t.Fatal(err)
    }
    
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status 200, got %d", resp.StatusCode)
    }
}
```

#### 运行测试
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/handlers

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 前端测试

#### 单元测试
```tsx
// components/UserProfile.test.tsx
import { render, screen } from '@testing-library/react';
import { UserProfile } from './UserProfile';

describe('UserProfile', () => {
    it('renders user information correctly', () => {
        const mockUser = {
            id: 1,
            username: 'testuser',
            email: 'test@example.com'
        };
        
        render(<UserProfile user={mockUser} />);
        
        expect(screen.getByText('testuser')).toBeInTheDocument();
        expect(screen.getByText('test@example.com')).toBeInTheDocument();
    });
});
```

#### 集成测试
```tsx
// tests/auth.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { LoginForm } from '../components/LoginForm';

describe('LoginForm Integration', () => {
    it('submits login form successfully', async () => {
        const mockOnLogin = jest.fn();
        
        render(<LoginForm onLogin={mockOnLogin} />);
        
        fireEvent.change(screen.getByLabelText('Email'), {
            target: { value: 'test@example.com' }
        });
        
        fireEvent.change(screen.getByLabelText('Password'), {
            target: { value: 'password123' }
        });
        
        fireEvent.click(screen.getByRole('button', { name: /login/i }));
        
        await waitFor(() => {
            expect(mockOnLogin).toHaveBeenCalledWith({
                email: 'test@example.com',
                password: 'password123'
            });
        });
    });
});
```

#### 运行测试
```bash
# 运行所有测试
npm test

# 运行测试并监听文件变化
npm run test:watch

# 运行测试并生成覆盖率报告
npm run test:coverage

# 运行 E2E 测试
npm run test:e2e
```

## 调试指南

### 后端调试

#### 使用 Delve 调试器
```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 启动调试
dlv debug cmd/main.go

# 常用调试命令
break main.go:20        # 设置断点
continue               # 继续执行
next                   # 下一步
step                   # 步入
print variable         # 打印变量
vars                   # 查看所有变量
```

#### 日志调试
```go
// 使用结构化日志
log.WithFields(log.Fields{
    "user_id": userID,
    "action":  "login",
    "ip":      clientIP,
}).Info("User login attempt")

// 使用不同日志级别
log.Debug("Debug information")
log.Info("Info message")
log.Warn("Warning message")
log.Error("Error message")
```

### 前端调试

#### 使用浏览器开发者工具
```javascript
// 在代码中添加断点
debugger;

// 使用 console 调试
console.log('Variable:', variable);
console.table(arrayData);
console.group('Group name');
console.groupEnd();
```

#### 使用 React DevTools
```bash
# 安装 React DevTools 浏览器扩展
# Chrome: React Developer Tools
# Firefox: React Developer Tools
```

#### 使用 Redux DevTools
```javascript
// 配置 Redux DevTools
const store = configureStore({
    reducer: rootReducer,
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware().concat(logger),
    devTools: process.env.NODE_ENV !== 'production',
});
```

## 性能优化

### 后端性能优化

#### 数据库优化
```go
// 使用连接池
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    PrepareStmt: true,
    Logger: logger.Default.LogMode(logger.Info),
})

sqlDB, err := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)

// 使用索引
// 在模型中添加索引标签
type Application struct {
    gorm.Model
    UserID uint `gorm:"index"`
    LabID  uint `gorm:"index"`
    Status string `gorm:"index"`
}
```

#### 缓存优化
```go
// 使用 Redis 缓存
func GetLabByID(id uint) (*Lab, error) {
    // 先从缓存获取
    cacheKey := fmt.Sprintf("lab:%d", id)
    if cached, err := redis.Get(ctx, cacheKey).Result(); err == nil {
        var lab Lab
        json.Unmarshal([]byte(cached), &lab)
        return &lab, nil
    }
    
    // 从数据库获取
    var lab Lab
    if err := db.First(&lab, id).Error; err != nil {
        return nil, err
    }
    
    // 存入缓存
    if data, err := json.Marshal(lab); err == nil {
        redis.Set(ctx, cacheKey, data, time.Hour)
    }
    
    return &lab, nil
}
```

### 前端性能优化

#### 代码分割
```tsx
// 路由级别的代码分割
const LabList = lazy(() => import('./pages/LabList'));
const LabDetail = lazy(() => import('./pages/LabDetail'));

// 组件级别的代码分割
const HeavyComponent = lazy(() => import('./components/HeavyComponent'));
```

#### 虚拟化长列表
```tsx
import { FixedSizeList as List } from 'react-window';

const VirtualizedList = ({ items }) => (
    <List
        height={400}
        itemCount={items.length}
        itemSize={50}
        itemData={items}
    >
        {({ index, style, data }) => (
            <div style={style}>
                {data[index].name}
            </div>
        )}
    </List>
);
```

#### 图片优化
```tsx
// 使用懒加载
import { LazyLoadImage } from 'react-lazy-load-image-component';

const LabCard = ({ lab }) => (
    <LazyLoadImage
        src={lab.coverImage}
        alt={lab.name}
        effect="blur"
        placeholderSrc="/placeholder.jpg"
    />
);
```

## 部署检查清单

### 后端部署检查
- [ ] 环境变量配置正确
- [ ] 数据库连接正常
- [ ] 日志目录权限正确
- [ ] 文件上传目录权限正确
- [ ] 防火墙配置正确
- [ ] SSL 证书配置正确
- [ ] 监控和告警配置正确

### 前端部署检查
- [ ] 环境变量配置正确
- [ ] API 地址配置正确
- [ ] 静态资源路径正确
- [ ] 缓存策略配置正确
- [ ] 错误页面配置正确
- [ ] 性能监控配置正确

## 常见问题

### 后端常见问题

1. **数据库连接失败**
   ```bash
   # 检查数据库服务状态
   sudo systemctl status mysql
   
   # 检查连接配置
   mysql -u username -p -h hostname database
   ```

2. **端口被占用**
   ```bash
   # 查看端口占用
   sudo netstat -tlnp | grep :8080
   
   # 杀死进程
   sudo kill -9 <PID>
   ```

3. **权限问题**
   ```bash
   # 修复文件权限
   sudo chown -R user:group /path/to/project
   sudo chmod -R 755 /path/to/project
   ```

### 前端常见问题

1. **依赖安装失败**
   ```bash
   # 清除缓存
   npm cache clean --force
   
   # 删除 node_modules 重新安装
   rm -rf node_modules package-lock.json
   npm install
   ```

2. **构建失败**
   ```bash
   # 检查 TypeScript 错误
   npm run type-check
   
   # 检查 ESLint 错误
   npm run lint
   ```

3. **API 请求失败**
   ```bash
   # 检查网络连接
   curl http://localhost:8080/api/health
   
   # 检查 CORS 配置
   # 确保后端允许前端域名访问
   ```

## 贡献指南

### 提交代码前检查
- [ ] 代码通过所有测试
- [ ] 代码符合编码规范
- [ ] 添加了必要的注释
- [ ] 更新了相关文档
- [ ] 提交信息符合规范

### 代码审查要点
- [ ] 功能实现是否正确
- [ ] 代码结构是否合理
- [ ] 是否有安全隐患
- [ ] 性能是否可接受
- [ ] 测试覆盖率是否足够

### 发布流程
1. 创建发布分支
2. 更新版本号
3. 更新更新日志
4. 创建发布标签
5. 部署到测试环境
6. 测试验证
7. 部署到生产环境 