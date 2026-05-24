import type { StreamEvent } from '@/generated/events'

interface SSEOptions {
  url: string
  params?: Record<string, string>
  onEvent: (event: StreamEvent) => void
  onError?: (error: Error) => void
  /** Enable auto-reconnect with exponential backoff on the same URL/params */
  reconnect?: boolean
  /** Max backoff in ms (default 10000) */
  maxBackoff?: number
}

export function connectSSE({
  url,
  params,
  onEvent,
  onError,
  reconnect = false,
  maxBackoff = 10000,
}: SSEOptions): AbortController {
  const controller = new AbortController()

  let attempt = 0

  const connect = async () => {
    try {
      const qs = params ? '?' + new URLSearchParams(params).toString() : ''
      const fullUrl = `${url}${qs}`

      const res = await fetch(fullUrl, {
        method: 'GET',
        headers: { Accept: 'text/event-stream' },
        signal: controller.signal,
      })

      if (!res.ok) {
        throw new Error(`HTTP ${res.status}: ${res.statusText}`)
      }

      // Reset attempt on successful connection
      attempt = 0

      const reader = res.body?.getReader()
      if (!reader) throw new Error('No readable stream')

      const decoder = new TextDecoder()
      let buffer = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        buffer += decoder.decode(value, { stream: true })

        const parts = buffer.split('\n\n')
        buffer = parts.pop() ?? ''

        for (const part of parts) {
          for (const line of part.split('\n')) {
            if (!line.startsWith('data: ')) continue
            const jsonStr = line.slice(6)
            if (!jsonStr.trim()) continue
            try {
              const event: StreamEvent = JSON.parse(jsonStr)
              onEvent(event)
            } catch {
              console.warn('Failed to parse SSE event:', jsonStr)
            }
          }
        }
      }
    } catch (err) {
      if ((err as Error).name === 'AbortError') return

      if (reconnect && !controller.signal.aborted) {
        attempt++
        const delay = Math.min(1000 * Math.pow(2, attempt - 1) + Math.random() * 500, maxBackoff)
        console.warn(
          `SSE disconnected, reconnecting in ${Math.round(delay)}ms (attempt ${attempt})`,
        )
        setTimeout(() => {
          if (!controller.signal.aborted) connect()
        }, delay)
        return
      }

      onError?.(err as Error)
    }
  }

  connect()

  return controller
}
