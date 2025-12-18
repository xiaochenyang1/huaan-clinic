export function pad2(n) {
  return String(n).padStart(2, '0')
}

export function formatYMD(d) {
  const date = d instanceof Date ? d : new Date(d)
  return `${date.getFullYear()}-${pad2(date.getMonth() + 1)}-${pad2(date.getDate())}`
}

export function nextDaysRange(days = 7) {
  const start = new Date()
  const end = new Date()
  end.setDate(end.getDate() + days)
  return { startDate: formatYMD(start), endDate: formatYMD(end) }
}

