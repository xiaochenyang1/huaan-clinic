import { http } from '../utils/request'

export function createAppointment({ idempotent_token, schedule_id, patient_id, symptom }) {
  return http.post('/appointments', { idempotent_token, schedule_id, patient_id, symptom })
}

export function listAppointments({ status } = {}) {
  return http.get('/appointments', { params: { status } })
}

export function getAppointment(id) {
  return http.get(`/appointments/${id}`)
}

export function cancelAppointment(id, { reason }) {
  return http.put(`/appointments/${id}/cancel`, { reason })
}

export function checkinAppointment(id) {
  return http.post(`/appointments/${id}/checkin`, {})
}

