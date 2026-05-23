const API_BASE = 'http://localhost:8080/api'

export interface Session {
  id: number
  session_id: string
  task_id: string
  agent_type: string
  status: string
  created_at: string
  updated_at: string
}

export interface TaskDetail {
  task: {
    id: number
    task_id: string
    title: string
    repo_path: string
    status: string
    created_at: string
    updated_at: string
  }
  sessions: Session[]
}

export async function fetchTasks(): Promise<{ task_id: string; title: string }[]> {
  const res = await fetch(`${API_BASE}/tasks`)
  const json = await res.json()
  return json.data
}

export async function fetchTask(taskId: string): Promise<TaskDetail> {
  const res = await fetch(`${API_BASE}/tasks/${taskId}`)
  const json = await res.json()
  return json.data
}

export async function patchSession(sessionId: string, status: string): Promise<void> {
  const res = await fetch(`${API_BASE}/sessions/${sessionId}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ status }),
  })
  if (!res.ok) {
    const json = await res.json()
    throw new Error(json.msg || 'Failed to update session')
  }
}
