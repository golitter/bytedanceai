import type { AgentType } from '@/generated/request'
import { useChatStream } from '@/hooks/use-chat-stream'

import { AgentAvatar } from './AgentAvatar'
import { MessageInput } from './MessageInput'
import { MessageList } from './MessageList'

const AGENT_NAMES: Record<string, string> = {
  'claude-code': 'Claude Code',
  opencode: 'OpenCode',
  orchestrator: 'Orchestrator',
}

interface ChatAreaProps {
  taskId: string
  sessionId: string
  agentType?: AgentType
  agentName?: string
}

export function ChatArea({
  taskId,
  sessionId,
  agentType = 'claude-code',
  agentName,
}: ChatAreaProps) {
  const { state, sendMessage } = useChatStream(taskId)
  const isStreaming = ['loading', 'streaming', 'tool_running'].includes(state.status)

  const handleSend = (message: string) => {
    sendMessage(message, sessionId, agentType)
  }

  const displayName = agentName ?? AGENT_NAMES[agentType] ?? agentType

  return (
    <div className="flex h-full flex-col" style={{ backgroundColor: 'var(--bg-canvas)' }}>
      {/* Header — IM style: agent avatar + name */}
      <div
        className="flex h-12 shrink-0 items-center gap-3 border-b px-6"
        style={{ borderColor: 'var(--divider)' }}
      >
        <AgentAvatar agentType={agentType} status={isStreaming ? 'running' : 'ready'} size={28} />
        <div>
          <h2 className="text-sm font-medium" style={{ color: 'var(--text-primary)' }}>
            {displayName}
          </h2>
          {isStreaming && (
            <p className="text-[11px]" style={{ color: 'var(--color-brand)' }}>
              正在回复...
            </p>
          )}
        </div>
      </div>

      {/* Messages */}
      {state.messages.length === 0 && !isStreaming ? (
        <div className="flex flex-1 flex-col items-center justify-center gap-2">
          <AgentAvatar agentType={agentType} status="ready" size={48} />
          <p className="mt-2 text-sm font-medium" style={{ color: 'var(--text-primary)' }}>
            {displayName}
          </p>
          <p className="text-xs" style={{ color: 'var(--text-tertiary)' }}>
            发送消息开始对话
          </p>
        </div>
      ) : (
        <MessageList
          messages={state.messages}
          streamingContent={state.streamingContent}
          streamingAgentType={state.streamingAgentType}
          isStreaming={isStreaming}
        />
      )}

      {/* Input */}
      <MessageInput
        onSend={handleSend}
        disabled={isStreaming}
        placeholder={`发消息给 ${displayName}...`}
      />
    </div>
  )
}
