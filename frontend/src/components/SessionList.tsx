import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { fetchTask, patchSession, type Session } from '@/lib/api'

function StatusBadge({ status }: { status: Session['status'] }) {
  const isInactive = status === 'inactive'
  return (
    <span
      className={`inline-block rounded px-2 py-0.5 text-xs font-medium ${
        isInactive
          ? 'bg-gray-200 text-gray-500 dark:bg-gray-700 dark:text-gray-400'
          : 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300'
      }`}
    >
      {status}
    </span>
  )
}

function SessionRow({ session }: { session: Session }) {
  const queryClient = useQueryClient()
  const isInactive = session.status === 'inactive'

  const deactivate = useMutation({
    mutationFn: () => patchSession(session.session_id, 'inactive'),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['task', session.task_id] }),
  })

  return (
    <div className="flex items-center justify-between py-2">
      <div className="flex items-center gap-3">
        <StatusBadge status={session.status} />
        <span className="text-sm text-muted-foreground">{session.agent_type}</span>
        <span className="text-xs text-muted-foreground">
          {new Date(session.created_at).toLocaleString()}
        </span>
      </div>
      {!isInactive && (
        <Button
          variant="outline"
          size="sm"
          onClick={() => deactivate.mutate()}
          disabled={deactivate.isPending}
        >
          停用
        </Button>
      )}
    </div>
  )
}

export function SessionList({ taskId }: { taskId: string }) {
  const { data, isLoading, error } = useQuery({
    queryKey: ['task', taskId],
    queryFn: () => fetchTask(taskId),
  })

  if (isLoading) return <p className="text-muted-foreground">Loading sessions...</p>
  if (error) return <p className="text-destructive">Failed to load sessions</p>
  if (!data?.sessions?.length) return <p className="text-muted-foreground">No sessions</p>

  return (
    <Card>
      <CardHeader>
        <CardTitle>Sessions</CardTitle>
      </CardHeader>
      <CardContent className="divide-y">
        {data.sessions.map((s) => (
          <SessionRow key={s.session_id} session={s} />
        ))}
      </CardContent>
    </Card>
  )
}
