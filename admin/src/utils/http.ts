import axios, { type AxiosResponse } from 'axios'
import { message } from 'antd'

const http = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

function isBinaryResponse(response: AxiosResponse) {
  const rt = response.config?.responseType
  if (!rt) return false
  return rt === 'blob' || rt === 'arraybuffer'
}

// 请求拦截器
http.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
http.interceptors.response.use(
  async (response) => {
    // 导出/下载等二进制响应不走业务 code 校验
    if (isBinaryResponse(response)) {
      const contentType = String(response.headers?.['content-type'] || '')
      // 如果后端返回 JSON（通常是业务错误），尝试解析并提示
      if (contentType.includes('application/json') && response.data?.text) {
        try {
          const text = await response.data.text()
          const parsed = JSON.parse(text)
          const code = parsed?.code
          const msg = parsed?.message
          if (code && code !== 200000) {
            message.error(msg || '请求失败')
            return Promise.reject(new Error(msg || '请求失败'))
          }
        } catch {
          // ignore parse error, fallback to treating as binary
        }
      }
      return response
    }

    // 后端统一返回格式: { code: 200000, message: "success", data: {...} }
    const code = response.data?.code
    const msg = response.data?.message

    // 成功码为 200000
    if (code === 200000) {
      return response
    }

    // 业务错误
    message.error(msg || '请求失败')
    return Promise.reject(new Error(msg || '请求失败'))
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response
      switch (status) {
        case 401:
          message.error('未授权，请重新登录')
          localStorage.removeItem('token')
          window.location.href = '/login'
          break
        case 403:
          message.error('没有权限访问')
          break
        case 404:
          message.error('请求的资源不存在')
          break
        case 500:
          message.error('服务器错误')
          break
        default:
          message.error(data?.message || '请求失败')
      }
    } else {
      message.error('网络错误，请检查网络连接')
    }
    return Promise.reject(error)
  }
)

export default http
