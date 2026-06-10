## ADDED Requirements

### Requirement: Block Reducer parses aka_yhy blocks from markdown
`reduceEventToBlocks(fullText: string) → MessageBlock[]` SHALL parse the full markdown text, detect ` ```aka_yhy ` code blocks, extract the `type` field, and produce typed Block objects. Non-aka_yhy content SHALL be grouped into TextBlock.

#### Scenario: Parse html-render block
- **WHEN** `reduceEventToBlocks` receives text containing ` ```aka_yhy\ntype: html-render\n<div>...</div>\n``` `
- **THEN** SHALL return `[TextBlock, HtmlBlock{content: "<div>...</div>"}]`

#### Scenario: Parse image block
- **WHEN** `reduceEventToBlocks` receives text containing ` ```aka_yhy\ntype: image\npath: chart.png\n``` `
- **THEN** SHALL return `[TextBlock, ImageBlock{path: "chart.png"}]`

#### Scenario: Parse attachment block
- **WHEN** `reduceEventToBlocks` receives text containing ` ```aka_yhy\ntype: attachment\npath: report.pdf\n``` `
- **THEN** SHALL return `[TextBlock, AttachmentBlock{path: "report.pdf"}]`

#### Scenario: Parse diff block
- **WHEN** `reduceEventToBlocks` receives text containing ` ```aka_yhy\ntype: diff\n``` `
- **THEN** SHALL return `[TextBlock, DiffBlock]`

#### Scenario: Parse preview block
- **WHEN** `reduceEventToBlocks` receives text containing ` ```aka_yhy\ntype: preview\nurl: http://localhost:3928/index.html\n``` `
- **THEN** SHALL return `[TextBlock, PreviewBlock{url: "http://localhost:3928/index.html"}]`

#### Scenario: Unknown aka_yhy type degrades gracefully
- **WHEN** `reduceEventToBlocks` receives text containing ` ```aka_yhy\ntype: unknown\n``` `
- **THEN** SHALL return `[TextBlock, TextBlock{content: "```aka_yhy\ntype: unknown\n```"}]`

#### Scenario: Plain text without aka_yhy blocks
- **WHEN** `reduceEventToBlocks` receives plain text with no aka_yhy blocks
- **THEN** SHALL return `[TextBlock{content: fullText}]`

### Requirement: MessageBlock types are discriminated union
MessageBlock SHALL be a TypeScript discriminated union with `type` as the discriminant field. Each block type SHALL have its own fields.

#### Scenario: TypeScript type narrowing works
- **WHEN** a MessageBlock is used in a switch statement on `block.type`
- **THEN** TypeScript SHALL correctly narrow the type to access block-specific fields

### Requirement: TextBlock contains markdown content
TextBlock SHALL have a `content: string` field containing markdown text to be rendered by MarkdownRenderer.

#### Scenario: TextBlock content rendered as markdown
- **WHEN** a TextBlock with `content: "# Hello\n\nWorld"` is rendered
- **THEN** MarkdownRenderer SHALL produce HTML heading + paragraph
