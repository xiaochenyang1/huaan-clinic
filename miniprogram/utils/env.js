import { API_BASE_URL } from './config'

export function getAppConfig() {
  const runtime = globalThis?.__APP_CONFIG__ || {}
  return {
    API_BASE_URL,
    ...runtime,
  }
}

