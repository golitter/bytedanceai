export interface PlanTask {
  task_id: string
  agent: string
  title: string
  status: 'pending' | 'running' | 'completed' | 'failed'
}

export interface CoordMessage {
  from: string
  to: string
  text: string
  round: number
}

export type MessageBlock =
  | { type: 'text'; id: string; content: string }
  | { type: 'html-render'; id: string; content: string }
  | { type: 'image'; id: string; path: string }
  | { type: 'attachment'; id: string; path: string }
  | { type: 'diff'; id: string; snapshotId: string }
  | { type: 'preview'; id: string; url: string }
  | { type: 'plan'; id: string; overview: string; tasks: PlanTask[] }
  | {
      type: 'runtime_status'
      id: string
      task_id: string
      agent: string
      status: string
      streamingText?: string
    }
  | {
      type: 'coordination'
      id: string
      messages: CoordMessage[]
      closed: boolean
      summary?: string
    }
  | { type: 'tool_call'; id: string; name: string; input?: string }
  | { type: 'tool_result'; id: string; output?: string }
