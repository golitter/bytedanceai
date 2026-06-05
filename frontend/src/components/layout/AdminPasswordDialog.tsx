import { Lock } from 'lucide-react'
import { useState } from 'react'

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { adminAuth } from '@/lib/api'
import { UI_ACTIONS, UI_LABELS, UI_MESSAGES, UI_PLACEHOLDERS, UI_STATUS } from '@/lib/ui-text'
import { useAdminStore } from '@/stores/admin'

export function AdminPasswordDialog() {
  const {
    showPasswordDialog,
    passwordDialogPurpose,
    hidePasswordDialog,
    setAdminToken,
    setIsAuthenticated,
  } = useAdminStore()
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!password) return
    setLoading(true)
    setError('')
    try {
      const res = await adminAuth(password)
      setAdminToken(res.token)
      setIsAuthenticated(true)
      hidePasswordDialog()
      setPassword('')
    } catch {
      setError(UI_MESSAGES.PASSWORD_ERROR)
    } finally {
      setLoading(false)
    }
  }

  const handleOpenChange = (open: boolean) => {
    if (!open) {
      hidePasswordDialog()
      setPassword('')
      setError('')
    }
  }

  return (
    <Dialog open={showPasswordDialog} onOpenChange={handleOpenChange}>
      <DialogContent className="max-w-[360px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Lock className="h-4 w-4" strokeWidth={1.25} />
            {passwordDialogPurpose === 'login'
              ? UI_LABELS.ADMIN_VERIFY
              : UI_LABELS.SENSITIVE_CONFIRM}
          </DialogTitle>
          <DialogDescription className="sr-only">
            {passwordDialogPurpose === 'login'
              ? UI_LABELS.ADMIN_VERIFY
              : UI_LABELS.SENSITIVE_CONFIRM}
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="flex flex-col gap-3">
          <p className="text-sm text-text-secondary">
            {passwordDialogPurpose === 'login'
              ? UI_LABELS.ADMIN_VERIFY
              : UI_LABELS.SENSITIVE_CONFIRM}
          </p>
          <input
            type="password"
            value={password}
            onChange={(e) => {
              setPassword(e.target.value)
              setError('')
            }}
            placeholder={UI_PLACEHOLDERS.PASSWORD}
            className="h-9 rounded-md px-3 text-sm outline-none"
            style={{
              border: '1px solid var(--border)',
              background: 'var(--bg-canvas)',
              color: 'var(--text-primary)',
            }}
            autoFocus
          />
          {error && (
            <p className="text-xs" style={{ color: 'var(--color-error)' }}>
              {error}
            </p>
          )}
          <button
            type="submit"
            disabled={loading || !password}
            className="h-9 rounded-md text-sm font-medium transition-[transform,opacity] disabled:opacity-50"
            style={{ background: 'var(--color-brand)', color: 'var(--primary-foreground)' }}
          >
            {loading ? UI_STATUS.VERIFYING : UI_ACTIONS.CONFIRM}
          </button>
        </form>
      </DialogContent>
    </Dialog>
  )
}
