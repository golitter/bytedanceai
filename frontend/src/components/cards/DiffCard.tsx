import 'react-diff-view/style/index.css'

import { Check, Pencil, RotateCcw, X } from 'lucide-react'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'

import { DiffFileEditor } from '@/components/diff/DiffFileEditor'
import { DiffFileTabs } from '@/components/diff/DiffFileTabs'
import { DiffFileView } from '@/components/diff/DiffFileView'
import type { ParsedDiffFile } from '@/lib/diff-parser'
import { getFileName, parseUnifiedDiff } from '@/lib/diff-parser'

const API_BASE = '/api'

interface DiffCardProps {
  snapshotId: string
  sessionId?: string
}

type SnapshotStatus = 'pending' | 'committed' | 'reverted' | 'cancelled'

export function DiffCard({ snapshotId, sessionId }: DiffCardProps) {
  const [diff, setDiff] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeFileIndex, setActiveFileIndex] = useState(0)
  const [editingFile, setEditingFile] = useState(false)
  const [actionStatus, setActionStatus] = useState<'idle' | 'committing' | 'reverting'>('idle')
  const [snapshotStatus, setSnapshotStatus] = useState<SnapshotStatus | null>(null)
  const initialized = useRef(false)

  const isSettled =
    snapshotStatus === 'committed' ||
    snapshotStatus === 'reverted' ||
    snapshotStatus === 'cancelled'

  const parsed = useMemo(() => parseUnifiedDiff(diff ?? ''), [diff])
  const activeFile: ParsedDiffFile | undefined = parsed.files[activeFileIndex]

  // Snapshot-first load: GET snapshot → 404 → workspace diff → PUT pending
  useEffect(() => {
    if (initialized.current) return
    initialized.current = true
    ;(async () => {
      setLoading(true)
      setError(null)
      try {
        // Try fetching existing snapshot
        const snapRes = await fetch(`${API_BASE}/diff-snapshots/${snapshotId}`)
        if (snapRes.ok) {
          const snap = await snapRes.json()
          const data = snap?.data ?? snap
          setDiff(data.diff_content ?? data.diff ?? '')
          setSnapshotStatus(data.status ?? 'pending')
          setLoading(false)
          return
        }

        // Snapshot not found — fetch workspace diff and create pending snapshot
        if (!sessionId) {
          setLoading(false)
          return
        }

        const wsRes = await fetch(`${API_BASE}/session/${sessionId}/diff`)
        let diffText = ''
        if (wsRes.ok) {
          diffText = await wsRes.text()
        }

        if (!diffText?.trim()) {
          setLoading(false)
          return
        }

        // Create pending snapshot
        await fetch(`${API_BASE}/diff-snapshots/${snapshotId}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ session_id: sessionId, diff: diffText, status: 'pending' }),
        })
        setDiff(diffText)
        setSnapshotStatus('pending')
      } catch (e) {
        setError(e instanceof Error ? e.message : 'Failed to load diff')
      } finally {
        setLoading(false)
      }
    })()
  }, [snapshotId, sessionId])

  const refresh = useCallback(async () => {
    if (!sessionId) return
    setLoading(true)
    setError(null)
    try {
      const res = await fetch(`${API_BASE}/session/${sessionId}/diff`)
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      setDiff(await res.text())
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load diff')
    } finally {
      setLoading(false)
    }
  }, [sessionId])

  const handleAccept = async () => {
    if (!sessionId || actionStatus !== 'idle' || !diff) return
    setActionStatus('committing')
    try {
      await fetch(`${API_BASE}/session/${sessionId}/commit`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ message: 'auto commit' }),
      })
      await fetch(`${API_BASE}/diff-snapshots/${snapshotId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ session_id: sessionId, diff, status: 'committed' }),
      })
      setSnapshotStatus('committed')
    } catch {
      // ignore
    } finally {
      setActionStatus('idle')
    }
  }

  const handleReject = async () => {
    if (!sessionId || actionStatus !== 'idle' || !diff) return
    setActionStatus('reverting')
    try {
      await fetch(`${API_BASE}/session/${sessionId}/revert`, { method: 'POST' })
      await fetch(`${API_BASE}/diff-snapshots/${snapshotId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ session_id: sessionId, diff, status: 'reverted' }),
      })
      setSnapshotStatus('reverted')
    } catch {
      // ignore
    } finally {
      setActionStatus('idle')
    }
  }

  const handleEditSave = async (content: string) => {
    if (!sessionId || !activeFile) return
    const filePath = activeFile.newPath
    await fetch(`${API_BASE}/session/${sessionId}/files/${encodeURIComponent(filePath)}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'text/plain' },
      body: content,
    })
    setEditingFile(false)
    await refresh()
  }

  if (loading) {
    return (
      <div className="my-2 rounded-lg border border-border bg-card px-4 py-3 text-sm text-muted-foreground">
        Loading diff...
      </div>
    )
  }

  if (error) {
    return (
      <div className="my-2 rounded-lg border border-destructive/50 bg-card px-4 py-3 text-sm text-destructive">
        {error}
      </div>
    )
  }

  if (!diff?.trim() || parsed.files.length === 0) {
    return null
  }

  const { summary } = parsed

  const badgeConfig: Record<string, { icon: React.ReactNode; label: string; className: string }> = {
    committed: {
      icon: <Check className="h-3 w-3" />,
      label: '已接受',
      className: 'bg-green-500/10 text-green-600',
    },
    reverted: {
      icon: <RotateCcw className="h-3 w-3" />,
      label: '已拒绝',
      className: 'bg-muted text-muted-foreground',
    },
    cancelled: {
      icon: <X className="h-3 w-3" />,
      label: '已取消',
      className: 'bg-muted text-muted-foreground',
    },
  }

  return (
    <div className="diff-card my-2 overflow-hidden rounded-lg border border-border">
      {/* Header */}
      <div className="flex items-center justify-between border-b border-border bg-muted/50 px-3 py-1.5">
        <span className="text-xs text-muted-foreground">
          {summary.filesChanged} file{summary.filesChanged !== 1 ? 's' : ''} changed,{' '}
          <span className="text-green-500">+{summary.additions}</span>{' '}
          <span className="text-red-500">-{summary.deletions}</span>
        </span>
        <div className="flex items-center gap-1">
          {snapshotStatus && badgeConfig[snapshotStatus] && (
            <span
              className={`mr-2 inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium ${badgeConfig[snapshotStatus].className}`}
            >
              {badgeConfig[snapshotStatus].icon} {badgeConfig[snapshotStatus].label}
            </span>
          )}
          {!isSettled && sessionId && (
            <button
              onClick={() => setEditingFile(true)}
              className="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
            >
              <Pencil className="h-3 w-3" />
              编辑
            </button>
          )}
          {!isSettled && (
            <>
              <button
                onClick={handleAccept}
                disabled={actionStatus !== 'idle'}
                className="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground disabled:opacity-50"
              >
                <Check className="h-3 w-3" />
                {actionStatus === 'committing' ? '提交中...' : '接受变更'}
              </button>
              <button
                onClick={handleReject}
                disabled={actionStatus !== 'idle'}
                className="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground disabled:opacity-50"
              >
                <RotateCcw className="h-3 w-3" />
                {actionStatus === 'reverting' ? '撤销中...' : '拒绝变更'}
              </button>
            </>
          )}
        </div>
      </div>

      {/* File tabs */}
      <DiffFileTabs
        files={parsed.files}
        activeIndex={activeFileIndex}
        onSelect={setActiveFileIndex}
      />

      {/* Content */}
      {activeFile && (
        <div className={`max-h-96 overflow-auto text-xs${isSettled ? ' opacity-60' : ''}`}>
          {editingFile ? (
            <DiffFileEditor
              oldContent={activeFile.oldContent}
              newContent={activeFile.newContent}
              fileName={getFileName(activeFile.newPath)}
              onSave={handleEditSave}
              onCancel={() => setEditingFile(false)}
            />
          ) : (
            <DiffFileView file={activeFile} />
          )}
        </div>
      )}
    </div>
  )
}
