import { useQuery } from '@tanstack/react-query'
import { RefreshCw, Trash2 } from 'lucide-react'
import { useState } from 'react'

import { type Conversation, deleteAdminSessions, fetchConversations } from '@/lib/api'
import { cn } from '@/lib/utils'

const STATUS_COLORS: Record<string, { bg: string; text: string }> = {
  running: {
    bg: 'color-mix(in srgb, var(--color-success) 10%, transparent)',
    text: 'var(--color-success)',
  },
  streaming: {
    bg: 'color-mix(in srgb, var(--color-success) 10%, transparent)',
    text: 'var(--color-success)',
  },
  loading: {
    bg: 'color-mix(in srgb, var(--color-warning) 10%, transparent)',
    text: 'var(--color-warning)',
  },
  done: {
    bg: 'color-mix(in srgb, var(--color-brand) 10%, transparent)',
    text: 'var(--color-brand)',
  },
  error: {
    bg: 'color-mix(in srgb, var(--color-error) 10%, transparent)',
    text: 'var(--color-error)',
  },
  idle: { bg: 'var(--bg-hover)', text: 'var(--text-secondary)' },
}

export function SessionCleanupPage() {
  const {
    data: sessions,
    isLoading,
    refetch,
    isRefetching,
  } = useQuery<Conversation[]>({
    queryKey: ['admin-sessions'],
    queryFn: fetchConversations,
    staleTime: 30_000,
  })
  const [selected, setSelected] = useState<Set<string>>(new Set())
  const [deleting, setDeleting] = useState(false)
  const [filter, setFilter] = useState('')

  const toggleSelect = (id: string) => {
    const next = new Set(selected)
    if (next.has(id)) next.delete(id)
    else next.add(id)
    setSelected(next)
  }

  const handleDelete = async () => {
    if (selected.size === 0) return
    if (!confirm(`确认清理 ${selected.size} 个会话？此操作不可恢复。`)) return
    setDeleting(true)
    try {
      await deleteAdminSessions(Array.from(selected))
      setSelected(new Set())
      refetch()
    } catch {
      /* ignore */
    }
    setDeleting(false)
  }

  const allSessions = sessions ?? []
  const filtered = filter ? allSessions.filter((s) => s.agentType === filter) : allSessions
  const agentTypes = [...new Set(allSessions.map((s) => s.agentType))]

  return (
    <div className="p-6">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-foreground">会话清理</h2>
        <div className="flex items-center gap-2">
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            className="h-8 rounded-md border border-border bg-card px-2 text-[13px] text-text-secondary"
          >
            <option value="">全部类型</option>
            {agentTypes.map((t) => (
              <option key={t} value={t}>
                {t}
              </option>
            ))}
          </select>
          <button
            onClick={() => refetch()}
            disabled={isLoading}
            className="flex items-center gap-1.5 rounded-md border border-border px-3 py-1.5 text-[13px] text-text-secondary transition-[transform,opacity] hover:bg-hover"
          >
            <RefreshCw
              className={cn('h-3.5 w-3.5', isRefetching && 'animate-spin')}
              strokeWidth={1.25}
            />
            刷新
          </button>
        </div>
      </div>

      {selected.size > 0 && (
        <div className="mb-3 flex items-center gap-2">
          <span className="text-[13px] text-text-secondary">已选 {selected.size} 项</span>
          <button
            onClick={handleDelete}
            disabled={deleting}
            className="flex items-center gap-1 rounded-md bg-error px-3 py-1 text-[12px] text-primary-foreground transition-[transform,opacity]"
          >
            <Trash2 className="h-3 w-3" strokeWidth={1.25} />
            {deleting ? '清理中...' : '批量清理'}
          </button>
        </div>
      )}

      <div className="rounded-lg overflow-hidden border border-border">
        <table className="w-full text-[13px]">
          <thead>
            <tr className="border-b border-border bg-hover">
              <th className="px-3 py-2 text-left font-medium text-foreground">选择</th>
              <th className="px-3 py-2 text-left font-medium text-foreground">会话 ID</th>
              <th className="px-3 py-2 text-left font-medium text-foreground">Agent</th>
              <th className="px-3 py-2 text-left font-medium text-foreground">类型</th>
              <th className="px-3 py-2 text-left font-medium text-foreground">任务</th>
              <th className="px-3 py-2 text-left font-medium text-foreground">状态</th>
            </tr>
          </thead>
          <tbody>
            {filtered.map((s) => {
              const sc = STATUS_COLORS[s.status] ?? {
                bg: 'var(--bg-hover)',
                text: 'var(--text-secondary)',
              }
              return (
                <tr key={s.sessionId} className="border-b border-border">
                  <td className="px-3 py-2">
                    <input
                      type="checkbox"
                      checked={selected.has(s.sessionId)}
                      onChange={() => toggleSelect(s.sessionId)}
                      className="h-3.5 w-3.5 rounded"
                    />
                  </td>
                  <td className="px-3 py-2 font-mono text-xs text-text-secondary">
                    {s.sessionId.slice(0, 8)}
                  </td>
                  <td className="px-3 py-2 text-foreground">{s.agentName || s.agentType}</td>
                  <td className="px-3 py-2 text-text-secondary">{s.agentType}</td>
                  <td className="px-3 py-2 text-text-secondary">{s.taskTitle}</td>
                  <td className="px-3 py-2">
                    <span
                      className="inline-block rounded-full px-2 py-0.5 text-[11px]"
                      style={{ background: sc.bg, color: sc.text }}
                    >
                      {s.status}
                    </span>
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
        {filtered.length === 0 && (
          <div className="py-8 text-center text-[13px] text-tertiary">暂无会话</div>
        )}
      </div>
    </div>
  )
}
