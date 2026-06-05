import { Component, type ErrorInfo, type ReactNode } from 'react'

import { UI_ACTIONS, UI_MESSAGES } from '@/lib/ui-text'

interface ErrorBoundaryProps {
  children: ReactNode
  fallback?: ReactNode
}

interface ErrorBoundaryState {
  hasError: boolean
  error: Error | null
}

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props)
    this.state = { hasError: false, error: null }
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error('[ErrorBoundary]', error, info.componentStack)
  }

  handleRetry = () => {
    this.setState({ hasError: false, error: null })
  }

  render() {
    if (this.state.hasError) {
      if (this.props.fallback) return this.props.fallback
      return (
        <div className="flex flex-col items-center justify-center gap-3 p-6">
          <p className="text-sm text-destructive">
            {this.state.error?.message || UI_MESSAGES.RENDER_ERROR}
          </p>
          <button
            className="rounded-[6px] bg-primary px-3 py-1.5 text-xs text-primary-foreground"
            onClick={this.handleRetry}
          >
            {UI_ACTIONS.RETRY}
          </button>
        </div>
      )
    }
    return this.props.children
  }
}
