import { Link, useParams } from 'react-router'

import { SessionList } from '@/components/SessionList'
import { Button } from '@/components/ui/button'

export function TaskDetail() {
  const { taskId } = useParams<{ taskId: string }>()

  if (!taskId) return <p>Missing taskId</p>

  return (
    <div className="mx-auto max-w-2xl p-8">
      <div className="mb-4 flex items-center gap-3">
        <Link to="/">
          <Button variant="ghost" size="sm">
            &larr; Back
          </Button>
        </Link>
        <h1 className="text-xl font-bold">Task: {taskId}</h1>
      </div>
      <SessionList taskId={taskId} />
    </div>
  )
}
