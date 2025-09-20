# 德州扑克在线游戏

一款基于Web的德州扑克游戏，支持多人在线对战，采用现代化技术栈构建。

## 项目特性

### 🎮 游戏功能
- **完整德州扑克规则**：标准德州扑克游戏流程
- **多房间支持**：低、中、高级别房间
- **实时对战**：WebSocket实现实时游戏同步
- **智能AI评估**：完整的手牌强度评估算法
- **移动端适配**：响应式设计，支持各种设备

### 🔐 用户系统
- **安全认证**：JWT Token + BCrypt密码加密
- **用户管理**：注册、登录、个人资料管理
- **筹码系统**：虚拟筹码，无真实货币交易
- **游戏统计**：胜率、总局数等数据统计

### 🛠 技术架构
- **后端**：Go + Gin + WebSocket + Goroutine Pool
- **前端**：Vue3 + Quasar + Pinia + TypeScript
- **数据库**：MySQL + Redis
- **部署**：Docker + Docker Compose + Nginx

## 快速开始

### 环境要求

- Docker 20.0+
- Docker Compose 2.0+
- Git

### 一键部署

```bash
# 1. 克隆项目
git clone <repository-url>
cd texas-poker

# 2. 配置环境变量
cp env.example .env
# 编辑 .env 文件，配置数据库密码等

# 3. 一键部署
chmod +x deploy.sh
./deploy.sh deploy
```

### 访问地址

部署完成后，可通过以下地址访问：

- **游戏主页**：http://localhost
- **后端API**：http://localhost/api
- **管理后台**：http://localhost/admin

## 开发指南

### 项目结构

```
texas-poker/
├── texas-poker-backend/     # Go后端
│   ├── cmd/                 # 主程序入口
│   ├── internal/            # 内部包
│   │   ├── handlers/        # HTTP处理器
│   │   ├── models/          # 数据模型
│   │   ├── database/        # 数据库连接
│   │   ├── cache/           # 缓存层
│   │   ├── websocket/       # WebSocket管理
│   │   ├── game/            # 游戏逻辑
│   │   └── utils/           # 工具函数
│   └── sql/                 # 数据库脚本
├── texas-poker-frontend/    # Vue3前端
│   ├── src/
│   │   ├── components/      # Vue组件
│   │   ├── pages/           # 页面组件
│   │   ├── stores/          # Pinia状态管理
│   │   ├── api/             # API接口
│   │   └── utils/           # 工具函数
│   └── public/              # 静态资源
├── nginx/                   # Nginx配置
├── docker-compose.yml       # Docker编排
└── deploy.sh               # 部署脚本
```

### 本地开发

#### 后端开发

```bash
cd texas-poker-backend

# 安装依赖
go mod download

# 启动数据库 (Docker)
docker-compose up -d mysql redis

# 配置环境变量
export DATABASE_URL="mysql://user:password@localhost:3306/texas_poker"
export REDIS_URL="redis://localhost:6379/0"
export JWT_SECRET="your-secret-key"

# 运行后端
go run cmd/main.go
```

#### 前端开发

```bash
cd texas-poker-frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

### 游戏规则

#### 德州扑克基本规则

1. **发牌**：每位玩家获得2张底牌
2. **翻牌前**：第一轮下注
3. **翻牌**：发出3张公共牌，第二轮下注
4. **转牌**：发出第4张公共牌，第三轮下注
5. **河牌**：发出第5张公共牌，第四轮下注
6. **摊牌**：比较手牌大小，最大者获胜

#### 手牌大小排序

1. 皇家同花顺 (Royal Flush)
2. 同花顺 (Straight Flush)
3. 四条 (Four of a Kind)
4. 葫芦 (Full House)
5. 同花 (Flush)
6. 顺子 (Straight)
7. 三条 (Three of a Kind)
8. 两对 (Two Pair)
9. 一对 (One Pair)
10. 高牌 (High Card)

### API文档

#### 认证接口

```http
POST /api/auth/register
POST /api/auth/login
GET  /api/auth/profile
PUT  /api/auth/profile
```

#### 房间接口

```http
GET  /api/rooms
POST /api/rooms
GET  /api/rooms/:id
```

#### WebSocket接口

```javascript
// 连接WebSocket
ws://localhost:8080/ws?token=<jwt_token>

// 消息格式
{
  "type": "player_action",
  "action": "call",
  "amount": 100
}
```

## 部署说明

### 生产环境部署

#### 1. 服务器要求

- **配置**：2核2GB内存以上
- **系统**：Linux (推荐Ubuntu 20.04+)
- **软件**：Docker, Docker Compose

#### 2. 域名配置

```bash
# 修改 .env 文件
DOMAIN=your-domain.com
API_DOMAIN=api.your-domain.com
VITE_API_BASE_URL=https://api.your-domain.com
```

#### 3. SSL证书配置

```bash
# 将证书文件放入 ssl/ 目录
ssl/
├── cert.pem
└── key.pem
```

#### 4. 启用HTTPS

编辑 `nginx/conf.d/default.conf`，取消HTTPS配置的注释。

### 监控和维护

#### 查看服务状态

```bash
./deploy.sh status
```

#### 查看日志

```bash
./deploy.sh logs
```

#### 数据备份

```bash
./deploy.sh backup
```

#### 健康检查

```bash
./deploy.sh health
```

## 性能优化

### 缓存策略

- **Redis缓存**：用户信息、房间状态、游戏数据
- **本地缓存**：热点数据本地缓存
- **CDN**：静态资源CDN加速

### 数据库优化

- **连接池**：MySQL连接池配置
- **索引优化**：关键字段索引
- **查询优化**：避免N+1查询

### WebSocket优化

- **连接池**：Goroutine连接池
- **心跳检测**：自动重连机制
- **消息队列**：异步消息处理

## 安全考虑

### 数据安全

- **密码加密**：BCrypt哈希算法
- **JWT Token**：有状态token认证
- **数据传输**：HTTPS加密传输

### 防护措施

- **SQL注入**：参数化查询
- **XSS攻击**：输入验证和输出编码
- **CSRF攻击**：Token验证
- **限流保护**：API访问频率限制

## 故障排除

### 常见问题

#### 1. 数据库连接失败

```bash
# 检查数据库状态
docker exec texas-poker-mysql mysqladmin ping

# 查看数据库日志
docker logs texas-poker-mysql
```

#### 2. Redis连接失败

```bash
# 检查Redis状态
docker exec texas-poker-redis redis-cli ping

# 查看Redis日志
docker logs texas-poker-redis
```

#### 3. WebSocket连接失败

检查防火墙设置，确保端口开放：

```bash
# 检查端口监听
netstat -tlnp | grep :8080
```

#### 4. 前端页面无法访问

检查Nginx配置：

```bash
# 测试Nginx配置
docker exec texas-poker-nginx nginx -t

# 重新加载配置
docker exec texas-poker-nginx nginx -s reload
```

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

- **项目地址**：[GitHub Repository](https://github.com/sh0TTvT/Texas-web)
- **问题反馈**：[Issues](https://github.com/sh0TTvT/Texas-web/issues)
- **邮箱**：shothollis@gmail.com

---

**注意**：本项目仅供学习和娱乐使用，不涉及真实货币交易。 