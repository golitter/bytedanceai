import { Send } from 'lucide-react'
import { useCallback, useEffect, useRef, useState } from 'react'

interface MessageInputProps {
  onSend: (message: string) => void
  disabled?: boolean
  sendDisabled?: boolean
  sendDisabledHint?: string
  placeholder?: string
}

export function MessageInput({
  onSend,
  disabled = false,
  sendDisabled = false,
  sendDisabledHint,
  placeholder = '输入消息...',
}: MessageInputProps) {
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const [hint, setHint] = useState<string | null>(null)
  const hintTimerRef = useRef<ReturnType<typeof setTimeout>>(null)

  useEffect(() => {
    return () => {
      if (hintTimerRef.current) clearTimeout(hintTimerRef.current)
    }
  }, [])

  const showHint = useCallback((message: string) => {
    setHint(message)
    if (hintTimerRef.current) clearTimeout(hintTimerRef.current)
    hintTimerRef.current = setTimeout(() => setHint(null), 3000)
  }, [])

  const adjustHeight = useCallback(() => {
    const el = textareaRef.current
    if (!el) return
    el.style.height = 'auto'
    el.style.height = `${Math.min(el.scrollHeight, 200)}px`
  }, [])

  const handleSend = useCallback(() => {
    if (sendDisabled) {
      if (sendDisabledHint) showHint(sendDisabledHint)
      return
    }
    const value = textareaRef.current?.value.trim()
    if (!value || disabled) return
    onSend(value)
    if (textareaRef.current) {
      textareaRef.current.value = ''
      textareaRef.current.style.height = '48px'
    }
  }, [onSend, disabled, sendDisabled, sendDisabledHint, showHint])

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
      if (e.key === 'Enter' && !e.shiftKey && !e.nativeEvent.isComposing) {
        e.preventDefault()
        handleSend()
      }
    },
    [handleSend],
  )

  const canSend = !disabled && !sendDisabled

  return (
    <div className="border-t border-border px-4 py-3">
      {hint && (
        <div className="mb-2 rounded-lg bg-muted px-3 py-1.5 text-xs text-tertiary">{hint}</div>
      )}
      <div className="flex items-end gap-2">
        <textarea
          ref={textareaRef}
          className="flex-1 resize-none rounded-lg bg-card px-3 py-2.5 text-sm text-foreground outline-none placeholder:text-tertiary disabled:opacity-50"
          style={{
            minHeight: 48,
            maxHeight: 200,
          }}
          placeholder={placeholder}
          disabled={disabled}
          rows={1}
          onChange={adjustHeight}
          onKeyDown={handleKeyDown}
        />
        <button
          className="flex h-[48px] w-[40px] shrink-0 items-center justify-center rounded-lg bg-primary transition-colors disabled:opacity-40"
          onClick={handleSend}
          disabled={!canSend}
        >
          <Send className="h-4 w-4 text-primary-foreground" strokeWidth={1.25} />
        </button>
      </div>
    </div>
  )
}
