import { useQuery } from '@tanstack/react-query'
import { RefreshCw, Trash2 } from 'lucide-react'

import { deleteAdminWorkspace, getAdminWorkspaces } from '@/lib/api'
import { UI_CONFIRMS, UI_MESSAGES } from '@/lib/ui-text'
import { cn } from '@/lib/utils'

export function WorkspacePage() {
  const { data, isLoading, refetch, isRefetching } = useQuery({
    queryKey: ['admin-workspaces'],
    queryFn: getAdminWorkspaces,
    staleTime: 30_000,
  })

  const workspaces = data?.workspaces ?? []
  const stats = data
    ? { total: data.total, active: data.active, cleaned: data.cleaned, totalDisk: data.totalDisk }
    : { total: 0, active: 0, cleaned: 0, totalDisk: 0 }

  const handleDelete = async (id: string) => {
    if (!confirm(UI_CONFIRMS.CLEAN_WORKSPACE)) return
    try {
      await deleteAdminWorkspace(id)
      refetch()
    } catch {
      /* ignore */
    }
  }

  return (
    <div className="p-6">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-foreground">工作区管理</h2>
        <button
          onClick={() => refetch()}
          disabled={isLoading}
          className="flex items-center gap-1.5 rounded-md border border-border px-3 py-1.5 text-[13px] text-text-secondary transition-[transform,opacity]"
          onMouseEnter={(e) => {
            e.currentTarget.style.background = 'var(--bg-hover)'
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.background = 'transparent'
          }}
        >
          <RefreshCw
            className={cn('h-3.5 w-3.5', isRefetching && 'animate-spin')}
            strokeWidth={1.25}
          />
          刷新
        </button>
      </div>

      <div className="mb-4 grid grid-cols-4 gap-3">
        {[
          { label: '总数', value: stats.total },
          { label: '活跃', value: stats.active },
          { label: '已清理', value: stats.cleaned },
          { label: '磁盘占用', value: `${stats.totalDisk.toFixed(1)} MB` },
        ].map((s) => (
          <div key={s.label} className="rounded-lg border border-border bg-card p-3 text-center">
            <div className="text-lg font-semibold text-foreground">{s.value}</div>
            <div className="text-[12px] text-tertiary">{s.label}</div>
          </div>
        ))}
      </div>

      <div className="rounded-lg overflow-hidden border border-border">
        <table className="w-full text-[13px]">
          <thead>
            <tr className="border-b border-border bg-hover">
              {['ID', '任务', 'Agent', '磁盘', '状态', '操作'].map((h) => (
                <th key={h} className="px-3 py-2 text-left font-medium text-text-secondary">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {workspaces.map((ws) => (
              <tr key={ws.id} className="border-b border-border">
                <td className="px-3 py-2 font-mono text-xs text-tertiary">{ws.id.slice(0, 8)}</td>
                <td className="px-3 py-2 text-text-secondary">{ws.task.slice(0, 8)}</td>
                <td className="px-3 py-2 text-foreground">{ws.agent}</td>
                <td className="px-3 py-2 text-text-secondary">{ws.disk_mb.toFixed(1)} MB</td>
                <td className="px-3 py-2">
                  <span
                    className={cn(
                      'inline-block rounded-full px-2 py-0.5 text-[11px]',
                      ws.status === 'active'
                        ? 'bg-success/10 text-success'
                        : 'bg-[var(--bg-hover)] text-[var(--text-tertiary)]',
                    )}
                  >
                    {ws.status}
                  </span>
                </td>
                <td className="px-3 py-2">
                  {ws.status === 'active' && (
                    <button onClick={() => handleDelete(ws.id)} className="text-error">
                      <Trash2 className="h-3.5 w-3.5" strokeWidth={1.25} />
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {workspaces.length === 0 && (
          <div className="py-8 text-center text-[13px] text-tertiary">
            {UI_MESSAGES.NO_WORKSPACES}
          </div>
        )}
      </div>
    </div>
  )
}
