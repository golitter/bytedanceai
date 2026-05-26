/**
 * Format a timestamp as a relative time string.
 * Rules: today "HH:mm", yesterday "昨天 HH:mm", 2-7 days "N天前",
 * this year "M月D日 HH:mm", cross-year "YYYY年M月D日".
 */
export function formatRelativeTime(timestamp: number): string {
  const date = new Date(timestamp)
  const now = new Date()

  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const yesterday = new Date(today.getTime() - 86400000)
  const dateDay = new Date(date.getFullYear(), date.getMonth(), date.getDate())

  const hh = String(date.getHours()).padStart(2, '0')
  const mm = String(date.getMinutes()).padStart(2, '0')
  const time = `${hh}:${mm}`

  if (dateDay.getTime() === today.getTime()) return time
  if (dateDay.getTime() === yesterday.getTime()) return `昨天 ${time}`

  const diffDays = Math.floor((today.getTime() - dateDay.getTime()) / 86400000)
  if (diffDays >= 2 && diffDays <= 7) return `${diffDays}天前`

  const month = date.getMonth() + 1
  const day = date.getDate()
  if (date.getFullYear() === now.getFullYear()) return `${month}月${day}日 ${time}`

  return `${date.getFullYear()}年${month}月${day}日`
}

/**
 * Determine whether to show a time separator between two messages.
 * Triggers: first message (no prev), gap >5min, or different calendar day.
 */
export function shouldShowTimeSeparator(
  prevTimestamp: number | undefined,
  currentTimestamp: number,
): boolean {
  if (prevTimestamp === undefined) return true

  const diff = currentTimestamp - prevTimestamp
  if (diff > 5 * 60 * 1000) return true

  const prevDate = new Date(prevTimestamp)
  const curDate = new Date(currentTimestamp)
  if (
    prevDate.getFullYear() !== curDate.getFullYear() ||
    prevDate.getMonth() !== curDate.getMonth() ||
    prevDate.getDate() !== curDate.getDate()
  ) {
    return true
  }

  return false
}
