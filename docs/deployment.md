# 部署指南

## 环境要求

### 系统要求
- **操作系统**: Linux (Ubuntu 20.04+ / CentOS 8+) 或 Windows Server 2019+
- **内存**: 最少 2GB，推荐 4GB+
- **存储**: 最少 20GB 可用空间
- **网络**: 稳定的网络连接

### 软件要求
- **Go**: 1.21+
- **MySQL**: 8.0+
- **Node.js**: 18+ (仅前端构建需要)
- **Nginx**: 1.18+ (生产环境)
- **Docker**: 20.10+ (可选)

## 部署方式

### 方式一：传统部署

#### 1. 后端部署

##### 1.1 环境准备
```bash
# 安装 Go
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 MySQL
sudo apt update
sudo apt install mysql-server
sudo mysql_secure_installation

# 创建数据库和用户
sudo mysql -u root -p
CREATE DATABASE lab_recruitment CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'lab_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON lab_recruitment.* TO 'lab_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

##### 1.2 项目部署
```bash
# 克隆项目
git clone https://github.com/your-org/lab-recruitment-platform.git
cd lab-recruitment-platform

# 进入后端目录
cd backend

# 安装依赖
go mod tidy

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件
```

##### 1.3 环境变量配置
```bash
# .env 文件内容
DB_HOST=localhost
DB_PORT=3306
DB_USER=lab_user
DB_PASSWORD=your_password
DB_NAME=lab_recruitment

JWT_SECRET=your_jwt_secret_key
JWT_EXPIRE_HOURS=24

SERVER_PORT=8080
SERVER_MODE=release

LOG_LEVEL=info
LOG_FILE=logs/app.log

UPLOAD_PATH=uploads
MAX_FILE_SIZE=10485760
```

##### 1.4 数据库迁移
```bash
# 运行数据库迁移
mysql -u lab_user -p lab_recruitment < ../migrations/001_initial_schema.sql

# 或者使用 Go 迁移工具
go run cmd/migrate/main.go
```

##### 1.5 启动服务
```bash
# 开发模式
go run cmd/main.go

# 生产模式
go build -o bin/server cmd/main.go
./bin/server

# 使用 systemd 管理服务
sudo tee /etc/systemd/system/lab-recruitment.service << EOF
[Unit]
Description=Lab Recruitment Backend
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/your/project/backend
ExecStart=/path/to/your/project/backend/bin/server
Restart=always
RestartSec=5
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable lab-recruitment
sudo systemctl start lab-recruitment
```

#### 2. 前端部署

##### 2.1 构建前端
```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件，设置 API 地址
VITE_API_BASE_URL=http://your-domain.com/api

# 构建生产版本
npm run build
```

##### 2.2 Nginx 配置
```nginx
# /etc/nginx/sites-available/lab-recruitment
server {
    listen 80;
    server_name your-domain.com;
    
    # 前端静态文件
    location / {
        root /path/to/your/project/frontend/dist;
        try_files $uri $uri/ /index.html;
        
        # 缓存配置
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }
    
    # API 代理
    location /api/ {
        proxy_pass http://localhost:8080/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时配置
        proxy_connect_timeout 30s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }
    
    # 文件上传
    location /uploads/ {
        alias /path/to/your/project/backend/uploads/;
        expires 1d;
    }
}
```

##### 2.3 启用站点
```bash
sudo ln -s /etc/nginx/sites-available/lab-recruitment /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 方式二：Docker 部署

#### 1. Docker Compose 配置
```yaml
# docker-compose.yml
version: '3.8'

services:
  backend:
    build: 
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=lab_user
      - DB_PASSWORD=lab_password
      - DB_NAME=lab_recruitment
      - JWT_SECRET=your_jwt_secret_key
      - JWT_EXPIRE_HOURS=24
      - SERVER_MODE=release
      - LOG_LEVEL=info
    volumes:
      - ./uploads:/app/uploads
      - ./logs:/app/logs
    depends_on:
      - mysql
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "3000:80"
    environment:
      - VITE_API_BASE_URL=http://your-domain.com/api
    depends_on:
      - backend
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=root_password
      - MYSQL_DATABASE=lab_recruitment
      - MYSQL_USER=lab_user
      - MYSQL_PASSWORD=lab_password
    volumes:
      - mysql_data:/var/lib/mysql
      - ./migrations:/docker-entrypoint-initdb.d
    ports:
      - "3306:3306"
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - frontend
      - backend
    restart: unless-stopped

volumes:
  mysql_data:
```

#### 2. Dockerfile 配置

##### 后端 Dockerfile
```dockerfile
# backend/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 复制二进制文件
COPY --from=builder /app/main .

# 创建必要的目录
RUN mkdir -p uploads logs

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
```

##### 前端 Dockerfile
```dockerfile
# frontend/Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app

# 复制 package 文件
COPY package*.json ./

# 安装依赖
RUN npm ci --only=production

# 复制源码
COPY . .

# 构建应用
RUN npm run build

# 运行阶段
FROM nginx:alpine

# 复制构建结果
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

#### 3. 部署命令
```bash
# 构建并启动服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f backend

# 停止服务
docker-compose down

# 更新服务
docker-compose pull
docker-compose up -d
```

### 方式三：Kubernetes 部署

#### 1. 命名空间配置
```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: lab-recruitment
```

#### 2. 配置映射
```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: lab-recruitment-config
  namespace: lab-recruitment
data:
  DB_HOST: mysql-service
  DB_PORT: "3306"
  DB_NAME: lab_recruitment
  SERVER_MODE: release
  LOG_LEVEL: info
```

#### 3. 密钥配置
```yaml
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: lab-recruitment-secret
  namespace: lab-recruitment
type: Opaque
data:
  DB_USER: bGFiX3VzZXI=  # lab_user
  DB_PASSWORD: bGFiX3Bhc3N3b3Jk  # lab_password
  JWT_SECRET: eW91cl9qd3Rfc2VjcmV0X2tleQ==  # your_jwt_secret_key
```

#### 4. 数据库部署
```yaml
# k8s/mysql.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
  namespace: lab-recruitment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - name: mysql
        image: mysql:8.0
        env:
        - name: MYSQL_ROOT_PASSWORD
          value: "root_password"
        - name: MYSQL_DATABASE
          value: "lab_recruitment"
        - name: MYSQL_USER
          valueFrom:
            secretKeyRef:
              name: lab-recruitment-secret
              key: DB_USER
        - name: MYSQL_PASSWORD
          valueFrom:
            secretKeyRef:
              name: lab-recruitment-secret
              key: DB_PASSWORD
        ports:
        - containerPort: 3306
        volumeMounts:
        - name: mysql-data
          mountPath: /var/lib/mysql
      volumes:
      - name: mysql-data
        persistentVolumeClaim:
          claimName: mysql-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: mysql-service
  namespace: lab-recruitment
spec:
  selector:
    app: mysql
  ports:
  - port: 3306
    targetPort: 3306
  type: ClusterIP
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pvc
  namespace: lab-recruitment
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

#### 5. 后端部署
```yaml
# k8s/backend.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
  namespace: lab-recruitment
spec:
  replicas: 2
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: your-registry/lab-recruitment-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: lab-recruitment-config
              key: DB_HOST
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: lab-recruitment-secret
              key: DB_USER
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: lab-recruitment-secret
              key: DB_PASSWORD
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: lab-recruitment-secret
              key: JWT_SECRET
        volumeMounts:
        - name: uploads
          mountPath: /app/uploads
      volumes:
      - name: uploads
        persistentVolumeClaim:
          claimName: uploads-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: backend-service
  namespace: lab-recruitment
spec:
  selector:
    app: backend
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP
```

#### 6. 前端部署
```yaml
# k8s/frontend.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
  namespace: lab-recruitment
spec:
  replicas: 2
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      labels:
        app: frontend
    spec:
      containers:
      - name: frontend
        image: your-registry/lab-recruitment-frontend:latest
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: frontend-service
  namespace: lab-recruitment
spec:
  selector:
    app: frontend
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP
```

#### 7. Ingress 配置
```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: lab-recruitment-ingress
  namespace: lab-recruitment
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: your-domain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-service
            port:
              number: 80
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: backend-service
            port:
              number: 8080
```

## SSL 证书配置

### Let's Encrypt 自动证书
```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo crontab -e
# 添加以下行
0 12 * * * /usr/bin/certbot renew --quiet
```

### 手动配置 SSL
```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /path/to/your/certificate.crt;
    ssl_certificate_key /path/to/your/private.key;
    
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    
    # 其他配置...
}
```

## 监控和日志

### 日志配置
```bash
# 创建日志目录
sudo mkdir -p /var/log/lab-recruitment
sudo chown www-data:www-data /var/log/lab-recruitment

# 配置 logrotate
sudo tee /etc/logrotate.d/lab-recruitment << EOF
/var/log/lab-recruitment/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 www-data www-data
    postrotate
        systemctl reload lab-recruitment
    endscript
}
EOF
```

### 监控配置
```bash
# 安装 Prometheus
wget https://github.com/prometheus/prometheus/releases/download/v2.45.0/prometheus-2.45.0.linux-amd64.tar.gz
tar xvf prometheus-*.tar.gz
cd prometheus-*

# 配置 Prometheus
cat > prometheus.yml << EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'lab-recruitment'
    static_configs:
      - targets: ['localhost:8080']
EOF

# 启动 Prometheus
./prometheus --config.file=prometheus.yml
```

## 备份策略

### 数据库备份
```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/backup/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="lab_recruitment"

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份数据库
mysqldump -u lab_user -p'lab_password' $DB_NAME > $BACKUP_DIR/${DB_NAME}_${DATE}.sql

# 压缩备份文件
gzip $BACKUP_DIR/${DB_NAME}_${DATE}.sql

# 删除7天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete

echo "Database backup completed: ${DB_NAME}_${DATE}.sql.gz"
```

### 文件备份
```bash
#!/bin/bash
# backup_files.sh

BACKUP_DIR="/backup/files"
DATE=$(date +%Y%m%d_%H%M%S)
UPLOAD_DIR="/path/to/your/project/backend/uploads"

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份上传文件
tar -czf $BACKUP_DIR/uploads_${DATE}.tar.gz -C $(dirname $UPLOAD_DIR) $(basename $UPLOAD_DIR)

# 删除30天前的备份
find $BACKUP_DIR -name "uploads_*.tar.gz" -mtime +30 -delete

echo "Files backup completed: uploads_${DATE}.tar.gz"
```

### 自动备份脚本
```bash
# 添加到 crontab
0 2 * * * /path/to/backup.sh
0 3 * * * /path/to/backup_files.sh
```

## 故障排除

### 常见问题

1. **数据库连接失败**
   ```bash
   # 检查 MySQL 服务状态
   sudo systemctl status mysql
   
   # 检查数据库连接
   mysql -u lab_user -p -h localhost lab_recruitment
   
   # 检查防火墙
   sudo ufw status
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
   sudo chown -R www-data:www-data /path/to/your/project
   sudo chmod -R 755 /path/to/your/project
   ```

4. **日志查看**
   ```bash
   # 查看应用日志
   sudo journalctl -u lab-recruitment -f
   
   # 查看 Nginx 日志
   sudo tail -f /var/log/nginx/access.log
   sudo tail -f /var/log/nginx/error.log
   ```

### 性能优化

1. **数据库优化**
   ```sql
   -- 添加索引
   CREATE INDEX idx_applications_status_created ON applications(status, created_at);
   CREATE INDEX idx_notifications_user_read ON notifications(user_id, is_read);
   
   -- 优化查询
   EXPLAIN SELECT * FROM applications WHERE status = 'pending' ORDER BY created_at DESC;
   ```

2. **应用优化**
   ```bash
   # 启用 Gzip 压缩
   # 在 Nginx 配置中添加
   gzip on;
   gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
   ```

3. **缓存配置**
   ```bash
   # 安装 Redis
   sudo apt install redis-server
   
   # 配置应用使用 Redis 缓存
   # 在应用配置中添加 Redis 连接信息
   ```

## 安全配置

### 防火墙配置
```bash
# 配置 UFW
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### 安全头配置
```nginx
# 在 Nginx 配置中添加
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header X-Content-Type-Options "nosniff" always;
add_header Referrer-Policy "no-referrer-when-downgrade" always;
add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;
```

### 定期安全更新
```bash
# 设置自动更新
sudo apt update && sudo apt upgrade -y

# 检查安全更新
sudo unattended-upgrades --dry-run --debug
``` 