#!/bin/bash
# 德州扑克游戏 - 自动化部署脚本
# 作用：一键部署应用到生产环境

set -e  # 遇到错误立即退出

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

# 检查环境
check_environment() {
    log_info "检查部署环境..."
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    
    # 检查Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    
    # 检查环境变量文件
    if [ ! -f ".env" ]; then
        log_warning ".env文件不存在，从env.example复制..."
        cp env.example .env
        log_warning "请编辑.env文件配置正确的环境变量"
        exit 1
    fi
    
    log_success "环境检查通过"
}

# 备份数据
backup_data() {
    log_info "备份现有数据..."
    
    # 创建备份目录
    BACKUP_DIR="./backups/$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$BACKUP_DIR"
    
    # 备份MySQL数据
    if docker ps | grep -q "texas-poker-mysql"; then
        log_info "备份MySQL数据..."
        docker exec texas-poker-mysql mysqldump -u root -p${MYSQL_ROOT_PASSWORD} texas_poker > "$BACKUP_DIR/mysql_backup.sql"
        log_success "MySQL数据备份完成"
    fi
    
    # 备份Redis数据
    if docker ps | grep -q "texas-poker-redis"; then
        log_info "备份Redis数据..."
        docker exec texas-poker-redis redis-cli --rdb "$BACKUP_DIR/redis_backup.rdb"
        log_success "Redis数据备份完成"
    fi
    
    log_success "数据备份完成: $BACKUP_DIR"
}

# 构建镜像
build_images() {
    log_info "构建应用镜像..."
    
    # 构建后端镜像
    log_info "构建后端镜像..."
    docker-compose build backend
    
    # 构建前端镜像
    log_info "构建前端镜像..."
    docker-compose build frontend
    
    log_success "镜像构建完成"
}

# 部署应用
deploy_application() {
    log_info "部署应用..."
    
    # 停止现有服务
    log_info "停止现有服务..."
    docker-compose down
    
    # 启动服务
    log_info "启动服务..."
    docker-compose up -d
    
    # 等待服务启动
    log_info "等待服务启动..."
    sleep 30
    
    log_success "应用部署完成"
}

# 健康检查
health_check() {
    log_info "执行健康检查..."
    
    # 检查后端服务
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        log_success "后端服务运行正常"
    else
        log_error "后端服务健康检查失败"
        return 1
    fi
    
    # 检查前端服务
    if curl -f http://localhost:3000/health > /dev/null 2>&1; then
        log_success "前端服务运行正常"
    else
        log_error "前端服务健康检查失败"
        return 1
    fi
    
    # 检查数据库连接
    if docker exec texas-poker-mysql mysqladmin ping -h localhost > /dev/null 2>&1; then
        log_success "MySQL数据库连接正常"
    else
        log_error "MySQL数据库连接失败"
        return 1
    fi
    
    # 检查Redis连接
    if docker exec texas-poker-redis redis-cli ping > /dev/null 2>&1; then
        log_success "Redis缓存连接正常"
    else
        log_error "Redis缓存连接失败"
        return 1
    fi
    
    log_success "所有服务健康检查通过"
}

# 初始化数据库
init_database() {
    log_info "初始化数据库..."
    
    # 等待MySQL启动
    log_info "等待MySQL启动..."
    while ! docker exec texas-poker-mysql mysqladmin ping -h localhost --silent; do
        sleep 2
    done
    
    # 执行数据库迁移
    log_info "执行数据库迁移..."
    docker exec texas-poker-mysql mysql -u root -p${MYSQL_ROOT_PASSWORD} texas_poker < texas-poker-backend/sql/schema.sql
    
    log_success "数据库初始化完成"
}

# 显示服务状态
show_status() {
    log_info "服务状态:"
    docker-compose ps
    
    echo ""
    log_info "访问地址:"
    echo "前端应用: http://localhost"
    echo "后端API: http://localhost/api"
    echo "管理后台: http://localhost/admin"
    echo "监控面板: http://localhost:3001 (如果启用)"
}

# 清理资源
cleanup() {
    log_info "清理未使用的资源..."
    docker system prune -f
    log_success "资源清理完成"
}

# 主函数
main() {
    echo "========================================"
    echo "     德州扑克游戏 - 自动化部署"
    echo "========================================"
    
    # 加载环境变量
    if [ -f ".env" ]; then
        set -a
        source .env
        set +a
    fi
    
    case "${1:-deploy}" in
        "check")
            check_environment
            ;;
        "backup")
            backup_data
            ;;
        "build")
            build_images
            ;;
        "deploy")
            check_environment
            backup_data
            build_images
            deploy_application
            init_database
            health_check
            show_status
            ;;
        "restart")
            log_info "重启服务..."
            docker-compose restart
            health_check
            show_status
            ;;
        "stop")
            log_info "停止服务..."
            docker-compose down
            log_success "服务已停止"
            ;;
        "logs")
            docker-compose logs -f
            ;;
        "status")
            show_status
            ;;
        "cleanup")
            cleanup
            ;;
        "init-db")
            init_database
            ;;
        "health")
            health_check
            ;;
        *)
            echo "使用方法: $0 {deploy|check|backup|build|restart|stop|logs|status|cleanup|init-db|health}"
            echo ""
            echo "命令说明:"
            echo "  deploy   - 完整部署流程"
            echo "  check    - 检查部署环境"
            echo "  backup   - 备份数据"
            echo "  build    - 构建镜像"
            echo "  restart  - 重启服务"
            echo "  stop     - 停止服务"
            echo "  logs     - 查看日志"
            echo "  status   - 查看状态"
            echo "  cleanup  - 清理资源"
            echo "  init-db  - 初始化数据库"
            echo "  health   - 健康检查"
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@" 