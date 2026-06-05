import { useState } from 'react'

/** Hook for collapsible section state persisted to localStorage. */
export function useCollapsible(key: string, defaultOpen = true): [boolean, () => void] {
  const lsKey = `sidebar-collapse-${key}`
  const [open, setOpen] = useState(() => {
    try {
      const stored = localStorage.getItem(lsKey)
      return stored === null ? defaultOpen : stored !== 'true'
    } catch {
      return defaultOpen
    }
  })
  const toggle = () => {
    setOpen((prev) => {
      try {
        localStorage.setItem(lsKey, String(!prev))
      } catch {
        /* ignore */
      }
      return !prev
    })
  }
  return [open, toggle]
}
