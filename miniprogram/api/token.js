import { http } from '../utils/request'

export function getIdempotentToken() {
  return http.get('/token/idempotent')
}

