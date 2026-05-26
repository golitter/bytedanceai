export type MessageBlock =
  | { type: 'text'; id: string; content: string }
  | { type: 'html-render'; id: string; content: string }
  | { type: 'image'; id: string; path: string }
  | { type: 'attachment'; id: string; path: string }
  | { type: 'diff'; id: string; snapshotId: string }
  | { type: 'preview'; id: string; url: string }
