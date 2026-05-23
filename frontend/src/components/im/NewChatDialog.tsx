import { useQuery } from '@tanstack/react-query'
import { useRef, useState } from 'react'

import { AgentAvatar } from '@/components/chat/AgentAvatar'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import type { AgentType } from '@/generated/request'
import { useCreateConversation } from '@/hooks/use-conversations'
import { fetchAgentTypes } from '@/lib/api'
import { useChatNav } from '@/stores/chat'

const AGENT_DESCRIPTIONS: Record<string, string> = {
  'claude-code': 'Anthropic 的 AI 编程助手，擅长代码生成、重构和调试',
  opencode: '开源 AI 编程工具，支持多种模型',
  orchestrator: '多 Agent 协调器，自动分派任务给合适的 Agent',
}

interface NewChatDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function NewChatDialog({ open, onOpenChange }: NewChatDialogProps) {
  const { data: agentTypes } = useQuery({
    queryKey: ['agent-types'],
    queryFn: fetchAgentTypes,
  })
  const createMutation = useCreateConversation()
  const { setCurrentSession } = useChatNav()
  const [expandedAgent, setExpandedAgent] = useState<string | null>(null)
  const nameRef = useRef<HTMLInputElement>(null)

  const handleSelect = (agentType: AgentType, agentName?: string) => {
    createMutation.mutate(
      { agentType, agentName },
      {
        onSuccess: (conversation) => {
          setCurrentSession(conversation.sessionId)
          onOpenChange(false)
        },
      },
    )
  }

  const types = agentTypes?.length
    ? agentTypes
    : (['claude-code', 'opencode', 'orchestrator'] as AgentType[]).map((t) => ({
        type: t,
        name: t,
        description: AGENT_DESCRIPTIONS[t] ?? '',
      }))

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        className="max-w-sm"
        style={{
          backgroundColor: 'var(--card)',
          borderColor: 'rgba(255,255,255,0.06)',
        }}
      >
        <DialogHeader>
          <DialogTitle style={{ color: 'var(--text-primary)' }}>新建对话</DialogTitle>
        </DialogHeader>
        <p className="mb-3 text-xs" style={{ color: 'var(--text-tertiary)' }}>
          选择一个 Agent 开始对话
        </p>
        <div className="flex flex-col gap-2">
          {types.map((agent) => (
            <div key={agent.type}>
              <button
                className="flex w-full items-center gap-3 rounded-lg px-3 py-3 text-left transition-colors"
                style={{ border: '1px solid rgba(255,255,255,0.06)' }}
                onClick={() => {
                  if (expandedAgent === agent.type) {
                    handleSelect(agent.type as AgentType, nameRef.current?.value || undefined)
                  } else {
                    setExpandedAgent(agent.type)
                  }
                }}
                onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = 'var(--bg-hover)')}
                onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}
                disabled={createMutation.isPending}
              >
                <AgentAvatar agentType={agent.type as AgentType} status="ready" />
                <div className="min-w-0 flex-1">
                  <p className="text-sm font-medium" style={{ color: 'var(--text-primary)' }}>
                    {agent.name}
                  </p>
                  <p className="mt-0.5 truncate text-xs" style={{ color: 'var(--text-tertiary)' }}>
                    {agent.description || AGENT_DESCRIPTIONS[agent.type]}
                  </p>
                </div>
              </button>
              {expandedAgent === agent.type && (
                <div className="flex items-center gap-2 px-3 pt-1 pb-2">
                  <input
                    ref={nameRef}
                    placeholder="自定义名称（可选）"
                    className="flex-1 rounded-md border px-2 py-1.5 text-xs outline-none"
                    style={{
                      borderColor: 'rgba(255,255,255,0.1)',
                      backgroundColor: 'var(--bg-canvas)',
                      color: 'var(--text-primary)',
                    }}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        handleSelect(agent.type as AgentType, nameRef.current?.value || undefined)
                      }
                    }}
                  />
                  <button
                    className="rounded-md px-3 py-1.5 text-xs font-medium"
                    style={{ backgroundColor: 'var(--color-brand)', color: '#fff' }}
                    onClick={() =>
                      handleSelect(agent.type as AgentType, nameRef.current?.value || undefined)
                    }
                    disabled={createMutation.isPending}
                  >
                    开始
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      </DialogContent>
    </Dialog>
  )
}
