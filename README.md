# 华安医疗预约系统

医疗预约管理平台，采用 Monorepo 架构，包含后端服务、微信小程序和后台管理系统。

## 项目结构

```
华安医疗/
├── backend/                 # 后端服务 (Go + Gin) ✅ 主要功能已完成
│   ├── cmd/                 # 程序入口
│   ├── internal/            # 内部代码（模型、路由、中间件）
│   ├── pkg/                 # 公共包（配置、工具、JWT等）
│   └── scripts/             # SQL脚本
├── miniprogram/             # 微信小程序 (uni-app + Vue3) 🔄 基础页面与请求封装已完成，联调中
│   ├── pages/               # 页面
│   ├── components/          # 组件
│   ├── api/                 # API封装
│   └── utils/               # 请求/鉴权/存储等工具
├── admin/                   # 后台管理系统 (React + Ant Design) 🔄 页面骨架与主要列表页已完成，联调中
│   └── src/                 # 源代码
├── 1.md                     # 需求规格说明书
├── 开发进度.md               # 开发进度跟踪
└── README.md
```

## 技术栈

| 子项目 | 技术栈 | 说明 |
|--------|--------|------|
| **backend** | Go + Gin + GORM + MySQL + Redis | RESTful API 服务 |
| **miniprogram** | uni-app + Vue3 + Pinia | 患者端微信小程序 |
| **admin** | React 18 + TypeScript + Ant Design + Vite | 管理员后台系统 |

## 各模块说明

### backend - 后端服务

提供 RESTful API 接口，处理业务逻辑、数据存储和微信认证。

**主要功能：**
- 用户认证（微信登录、JWT）
- 预约管理（号源、排班、冲突检测）
- 就诊人管理
- 医生/科室管理
- 数据统计

### miniprogram - 微信小程序

面向患者的移动端应用，提供预约挂号服务。

**主要功能：**
- 微信一键登录
- 在线预约挂号
- 就诊人管理
- 预约记录查询
- 就诊签到

### admin - 后台管理系统

面向医院管理员的 Web 管理平台。

**主要功能：**
- 预约管理
- 患者管理
- 医生/科室管理
- 排班管理
- 数据统计报表
- 系统配置

## 开发指南

### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+

### 快速开始

```bash
# 后端
cd backend
go mod download
go run cmd/main.go

# 后台管理
cd admin
pnpm install
pnpm dev

# 小程序（uni-app 项目：当前工程未提供 package.json，推荐用 HBuilderX 运行）
# 1) 打开 HBuilderX -> 导入项目 -> 选择 miniprogram 目录
# 2) 运行 -> 运行到小程序模拟器 -> 微信开发者工具
```

## 开发规范

- 后端遵循 Go 标准项目布局
- 前端遵循 ESLint + Prettier 代码规范
- Git 提交信息遵循 Conventional Commits
- API 接口遵循 RESTful 设计规范

## 文档

- [需求规格说明书](./1.md)
- [开发进度跟踪](./开发进度.md)
- 子项目文档：`backend/README.md`、`admin/README.md`、`miniprogram/README.md`

## 部署

各子项目独立部署：

| 模块 | 部署方式 |
|------|----------|
| backend | Docker / 二进制部署 |
| miniprogram | 微信开发者平台上传 |
| admin | Nginx 静态资源部署 |
