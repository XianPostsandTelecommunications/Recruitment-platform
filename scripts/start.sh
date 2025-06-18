#!/bin/bash

# 实验室招新平台启动脚本
# 作者: Lab Recruitment Team
# 版本: 1.0.0

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "检查系统依赖..."
    
    # 检查 Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    # 检查 Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    log_success "依赖检查完成"
}

# 创建必要的目录
create_directories() {
    log_info "创建必要的目录..."
    
    mkdir -p uploads
    mkdir -p logs
    mkdir -p ssl
    mkdir -p monitoring/grafana/provisioning
    
    log_success "目录创建完成"
}

# 生成自签名 SSL 证书 (开发环境)
generate_ssl_cert() {
    if [ ! -f "ssl/cert.pem" ] || [ ! -f "ssl/key.pem" ]; then
        log_info "生成自签名 SSL 证书..."
        
        openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
            -keyout ssl/key.pem \
            -out ssl/cert.pem \
            -subj "/C=CN/ST=Beijing/L=Beijing/O=LabRecruitment/OU=IT/CN=localhost"
        
        log_success "SSL 证书生成完成"
    else
        log_info "SSL 证书已存在，跳过生成"
    fi
}

# 设置环境变量
setup_environment() {
    log_info "设置环境变量..."
    
    # 检查 .env 文件
    if [ ! -f ".env" ]; then
        log_warning ".env 文件不存在，创建默认配置"
        cat > .env << EOF
# 数据库配置
DB_HOST=mysql
DB_PORT=3306
DB_USER=lab_user
DB_PASSWORD=lab_password
DB_NAME=lab_recruitment

# Redis 配置
REDIS_HOST=redis
REDIS_PORT=6379

# JWT 配置
JWT_SECRET=your_jwt_secret_key_change_in_production
JWT_EXPIRE_HOURS=24

# 服务器配置
SERVER_PORT=8080
SERVER_MODE=release
LOG_LEVEL=info

# 文件上传配置
UPLOAD_PATH=uploads
MAX_FILE_SIZE=10485760

# 前端配置
VITE_API_BASE_URL=http://localhost/api
VITE_APP_TITLE=实验室招新平台
EOF
        log_success "默认 .env 文件创建完成"
    fi
    
    # 加载环境变量
    export $(cat .env | grep -v '^#' | xargs)
}

# 启动服务
start_services() {
    log_info "启动服务..."
    
    # 构建并启动服务
    docker-compose up -d --build
    
    log_success "服务启动完成"
}

# 等待服务就绪
wait_for_services() {
    log_info "等待服务就绪..."
    
    # 等待 MySQL 就绪
    log_info "等待 MySQL 就绪..."
    timeout=60
    while ! docker-compose exec -T mysql mysqladmin ping -h"localhost" --silent; do
        if [ $timeout -le 0 ]; then
            log_error "MySQL 启动超时"
            exit 1
        fi
        sleep 1
        timeout=$((timeout - 1))
    done
    log_success "MySQL 就绪"
    
    # 等待后端服务就绪
    log_info "等待后端服务就绪..."
    timeout=60
    while ! curl -f http://localhost:8080/api/health &> /dev/null; do
        if [ $timeout -le 0 ]; then
            log_error "后端服务启动超时"
            exit 1
        fi
        sleep 2
        timeout=$((timeout - 2))
    done
    log_success "后端服务就绪"
    
    # 等待前端服务就绪
    log_info "等待前端服务就绪..."
    timeout=60
    while ! curl -f http://localhost:3000 &> /dev/null; do
        if [ $timeout -le 0 ]; then
            log_error "前端服务启动超时"
            exit 1
        fi
        sleep 2
        timeout=$((timeout - 2))
    done
    log_success "前端服务就绪"
}

# 显示服务状态
show_status() {
    log_info "服务状态:"
    docker-compose ps
    
    echo ""
    log_success "实验室招新平台启动完成！"
    echo ""
    echo "访问地址:"
    echo "  - 前端应用: http://localhost"
    echo "  - API 文档: http://localhost/api/docs"
    echo "  - 健康检查: http://localhost/health"
    echo ""
    echo "监控面板:"
    echo "  - Prometheus: http://localhost:9090"
    echo "  - Grafana: http://localhost:3001 (admin/admin)"
    echo ""
    echo "数据库:"
    echo "  - MySQL: localhost:3306"
    echo "  - Redis: localhost:6379"
    echo ""
    echo "管理命令:"
    echo "  - 查看日志: docker-compose logs -f"
    echo "  - 停止服务: docker-compose down"
    echo "  - 重启服务: docker-compose restart"
}

# 主函数
main() {
    echo "=========================================="
    echo "    实验室招新平台启动脚本"
    echo "=========================================="
    echo ""
    
    # 检查是否在项目根目录
    if [ ! -f "docker-compose.yml" ]; then
        log_error "请在项目根目录运行此脚本"
        exit 1
    fi
    
    # 执行启动步骤
    check_dependencies
    create_directories
    generate_ssl_cert
    setup_environment
    start_services
    wait_for_services
    show_status
}

# 处理命令行参数
case "${1:-}" in
    "stop")
        log_info "停止服务..."
        docker-compose down
        log_success "服务已停止"
        ;;
    "restart")
        log_info "重启服务..."
        docker-compose down
        docker-compose up -d
        log_success "服务已重启"
        ;;
    "logs")
        docker-compose logs -f
        ;;
    "status")
        docker-compose ps
        ;;
    "clean")
        log_warning "清理所有数据..."
        docker-compose down -v
        rm -rf uploads logs ssl
        log_success "清理完成"
        ;;
    "help"|"-h"|"--help")
        echo "用法: $0 [命令]"
        echo ""
        echo "命令:"
        echo "  start   启动服务 (默认)"
        echo "  stop    停止服务"
        echo "  restart 重启服务"
        echo "  logs    查看日志"
        echo "  status  查看状态"
        echo "  clean   清理数据"
        echo "  help    显示帮助"
        ;;
    *)
        main
        ;;
esac 