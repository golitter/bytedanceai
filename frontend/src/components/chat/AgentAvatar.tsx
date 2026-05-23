import type { AgentType } from '@/generated/request'

const AGENT_COLORS: Record<AgentType, string> = {
  'claude-code': 'var(--agent-claude)',
  opencode: 'var(--agent-opencode)',
  orchestrator: 'var(--agent-orchestrator)',
}

const AGENT_LABELS: Record<AgentType, string> = {
  'claude-code': 'Claude',
  opencode: 'OpenCode',
  orchestrator: 'Orchestrator',
}

type Status = 'ready' | 'running' | 'offline' | 'error'

const STATUS_COLORS: Record<Status, string> = {
  ready: '#22C55E',
  running: '#F59E0B',
  offline: '#5A6070',
  error: '#EF4444',
}

interface AgentAvatarProps {
  agentType: AgentType
  status?: Status
  size?: number
}

export function AgentAvatar({ agentType, status = 'offline', size = 32 }: AgentAvatarProps) {
  const color = AGENT_COLORS[agentType] ?? 'var(--color-brand)'
  const label = AGENT_LABELS[agentType] ?? agentType

  return (
    <div className="relative inline-flex shrink-0" title={label}>
      <div
        className="flex items-center justify-center rounded-lg text-xs font-semibold text-white"
        style={{
          width: size,
          height: size,
          backgroundColor: color,
          borderRadius: 8,
        }}
      >
        {label.charAt(0).toUpperCase()}
      </div>
      {status && (
        <span
          className="absolute -right-0.5 -bottom-0.5 block rounded-full border-2 border-[var(--bg-canvas)]"
          style={{
            width: 8,
            height: 8,
            backgroundColor: STATUS_COLORS[status],
            animation: status === 'running' ? 'pulse 1.5s ease-in-out infinite' : undefined,
          }}
        />
      )}
    </div>
  )
}
