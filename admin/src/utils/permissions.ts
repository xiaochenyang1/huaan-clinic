import http from '@/utils/http'

export type StoredPermissions = {
  isSuperAdmin: boolean
  permissions: string[]
}

const STORAGE_KEY = 'permissions'

export function getStoredPermissions(): StoredPermissions | null {
  const raw = localStorage.getItem(STORAGE_KEY)
  if (!raw) return null
  try {
    const parsed = JSON.parse(raw) as StoredPermissions
    if (!parsed || !Array.isArray(parsed.permissions)) return null
    return {
      isSuperAdmin: Boolean(parsed.isSuperAdmin),
      permissions: parsed.permissions.filter((p) => typeof p === 'string'),
    }
  } catch {
    return null
  }
}

export function setStoredPermissions(data: StoredPermissions) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify({
    isSuperAdmin: Boolean(data.isSuperAdmin),
    permissions: Array.isArray(data.permissions) ? data.permissions : [],
  }))
}

export function clearStoredPermissions() {
  localStorage.removeItem(STORAGE_KEY)
}

export function hasPermission(code: string) {
  const stored = getStoredPermissions()
  if (!stored) return false
  if (stored.isSuperAdmin) return true
  return stored.permissions.includes(code)
}

export function hasAnyPermission(codes: string[]) {
  if (!codes || codes.length === 0) return true
  const stored = getStoredPermissions()
  if (!stored) return false
  if (stored.isSuperAdmin) return true
  const set = new Set(stored.permissions)
  return codes.some((c) => set.has(c))
}

export async function refreshMyPermissions(): Promise<StoredPermissions> {
  const response = await http.get('/admin/permissions/me')
  const data = response.data?.data
  const next: StoredPermissions = {
    isSuperAdmin: Boolean(data?.is_super_admin),
    permissions: Array.isArray(data?.permissions) ? data.permissions : [],
  }
  setStoredPermissions(next)
  return next
}

