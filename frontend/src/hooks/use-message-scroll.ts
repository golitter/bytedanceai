import { useCallback, useEffect, useLayoutEffect, useRef, useState } from 'react'

const SCROLL_BOTTOM_THRESHOLD = 60

interface UseMessageScrollOptions {
  hasMore: boolean
  isLoadingMore: boolean
  onLoadMore: () => Promise<void>
  streamingContent: string
  messagesLength: number
}

export function useMessageScroll(
  parentRef: React.RefObject<HTMLDivElement | null>,
  options: UseMessageScrollOptions,
) {
  const { hasMore, isLoadingMore, onLoadMore, streamingContent, messagesLength } = options
  const [autoScroll, setAutoScroll] = useState(true)
  const loadingRef = useRef(false)
  const scrollRafRef = useRef<number | null>(null)
  const prevMsgLenRef = useRef(messagesLength)

  const scrollToBottom = useCallback(() => {
    if (!parentRef.current) return
    parentRef.current.scrollTop = parentRef.current.scrollHeight
  }, [parentRef])

  // Initial load or history loaded: scroll to bottom synchronously after DOM layout
  useLayoutEffect(() => {
    if (autoScroll) {
      scrollToBottom()
    }
    prevMsgLenRef.current = messagesLength
  }, [autoScroll, scrollToBottom, messagesLength])

  // Streaming content updates: throttle with rAF for smooth rendering
  useEffect(() => {
    if (!autoScroll || !streamingContent) return
    // Skip if this is an initial load (message count changed), already handled above
    if (messagesLength !== prevMsgLenRef.current) return
    if (scrollRafRef.current !== null) return
    scrollRafRef.current = requestAnimationFrame(() => {
      scrollRafRef.current = null
      // Double-rAF to ensure browser has completed layout after streaming content change
      requestAnimationFrame(() => {
        scrollToBottom()
      })
    })
  }, [autoScroll, streamingContent, messagesLength, scrollToBottom])

  useEffect(() => {
    return () => {
      if (scrollRafRef.current !== null) {
        cancelAnimationFrame(scrollRafRef.current)
      }
    }
  }, [])

  const handleScroll = useCallback(() => {
    const el = parentRef.current
    if (!el) return

    const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < SCROLL_BOTTOM_THRESHOLD
    setAutoScroll(atBottom)

    if (el.scrollTop === 0 && hasMore && !isLoadingMore && !loadingRef.current) {
      const oldScrollHeight = el.scrollHeight
      loadingRef.current = true
      onLoadMore().finally(() => {
        loadingRef.current = false
        requestAnimationFrame(() => {
          if (parentRef.current) {
            parentRef.current.scrollTop = parentRef.current.scrollHeight - oldScrollHeight
          }
        })
      })
    }
  }, [hasMore, isLoadingMore, onLoadMore, parentRef])

  const enableAutoScroll = useCallback(() => setAutoScroll(true), [])

  return { autoScroll, handleScroll, scrollToBottom, enableAutoScroll }
}
