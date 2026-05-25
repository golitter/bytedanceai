export type MessageBlock =
  | { type: 'text'; content: string }
  | { type: 'html-render'; content: string }
  | { type: 'image'; path: string }
  | { type: 'attachment'; path: string }
  | { type: 'diff'; snapshotId: string }
  | { type: 'preview'; url: string }
