import { ChevronDown } from 'lucide-react'

import type { AgentType } from '@/generated/request'
import { AGENT_NAMES } from '@/lib/constants'
import { UI_LABELS } from '@/lib/ui-text'

import { AgentHoverCard } from './AgentHoverCard'
import { useCollapsible } from './useCollapsible'

/** Single-chat agent info section — mirrors MembersSection layout but for one agent. */
export function AgentInfoSection({
  agentType,
  agentName,
  avatarUrl,
  sessionId,
}: {
  agentType?: AgentType
  agentName?: string
  avatarUrl?: string
  sessionId: string
}) {
  const [open, toggleOpen] = useCollapsible('agent-info')

  const displayName = agentName ?? (agentType ? AGENT_NAMES[agentType] : 'Agent')
  const typeLabel = agentType ?? ''
  const online = false // TODO: wire to real session status

  return (
    <div className="border-b border-sidebar-border">
      {/* Header */}
      <button
        type="button"
        className="flex w-full items-center justify-between px-4 py-3 pb-2.5 text-left user-select-none"
        onClick={toggleOpen}
      >
        <span className="flex items-center gap-1.5 text-xs font-semibold uppercase tracking-wide text-text-secondary transition-[transform,opacity] hover:text-foreground">
          {UI_LABELS.AGENT_INFO}
        </span>
        <ChevronDown
          className={`h-3.5 w-3.5 text-tertiary transition-transform ${open ? '' : '-rotate-90'}`}
          strokeWidth={1.25}
        />
      </button>

      {/* Body */}
      <div
        className={`overflow-hidden transition-[max-height] duration-200 ease-out ${open ? 'max-h-[600px] overflow-y-auto' : 'max-h-0'}`}
      >
        <div className="px-4 pb-3.5">
          <div className="flex items-center gap-2.5 rounded-md px-2 py-1.5">
            <AgentHoverCard
              agentType={typeLabel}
              agentName={displayName}
              sessionId={sessionId}
              avatarUrl={avatarUrl}
              status={online ? 'running' : 'offline'}
            />
            <div className="min-w-0 flex-1">
              <div className="text-[13px] font-medium">{displayName}</div>
              <div className="text-[11px] text-tertiary">{typeLabel}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
