import type { MessageBlock } from './block-types'

const BLOCK_MARKER = 'aka_yhy'

const BLOCK_RE = new RegExp('```' + BLOCK_MARKER + '\\n([\\s\\S]*?)```', 'g')

let _blockIdCounter = 0
function nextBlockId(): string {
  return `blk-${++_blockIdCounter}`
}

export function reduceEventToBlocks(fullText: string): MessageBlock[] {
  const blocks: MessageBlock[] = []
  let lastIndex = 0

  for (const match of fullText.matchAll(BLOCK_RE)) {
    const matchStart = match.index!
    if (matchStart > lastIndex) {
      const text = fullText.slice(lastIndex, matchStart)
      if (text) blocks.push({ type: 'text', id: nextBlockId(), content: text })
    }

    const inner = match[1].trim()
    const parsed = parseBlockContent(inner)
    if (parsed) {
      blocks.push(parsed)
    } else {
      blocks.push({ type: 'text', id: nextBlockId(), content: match[0] })
    }

    lastIndex = matchStart + match[0].length
  }

  if (lastIndex < fullText.length) {
    blocks.push({ type: 'text', id: nextBlockId(), content: fullText.slice(lastIndex) })
  }

  if (blocks.length === 0) {
    return [{ type: 'text', id: nextBlockId(), content: fullText }]
  }

  return blocks
}

function parseBlockContent(inner: string): MessageBlock | null {
  const lines = inner.split('\n')
  const typeLine = lines.find((l) => l.startsWith('type:'))
  if (!typeLine) return null

  const blockType = typeLine.slice('type:'.length).trim()

  switch (blockType) {
    case 'html-render': {
      const content = lines.slice(1).join('\n').trim()
      return { type: 'html-render', id: nextBlockId(), content }
    }
    case 'image': {
      const path = extractField(lines, 'path')
      return path ? { type: 'image', id: nextBlockId(), path } : null
    }
    case 'attachment': {
      const path = extractField(lines, 'path')
      return path ? { type: 'attachment', id: nextBlockId(), path } : null
    }
    case 'diff': {
      const snapshotId = extractField(lines, 'snapshotId')
      return snapshotId ? { type: 'diff', id: nextBlockId(), snapshotId } : null
    }
    case 'preview': {
      const url = extractField(lines, 'url')
      return url ? { type: 'preview', id: nextBlockId(), url } : null
    }
    default:
      return null
  }
}

function extractField(lines: string[], field: string): string | null {
  const line = lines.find((l) => l.startsWith(`${field}:`))
  return line ? line.slice(field.length + 1).trim() : null
}
