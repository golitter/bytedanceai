import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { fetchTasks } from '@/lib/api'

export function TaskList() {
  const { data: tasks, isLoading } = useQuery({
    queryKey: ['tasks'],
    queryFn: fetchTasks,
  })

  return (
    <div className="mx-auto max-w-2xl p-8">
      <h1 className="mb-6 text-2xl font-bold">Tasks</h1>
      {isLoading && <p className="text-muted-foreground">Loading...</p>}
      {tasks?.map((t) => (
        <Link key={t.task_id} to={`/tasks/${t.task_id}`}>
          <Card className="mb-3 transition-colors hover:bg-muted/50">
            <CardHeader className="py-3">
              <CardTitle className="text-base">{t.title}</CardTitle>
            </CardHeader>
            <CardContent className="pb-3 pt-0">
              <p className="text-xs text-muted-foreground">{t.task_id}</p>
            </CardContent>
          </Card>
        </Link>
      ))}
    </div>
  )
}
