# 华安医疗预约系统 - 后端服务

基于 Go + Gin + GORM + MySQL + Redis 的医疗预约系统后端 API 服务。

## 技术栈

- **语言**: Go 1.21+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0+
- **缓存**: Redis 6.0+
- **认证**: JWT (双Token机制)
- **日志**: Zap + Lumberjack
- **配置**: Viper
- **API文档**: Swagger

## 项目结构

```
backend/
├── cmd/
│   └── main.go              # 程序入口
├── internal/
│   ├── handler/             # HTTP处理器 (Controller层)
│   ├── service/             # 业务逻辑层
│   ├── repository/          # 数据访问层
│   ├── model/               # 数据模型
│   ├── middleware/          # 中间件
│   └── router/              # 路由配置
├── pkg/
│   ├── config/              # 配置管理
│   ├── database/            # 数据库连接
│   ├── redis/               # Redis连接
│   ├── jwt/                 # JWT工具
│   ├── logger/              # 日志工具
│   ├── response/            # 统一响应格式
│   ├── errorcode/           # 错误码定义
│   └── utils/               # 工具函数
├── docs/                    # Swagger文档
├── config.yaml              # 配置文件
├── config.example.yaml      # 配置文件示例
├── go.mod
└── README.md
```

## 快速开始

### 1. 环境准备

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

### 2. 配置

复制配置文件并修改：

```bash
cp config.example.yaml config.yaml
```

编辑 `config.yaml`，配置数据库、Redis、微信等信息。

### 3. 安装依赖

```bash
go mod download
```

### 4. 创建数据库

```sql
CREATE DATABASE huaan_medical DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. 运行

```bash
go run cmd/main.go
```

服务将在 http://localhost:8080 启动。

### 6. 访问API文档

http://localhost:8080/swagger/index.html

## 开发

### 生成Swagger文档

```bash
# 安装swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g cmd/main.go -o docs
```

### 热重载开发

```bash
# 安装air
go install github.com/cosmtrek/air@latest

# 运行
air
```

## API接口

### 公开接口

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 微信登录 | POST | /api/user/login | 微信小程序登录 |
| 刷新Token | POST | /api/auth/refresh | 刷新访问令牌 |
| 科室列表 | GET | /api/departments | 获取所有科室 |
| 医生列表 | GET | /api/doctors | 获取医生列表 |
| 医生详情 | GET | /api/doctors/:id | 获取医生详情 |
| 排班查询 | GET | /api/schedule | 查询排班信息 |

### 用户接口 (需认证)

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 用户信息 | GET | /api/user/info | 获取当前用户信息 |
| 就诊人列表 | GET | /api/user/patients | 获取就诊人列表 |
| 创建预约 | POST | /api/appointments | 创建预约 |
| 取消预约 | PUT | /api/appointments/:id/cancel | 取消预约 |

### 管理接口 (需管理员认证)

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 管理员登录 | POST | /api/admin/login | 管理员登录 |
| 预约列表 | GET | /api/admin/appointments | 预约管理列表 |
| 科室管理 | CRUD | /api/admin/departments | 科室增删改查 |
| 医生管理 | CRUD | /api/admin/doctors | 医生增删改查 |
| 排班管理 | CRUD | /api/admin/schedules | 排班增删改查 |

## 响应格式

```json
{
  "code": 200000,
  "message": "success",
  "data": {}
}
```

## 错误码

| 错误码 | 说明 |
|--------|------|
| 200000 | 成功 |
| 400xxx | 客户端错误 |
| 401xxx | 认证错误 |
| 403xxx | 权限错误 |
| 404xxx | 资源不存在 |
| 4xxxx | 业务错误 |
| 500xxx | 服务端错误 |

详见 `pkg/errorcode/errorcode.go`

## 部署

### Docker部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o huaan-medical cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/huaan-medical .
COPY config.yaml .
EXPOSE 8080
CMD ["./huaan-medical"]
```

### 编译部署

```bash
# 编译
CGO_ENABLED=0 GOOS=linux go build -o huaan-medical cmd/main.go

# 运行
./huaan-medical
```
