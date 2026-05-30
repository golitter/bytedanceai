import { describe, expect, it } from 'vitest'

import { reduceEventToBlocks } from '../block-reducer'

describe('reduceEventToBlocks', () => {
  it('returns text block for plain text', () => {
    const result = reduceEventToBlocks('Hello world')
    expect(result).toHaveLength(1)
    expect(result[0].type).toBe('text')
    if (result[0].type === 'text') expect(result[0].content).toBe('Hello world')
    expect(result[0].id).toMatch(/^blk-/)
  })

  it('parses html-render block', () => {
    const input = 'Before\n```aka_yhy\ntype: html-render\n<div>Hello</div>\n```\nAfter'
    const result = reduceEventToBlocks(input)
    expect(result).toHaveLength(3)
    expect(result[0].type).toBe('text')
    if (result[0].type === 'text') expect(result[0].content).toBe('Before\n')
    expect(result[1].type).toBe('html-render')
    if (result[1].type === 'html-render') expect(result[1].content).toBe('<div>Hello</div>')
    expect(result[2].type).toBe('text')
    if (result[2].type === 'text') expect(result[2].content).toBe('\nAfter')
  })

  it('parses image block', () => {
    const input = '```aka_yhy\ntype: image\npath: chart.png\n```'
    const result = reduceEventToBlocks(input)
    expect(result).toHaveLength(1)
    expect(result[0].type).toBe('image')
    if (result[0].type === 'image') expect(result[0].path).toBe('chart.png')
  })

  it('parses attachment block', () => {
    const input = '```aka_yhy\ntype: attachment\npath: report.pdf\n```'
    const result = reduceEventToBlocks(input)
    expect(result).toHaveLength(1)
    expect(result[0].type).toBe('attachment')
    if (result[0].type === 'attachment') expect(result[0].path).toBe('report.pdf')
  })

  it('parses diff block', () => {
    const input = '```aka_yhy\ntype: diff\nsnapshotId: snap-123\n```'
    const result = reduceEventToBlocks(input)
    expect(result).toHaveLength(1)
    expect(result[0].type).toBe('diff')
    if (result[0].type === 'diff') expect(result[0].snapshotId).toBe('snap-123')
  })

  it('parses preview block', () => {
    const input = '```aka_yhy\ntype: preview\nurl: http://localhost:3928/index.html\n```'
    const result = reduceEventToBlocks(input)
    expect(result).toHaveLength(1)
    expect(result[0].type).toBe('preview')
    if (result[0].type === 'preview') expect(result[0].url).toBe('http://localhost:3928/index.html')
  })

  it('degrades unknown type to text block', () => {
    const input = '```aka_yhy\ntype: unknown\n```'
    const result = reduceEventToBlocks(input)
    expect(result).toHaveLength(1)
    expect(result[0].type).toBe('text')
    if (result[0].type === 'text') expect(result[0].content).toBe('```aka_yhy\ntype: unknown\n```')
  })

  it('handles multiple blocks mixed with text', () => {
    const input =
      'Intro\n```aka_yhy\ntype: diff\nsnapshotId: snap-1\n```\nMiddle\n```aka_yhy\ntype: preview\nurl: http://localhost:3000\n```\nEnd'
    const result = reduceEventToBlocks(input)
    expect(result).toHaveLength(5)
    expect(result.map((b) => b.type)).toEqual(['text', 'diff', 'text', 'preview', 'text'])
  })

  it('assigns unique ids to each block', () => {
    const input =
      'A\n```aka_yhy\ntype: diff\n```\nB\n```aka_yhy\ntype: preview\nurl: http://x\n```\nC'
    const result = reduceEventToBlocks(input)
    const ids = result.map((b) => b.id)
    expect(new Set(ids).size).toBe(ids.length)
  })

  it('folds legacy bare runtime status lines', () => {
    const input =
      'type: runtime_status\n' +
      'json: {"agent":"执行者","status":"running","streamingText":"你好","task_id":"task-001"}\n' +
      'type: runtime_status\n' +
      'json: {"agent":"执行者","status":"running","streamingText":"！","task_id":"task-001"}\n' +
      'type: runtime_status\n' +
      'json: {"agent":"执行者","status":"completed","task_id":"task-001"}'
    const result = reduceEventToBlocks(input)

    expect(result).toHaveLength(1)
    expect(result[0].type).toBe('runtime_status')
    if (result[0].type === 'runtime_status') {
      expect(result[0].status).toBe('completed')
      expect(result[0].streamingText).toBe('你好！')
    }
  })

  it('folds legacy bare plan lines', () => {
    const input =
      '说明\n' +
      'type: plan\n' +
      'json: {"overview":"","tasks":[{"agent":"执行者","status":"pending","task_id":"task-001","title":"介绍自己"}]}'
    const result = reduceEventToBlocks(input)

    expect(result.map((block) => block.type)).toEqual(['text', 'plan'])
    expect(result[1].type).toBe('plan')
    if (result[1].type === 'plan') {
      expect(result[1].tasks[0].title).toBe('介绍自己')
    }
  })

  it('folds legacy fenced runtime status lines', () => {
    const input =
      '```yaml\n' +
      'type: runtime_status\n' +
      'json: {"agent":"执行者","status":"running","streamingText":"Cla","task_id":"task-001"}\n' +
      '```\n' +
      '```yaml\n' +
      'type: runtime_status\n' +
      'json: {"agent":"执行者","status":"running","streamingText":"ude","task_id":"task-001"}\n' +
      '```'
    const result = reduceEventToBlocks(input)

    expect(result).toHaveLength(1)
    expect(result[0].type).toBe('runtime_status')
    if (result[0].type === 'runtime_status') {
      expect(result[0].streamingText).toBe('Claude')
    }
  })
})
