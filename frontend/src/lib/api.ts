import type { StreamEvent } from '@/generated/events'
import type { AgentType } from '@/generated/request'

export interface Task {
  task_id: string
  title: string
  repo_path: string
  status: string
  created_at: string
  updated_at: string
}

export interface Session {
  id: number
  session_id: string
  task_id: string
  agent_type: AgentType
  agent_name: string
  status: string
  created_at: string
  updated_at: string
}

export interface TaskDetail {
  task: Task
  sessions: Session[]
}

export interface AgentTypeInfo {
  type: AgentType
  name: string
  description: string
}

const API_BASE = '/api'

export async function fetchTasks(): Promise<Task[]> {
  const res = await fetch(`${API_BASE}/tasks`)
  const json = await res.json()
  return json.data
}

export async function fetchTask(taskId: string): Promise<TaskDetail> {
  const res = await fetch(`${API_BASE}/tasks/${taskId}`)
  const json = await res.json()
  return json.data
}

export async function createTask(
  title: string,
  agents?: { type: AgentType; name?: string }[],
): Promise<Task> {
  const res = await fetch(`${API_BASE}/tasks`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title, agents }),
  })
  const json = await res.json()
  return json.data
}

export async function deleteTask(taskId: string): Promise<void> {
  const res = await fetch(`${API_BASE}/tasks/${taskId}`, { method: 'DELETE' })
  if (!res.ok) {
    const json = await res.json()
    throw new Error(json.msg || 'Failed to delete task')
  }
}

export async function fetchAgentTypes(): Promise<AgentTypeInfo[]> {
  const res = await fetch(`${API_BASE}/agent-types`)
  const json = await res.json()
  const data: unknown[] = json.data
  return data.map((item) =>
    typeof item === 'string'
      ? { type: item as AgentType, name: item, description: '' }
      : (item as AgentTypeInfo),
  )
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

export type { StreamEvent }

// IM Conversation — a flattened view of Session across Tasks
export interface Conversation {
  taskId: string
  sessionId: string
  agentType: AgentType
  agentName: string
  title: string
  lastActiveAt: string
  taskTitle: string
  status: string
}

export async function fetchConversations(): Promise<Conversation[]> {
  const tasks = await fetchTasks()
  const details = await Promise.all(tasks.map((t) => fetchTask(t.task_id)))
  const convos: Conversation[] = []
  for (const detail of details) {
    for (const s of detail.sessions) {
      convos.push({
        taskId: s.task_id,
        sessionId: s.session_id,
        agentType: s.agent_type,
        agentName: s.agent_name ?? '',
        title: s.agent_name || s.agent_type,
        lastActiveAt: s.updated_at,
        taskTitle: detail.task.title,
        status: s.status,
      })
    }
  }
  convos.sort((a, b) => new Date(b.lastActiveAt).getTime() - new Date(a.lastActiveAt).getTime())
  return convos
}

export async function createConversation(
  agentType: AgentType,
  agentName?: string,
  title?: string,
): Promise<Conversation> {
  const taskTitle = title ?? `Chat with ${agentName || agentType}`
  const task = await createTask(taskTitle, [{ type: agentType, name: agentName }])
  const detail = await fetchTask(task.task_id)
  const session = detail.sessions[0]
  return {
    taskId: task.task_id,
    sessionId: session.session_id,
    agentType,
    agentName: agentName ?? '',
    title: agentName || agentType,
    lastActiveAt: session.updated_at,
    taskTitle: task.title,
    status: session.status,
  }
}
