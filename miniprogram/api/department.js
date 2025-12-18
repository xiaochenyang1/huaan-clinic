import { http } from '../utils/request'

export function listDepartments() {
  return http.get('/departments')
}

