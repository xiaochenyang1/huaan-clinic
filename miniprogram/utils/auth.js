import { STORAGE_KEYS } from './config'
import { getStorage, removeStorage, setStorage } from './storage'
import { refreshToken as refreshTokenApi } from '../api/auth'

let refreshPromise = null

function base64UrlDecode(input) {
  const str = (input || '').replace(/-/g, '+').replace(/_/g, '/')
  const pad = str.length % 4
  const normalized = pad ? str + '='.repeat(4 - pad) : str
  try {
    return decodeURIComponent(
      atob(normalized)
        .split('')
        .map((c) => `%${`00${c.charCodeAt(0).toString(16)}`.slice(-2)}`)
        .join('')
    )
  } catch (e) {
    try {
      return atob(normalized)
    } catch (e2) {
      return ''
    }
  }
}

function parseJwtExp(token) {
  if (!token) return 0
  const parts = token.split('.')
  if (parts.length < 2) return 0
  try {
    const payload = JSON.parse(base64UrlDecode(parts[1]) || '{}')
    return Number(payload.exp || 0)
  } catch (e) {
    return 0
  }
}

export function getAccessToken() {
  return getStorage(STORAGE_KEYS.accessToken) || ''
}

export function getRefreshToken() {
  return getStorage(STORAGE_KEYS.refreshToken) || ''
}

export function setTokens({ accessToken, refreshToken }) {
  if (accessToken) setStorage(STORAGE_KEYS.accessToken, accessToken)
  if (refreshToken) setStorage(STORAGE_KEYS.refreshToken, refreshToken)
}

export function clearTokens() {
  removeStorage(STORAGE_KEYS.accessToken)
  removeStorage(STORAGE_KEYS.refreshToken)
}

export function setUser(user) {
  setStorage(STORAGE_KEYS.user, user || null)
}

export function getUser() {
  return getStorage(STORAGE_KEYS.user)
}

export function isLoggedIn() {
  return !!getAccessToken()
}

export async function ensureValidAccessToken() {
  const token = getAccessToken()
  if (token) {
    const exp = parseJwtExp(token)
    const now = Math.floor(Date.now() / 1000)
    // 还有 60 秒内过期则提前刷新
    if (!exp || exp - now > 60) return token
  }

  const rt = getRefreshToken()
  if (!rt) return ''

  if (!refreshPromise) {
    refreshPromise = (async () => {
      try {
        const data = await refreshTokenApi({ refresh_token: rt })
        setTokens({
          accessToken: data?.access_token,
          refreshToken: data?.refresh_token,
        })
        return getAccessToken()
      } finally {
        refreshPromise = null
      }
    })()
  }

  return refreshPromise
}

export function toLoginPage() {
  uni.reLaunch({ url: '/pages/login/index' })
}

export async function logout() {
  clearTokens()
  setUser(null)
  toLoginPage()
}

export async function bootstrapAuth() {
  await ensureValidAccessToken()
}
