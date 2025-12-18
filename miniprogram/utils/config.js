function getEnvApiBaseUrl() {
  // 优先使用运行时注入（如 HBuilderX 运行参数/自定义全局变量）
  const runtime = globalThis?.__APP_CONFIG__?.API_BASE_URL
  if (runtime) return runtime

  // 小程序端可通过编译条件控制；默认开发用 localhost，生产建议改为你的正式域名
  // #ifdef MP-WEIXIN
  // 注意：真机/微信开发者工具访问 localhost 会指向设备自身，需要改为局域网IP或域名
  // #endif
  return 'http://localhost:8080/api'
}

export const API_BASE_URL = getEnvApiBaseUrl()

export const STORAGE_KEYS = {
  accessToken: 'access_token',
  refreshToken: 'refresh_token',
  user: 'user',
  subscribe: 'subscribe_state',
  recent: 'recent_views',
}

export const SUCCESS_CODE = 200000

// 订阅消息模板ID（需要你在微信公众平台配置）
export const WECHAT_SUBSCRIBE_TEMPLATE_IDS = {
  appointmentReminder: '',
}
