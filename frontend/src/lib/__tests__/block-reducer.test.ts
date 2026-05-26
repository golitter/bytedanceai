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
})
