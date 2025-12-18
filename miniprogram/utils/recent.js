import { STORAGE_KEYS } from './config'
import { getStorage, setStorage } from './storage'

function readAll() {
  return getStorage(STORAGE_KEYS.recent) || { doctors: [], departments: [] }
}

function writeAll(all) {
  setStorage(STORAGE_KEYS.recent, all)
}

function upsert(list, item, limit) {
  const id = item?.id
  if (!id) return list
  const next = [item, ...list.filter((x) => x?.id !== id)]
  return next.slice(0, limit)
}

export function addRecentDepartment(dept) {
  const all = readAll()
  all.departments = upsert(all.departments || [], { id: dept.id, name: dept.name }, 8)
  writeAll(all)
}

export function addRecentDoctor(doctor) {
  const all = readAll()
  all.doctors = upsert(
    all.doctors || [],
    { id: doctor.id, name: doctor.name, title: doctor.title, department_name: doctor.department_name },
    8
  )
  writeAll(all)
}

export function getRecent() {
  return readAll()
}

export function clearRecent() {
  writeAll({ doctors: [], departments: [] })
}

