import { Send } from 'lucide-react'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'

import type { AgentSessionInfo } from '@/lib/api'

const MAX_INPUT_HEIGHT = 200
const MIN_INPUT_HEIGHT = 48
const HINT_DISPLAY_DURATION = 3000

interface MessageInputProps {
  onSend: (message: string) => void
  disabled?: boolean
  sendDisabled?: boolean
  sendDisabledHint?: string
  placeholder?: string
  mentionSessions?: AgentSessionInfo[]
}

export function MessageInput({
  onSend,
  disabled = false,
  sendDisabled = false,
  sendDisabledHint,
  placeholder = '输入消息...',
  mentionSessions,
}: MessageInputProps) {
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const [hint, setHint] = useState<string | null>(null)
  const [mentionOpen, setMentionOpen] = useState(false)
  const [mentionQuery, setMentionQuery] = useState('')
  const [mentionStart, setMentionStart] = useState(0)
  const [activeMentionIndex, setActiveMentionIndex] = useState(0)
  const hintTimerRef = useRef<ReturnType<typeof setTimeout>>(null)

  useEffect(() => {
    return () => {
      if (hintTimerRef.current) clearTimeout(hintTimerRef.current)
    }
  }, [])

  const showHint = useCallback((message: string) => {
    setHint(message)
    if (hintTimerRef.current) clearTimeout(hintTimerRef.current)
    hintTimerRef.current = setTimeout(() => setHint(null), HINT_DISPLAY_DURATION)
  }, [])

  const adjustHeight = useCallback(() => {
    const el = textareaRef.current
    if (!el) return
    el.style.height = 'auto'
    el.style.height = `${Math.min(el.scrollHeight, MAX_INPUT_HEIGHT)}px`
  }, [])

  const updateMentionState = useCallback(() => {
    const el = textareaRef.current
    if (!el || !mentionSessions?.length) {
      setMentionOpen(false)
      return
    }
    const cursor = el.selectionStart ?? el.value.length
    const beforeCursor = el.value.slice(0, cursor)
    const match = /(^|\s)@([^\s@]*)$/.exec(beforeCursor)
    if (!match) {
      setMentionOpen(false)
      return
    }
    const prefix = match[1] ?? ''
    setMentionStart(beforeCursor.length - match[0].length + prefix.length)
    setMentionQuery(match[2] ?? '')
    setMentionOpen(true)
    setActiveMentionIndex(0)
  }, [mentionSessions])

  const mentionOptions = useMemo(() => {
    if (!mentionSessions?.length) return []
    const query = mentionQuery.trim().toLowerCase()
    return mentionSessions
      .filter((session) => {
        if (!query) return true
        const values = [
          session.mentionLabel,
          session.routeId,
          session.agentName,
          session.agentType,
          ...(session.aliases ?? []),
        ]
        return values.some((value) => value.toLowerCase().includes(query))
      })
      .slice(0, 8)
  }, [mentionQuery, mentionSessions])

  const insertMention = useCallback(
    (session: AgentSessionInfo) => {
      const el = textareaRef.current
      if (!el) return
      const cursor = el.selectionStart ?? el.value.length
      const before = el.value.slice(0, mentionStart)
      const after = el.value.slice(cursor)
      const insertion = `@${session.mentionLabel} `
      el.value = `${before}${insertion}${after}`
      const nextCursor = before.length + insertion.length
      el.focus()
      el.setSelectionRange(nextCursor, nextCursor)
      adjustHeight()
      setMentionOpen(false)
    },
    [adjustHeight, mentionStart],
  )

  const handleSend = useCallback(() => {
    if (sendDisabled) {
      if (sendDisabledHint) showHint(sendDisabledHint)
      return
    }
    const value = textareaRef.current?.value.trim()
    if (!value || disabled) return
    onSend(value)
    setMentionOpen(false)
    if (textareaRef.current) {
      textareaRef.current.value = ''
      textareaRef.current.style.height = `${MIN_INPUT_HEIGHT}px`
    }
  }, [onSend, disabled, sendDisabled, sendDisabledHint, showHint])

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
      if (e.key === 'Enter' && !e.shiftKey && !e.nativeEvent.isComposing) {
        if (mentionOpen && mentionOptions.length > 0) {
          e.preventDefault()
          insertMention(mentionOptions[activeMentionIndex] ?? mentionOptions[0])
          return
        }
        e.preventDefault()
        handleSend()
        return
      }
      if (mentionOpen && mentionOptions.length > 0) {
        if (e.key === 'ArrowDown') {
          e.preventDefault()
          setActiveMentionIndex((idx) => (idx + 1) % mentionOptions.length)
        } else if (e.key === 'ArrowUp') {
          e.preventDefault()
          setActiveMentionIndex((idx) => (idx - 1 + mentionOptions.length) % mentionOptions.length)
        } else if (e.key === 'Tab') {
          e.preventDefault()
          insertMention(mentionOptions[activeMentionIndex] ?? mentionOptions[0])
        } else if (e.key === 'Escape') {
          e.preventDefault()
          setMentionOpen(false)
        }
      }
    },
    [activeMentionIndex, handleSend, insertMention, mentionOpen, mentionOptions],
  )

  const canSend = !disabled && !sendDisabled

  return (
    <div className="border-t border-border px-4 py-3">
      {hint && (
        <div className="mb-2 rounded-lg bg-muted px-3 py-1.5 text-xs text-tertiary">{hint}</div>
      )}
      <div className="relative flex items-end gap-2">
        {mentionOpen && mentionOptions.length > 0 && (
          <div className="absolute bottom-[calc(100%+8px)] left-0 z-20 w-[min(360px,calc(100vw-2rem))] min-w-[220px] overflow-hidden rounded-[8px] border border-border bg-popover py-1 shadow-lg">
            {mentionOptions.map((session, index) => (
              <button
                key={session.sessionId}
                type="button"
                className={`flex w-full min-w-0 items-center justify-between gap-3 px-3 py-2 text-left text-sm ${
                  index === activeMentionIndex ? 'bg-muted text-foreground' : 'text-foreground'
                }`}
                onMouseDown={(e) => {
                  e.preventDefault()
                  insertMention(session)
                }}
              >
                <span className="min-w-0 truncate font-medium">{session.mentionLabel}</span>
                <span className="shrink-0 text-[11px] text-tertiary">{session.agentType}</span>
              </button>
            ))}
          </div>
        )}
        <textarea
          ref={textareaRef}
          className="flex-1 resize-none break-words rounded-[8px] bg-card px-3 py-2.5 text-sm text-foreground outline-none placeholder:text-tertiary disabled:opacity-50"
          style={{
            minHeight: MIN_INPUT_HEIGHT,
            maxHeight: MAX_INPUT_HEIGHT,
          }}
          placeholder={placeholder}
          disabled={disabled}
          rows={1}
          onChange={() => {
            adjustHeight()
            updateMentionState()
          }}
          onClick={updateMentionState}
          onKeyDown={handleKeyDown}
        />
        <button
          className="flex w-[40px] shrink-0 items-center justify-center rounded-[6px] bg-primary disabled:opacity-40"
          style={{ height: MIN_INPUT_HEIGHT }}
          onClick={handleSend}
          disabled={!canSend}
        >
          <Send className="h-4 w-4 text-primary-foreground" strokeWidth={1.25} />
        </button>
      </div>
    </div>
  )
}
