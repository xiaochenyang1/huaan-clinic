import type { ReactElement } from 'react'
import { Navigate, useLocation } from 'react-router-dom'

type Props = {
  children: ReactElement
}

export default function RequireAuth({ children }: Props) {
  const location = useLocation()
  const token = localStorage.getItem('token')

  if (!token) {
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  return children
}

