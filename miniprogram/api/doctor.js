import { http } from '../utils/request'

export function listDoctors({ department_id, keyword } = {}) {
  return http.get('/doctors', { params: { department_id, keyword } })
}

export function getDoctor(id) {
  return http.get(`/doctors/${id}`)
}

