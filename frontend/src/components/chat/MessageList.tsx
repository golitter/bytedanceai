import { useVirtualizer } from '@tanstack/react-virtual'
import { ArrowDown, Loader2 } from 'lucide-react'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'

import type { AgentType } from '@/generated/request'
import type { ChatMessage } from '@/stores/chat'
import { shouldShowTimeSeparator } from '@/utils/time'

import { MarkdownRenderer } from '../markdown/MarkdownRenderer'
import { MessageBubble } from './MessageBubble'
import { TimeDivider } from './TimeDivider'

interface MessageListProps {
  messages: ChatMessage[]
  streamingContent: string
  streamingAgentType?: string
  isStreaming: boolean
  avatarUrl?: string
  agentName?: string
  sessionId?: string
  hasMore: boolean
  isLoadingMore: boolean
  onLoadMore: () => Promise<void>
}

type DisplayItem =
  | { type: 'message'; msg: ChatMessage; isStreamingMsg: boolean }
  | { type: 'time-divider'; timestamp: number }

const VIRTUALIZE_THRESHOLD = 50
const SCROLL_BOTTOM_THRESHOLD = 60

export function MessageList({
  messages,
  streamingContent,
  streamingAgentType,
  isStreaming,
  avatarUrl,
  agentName,
  sessionId,
  hasMore,
  isLoadingMore,
  onLoadMore,
}: MessageListProps) {
  const parentRef = useRef<HTMLDivElement>(null)
  const [autoScroll, setAutoScroll] = useState(true)
  const loadingRef = useRef(false)

  const displayItems = useMemo<DisplayItem[]>(() => {
    const allMsgs =
      isStreaming && streamingContent
        ? [
            ...messages,
            {
              id: 'streaming',
              role: 'agent' as const,
              content: streamingContent,
              agentType: streamingAgentType as AgentType | undefined,
              timestamp: Date.now(),
            },
          ]
        : messages
    const items: DisplayItem[] = []
    for (let i = 0; i < allMsgs.length; i++) {
      const msg = allMsgs[i]
      const prevMsg = i > 0 ? allMsgs[i - 1] : undefined
      if (shouldShowTimeSeparator(prevMsg?.timestamp, msg.timestamp)) {
        items.push({ type: 'time-divider', timestamp: msg.timestamp })
      }
      items.push({
        type: 'message',
        msg,
        isStreamingMsg: isStreaming && msg.id === 'streaming',
      })
    }
    return items
  }, [messages, isStreaming, streamingContent, streamingAgentType])

  const useVirtual = displayItems.length > VIRTUALIZE_THRESHOLD

  // eslint-disable-next-line react-hooks/incompatible-library
  const virtualizer = useVirtualizer({
    count: displayItems.length,
    getScrollElement: () => parentRef.current,
    estimateSize: (index) => {
      const item = displayItems[index]
      if (!item) return 60
      if (item.type === 'time-divider') return 40
      return item.msg.content.length > 200 ? 200 : 80
    },
    overscan: 5,
    enabled: useVirtual,
  })

  const scrollToBottom = useCallback(() => {
    if (!parentRef.current) return
    parentRef.current.scrollTop = parentRef.current.scrollHeight
  }, [])

  useEffect(() => {
    if (autoScroll) {
      scrollToBottom()
    }
  }, [autoScroll, scrollToBottom, streamingContent, messages.length])

  const handleScroll = useCallback(() => {
    const el = parentRef.current
    if (!el) return

    const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < SCROLL_BOTTOM_THRESHOLD
    setAutoScroll(atBottom)

    // Pull-up load: scrollTop === 0 and has more
    if (el.scrollTop === 0 && hasMore && !isLoadingMore && !loadingRef.current) {
      const oldScrollHeight = el.scrollHeight
      loadingRef.current = true
      onLoadMore().finally(() => {
        loadingRef.current = false
        // Restore scroll position after prepending messages
        requestAnimationFrame(() => {
          if (parentRef.current) {
            parentRef.current.scrollTop = parentRef.current.scrollHeight - oldScrollHeight
          }
        })
      })
    }
  }, [hasMore, isLoadingMore, onLoadMore])

  const renderItem = (item: DisplayItem) => {
    if (item.type === 'time-divider') {
      return (
        <div className="px-6 py-2">
          <TimeDivider timestamp={item.timestamp} />
        </div>
      )
    }
    return (
      <div className="px-6 py-2">
        <MessageRenderer
          msg={item.msg}
          isStreaming={item.isStreamingMsg}
          avatarUrl={avatarUrl}
          agentName={agentName}
          sessionId={sessionId}
        />
      </div>
    )
  }

  return (
    <div className="relative flex-1 overflow-hidden">
      {isLoadingMore && (
        <div className="absolute left-0 right-0 top-0 z-10 flex justify-center py-2">
          <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
        </div>
      )}
      <div ref={parentRef} className="h-full overflow-y-auto" onScroll={handleScroll}>
        {useVirtual ? (
          <div
            style={{
              height: `${virtualizer.getTotalSize()}px`,
              width: '100%',
              position: 'relative',
            }}
          >
            {virtualizer.getVirtualItems().map((virtualRow) => {
              const item = displayItems[virtualRow.index]
              if (!item) return null
              return (
                <div
                  key={item.type === 'time-divider' ? `divider-${virtualRow.index}` : item.msg.id}
                  style={{
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    width: '100%',
                    transform: `translateY(${virtualRow.start}px)`,
                  }}
                  ref={virtualizer.measureElement}
                  data-index={virtualRow.index}
                >
                  {renderItem(item)}
                </div>
              )
            })}
          </div>
        ) : (
          <div className="py-4">
            {displayItems.map((item, i) => {
              const key = item.type === 'time-divider' ? `divider-${i}` : item.msg.id
              return <div key={key}>{renderItem(item)}</div>
            })}
          </div>
        )}
      </div>

      {!autoScroll && (
        <button
          className="absolute bottom-4 right-6 flex h-8 w-8 items-center justify-center rounded-lg bg-accent transition-colors"
          onClick={() => {
            scrollToBottom()
            setAutoScroll(true)
          }}
        >
          <ArrowDown className="h-4 w-4 text-muted-foreground" strokeWidth={1.25} />
        </button>
      )}
    </div>
  )
}

function MessageRenderer({
  msg,
  isStreaming,
  avatarUrl,
  agentName,
  sessionId,
}: {
  msg: ChatMessage
  isStreaming: boolean
  avatarUrl?: string
  agentName?: string
  sessionId?: string
}) {
  if (msg.role === 'user') {
    return <MessageBubble variant="user">{msg.content}</MessageBubble>
  }

  if (msg.role === 'agent') {
    return (
      <MessageBubble
        variant="agent"
        agentType={msg.agentType ?? 'claude-code'}
        avatarUrl={avatarUrl}
        agentName={agentName}
        status={isStreaming ? 'running' : 'ready'}
        isStreaming={isStreaming}
        blocks={msg.blocks}
        sessionId={sessionId}
      >
        <MarkdownRenderer content={msg.content} />
      </MessageBubble>
    )
  }

  return <MessageBubble variant="system">{msg.content}</MessageBubble>
}
