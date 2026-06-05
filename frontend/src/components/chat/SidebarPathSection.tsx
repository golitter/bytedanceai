import { ChevronRight, FolderOpen, FolderSync } from 'lucide-react'

import { UI_LABELS, UI_MESSAGES } from '@/lib/ui-text'

import { showCopyToast } from './SidebarActions'
import { useCollapsible } from './useCollapsible'

export function SidebarPathSection({ repoPath, taskId }: { repoPath: string; taskId: string }) {
  const [pathsOpen, togglePaths] = useCollapsible('paths', false)

  return (
    <div className="border-b border-sidebar-border px-4 py-3">
      <button
        type="button"
        className="flex w-full items-center gap-1 text-[11px] font-medium uppercase tracking-wider text-tertiary"
        onClick={togglePaths}
      >
        <ChevronRight
          className={`h-3 w-3 transition-transform ${pathsOpen ? 'rotate-90' : ''}`}
          strokeWidth={1.25}
        />
        {UI_LABELS.PATH_INFO}
      </button>
      {pathsOpen && (
        <div className="mt-2 flex flex-col gap-2">
          {/* Repo path */}
          <div>
            <span className="mb-0.5 block text-[11px] text-tertiary">{UI_LABELS.REPO_PATH}</span>
            <p
              className="flex select-none truncate rounded-md bg-bg-subtle px-2.5 py-1.5 text-xs text-muted-foreground transition-[transform,opacity] hover:bg-bg-hover hover:text-foreground"
              title={`${repoPath} — ${UI_MESSAGES.DOUBLE_CLICK_TO_COPY}`}
              onDoubleClick={() => {
                navigator.clipboard.writeText(repoPath)
                showCopyToast()
              }}
            >
              <FolderOpen
                className="mr-1.5 h-3.5 w-3.5 shrink-0 text-tertiary"
                strokeWidth={1.25}
              />
              <span className="truncate">{repoPath}</span>
            </p>
          </div>
          {/* Task path — worktrees live in <repo_parent>/worktrees/<taskId> */}
          <div>
            <span className="mb-0.5 block text-[11px] text-tertiary">{UI_LABELS.TASK_PATH}</span>
            <p
              className="flex select-none truncate rounded-md bg-bg-subtle px-2.5 py-1.5 text-xs text-muted-foreground transition-[transform,opacity] hover:bg-bg-hover hover:text-foreground"
              title={`${repoPath.replace(/\/[^/]+$/, '')}/worktrees/${taskId} — ${UI_MESSAGES.DOUBLE_CLICK_TO_COPY}`}
              onDoubleClick={() => {
                navigator.clipboard.writeText(
                  `${repoPath.replace(/\/[^/]+$/, '')}/worktrees/${taskId}`,
                )
                showCopyToast()
              }}
            >
              <FolderSync
                className="mr-1.5 h-3.5 w-3.5 shrink-0 text-tertiary"
                strokeWidth={1.25}
              />
              <span className="truncate">{`${repoPath.replace(/\/[^/]+$/, '')}/worktrees/${taskId}`}</span>
            </p>
          </div>
        </div>
      )}
    </div>
  )
}
