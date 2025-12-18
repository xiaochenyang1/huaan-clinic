import { http } from '../utils/request'

export function listPatients() {
  return http.get('/user/patients')
}

export function getPatient(id) {
  return http.get(`/user/patients/${id}`)
}

export function createPatient(payload) {
  return http.post('/user/patients', payload)
}

export function updatePatient(id, payload) {
  return http.put(`/user/patients/${id}`, payload)
}

export function deletePatient(id) {
  return http.del(`/user/patients/${id}`)
}

