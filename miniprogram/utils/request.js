import { API_BASE_URL, SUCCESS_CODE } from './config'
import { ensureValidAccessToken, getRefreshToken, logout, setTokens } from './auth'
import { refreshToken as refreshTokenApi } from '../api/auth'

function buildUrl(path) {
  if (!path) return API_BASE_URL
  if (path.startsWith('http://') || path.startsWith('https://')) return path
  if (path.startsWith('/')) return `${API_BASE_URL}${path}`
  return `${API_BASE_URL}/${path}`
}

function requestRaw({ method, url, data, header }) {
  return new Promise((resolve, reject) => {
    uni.request({
      url,
      method,
      data,
      header,
      timeout: 10000,
      success: (res) => resolve(res),
      fail: (err) => reject(err),
    })
  })
}

async function tryRefreshToken() {
  const rt = getRefreshToken()
  if (!rt) return ''
  const data = await refreshTokenApi({ refresh_token: rt })
  setTokens({ accessToken: data?.access_token, refreshToken: data?.refresh_token })
  return data?.access_token || ''
}

export async function request({ method = 'GET', path, data, params, headers } = {}) {
  const token = await ensureValidAccessToken()
  const url = buildUrl(path)

  const header = {
    'Content-Type': 'application/json',
    ...(headers || {}),
  }
  if (token) header.Authorization = `Bearer ${token}`

  const payload = params ? { ...(data || {}), ...(params || {}) } : data

  try {
    let res = await requestRaw({ method, url, data: payload, header })
    const httpStatus = res.statusCode
    const body = res.data || {}

    if (httpStatus === 401) {
      // 尝试刷新一次并重试
      await tryRefreshToken()
      const newToken = await ensureValidAccessToken()
      if (!newToken) {
        await logout()
        throw new Error('未授权')
      }
      const retryHeader = { ...header, Authorization: `Bearer ${newToken}` }
      res = await requestRaw({ method, url, data: payload, header: retryHeader })
      if (res.statusCode === 401) {
        await logout()
        throw new Error('未授权')
      }
    }

    const finalBody = res.data || {}
    if (typeof finalBody === 'object' && finalBody && 'code' in finalBody) {
      if (finalBody.code === SUCCESS_CODE) return finalBody.data
      if (String(finalBody.code).startsWith('401')) {
        await logout()
        throw new Error(finalBody.message || '登录已过期')
      }
      throw new Error(finalBody.message || '请求失败')
    }

    return finalBody
  } catch (e) {
    const msg = e?.message || '网络错误'
    uni.showToast({ title: msg, icon: 'none' })
    throw e
  }
}

export const http = {
  get: (path, options) => request({ method: 'GET', path, ...(options || {}) }),
  post: (path, data, options) => request({ method: 'POST', path, data, ...(options || {}) }),
  put: (path, data, options) => request({ method: 'PUT', path, data, ...(options || {}) }),
  del: (path, options) => request({ method: 'DELETE', path, ...(options || {}) }),
}
