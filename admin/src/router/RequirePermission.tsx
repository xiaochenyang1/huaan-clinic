import type { ReactElement } from 'react'
import Forbidden from '@/pages/Forbidden'
import { hasAnyPermission } from '@/utils/permissions'

type Props = {
  any: string[]
  children: ReactElement
}

export default function RequirePermission({ any, children }: Props) {
  if (!any || any.length === 0) return children
  if (!hasAnyPermission(any)) return <Forbidden />
  return children
}

