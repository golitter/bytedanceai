import { MessageSquare, Plus, Search } from 'lucide-react'
import { useState } from 'react'

import { useConversations } from '@/hooks/use-conversations'
import { useChatNav } from '@/stores/chat'

import { ConversationItem } from './ConversationItem'
import { NewChatDialog } from './NewChatDialog'

export function ConversationList() {
  const [search, setSearch] = useState('')
  const [showNewChat, setShowNewChat] = useState(false)
  const { data: conversations, isLoading } = useConversations()
  const { currentSessionId, setCurrentSession } = useChatNav()

  const filtered = conversations?.filter((c) => {
    if (!search) return true
    const q = search.toLowerCase()
    return c.agentType.includes(q) || c.taskTitle.toLowerCase().includes(q)
  })

  return (
    <div
      className="flex h-full w-[280px] shrink-0 flex-col border-r"
      style={{
        backgroundColor: 'var(--bg-sidebar)',
        borderColor: 'rgba(255,255,255,0.06)',
      }}
    >
      {/* Header */}
      <div
        className="flex shrink-0 items-center gap-2 border-b px-4 py-3"
        style={{ borderColor: 'rgba(255,255,255,0.06)' }}
      >
        <div
          className="flex h-6 w-6 items-center justify-center rounded-md text-xs font-bold text-white"
          style={{ backgroundColor: 'var(--color-brand)' }}
        >
          A
        </div>
        <span className="flex-1 text-sm font-semibold" style={{ color: 'var(--text-primary)' }}>
          AgentHub
        </span>
        <button
          className="flex h-7 w-7 items-center justify-center rounded-md transition-colors"
          style={{ color: 'var(--text-secondary)' }}
          onClick={() => setShowNewChat(true)}
          onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = 'var(--bg-hover)')}
          onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}
        >
          <Plus className="h-4 w-4" strokeWidth={1.5} />
        </button>
      </div>

      {/* Search */}
      <div className="shrink-0 px-3 py-2">
        <div
          className="flex items-center gap-2 rounded-lg px-3 py-1.5"
          style={{ backgroundColor: 'var(--bg-hover)' }}
        >
          <Search
            className="h-3.5 w-3.5 shrink-0"
            style={{ color: 'var(--text-tertiary)' }}
            strokeWidth={1.5}
          />
          <input
            type="text"
            placeholder="搜索对话..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full bg-transparent text-xs outline-none"
            style={{ color: 'var(--text-primary)' }}
          />
        </div>
      </div>

      {/* Conversation list */}
      <div className="flex-1 overflow-y-auto px-2">
        {isLoading ? (
          <div className="px-2 py-4 text-xs" style={{ color: 'var(--text-tertiary)' }}>
            加载中...
          </div>
        ) : !filtered?.length ? (
          <div className="flex flex-col items-center gap-2 px-2 py-8">
            <MessageSquare
              className="h-6 w-6"
              style={{ color: 'var(--text-tertiary)' }}
              strokeWidth={1.5}
            />
            <p className="text-center text-xs" style={{ color: 'var(--text-tertiary)' }}>
              {search ? '没有找到对话' : '还没有对话，点击 + 开始新对话'}
            </p>
          </div>
        ) : (
          filtered.map((c) => (
            <ConversationItem
              key={c.sessionId}
              conversation={c}
              isActive={c.sessionId === currentSessionId}
              onClick={() => setCurrentSession(c.sessionId)}
            />
          ))
        )}
      </div>

      <NewChatDialog open={showNewChat} onOpenChange={setShowNewChat} />
    </div>
  )
}
