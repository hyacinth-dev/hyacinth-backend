# Hyacinth Backend

Hyacinth Backend 是一个基于 Go 语言和 Nunu 框架构建的现代化 Web API 服务，提供虚拟网络管理、用户认证、流量统计等功能。项目采用 Clean Architecture 架构模式，具有良好的可维护性和可扩展性。

## 🚀 特性

- **用户管理**: 用户注册、登录、信息管理、套餐购买
- **虚拟网络**: 虚拟网络创建、配置、管理和监控
- **流量统计**: 实时流量监控和历史数据分析
- **权限控制**: 基于 JWT 的身份认证和授权
- **API 文档**: 集成 Swagger 自动生成 API 文档
- **数据库支持**: 支持 MySQL、PostgreSQL、SQLite
- **缓存**: Redis 缓存支持
- **定时任务**: 基于 Cron 的定时任务调度
- **日志**: 结构化日志记录
- **测试**: 完整的单元测试和集成测试

## 🏗️ 技术栈

- **框架**: [Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP Web 框架
- **ORM**: [GORM](https://gorm.io/) - 开发者友好的 ORM 库
- **数据库**: MySQL / PostgreSQL / SQLite
- **缓存**: [Redis](https://redis.io/)
- **认证**: [JWT](https://github.com/golang-jwt/jwt)
- **配置管理**: [Viper](https://github.com/spf13/viper)
- **日志**: [Zap](https://github.com/uber-go/zap)
- **API 文档**: [Swagger](https://github.com/swaggo/swag)
- **定时任务**: [Gocron](https://github.com/go-co-op/gocron)
- **依赖注入**: [Wire](https://github.com/google/wire)

## 📁 项目结构

```
├── api/                  # API 接口定义
│   └── v1/               # v1 版本 API
├── cmd/                  # 应用程序入口
│   ├── migration/        # 数据库迁移
│   ├── server/           # HTTP 服务器
│   └── task/             # 定时任务
├── config/               # 配置文件
├── docs/                 # Swagger 文档
├── internal/             # 内部代码
│   ├── handler/          # HTTP 处理器
│   ├── middleware/       # 中间件
│   ├── model/            # 数据模型
│   ├── repository/       # 数据访问层
│   └── service/          # 业务逻辑层
├── pkg/                  # 公共库
└── test/                 # 测试文件
```

## 🚀 快速开始

### 环境要求

- Go 1.19+
- MySQL 5.7+ / PostgreSQL 12+ / SQLite 3

### 安装

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd hyacinth-backend
   ```

2. **安装依赖**
   ```bash
   make init
   ```

3. **配置环境**
   ```bash
   # 编辑默认配置文件，填入数据库和 Redis 连接信息
   vim config/local.yml
   ```

4. **运行数据库迁移**
   ```bash
   go run cmd/migration/main.go
   ```

5. **启动服务**
   ```bash
   go run cmd/server/main.go
   ```

服务启动后，可以通过以下地址访问：
- API 服务: http://localhost:8000
- Swagger 文档: http://localhost:8000/swagger/index.html

## 📖 API 文档

项目集成了 Swagger 自动生成 API 文档，启动服务后访问 `/swagger/index.html` 查看完整的 API 文档。

### 主要 API 端点

- **用户相关**
  - `POST /register` - 用户注册
  - `POST /login` - 用户登录
  - `GET /user` - 获取用户信息
  - `PUT /user` - 更新用户信息
  - `PUT /user/password` - 修改密码

- **虚拟网络**
  - `GET /vnet` - 获取虚拟网络列表
  - `POST /vnet` - 创建虚拟网络
  - `PUT /vnet/{id}` - 更新虚拟网络
  - `DELETE /vnet/{id}` - 删除虚拟网络

- **流量统计**
  - `GET /usage` - 获取流量使用统计

## 🧪 测试

运行单元测试：
```bash
go test ./...
```

运行测试并生成覆盖率报告：
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 🔧 配置

配置文件位于 `config/` 目录，支持不同环境的配置：

- `local.yml` - 本地开发环境
- `prod.yml` - 生产环境

### 配置示例

```yaml
app:
  name: "hyacinth-backend"
  port: 8000
  mode: "debug"

database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  database: "hyacinth"
  username: "root"
  password: "password"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key"
  expire: 7200
```

## 🚀 部署

### Docker 部署

1. **构建镜像**
   ```bash
   docker build -t hyacinth-backend .
   ```

2. **使用 Docker Compose**
   ```bash
   cd deploy/docker-compose
   docker-compose up -d
   ```

### 传统部署

1. **编译应用**
   ```bash
   make build
   ```

2. **上传到服务器并运行**
   ```bash
   ./hyacinth-backend -conf config/prod.yml
   ```

## 📝 开发指南

### 添加新的 API 端点

1. 在 `api/v1/` 中定义请求和响应结构体
2. 在 `internal/handler/` 中实现处理器
3. 在 `internal/service/` 中实现业务逻辑
4. 在 `internal/repository/` 中实现数据访问
5. 更新路由配置

### 数据库迁移

使用 GORM 的 AutoMigrate 功能进行数据库迁移：

```go
// 在 cmd/migration/main.go 中添加新的模型
db.AutoMigrate(&model.NewModel{})
```

## 🤝 贡献

欢迎提交 Pull Request 或 Issue。在贡献代码前，请确保：

1. 代码通过所有测试
2. 遵循项目的代码规范
3. 添加必要的测试用例
4. 更新相关文档

## 📄 许可证

本项目采用 MIT 许可证。

---

**注意**: 这是一个后端 API 服务，需要配合前端应用使用。