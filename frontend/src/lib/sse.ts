import type { StreamEvent } from '@/generated/events'

interface SSEOptions {
  url: string
  body: Record<string, unknown>
  onEvent: (event: StreamEvent) => void
  onError?: (error: Error) => void
}

export function connectSSE({ url, body, onEvent, onError }: SSEOptions): AbortController {
  const controller = new AbortController()

  ;(async () => {
    try {
      const res = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
        signal: controller.signal,
      })

      if (!res.ok) {
        throw new Error(`HTTP ${res.status}: ${res.statusText}`)
      }

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
      if ((err as Error).name !== 'AbortError') {
        onError?.(err as Error)
      }
    }
  })()

  return controller
}
