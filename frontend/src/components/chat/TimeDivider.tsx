import { formatRelativeTime } from '@/utils/time'

export function TimeDivider({ timestamp }: { timestamp: number }) {
  return (
    <div className="flex items-center gap-2 px-6 py-2">
      <div className="h-px flex-1 bg-border" />
      <span className="shrink-0 text-xs text-muted-foreground">
        {formatRelativeTime(timestamp)}
      </span>
      <div className="h-px flex-1 bg-border" />
    </div>
  )
}
