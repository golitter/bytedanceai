import { useState } from 'react'

import { validateRepoPath } from '@/lib/api'
import { UI_ACTIONS, UI_ERRORS, UI_LABELS, UI_STATUS } from '@/lib/ui-text'

interface RepoPathInputProps {
  onValidationChange: (path: string, validated: boolean) => void
}

export function RepoPathInput({ onValidationChange }: RepoPathInputProps) {
  const [repoPath, setRepoPath] = useState('')
  const [validated, setValidated] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [validating, setValidating] = useState(false)

  const handleValidate = async () => {
    const path = repoPath.trim()
    if (!path) {
      setError(UI_ERRORS.REPO_PATH_REQUIRED)
      setValidated(false)
      onValidationChange('', false)
      return
    }
    setValidating(true)
    setError(null)
    try {
      const result = await validateRepoPath(path)
      if (result.valid) {
        setValidated(true)
        setError(null)
        onValidationChange(path, true)
      } else {
        setValidated(false)
        setError(result.errors.join('; '))
        onValidationChange(path, false)
      }
    } catch {
      setValidated(false)
      setError(UI_ERRORS.VALIDATE_FAILED)
      onValidationChange(path, false)
    } finally {
      setValidating(false)
    }
  }

  return (
    <div className="mb-3">
      <label className="mb-1 block text-xs font-medium text-muted-foreground">
        {UI_LABELS.REPO_PATH}
      </label>
      <div className="flex items-center gap-2">
        <input
          value={repoPath}
          placeholder="/path/to/repo"
          className="flex-1 rounded-md border bg-background px-2 py-1.5 text-xs text-foreground outline-none"
          style={{
            borderColor: error
              ? 'var(--destructive)'
              : validated
                ? 'var(--color-success)'
                : 'var(--border)',
          }}
          onChange={(e) => {
            setRepoPath(e.target.value)
            setValidated(false)
            setError(null)
            onValidationChange(e.target.value.trim(), false)
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
          {validating ? UI_STATUS.VALIDATING : UI_ACTIONS.VALIDATE}
        </button>
      </div>
      {error && <p className="mt-1 text-xs text-destructive">{error}</p>}
      {validated && (
        <p className="mt-1 text-xs" style={{ color: 'var(--color-success)' }}>
          路径校验通过
        </p>
      )}
    </div>
  )
}
