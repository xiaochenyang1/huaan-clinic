# 华安医疗微信小程序（uni-app + Vue3）

患者端小程序，包含登录、预约挂号、就诊人管理、就诊记录、个人中心等页面，并已提供基础的请求封装与 Token 刷新逻辑。

## 目录结构

```
miniprogram/
├── api/                  # 后端 API 封装（按业务模块拆分）
├── pages/                # 页面（index/appointment/login/record/user/legal）
├── static/               # 静态资源
├── utils/                # request/auth/storage/config 等
├── App.vue
├── main.js
├── manifest.json
├── pages.json
└── uni.scss
```

## 运行方式

当前工程未提供 `package.json`（无 npm scripts），推荐使用 HBuilderX 运行：

1. 打开 HBuilderX -> 导入项目 -> 选择本目录 `miniprogram/`
2. 运行 -> 运行到小程序模拟器 -> 微信开发者工具

## 后端地址配置

默认 API 地址在 `miniprogram/utils/config.js`：

- 默认值：`http://localhost:8080/api`
- 注意：真机/微信开发者工具访问 `localhost` 会指向设备自身，需要改为局域网 IP 或域名。

