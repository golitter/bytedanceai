import { useQuery } from '@tanstack/react-query'
import { useState } from 'react'

import { AgentAvatar } from '@/components/chat/AgentAvatar'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import type { AgentType } from '@/generated/request'
import { useCreateConversation } from '@/hooks/use-conversations'
import { useHoverStyle } from '@/hooks/use-hover-style'
import { fetchAgentTypes, validateRepoPath } from '@/lib/api'
import { AGENT_DESCRIPTIONS } from '@/lib/constants'
import { useChatNav } from '@/stores/chat'

interface AgentEntry {
  type: AgentType
  name: string
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
  const agentHover = useHoverStyle()

  const [repoPath, setRepoPath] = useState('')
  const [repoPathValidated, setRepoPathValidated] = useState(false)
  const [repoPathError, setRepoPathError] = useState<string | null>(null)
  const [validating, setValidating] = useState(false)
  const [groupTitle, setGroupTitle] = useState('')
  const [groupTitleError, setGroupTitleError] = useState(false)

  const [agents, setAgents] = useState<AgentEntry[]>([])
  const [addingType, setAddingType] = useState<AgentType | null>(null)
  const [inputName, setInputName] = useState('')
  const [nameError, setNameError] = useState(false)

  const [prevOpen, setPrevOpen] = useState(open)
  if (prevOpen !== open) {
    setPrevOpen(open)
    if (open) {
      setRepoPath('')
      setRepoPathValidated(false)
      setRepoPathError(null)
      setValidating(false)
      setAgents([])
      setAddingType(null)
      setInputName('')
      setNameError(false)
      setGroupTitle('')
      setGroupTitleError(false)
    }
  }

  const handleValidate = async () => {
    const path = repoPath.trim()
    if (!path) {
      setRepoPathError('请输入仓库路径')
      setRepoPathValidated(false)
      return
    }
    setValidating(true)
    setRepoPathError(null)
    try {
      const result = await validateRepoPath(path)
      if (result.valid) {
        setRepoPathValidated(true)
        setRepoPathError(null)
      } else {
        setRepoPathValidated(false)
        setRepoPathError(result.errors.join('; '))
      }
    } catch {
      setRepoPathValidated(false)
      setRepoPathError('校验失败，请检查 Agent 服务是否可用')
    } finally {
      setValidating(false)
    }
  }

  const handleAddAgent = () => {
    const trimmed = inputName.trim()
    if (!trimmed) {
      setNameError(true)
      return
    }
    if (agents.some((a) => a.name === trimmed)) {
      setNameError(true)
      return
    }
    setAgents((prev) => [...prev, { type: addingType!, name: trimmed }])
    setInputName('')
    setAddingType(null)
    setNameError(false)
  }

  const handleRemoveAgent = (index: number) => {
    setAgents((prev) => prev.filter((_, i) => i !== index))
  }

  const handleSubmit = () => {
    if (agents.length === 0) return
    if (agents.length >= 2 && !groupTitle.trim()) {
      setGroupTitleError(true)
      return
    }
    createMutation.mutate(
      {
        agents,
        repoPath: repoPath.trim(),
        title: agents.length >= 2 ? groupTitle.trim() : undefined,
      },
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
    : (['claude-code', 'opencode', 'orchestrator', 'codex'] as AgentType[]).map((t) => ({
        type: t,
        name: t,
        description: AGENT_DESCRIPTIONS[t] ?? '',
      }))

  const canSubmit = agents.length > 0 && repoPathValidated && !createMutation.isPending

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-sm bg-card border-border">
        <DialogHeader>
          <DialogTitle className="text-foreground">新建对话</DialogTitle>
        </DialogHeader>

        {/* Repo path */}
        <div className="mb-3">
          <label className="mb-1 block text-xs font-medium text-muted-foreground">仓库路径</label>
          <div className="flex items-center gap-2">
            <input
              value={repoPath}
              placeholder="/path/to/repo"
              className="flex-1 rounded-md border bg-background px-2 py-1.5 text-xs text-foreground outline-none"
              style={{
                borderColor: repoPathError
                  ? 'var(--destructive)'
                  : repoPathValidated
                    ? 'var(--color-success)'
                    : 'var(--border)',
              }}
              onChange={(e) => {
                setRepoPath(e.target.value)
                setRepoPathValidated(false)
                setRepoPathError(null)
              }}
              onKeyDown={(e) => {
                if (e.key === 'Enter') handleValidate()
              }}
            />
            <button
              className="shrink-0 rounded-md bg-primary px-3 py-1.5 text-xs font-medium text-primary-foreground"
              style={{ opacity: validating ? 0.6 : 1 }}
              onClick={handleValidate}
              disabled={validating}
            >
              {validating ? '校验中...' : '校验'}
            </button>
          </div>
          {repoPathError && <p className="mt-1 text-xs text-destructive">{repoPathError}</p>}
          {repoPathValidated && (
            <p className="mt-1 text-xs" style={{ color: 'var(--color-success)' }}>
              路径校验通过
            </p>
          )}
        </div>

        {/* Added agents list */}
        {agents.length > 0 && (
          <div className="mb-3 flex flex-col gap-1.5">
            <p className="text-xs font-medium text-muted-foreground">
              已选 Agent（{agents.length}）
            </p>
            {agents.map((agent, i) => (
              <div
                key={i}
                className="flex items-center gap-2 rounded-lg border border-border bg-background px-2.5 py-2"
              >
                <AgentAvatar agentType={agent.type} status="ready" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs font-medium text-foreground truncate">{agent.name}</p>
                  <p className="text-[11px] text-muted-foreground">{agent.type}</p>
                </div>
                <button
                  className="text-xs text-muted-foreground hover:text-destructive"
                  onClick={() => handleRemoveAgent(i)}
                >
                  ✕
                </button>
              </div>
            ))}
          </div>
        )}

        {/* Add agent section */}
        {repoPathValidated && (
          <div className="mb-3">
            {addingType ? (
              <div className="flex items-center gap-2 rounded-lg border border-primary-border bg-primary-soft px-3 py-2">
                <AgentAvatar agentType={addingType} status="ready" />
                <input
                  value={inputName}
                  placeholder="输入 Agent 名称"
                  className="flex-1 rounded-md border bg-background px-2 py-1.5 text-xs text-foreground outline-none"
                  style={{
                    borderColor: nameError ? 'var(--destructive)' : 'var(--border)',
                    animation: nameError ? 'shake 0.4s ease' : undefined,
                  }}
                  onChange={(e) => {
                    setInputName(e.target.value)
                    setNameError(false)
                  }}
                  onAnimationEnd={() => setNameError(false)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                      e.preventDefault()
                      handleAddAgent()
                    }
                  }}
                />
                <button
                  className="rounded-md bg-primary px-2.5 py-1.5 text-xs font-medium text-primary-foreground"
                  onClick={handleAddAgent}
                >
                  添加
                </button>
                <button
                  className="text-xs text-muted-foreground hover:text-foreground"
                  onClick={() => {
                    setAddingType(null)
                    setInputName('')
                    setNameError(false)
                  }}
                >
                  取消
                </button>
              </div>
            ) : (
              <div className="flex flex-wrap gap-1.5">
                {types.map((agent) => (
                  <button
                    key={agent.type}
                    className="flex items-center gap-1.5 rounded-lg border border-border px-2.5 py-1.5 text-xs text-foreground transition-colors hover:bg-accent"
                    onMouseEnter={agentHover.onMouseEnter}
                    onMouseLeave={agentHover.onMouseLeave}
                    onClick={() => {
                      setAddingType(agent.type as AgentType)
                      setInputName('')
                    }}
                    disabled={createMutation.isPending}
                  >
                    <AgentAvatar agentType={agent.type as AgentType} status="ready" />
                    <span>{agent.name}</span>
                  </button>
                ))}
              </div>
            )}
            {nameError && (
              <p className="mt-1 text-xs text-destructive">请输入不重复的 Agent 名称</p>
            )}
          </div>
        )}

        {/* Group title (required when 2+ agents) */}
        {agents.length >= 2 && (
          <div className="mb-3">
            <label className="mb-1 block text-xs font-medium text-muted-foreground">
              群聊名称 <span className="text-destructive">*</span>
            </label>
            <input
              value={groupTitle}
              placeholder="为群聊起个名字"
              className="w-full rounded-md border bg-background px-2 py-1.5 text-xs text-foreground outline-none"
              style={{
                borderColor: groupTitleError ? 'var(--destructive)' : 'var(--border)',
                animation: groupTitleError ? 'shake 0.4s ease' : undefined,
              }}
              onChange={(e) => {
                setGroupTitle(e.target.value)
                setGroupTitleError(false)
              }}
              onAnimationEnd={() => setGroupTitleError(false)}
              onKeyDown={(e) => {
                if (e.key === 'Enter') handleSubmit()
              }}
            />
            {groupTitleError && <p className="mt-1 text-xs text-destructive">群聊必须填写名称</p>}
          </div>
        )}

        {/* Submit */}
        <button
          className="w-full rounded-md py-2 text-sm font-medium transition-colors"
          style={{
            backgroundColor: canSubmit ? 'var(--primary)' : 'var(--muted)',
            color: canSubmit ? 'var(--primary-foreground)' : 'var(--muted-foreground)',
            opacity: createMutation.isPending ? 0.6 : 1,
          }}
          onClick={handleSubmit}
          disabled={!canSubmit}
        >
          {createMutation.isPending
            ? '创建中...'
            : agents.length > 1
              ? `创建群聊（${agents.length} 个 Agent）`
              : agents.length === 1
                ? '开始对话'
                : '请添加 Agent'}
        </button>
      </DialogContent>
    </Dialog>
  )
}
