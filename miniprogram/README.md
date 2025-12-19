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

当前工程已补充 `package.json`（用于 Pinia 等依赖）。推荐使用 HBuilderX 运行：

1. 打开 HBuilderX -> 导入项目 -> 选择本目录 `miniprogram/`
2. 运行 -> 运行到小程序模拟器 -> 微信开发者工具

如需安装依赖（可选）：在 `miniprogram/` 目录执行 `pnpm install`（或在 HBuilderX 中启用 npm）。

## 后端地址配置

默认 API 地址在 `miniprogram/utils/config.js`：

- 默认值：`http://localhost:8080/api`
- 注意：真机/微信开发者工具访问 `localhost` 会指向设备自身，需要改为局域网 IP 或域名。

## 可交付运行所需配置（建议）

### 1) API_BASE_URL

将请求地址改为可访问的后端地址（局域网 IP 或域名）。你可以：

- 直接修改 `miniprogram/utils/config.js` 的默认值
- 或运行时注入 `globalThis.__APP_CONFIG__`：

```js
globalThis.__APP_CONFIG__ = {
  API_BASE_URL: 'http://192.168.1.100:8080/api',
}
```

### 2) 订阅消息模板ID（就诊提醒）

页面 `pages/user/subscribe` 需要配置模板ID，否则会提示“未配置模板ID”。可通过运行时注入：

```js
globalThis.__APP_CONFIG__ = {
  WECHAT_APPOINTMENT_REMINDER_TEMPLATE_ID: '你的模板ID',
}
```
