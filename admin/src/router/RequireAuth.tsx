import { useEffect, useState, type ReactElement } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { Spin } from 'antd'
import { getStoredPermissions, refreshMyPermissions, setStoredPermissions } from '@/utils/permissions'

type Props = {
  children: ReactElement
}

export default function RequireAuth({ children }: Props) {
  const location = useLocation()
  const token = localStorage.getItem('token')
  const [ready, setReady] = useState(false)

  useEffect(() => {
    if (!token) return

    let active = true
    const stored = getStoredPermissions()
    if (stored) {
      setReady(true)
      return () => {
        active = false
      }
    }

    refreshMyPermissions()
      .catch(() => {
        setStoredPermissions({ isSuperAdmin: false, permissions: [] })
      })
      .finally(() => {
        if (active) setReady(true)
      })

    return () => {
      active = false
    }
  }, [token])

  if (!token) {
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  if (!ready) {
    return (
      <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <Spin />
      </div>
    )
  }

  return children
}
