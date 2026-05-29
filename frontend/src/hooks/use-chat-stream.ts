import { useCallback, useEffect, useRef } from 'react'

import type { StreamEvent } from '@/generated/events'
import { EventTypeValues } from '@/generated/events'
import type { AgentType } from '@/generated/request'
import { getTaskMessages, submitMessage } from '@/lib/api'
import { connectSSE } from '@/lib/sse'
import { type ChatMessage, useChatStore } from '@/stores/chat'

// Re-export ChatMessage for consumers
export type { ChatMessage }

export function useChatStream(
  taskId: string,
  sessionId: string,
  agentType: AgentType = 'claude-code',
) {
  const store = useChatStore()
  const abortRef = useRef<AbortController | null>(null)
  const session = store.getSession(sessionId)

  const connectToStream = useCallback(
    (messageId: string) => {
      abortRef.current?.abort()

      store.streamStart(sessionId, agentType)

      const controller = connectSSE({
        url: `/api/tasks/${taskId}/stream`,
        params: { session_id: sessionId, message_id: messageId },
        reconnect: true,
        onEvent: (event: StreamEvent) => {
          switch (event.type) {
            case EventTypeValues.Init:
              break
            case EventTypeValues.Text: {
              const textAgent = event.content?.agent as string | undefined
              const textAgentType = event.content?.agent_type as AgentType | undefined
              if (textAgent && textAgentType) {
                store.streamAgentUpdate(sessionId, textAgentType, textAgent)
              }
              store.streamText(sessionId, (event.content?.text as string) ?? '')
              break
            }
            case EventTypeValues.ToolCall:
              store.streamToolCall(sessionId, (event.content?.name as string) ?? 'unknown')
              break
            case EventTypeValues.ToolResult:
              store.streamToolResult(sessionId)
              break
            case EventTypeValues.Done:
              store.streamDone(sessionId)
              break
            case EventTypeValues.Error:
              store.streamError(
                sessionId,
                new Error(
                  (event.content?.error as string) ||
                    (event.content?.message as string) ||
                    'Unknown error',
                ),
              )
              break
            case EventTypeValues.RuntimeExecuting:
              store.streamRuntimeEvent(sessionId, {
                task_id: (event.content?.task_id as string) ?? '',
                agent: (event.content?.agent as string) ?? '',
                status: 'running',
              })
              break
            case EventTypeValues.RuntimeCompleted: {
              const success = event.content?.success ?? false
              store.streamRuntimeEvent(sessionId, {
                task_id: (event.content?.task_id as string) ?? '',
                agent: (event.content?.agent as string) ?? '',
                status: success ? 'completed' : 'failed',
              })
              break
            }
            case EventTypeValues.RuntimeText: {
              store.streamRuntimeText(sessionId, {
                task_id: (event.content?.task_id as string) ?? '',
                agent: (event.content?.agent as string) ?? '',
                text: (event.content?.text as string) ?? '',
              })
              break
            }
            case EventTypeValues.Planning: {
              const node = event.content?.node as string
              if (node === 'dispatch') {
                const dispatch = event.content?.dispatch as
                  | { task_id?: string; agent?: string; content?: string }
                  | undefined
                if (dispatch) {
                  store.streamPlanEvent(
                    sessionId,
                    [
                      {
                        task_id: dispatch.task_id ?? '',
                        agent: dispatch.agent ?? '',
                        title: (dispatch.content ?? '').slice(0, 80),
                        status: 'pending',
                      },
                    ],
                    '',
                  )
                }
              }
              break
            }
            case EventTypeValues.CoordinationStart:
              // coordination channel opens — no action needed, messages will follow
              break
            case EventTypeValues.CoordinationMessage:
              store.streamCoordinationEvent(sessionId, {
                from: (event.content?.from as string) ?? '',
                to: (event.content?.to as string) ?? '',
                text: (event.content?.text as string) ?? '',
                round: (event.content?.round as number) ?? 1,
              })
              break
            case EventTypeValues.CoordinationDone: {
              const decisions = event.content?.decisions as string[] | undefined
              store.streamCoordinationDone(sessionId, decisions?.join('\n') ?? '')
              break
            }
            default:
              break
          }
        },
        onError: (error) => {
          store.streamError(sessionId, error)
        },
      })

      abortRef.current = controller
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [taskId, sessionId],
  )

  const sendMessage = useCallback(
    async (message: string, agentType: AgentType = 'claude-code') => {
      const userMessage: ChatMessage = {
        id: `user-${Date.now()}`,
        role: 'user',
        content: message,
        timestamp: Date.now(),
      }

      store.sendMessage(sessionId, userMessage, {
        messageId: '',
        sessionId,
      })

      try {
        const result = await submitMessage(taskId, {
          message,
          session_id: sessionId,
          agent_type: agentType,
        })

        connectToStream(result.message_id)
      } catch (err) {
        store.streamError(sessionId, err instanceof Error ? err : new Error('发送失败'))
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [taskId, sessionId, connectToStream],
  )

  // Load history on mount; auto-reconnect if streaming message found
  useEffect(() => {
    let cancelled = false

    getTaskMessages(taskId, { limit: 20, sessionId })
      .then((res) => {
        if (cancelled || res.data.length === 0) return
        const chatMessages: ChatMessage[] = res.data.map((m) => ({
          id: `${m.role}-${m.id}`,
          dbId: m.id,
          role: m.role,
          content: m.content,
          agentType: m.agent_type as AgentType | undefined,
          agentName: m.agent_name || undefined,
          sessionId: m.session_id || undefined,
          timestamp: new Date(m.created_at).getTime(),
          messageId: m.message_id,
          status: m.status,
        }))
        store.loadHistory(sessionId, chatMessages, res.has_more)

        const streaming = res.data.find((m) => m.role === 'agent' && m.status === 'streaming')
        if (streaming && streaming.message_id) {
          connectToStream(streaming.message_id)
        }
      })
      .catch(() => {
        // silently ignore — empty state is fine
      })

    return () => {
      cancelled = true
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [taskId, sessionId, connectToStream])

  const abort = useCallback(() => {
    abortRef.current?.abort()
    abortRef.current = null
  }, [])

  return { state: session, sendMessage, abort }
}
