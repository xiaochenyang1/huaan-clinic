import { http } from '../utils/request'

export function listAvailableSchedules({ doctor_id, department_id, start_date, end_date }) {
  return http.get('/schedule/available', {
    params: { doctor_id, department_id, start_date, end_date },
  })
}

export function listDoctorSchedules({ doctor_id, start_date, end_date }) {
  return http.get('/schedule', { params: { doctor_id, start_date, end_date } })
}

