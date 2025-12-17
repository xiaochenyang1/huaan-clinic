# 华安医疗后台管理系统

基于 Vite + React + TypeScript + Ant Design 构建的医疗预约系统后台管理界面。

## 技术栈

- **框架**: React 18
- **构建工具**: Vite 7
- **语言**: TypeScript
- **UI 组件库**: Ant Design 5
- **路由**: React Router v6
- **HTTP 客户端**: Axios
- **日期处理**: Day.js

## 项目结构

```
admin/
├── src/
│   ├── layouts/          # 布局组件
│   │   └── MainLayout.tsx
│   ├── pages/            # 页面组件
│   │   ├── Dashboard/    # 仪表盘
│   │   ├── Login/        # 登录页
│   │   ├── Department/   # 科室管理
│   │   ├── Doctor/       # 医生管理
│   │   ├── Appointment/  # 预约管理
│   │   └── Patient/      # 患者管理
│   ├── router/           # 路由配置
│   ├── styles/           # 全局样式
│   ├── types/            # TypeScript 类型定义
│   ├── utils/            # 工具函数
│   │   └── http.ts       # Axios 封装
│   ├── App.tsx           # 根组件
│   ├── main.tsx          # 入口文件
│   └── vite-env.d.ts     # Vite 类型声明
├── index.html            # HTML 模板
├── vite.config.ts        # Vite 配置
├── tsconfig.json         # TypeScript 配置
├── tsconfig.node.json    # TypeScript Node 配置
├── eslint.config.js      # ESLint 配置
└── package.json          # 项目依赖

```

## 功能模块

- **登录认证**: 用户登录与权限验证
- **仪表盘**: 数据概览与统计信息
- **科室管理**: 科室信息的增删改查
- **医生管理**: 医生信息的增删改查
- **预约管理**: 预约信息查看与状态管理
- **患者管理**: 患者信息查看

## 快速开始

### 安装依赖

```bash
pnpm install
```

### 开发模式

```bash
pnpm dev
```

项目将在 http://localhost:3000 启动

### 构建生产版本

```bash
pnpm build
```

### 预览生产构建

```bash
pnpm preview
```

### 代码检查

```bash
pnpm lint
```

## API 配置

后端 API 地址配置在 `vite.config.ts` 中：

```typescript
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:3030',  // 后端服务地址
      changeOrigin: true,
    },
  },
}
```

## 开发说明

### HTTP 请求

项目使用封装的 axios 实例 (`src/utils/http.ts`)，已配置：

- 请求拦截：自动添加 Authorization token
- 响应拦截：统一处理错误状态码
- 超时设置：10 秒

### 路由配置

路由配置在 `src/router/index.tsx`，使用 React Router v6 的钩子方式。

### 状态管理

当前使用 React 本地状态管理，如需全局状态管理可以引入 Redux 或 Zustand。

### 样式方案

- 全局样式：`src/styles/index.css`
- 组件样式：使用 Ant Design 的主题定制和内联样式

## License

MIT
