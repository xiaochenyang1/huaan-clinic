import { http } from '../utils/request'

export function listRecords() {
  return http.get('/records')
}

export function getRecord(id) {
  return http.get(`/records/${id}`)
}

