import { http } from '../utils/request'

export function wechatLogin({ code }) {
  return http.post('/user/login', { code })
}

export function passwordLogin({ username, password }) {
  return http.post('/user/login/password', { username, password })
}

export function sendSmsCode({ phone }) {
  return http.post('/sms/send', { phone })
}

export function phoneLogin({ phone, code }) {
  return http.post('/user/login/phone', { phone, code })
}

export function refreshToken({ refresh_token }) {
  return http.post('/auth/refresh', { refresh_token })
}

export function getUserInfo() {
  return http.get('/user/info')
}

export function updateUserInfo(payload) {
  return http.put('/user/info', payload)
}
